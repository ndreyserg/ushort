package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ndreyserg/ushort/internal/app/logger"
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
	r.Use(logger.LoggerMiddleware)
	r.Get("/{id}", makeGetHandler(s))
	r.Post("/", makePostHandler(s, baseURL))
	r.Post("/api/shorten", MakePostJSONHandler(s, baseURL))
	r.MethodNotAllowed(errHandler)
	r.NotFound(errHandler)
	return r

}
