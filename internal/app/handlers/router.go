// Package handlers с обработчиками endpoints
package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

// Queue интерфейс очереди для асинхронного удаления
type Queue interface {
	AddTask(IDs []string, userID string)
}

// MakeRouter создает роутер с прописанными endpoint-им и обработчиками
func MakeRouter(s storage.Storage, baseURL string, session auth.Session, q Queue) http.Handler {

	errHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusBadRequest)
	}

	r := chi.NewRouter()

	r.Use(gzipMiddleware)
	r.Use(logger.LoggerMiddleware)
	r.Get("/ping", MakePingHandler(s))
	r.Get("/{id}", MakeGetHandler(s))
	r.Post("/", MakePostHandler(s, baseURL, session))
	r.Post("/api/shorten", MakePostJSONHandler(s, baseURL, session))
	r.Post("/api/shorten/batch", MakePostBatchHandler(s, baseURL, session))
	r.Get("/api/user/urls", MakeGetUserUrlsHandler(s, baseURL, session))
	r.Delete("/api/user/urls", MakeDeleteHandler(q, baseURL, session))
	r.MethodNotAllowed(errHandler)
	r.NotFound(errHandler)

	return r

}
