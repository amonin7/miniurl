package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type HttpHandler struct {
	storageMu sync.RWMutex
	storage   map[string]string
}

func getRandomKey() string {
	alphaBet := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Shuffle(len(alphaBet), func(i, j int) {
		alphaBet[i], alphaBet[j] = alphaBet[j], alphaBet[i]
	})
	id := string(alphaBet[:5])
	return id
}

type PutRequestData struct {
	Url string `json:"url"`
}

type PutResponseData struct {
	Key string `json:"key"`
}

func (h *HttpHandler) handlePutUrl(w http.ResponseWriter, r *http.Request) {
	var data PutRequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUrlKey := getRandomKey()

	h.storageMu.Lock()
	h.storage[newUrlKey] = data.Url
	h.storageMu.Unlock()

	response := PutResponseData{
		Key: newUrlKey,
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

func (h *HttpHandler) handleGetUrl(w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
	h.storageMu.RLock()
	url, found := h.storage[key]
	h.storageMu.RUnlock()

	if !found {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from Server!"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
}

func NewServer() *http.Server {
	r := mux.NewRouter()

	handler := &HttpHandler{
		storage: make(map[string]string),
	}

	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/{shorturl:\\w{5}}", handler.handleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.handlePutUrl).Methods(http.MethodPost)

	return &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func main() {
	srv := NewServer()
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
