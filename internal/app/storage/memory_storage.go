package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/ndreyserg/ushort/internal/app/models"
)

type memoryStorage struct {
	mt    *sync.Mutex
	byKey map[string]string
	byVal map[string]string
}

func (s *memoryStorage) Set(ctx context.Context, val string) (string, error) {
	s.mt.Lock()
	defer s.mt.Unlock()
	if s.byVal[val] != "" {
		return s.byVal[val], nil
	}
	key := getUniqKey()
	s.byVal[val] = key
	s.byKey[key] = val
	return key, nil
}

func (s *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	val := s.byKey[key]
	if val == "" {
		return "", errors.New("")
	}
	return val, nil
}

func (s *memoryStorage) Close() error {
	return nil
}

func (s *memoryStorage) Check(ctx context.Context) error {
	return errors.New("memory storage has no db")
}

func (s *memoryStorage) SetBatch(ctx context.Context, batch models.BatchRequest) (models.BatchResult, error) {
	result := make(models.BatchResult, 0, len(batch))

	for _, item := range batch {
		short, err := s.Set(ctx, item.Original)

		if err != nil {
			return nil, err
		}
		resultItem := models.BatchResultItem{
			ID:    item.ID,
			Short: short,
		}

		result = append(result, resultItem)
	}
	return result, nil
}

func NewMemoryStorage() Storage {
	byKey := map[string]string{}
	byVal := map[string]string{}

	return &memoryStorage{
		mt:    &sync.Mutex{},
		byKey: byKey,
		byVal: byVal,
	}
}
