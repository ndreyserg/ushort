package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func makePostHandler(s Repositiry, urlAddr string) http.HandlerFunc {
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
		w.WriteHeader(http.StatusCreated)
		urlID := s.Set(strings.Trim(string(b), " "))
		short := fmt.Sprintf("http://%s/%s", urlAddr, urlID)
		w.Write([]byte(short))
	}
}
