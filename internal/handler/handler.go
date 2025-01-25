package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"url-shortener/internal/service"

	"github.com/gorilla/mux"
)

type ShortenRequest struct {
	URL string `json:"url"`
}
type Handler struct {
	service *service.URLService
}

func NewHandler(s *service.URLService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error:problem with body")
		http.Error(w, "missing body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req ShortenRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		http.Error(w, "invalid link 1", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.ShortenURL(&req.URL)
	if err != nil {
		http.Error(w, err.Error()+shortURL, http.StatusBadRequest)
		return
	}
	if shortURL == "" {
		http.Error(w, "invalid link 3", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": %s}`, shortURL)

}

func (h *Handler) RedirectURL(w http.ResponseWriter, r *http.Request) {

	var req ShortenRequest
	url := r.URL.Path
	req.URL = url
	if req.URL == "" {
		http.Error(w, "invalid url", http.StatusBadRequest)
	}
	originalURL, err := h.service.Redirect(&req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)

}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/shorten", h.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortened}", h.RedirectURL).Methods("GET")
}
