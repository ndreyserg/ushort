package main

import (
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/config"
	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/queue"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func main() {
	cfg := config.MakeConfig()
	session := auth.NewJWTSession(cfg.Secret)
	err := logger.Initialize(cfg.LogLevel)

	if err != nil {
		panic(err)
	}

	st, err := storage.NewStorage(cfg.DSN, cfg.StoragePath)

	if err != nil {
		panic(err)
	}
	defer st.Close()

	queue := queue.NewQueue(st)
	go queue.Listen()

	err = http.ListenAndServe(cfg.ServerAddr, handlers.MakeRouter(st, cfg.BaseURL, session, queue))

	if err != nil {
		panic(err)
	}
}
