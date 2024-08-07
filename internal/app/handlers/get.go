package handlers

import (
	"net/http"
	"strings"
)

func get(s Repositiry, w http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")
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
