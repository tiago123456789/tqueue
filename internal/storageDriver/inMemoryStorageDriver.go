package storageDriver

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tiago123456789/tqueue/pkg/types"
)

type InMemoryStorageDriver struct {
	mu                sync.Mutex
	messagesUnavaible *LinkedList
	messages          *LinkedList
}

func NewInMemoryStorageDriver() *InMemoryStorageDriver {
	return &InMemoryStorageDriver{
		messagesUnavaible: &LinkedList{},
		messages:          &LinkedList{},
	}
}

func (i *InMemoryStorageDriver) GetByIdFromUnavaible(id string) *types.QueueItem {
	i.mu.Lock()
	defer i.mu.Unlock()
	queueItem := i.messagesUnavaible.GetById(types.QueueItem{Id: id}).Value.(types.QueueItem)
	return &queueItem
}

func (i *InMemoryStorageDriver) PushToUnavaible(message string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.messagesUnavaible.Append(types.QueueItem{Message: message, Id: uuid.New().String()})
}

func (i *InMemoryStorageDriver) Push(message string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.messages.Append(types.QueueItem{Message: message, Id: uuid.New().String()})
}

func (i *InMemoryStorageDriver) Pop() *types.QueueItem {
	if i.messages.Size == 0 {
		return nil
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	message := i.messages.Get()
	if message != nil {
		queueItem, ok := message.Value.(types.QueueItem)
		if ok {
			queueItem.AvailableAt = time.Now().Add(time.Second * 30)
			i.messagesUnavaible.Append(queueItem)
			return &queueItem
		}
	}
	return nil
}

func (i *InMemoryStorageDriver) TotalMessages() int {
	return i.messages.Size
}

func (i *InMemoryStorageDriver) RequeueUnavailableMessages() {
	for i.messagesUnavaible.Size > 0 {
		i.mu.Lock()
		message := i.messagesUnavaible.Get()
		if message != nil && message.Value.(types.QueueItem).AvailableAt.Before(time.Now()) {
			i.messages.Append(message.Value.(types.QueueItem))
			i.messagesUnavaible.GetById(message.Value.(types.QueueItem))
		} else if message != nil {
			i.messagesUnavaible.Append(message.Value.(types.QueueItem))
		}
		i.mu.Unlock()
	}
}
