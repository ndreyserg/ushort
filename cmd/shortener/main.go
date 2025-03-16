package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	cfg := config.MakeConfig()
	session := auth.NewJWTSession(cfg.Secret)
	err := logger.Initialize(cfg.LogLevel)

	if err != nil {
		return err
	}

	st, err := storage.NewStorage(cfg.DSN, cfg.StoragePath)

	if err != nil {
		return err
	}
	defer st.Close()

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	queue := queue.NewQueue(st)
	go queue.Listen()
	handler := handlers.MakeRouter(st, cfg.BaseURL, session, queue)

	var srv = http.Server{Addr: cfg.ServerAddr, Handler: handler}

	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Log.Info("HTTP server Shutdown: %v", err)
		}
		queue.Wait()
		close(idleConnsClosed)
	}()

	var serverError error
	if cfg.EnableHTTPS {
		logger.Log.Info("start https server")
		certPair, err := certs.GetCertFiles()
		if err != nil {
			return err
		}
		serverError = srv.ListenAndServeTLS(certPair.CertPem, certPair.PrivateKeyPEM)

	} else {
		serverError = srv.ListenAndServe()
	}

	if serverError != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	return nil
}
