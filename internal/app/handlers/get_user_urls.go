package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

type responseItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func makeGetUserUrlsHandler(s storage.Storage, baseURL string, session auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := session.GetID(r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		urls, err := s.GetUserUrls(r.Context(), userID)

		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		res := make([]responseItem, 0, len(urls))

		for _, i := range urls {
			res = append(res, responseItem{
				ShortURL:    fmt.Sprintf("%s/%s", baseURL, i.Short),
				OriginalURL: i.Original,
			})
		}
		w.Header().Set("content-type", "application/json")
		enc := json.NewEncoder(w)

		if err := enc.Encode(res); err != nil {
			return
		}
	}
}
