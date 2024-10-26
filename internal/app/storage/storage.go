package storage

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/ndreyserg/ushort/internal/app/logger"
	"github.com/ndreyserg/ushort/internal/app/models"
)

func getUniqKey() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

var ErrConflict = errors.New("url exists")
var ErrIsGone = errors.New("url is gone")

type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	GetUserUrls(ctx context.Context, userID string) ([]StorageItem, error)
	Set(ctx context.Context, val string, userID string) (string, error)
	SetBatch(ctx context.Context, batch models.BatchRequest, userID string) (models.BatchResult, error)
	Check(ctx context.Context) error
	Close() error
	DeleteUserData(ctx context.Context, ids []string, userID string) error
}

type StorageItem struct {
	Original  string `json:"original"`
	Short     string `json:"short"`
	UserID    string `json:"user_id"`
	IsDeleted bool
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
