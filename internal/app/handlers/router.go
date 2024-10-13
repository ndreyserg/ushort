package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func MakeRouter(s storage.Storage, baseURL string) http.Handler {

	errHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}

	r := chi.NewRouter()
	r.Use(gzipMiddleware)
	r.Use(logger.LoggerMiddleware)
	r.Get("/ping", makePingHandler(s))
	r.Get("/{id}", makeGetHandler(s))
	r.Post("/", makePostHandler(s, baseURL))
	r.Post("/api/shorten", MakePostJSONHandler(s, baseURL))
	r.Post("/api/shorten/batch", MakePostBatchHandler(s, baseURL))
	r.MethodNotAllowed(errHandler)
	r.NotFound(errHandler)
	return r

}
