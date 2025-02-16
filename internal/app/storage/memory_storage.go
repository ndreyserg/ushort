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

// Set - сохранение ссылки в хранилище и получение ee ID.
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

// Get получение оригинала ссылки по ее ID.
func (s *memoryStorage) Get(ctx context.Context, key string) (string, error) {
	si, ok := s.byKey[key]

	if !ok {
		return "", errors.New("not found")
	}
	return si.Original, nil
}

// Close закрытие хранилища.
func (s *memoryStorage) Close() error {
	return nil
}

// Check проверка доступности БД.
func (s *memoryStorage) Check(ctx context.Context) error {
	return errors.New("memory storage has no db")
}

// SetBatch сохранение массива ссылок хранилище и получение их ID.
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

// GetUserUrls получения сохраненных ссылок по пользователю.
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

// DeleteUserData удаление ссылок по пользователя.
func (s *memoryStorage) DeleteUserData(ctx context.Context, ids []string, userID string) error {
	return nil
}

// NewMemoryStorage создание inmemory хранилища
func NewMemoryStorage() Storage {
	byKey := map[string]StorageItem{}
	byVal := map[string]StorageItem{}

	return &memoryStorage{
		mt:    &sync.Mutex{},
		byKey: byKey,
		byVal: byVal,
	}
}
