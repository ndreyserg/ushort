package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ndreyserg/ushort/internal/app/storage"
)

func makePingHandler(st storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()
		err := st.Check(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
