package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/models"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func MakePostJSONHandler(s storage.Storage, baseURL string, session auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Request

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&req); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		userID, err := session.Open(w, r)

		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")

		urlID, err := s.Set(r.Context(), req.URL, userID)

		if err != nil && !errors.Is(err, storage.ErrConflict) {
			logger.Log.Error(err)
			http.Error(w, "", http.StatusBadRequest)
			return

		}

		if err != nil && errors.Is(err, storage.ErrConflict) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		resp := models.Response{
			Result: fmt.Sprintf("%s/%s", baseURL, urlID),
		}

		enc := json.NewEncoder(w)

		if err := enc.Encode(resp); err != nil {
			return
		}
		logger.Log.Info("success")
	}
}
