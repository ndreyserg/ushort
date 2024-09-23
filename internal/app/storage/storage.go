package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
)

func getUniqKey() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

type Storage struct {
	byKey map[string]string
	byVal map[string]string
}

func (s Storage) Set(val string) string {
	if s.byVal[val] != "" {
		return s.byVal[val]
	}
	key := getUniqKey()
	s.byVal[val] = key
	s.byKey[key] = val
	return key
}

func (s *Storage) Get(key string) (string, error) {
	val := s.byKey[key]
	if val == "" {
		return "", errors.New("")
	}
	return val, nil
}

func NewStorage() *Storage {
	return &Storage{
		byKey: map[string]string{},
		byVal: map[string]string{},
	}
}
