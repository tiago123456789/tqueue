package queue

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tiago123456789/tqueue/pkg/types"
)

type IQueue interface {
	Push(message string)
	Pop() types.QueueItem
	TotalMessages() int
	RequeueUnavailableMessages()
}

type Queue struct {
	mu                sync.Mutex
	messagesUnavaible *LinkedList
	messages          *LinkedList
}

func (q *Queue) RequeueUnavailableMessages() {
	for q.messagesUnavaible.Size > 0 {
		q.mu.Lock()
		message := q.messagesUnavaible.Get()
		if message != nil && message.Value.(types.QueueItem).AvailableAt.Before(time.Now()) {
			q.messages.Append(message.Value.(types.QueueItem))
			q.messagesUnavaible.GetById(message.Value.(types.QueueItem))
		} else {
			q.messagesUnavaible.Append(message.Value.(types.QueueItem))
		}
		q.mu.Unlock()
	}
}
func (q *Queue) TotalMessages() int {
	return q.messages.Size
}

func (q *Queue) Push(message string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages.Append(types.QueueItem{Message: message, Id: uuid.New().String()})
}

func (q *Queue) Pop() types.QueueItem {
	if q.TotalMessages() == 0 {
		return types.QueueItem{}
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	message := q.messages.Get()
	log.Println("Message: ", message)
	if queueItem, ok := message.Value.(types.QueueItem); ok {
		queueItem.AvailableAt = time.Now().Add(time.Second * 30)
		q.messagesUnavaible.Append(queueItem)
		return queueItem
	}
	return types.QueueItem{}
}
