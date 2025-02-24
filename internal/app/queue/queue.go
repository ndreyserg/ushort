// Package queue с реализацией простой очереди для асинхронного удаления
package queue

import (
	"context"

	"github.com/ndreyserg/ushort/internal/app/logger"
)

type storage interface {
	DeleteUserData(ctx context.Context, ids []string, userID string) error
}

type deleteTask struct {
	IDs    []string
	UserID string
}

// Queue структура очереди
type Queue struct {
	s  storage
	ch chan deleteTask
}

// AddTask добавляет задачу в очередь
func (q *Queue) AddTask(IDs []string, userID string) {
	q.ch <- deleteTask{IDs: IDs, UserID: userID}
}

// Listen запускает очередь в обработку
func (q *Queue) Listen() {
	for t := range q.ch {
		go func(t deleteTask) {
			err := q.s.DeleteUserData(context.Background(), t.IDs, t.UserID)
			if err != nil {
				logger.Log.Error(err)
			}
		}(t)
	}
}

// NewQueue возвращает объект очереди
func NewQueue(s storage) *Queue {
	return &Queue{
		s:  s,
		ch: make(chan deleteTask, 1),
	}
}
