package storage

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type storageItem struct {
	Key   string `json:"short"`
	Value string `json:"original"`
}

func getUniqKey() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

type Storage struct {
	mt      *sync.Mutex
	byKey   map[string]string
	byVal   map[string]string
	file    *os.File
	encoder *json.Encoder
}

func (s *Storage) Set(val string) (string, error) {
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

func (s *Storage) Get(key string) (string, error) {
	val := s.byKey[key]
	if val == "" {
		return "", errors.New("")
	}
	return val, nil
}

func (s *Storage) Close() error {
	return s.file.Close()
}

func NewStorage(filepath string) (*Storage, error) {
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

	return &Storage{
		mt:      &sync.Mutex{},
		file:    file,
		encoder: json.NewEncoder(file),
		byKey:   byKey,
		byVal:   byVal,
	}, nil
}
