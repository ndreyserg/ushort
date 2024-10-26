package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/ndreyserg/ushort/internal/app/models"
)

type memoryStorage struct {
	mt    *sync.Mutex
	byKey map[string]StorageItem
	byVal map[string]StorageItem
}

func (s *memoryStorage) Set(ctx context.Context, val string, userID string) (string, error) {
	s.mt.Lock()
	defer s.mt.Unlock()

	si, ok := s.byVal[val]

	if ok {
		return si.Short, nil
	}
	si = StorageItem{
		Short:    getUniqKey(),
		Original: val,
		UserID:   userID,
	}
	s.byVal[si.Original] = si
	s.byKey[si.Short] = si
	return si.Short, nil
}

func (s *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	si, ok := s.byKey[key]

	if !ok {
		return "", errors.New("not found")
	}
	return si.Original, nil
}

func (s *memoryStorage) Close() error {
	return nil
}

func (s *memoryStorage) Check(ctx context.Context) error {
	return errors.New("memory storage has no db")
}

func (s *memoryStorage) SetBatch(ctx context.Context, batch models.BatchRequest, userID string) (models.BatchResult, error) {
	result := make(models.BatchResult, 0, len(batch))

	for _, item := range batch {
		short, err := s.Set(ctx, item.Original, userID)

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

func (s *memoryStorage) GetUserUrls(ctx context.Context, userID string) ([]StorageItem, error) {
	s.mt.Lock()
	defer s.mt.Unlock()
	res := []StorageItem{}

	for _, si := range s.byKey {
		if si.UserID == userID {
			res = append(res, si)
		}
	}
	return res, nil
}

func (s *memoryStorage) DeleteUserData(ctx context.Context, ids []string, userID string) error {
	return nil
}

func NewMemoryStorage() Storage {
	byKey := map[string]StorageItem{}
	byVal := map[string]StorageItem{}

	return &memoryStorage{
		mt:    &sync.Mutex{},
		byKey: byKey,
		byVal: byVal,
	}
}
