package services

import (
	"context"
)

type Queue interface {
	Add(data string)
}

const queueBufCapacity = 16

type QueueService struct {
	dataReceiver chan string
	queue        Queue
}

func NewQueueService(queue Queue) *QueueService {
	return &QueueService{
		dataReceiver: make(chan string, queueBufCapacity),
		queue:        queue,
	}
}

func (s *QueueService) Add(data []string) {
	for _, v := range data {
		s.dataReceiver <- v
	}
}

func (s *QueueService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-s.dataReceiver:
			s.queue.Add(data)
		}
	}
}
