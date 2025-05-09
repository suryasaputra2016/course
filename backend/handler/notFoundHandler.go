package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type NotFoundHandler struct{}

func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (nfh NotFoundHandler) PageNotFound(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{"message": "page not found"})
	if err != nil {
		log.Printf("encoding message: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (nfh NotFoundHandler) Home(w http.ResponseWriter, r *http.Request) {
	home := struct {
		Title string
		Body  string
	}{
		Title: "Welome Home",
		Body:  "This is the body of home page.",
	}
	err := json.NewEncoder(w).Encode(home)
	if err != nil {
		log.Printf("encoding message: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
