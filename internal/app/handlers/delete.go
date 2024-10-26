package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/logger"
)

func MakeDeleteHandler(q Queue, baseURL string, session auth.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var ids []string

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&ids); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		userID, err := session.Open(w, r)

		if err != nil {
			logger.Log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		q.AddTask(ids, userID)
		w.WriteHeader(http.StatusAccepted)
	}
}
