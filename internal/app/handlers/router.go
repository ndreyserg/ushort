package handlers

import (
	"net/http"
)

type Repositiry interface {
	Get(key string) (string, error)
	Set(val string) string
}

func MakeRouter(s Repositiry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			get(s, w, r)
		case http.MethodPost:
			post(s, w, r)
		default:
			http.Error(w, "", http.StatusBadRequest)
		}
	}
}
