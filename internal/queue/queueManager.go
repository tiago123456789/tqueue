package queue

import (
	"errors"
	"log"
	"time"

	"github.com/tiago123456789/tqueue/pkg/types"
)

type IQueueManager interface {
	CreateQueue(queueName string)
	Push(queueName string, message string) error
	Pop(queueName string) (*types.QueueItem, error)
	GetQueue(queueName string) (*Queue, error)
	GetQueueProducerConnected(connectionaAddress string) (string, error)
	SetQueueProducerConnected(connectionaAddress string, queueName string)
	RemoveQueueProducerConnected(connectionaAddress string)

	GetQueueConsumerConnected(connectionaAddress string) (string, error)
	SetQueueConsumerConnected(connectionaAddress string, queueName string)
	RemoveQueueConsumerConnected(connectionaAddress string)
	RequeueUnavailableMessages()
	RemoveAvailableMessageById(queueName string, id string)
}

type QueueManager struct {
	storageMode   string
	queues        map[string]*Queue
	producerQueue map[string]string
	consumerQueue map[string]string
}

func NewQueueManager(storageMode string) *QueueManager {
	return &QueueManager{
		storageMode:   storageMode,
		queues:        make(map[string]*Queue),
		producerQueue: make(map[string]string),
		consumerQueue: make(map[string]string),
	}
}

func (q *QueueManager) RemoveAvailableMessageById(queueName string, id string) {
	if q.queues[queueName] != nil {
		(*q.queues[queueName].storageDriver).GetByIdFromUnavaible(id)
	}
}

func (q *QueueManager) RequeueUnavailableMessages() {
	log.Println("Started requeue messages in queue")
	for {
		for _, queue := range q.queues {
			time.Sleep(time.Second * 1)
			queue.RequeueUnavailableMessages()
		}
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
		var storageDriver IStorageDriver
		if q.storageMode == "inmemory" {
			storageDriver = NewInMemoryStorageDriver()
		}
		q.queues[queueName] = &Queue{
			storageDriver: &storageDriver,
		}
	}
}

func (q *QueueManager) Push(queueName string, message string) error {
	if q.queues[queueName] != nil {
		q.queues[queueName].Push(message)
		log.Println("Message added, so total messages are:", (*q.queues[queueName].storageDriver).TotalMessages())
		return nil
	}

	return errors.New("Queue not found")
}

func (q *QueueManager) Pop(queueName string) (*types.QueueItem, error) {
	if q.queues[queueName] != nil {
		message := q.queues[queueName].Pop()
		message.AvailableAt = time.Now().Add(time.Second * 30)
		(*q.queues[queueName].storageDriver).PushToUnavaible(message.Message)
		return message, nil
	}
	return nil, errors.New("Queue not found")
}

func (q *QueueManager) GetQueue(queueName string) (*Queue, error) {
	if q.queues[queueName] != nil {
		return q.queues[queueName], nil
	}
	return nil, errors.New("Queue not found")
}
