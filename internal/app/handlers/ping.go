package handlers

import (
	"context"
	"net/http"
	"time"
)

func makePingHandler(db DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		err := db.PingContext(ctx)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
