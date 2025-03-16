package main

import (
	"fmt"
	"net/http"

	"github.com/ndreyserg/ushort/internal/app/auth"
	"github.com/ndreyserg/ushort/internal/app/certs"
	"github.com/ndreyserg/ushort/internal/app/config"
	"github.com/ndreyserg/ushort/internal/app/handlers"
	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/queue"
	"github.com/ndreyserg/ushort/internal/app/storage"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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

	handler := handlers.MakeRouter(st, cfg.BaseURL, session, queue)

	if cfg.EnableHTTPS {
		logger.Log.Info("start https server")
		certPair, err := certs.GetCertFiles()
		if err != nil {
			panic(err)
		}
		err = http.ListenAndServeTLS(cfg.ServerAddr, certPair.CertPem, certPair.PrivateKeyPEM, handler)
		if err != nil {
			panic(err)
		}

	} else {
		logger.Log.Info("start http server")
		err = http.ListenAndServe(cfg.ServerAddr, handler)
		if err != nil {
			panic(err)
		}
	}
}
