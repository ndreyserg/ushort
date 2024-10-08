package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func makeGetHandler(s Repositiry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "id")
		if key == "" {
			http.Error(w, "empty key", http.StatusBadRequest)
			return
		}
		url, err := s.Get(key)
		if err != nil {
			http.Error(w, "key not found", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
