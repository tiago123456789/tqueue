package queue

import (
	"errors"
	"log"
	"sync"
)

type IQueueManager interface {
	CreateQueue(queueName string)
	Push(queueName string, message string) error
	Pop(queueName string) (string, error)
	GetQueue(queueName string) (*Queue, error)
	GetQueueProducerConnected(connectionaAddress string) (string, error)
	SetQueueProducerConnected(connectionaAddress string, queueName string)
	RemoveQueueProducerConnected(connectionaAddress string)

	GetQueueConsumerConnected(connectionaAddress string) (string, error)
	SetQueueConsumerConnected(connectionaAddress string, queueName string)
	RemoveQueueConsumerConnected(connectionaAddress string)
}

type QueueManager struct {
	queues        map[string]*Queue
	producerQueue map[string]string
	consumerQueue map[string]string
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues:        make(map[string]*Queue),
		producerQueue: make(map[string]string),
		consumerQueue: make(map[string]string),
	}
}

func (q *QueueManager) RemoveQueueProducerConnected(connectionaAddress string) {
	delete(q.producerQueue, connectionaAddress)
}
func (q *QueueManager) SetQueueProducerConnected(connectionaAddress string, queueName string) {
	q.producerQueue[connectionaAddress] = queueName
}

func (q *QueueManager) GetQueueProducerConnected(connectionaAddress string) (string, error) {
	if q.producerQueue[connectionaAddress] != "" {
		return q.producerQueue[connectionaAddress], nil
	}
	return "", errors.New("Queue not found")
}

func (q *QueueManager) RemoveQueueConsumerConnected(connectionaAddress string) {
	delete(q.consumerQueue, connectionaAddress)
}
func (q *QueueManager) SetQueueConsumerConnected(connectionaAddress string, queueName string) {
	q.consumerQueue[connectionaAddress] = queueName
}

func (q *QueueManager) GetQueueConsumerConnected(connectionaAddress string) (string, error) {
	if q.consumerQueue[connectionaAddress] != "" {
		return q.consumerQueue[connectionaAddress], nil
	}
	return "", errors.New("Queue not found")
}
func (q *QueueManager) CreateQueue(queueName string) {
	if q.queues[queueName] == nil {
		q.queues[queueName] = &Queue{
			mu:       sync.Mutex{},
			messages: []string{},
		}
	}
}

func (q *QueueManager) Push(queueName string, message string) error {
	if q.queues[queueName] != nil {
		q.queues[queueName].Push(message)
		log.Println("Message added, so total messages are:", q.queues[queueName].TotalMessages())
		return nil
	}

	return errors.New("Queue not found")
}

func (q *QueueManager) Pop(queueName string) (string, error) {
	if q.queues[queueName] != nil {
		return q.queues[queueName].Pop(), nil
	}
	return "", errors.New("Queue not found")
}

func (q *QueueManager) GetQueue(queueName string) (*Queue, error) {
	if q.queues[queueName] != nil {
		return q.queues[queueName], nil
	}
	return nil, errors.New("Queue not found")
}
