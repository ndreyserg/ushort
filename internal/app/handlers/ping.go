package handlers

import (
	"context"
	"net/http"
	"time"
)

func makePingHandler(st Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		err := st.Check(ctx)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
