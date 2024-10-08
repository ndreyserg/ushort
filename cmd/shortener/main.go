package main

import (
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/config"
	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func main() {
	cfg := config.MakeConfig()
	st := storage.NewStorage()
	err := http.ListenAndServe(cfg.ServerAddr, handlers.MakeRouter(st, cfg.BaseURL))

	if err != nil {
		panic(err)
	}
}
