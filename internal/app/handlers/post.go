package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func makePostHandler(s storage.Storage, baseURL string, session auth.Session) http.HandlerFunc {
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

		userID, err := session.Open(w, r)
		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		urlID, err := s.Set(
			r.Context(),
			strings.Trim(string(b), " "),
			userID,
		)

		if err != nil {
			if errors.Is(err, storage.ErrConflict) {
				w.WriteHeader(http.StatusConflict)
			} else {
				logger.Log.Error(err)
				http.Error(w, "storage error", http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusCreated)
		}

		short := fmt.Sprintf("%s/%s", baseURL, urlID)
		w.Write([]byte(short))
	}
}
