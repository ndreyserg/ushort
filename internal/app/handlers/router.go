package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Repositiry interface {
	Get(key string) (string, error)
	Set(val string) string
}

func MakeRouter(s Repositiry, baseURL string) http.Handler {
	errHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}

	r := chi.NewRouter()
	r.Get("/{id}", makeGetHandler(s))
	r.Post("/", makePostHandler(s, baseURL))
	r.MethodNotAllowed(errHandler)
	r.NotFound(errHandler)
	return r

}
