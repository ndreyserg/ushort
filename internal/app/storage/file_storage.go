package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ndreyserg/ushort/internal/app/models"
)

type fileStorage struct {
	mt      *sync.Mutex
	byKey   map[string]StorageItem
	byVal   map[string]StorageItem
	file    *os.File
	encoder *json.Encoder
}

func (s *fileStorage) Set(ctx context.Context, val string, userID string) (string, error) {
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

	err := s.encoder.Encode(&si)

	if err != nil {
		return "", err
	}

	s.byVal[si.Original] = si
	s.byKey[si.Short] = si
	return si.Short, nil
}

func (s *fileStorage) Get(ctx context.Context, key string) (string, error) {
	si, ok := s.byKey[key]
	if !ok {
		return "", errors.New("")
	}
	return si.Original, nil
}

func (s *fileStorage) Close() error {
	return s.file.Close()
}

func (s *fileStorage) Check(ctx context.Context) error {
	return errors.New("file storage has no db")
}

func (s *fileStorage) SetBatch(ctx context.Context, batch models.BatchRequest, userID string) (models.BatchResult, error) {
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

func (s *fileStorage) GetUserUrls(ctx context.Context, userID string) ([]StorageItem, error) {
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

func NewFileStorage(filepath string) (Storage, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	byKey := map[string]StorageItem{}
	byVal := map[string]StorageItem{}

	for decoder.More() {
		si := StorageItem{}
		err := decoder.Decode(&si)
		if err != nil {
			return nil, err
		}
		byKey[si.Short] = si
		byVal[si.Original] = si
	}

	return &fileStorage{
		mt:      &sync.Mutex{},
		file:    file,
		encoder: json.NewEncoder(file),
		byKey:   byKey,
		byVal:   byVal,
	}, nil
}
