package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/models"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func MakePostBatchHandler(s storage.Storage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.BatchRequest

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&req); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if len(req) == 0 {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		res, err := s.SetBatch(r.Context(), req)

		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			logger.Log.Error(err)
			return
		}

		for key := range res {
			res[key].Short = fmt.Sprintf("%s/%s", baseURL, res[key].Short)
		}

		enc := json.NewEncoder(w)
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := enc.Encode(res); err != nil {
			return
		}
		logger.Log.Info("success")
	}
}
