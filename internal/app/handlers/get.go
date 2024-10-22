package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func makeGetHandler(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "id")
		if key == "" {
			http.Error(w, "empty key", http.StatusBadRequest)
			return
		}
		url, err := s.Get(r.Context(), key)
		if err != nil {
			http.Error(w, "key not found", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
