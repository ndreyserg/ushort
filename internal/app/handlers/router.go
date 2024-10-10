package handlers

import (
	"net/http"

	"context"

	"github.com/go-chi/chi/v5"
	"github.com/ndreyserg/ushort/internal/app/logger"
)

type Repositiry interface {
	Get(key string) (string, error)
	Set(val string) (string, error)
}

type DBConnection interface {
	PingContext(ctx context.Context) error
}


func MakeRouter(s Repositiry, baseURL string, db DBConnection) http.Handler {

	errHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}

	r := chi.NewRouter()
	r.Use(gzipMiddleware)
	r.Use(logger.LoggerMiddleware)
	r.Get("/ping", makePingHandler(db))
	r.Get("/{id}", makeGetHandler(s))
	r.Post("/", makePostHandler(s, baseURL))
	r.Post("/api/shorten", MakePostJSONHandler(s, baseURL))
	r.MethodNotAllowed(errHandler)
	r.NotFound(errHandler)
	return r

}
