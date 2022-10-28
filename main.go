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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			return
		case "POST":
			return
	}
	_, err := w.Write([]byte("Hello world\n"))
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Set("Content-Type", "plain/text")
}

type HttpHandler struct {
	storage map[string]string
}

func generateRandomKey() string {
	alphabet := []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789")
	rand.Shuffle(len(alphabet), func(i, j int) {
		alphabet[i], alphabet[j] = alphabet[j], alphabet[i]
	})
	return string(alphabet[:8])
}

type PutRequestData struct {
	Url string `json:"url,omitempty"`
}

type PutResponse struct {
	Key string `json:"key,omitempty"`
}

func (h *HttpHandler) handleCreateUrl(w http.ResponseWriter, r *http.Request) {
	var data PutRequestData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUrlKey := generateRandomKey()
	h.storage[newUrlKey] = data.Url
	response := PutResponse {
		Key: newUrlKey,
	}
	rawResponse, _ := json.Marshal(response)

	_, err = w.Write(rawResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (h *HttpHandler) handleGetUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"]

	url, ok := h.storage[shortUrl]
	if !ok {
		http.Error(w, "No such url", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, 302)
}

func main() {
	r := mux.NewRouter()

	handler := &HttpHandler{
		storage: make(map[string]string),
	}

	r.HandleFunc("/", handleRoot).Methods("GET")
	r.HandleFunc("/{shortUrl:\\w{8}}", handler.handleGetUrl).Methods("GET")
	r.HandleFunc("/api/urls", handler.handleCreateUrl).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr: "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}

	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

