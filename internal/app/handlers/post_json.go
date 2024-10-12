package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/models"
)

func MakePostJSONHandler(s Storage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Requst

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&req); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		urlID, err := s.Set(r.Context(), req.URL)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		resp := models.Response{
			Result: fmt.Sprintf("%s/%s", baseURL, urlID),
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		enc := json.NewEncoder(w)

		if err := enc.Encode(resp); err != nil {
			return
		}
		logger.Log.Info("success")
	}
}
