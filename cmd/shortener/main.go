package main

import (
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/config"
	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func main() {
	cfg := config.MakeConfig()

	err := logger.Initialize(cfg.LogLevel)

	if err != nil {
		panic(err)
	}

	st, err := storage.NewStorage(cfg.DSN, cfg.StoragePath)

	if err != nil {
		panic(err)
	}
	defer st.Close()

	err = http.ListenAndServe(cfg.ServerAddr, handlers.MakeRouter(st, cfg.BaseURL))

	if err != nil {
		panic(err)
	}
}
