package main

import (
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

func main() {
	st := storage.NewStorage()
	err := http.ListenAndServe(":8080", handlers.MakeRouter(st))

	if err != nil {
		panic(err)
	}
}
