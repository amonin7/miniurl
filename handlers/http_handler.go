package handlers

import (
	"encoding/json"
	"fmt"
	"miniurl/ratelimit"
	"miniurl/storage"
	"net/http"
	"strings"
	"time"
)

func NewHttpHandler(
	storage storage.Storage,
	limiterFactory *ratelimit.Factory,
	indexMaintainers []storage.IndexMaintainer,
) *HttpHandler {
	return &HttpHandler{
		Storage:          storage,
		indexMaintainers: indexMaintainers,
		postLimit:        limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		getLimit:         limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
	}
}

type HttpHandler struct {
	Storage          storage.Storage
	indexMaintainers []storage.IndexMaintainer

	postLimit *ratelimit.Limiter
	getLimit  *ratelimit.Limiter
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HttpHandler) HandlePutUrl(w http.ResponseWriter, r *http.Request) {
	canDo, err := h.postLimit.CanDoAt(r.Context(), time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !canDo {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var data PutRequestData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUrlKey, err := h.Storage.PutURL(r.Context(), storage.ShortedURL(data.Url))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := PutResponseData{
		Key: string(newUrlKey),
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(rawResponse)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (h *HttpHandler) HandleGetUrl(w http.ResponseWriter, r *http.Request) {
	canDo, err := h.getLimit.CanDoAt(r.Context(), time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !canDo {
		http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	key := strings.Trim(r.URL.Path, "/")

	url, err := h.Storage.GetURL(r.Context(), storage.URLKey(key))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, string(url), http.StatusPermanentRedirect)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from Server!"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
}

func (h *HttpHandler) CreateIndices(rw http.ResponseWriter, r *http.Request) {
	for _, maintainer := range h.indexMaintainers {
		if err := maintainer.EnsureIndices(r.Context()); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_, _ = rw.Write([]byte("All indices are successfully created"))
}
