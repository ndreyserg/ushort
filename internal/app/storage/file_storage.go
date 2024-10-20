package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/ndreyserg/ushort/internal/app/models"
)

type storageItem struct {
	Key   string `json:"short"`
	Value string `json:"original"`
}

type fileStorage struct {
	mt      *sync.Mutex
	byKey   map[string]string
	byVal   map[string]string
	file    *os.File
	encoder *json.Encoder
}

func (s *fileStorage) Set(ctx context.Context, val string) (string, error) {
	s.mt.Lock()
	defer s.mt.Unlock()
	if s.byVal[val] != "" {
		return s.byVal[val], nil
	}
	key := getUniqKey()
	si := storageItem{
		Key:   key,
		Value: val,
	}
	err := s.encoder.Encode(&si)
	if err != nil {
		return "", err
	}

	s.byVal[val] = key
	s.byKey[key] = val
	return key, nil
}

func (s *fileStorage) Get(ctx context.Context, key string) (string, error) {
	val := s.byKey[key]
	if val == "" {
		return "", errors.New("")
	}
	return val, nil
}

func (s *fileStorage) Close() error {
	return s.file.Close()
}

func (s *fileStorage) Check(ctx context.Context) error {
	return errors.New("file storage has no db")
}

func (s *fileStorage) SetBatch(ctx context.Context, batch models.BatchRequest) (models.BatchResult, error) {
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

func NewFileStorage(filepath string) (Storage, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)

	byKey := map[string]string{}
	byVal := map[string]string{}

	for decoder.More() {
		si := &storageItem{}
		err := decoder.Decode(&si)
		if err != nil {
			return nil, err
		}
		byKey[si.Key] = si.Value
		byVal[si.Value] = si.Key
	}

	return &fileStorage{
		mt:      &sync.Mutex{},
		file:    file,
		encoder: json.NewEncoder(file),
		byKey:   byKey,
		byVal:   byVal,
	}, nil
}
