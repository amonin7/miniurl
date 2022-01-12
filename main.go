package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type HttpHandler struct {
	storage map[string]string
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
	h.storage[newUrlKey] = data.Url
	response := PutResponseData{
		Key: newUrlKey,
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(rawResponse)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *HttpHandler) handleGetUrl(w http.ResponseWriter, r *http.Request) {
	var data PutRequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUrlKey := getRandomKey()
	h.storage[newUrlKey] = data.Url
	response := PutResponseData{
		Key: newUrlKey,
	}
	rawResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(rawResponse)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from Server!"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain")
}

func main() {
	r := mux.NewRouter()

	handler := &HttpHandler{
		storage: make(map[string]string),
	}

	r.HandleFunc("/", handleRoot)
	r.HandleFunc("/{shorturl:\\w{5}}", handler.handleGetUrl).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.handlePutUrl).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
