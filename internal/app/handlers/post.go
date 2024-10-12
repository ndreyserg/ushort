package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func makePostHandler(s Storage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if string(b) == "" {
			http.Error(w, "empty request body", http.StatusBadRequest)
			return
		}
		urlID, err := s.Set(r.Context(), strings.Trim(string(b), " "))

		if err != nil {
			http.Error(w, "storage error", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		short := fmt.Sprintf("%s/%s", baseURL, urlID)
		w.Write([]byte(short))
	}
}
