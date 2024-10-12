package storage

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/ndreyserg/ushort/internal/app/logger"
)

func getUniqKey() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, val string) (string, error)
	Check(ctx context.Context) error
	Close() error
}

func NewStorage(dsn, fileName string) (Storage, error) {
	if dsn != "" {
		logger.Log.Info("database storage used")
		return NewDBStorage(dsn)
	}

	if fileName != "" {
		logger.Log.Info("file storage used")
		return NewFileStorage(fileName)
	}
	logger.Log.Info("memory storage used")
	return NewMemoryStorage(), nil
}
