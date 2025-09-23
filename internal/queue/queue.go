package queue

import (
	"github.com/tiago123456789/tqueue/pkg/types"
)

type IQueue interface {
	Push(message string)
	Pop() types.QueueItem
	TotalMessages() int
	RequeueUnavailableMessages()
}

type Queue struct {
	storageDriver *IStorageDriver
}

func (q *Queue) RequeueUnavailableMessages() {
	(*q.storageDriver).RequeueUnavailableMessages()
}
func (q *Queue) TotalMessages() int {
	return (*q.storageDriver).TotalMessages()
}

func (q *Queue) Push(message string) {
	(*q.storageDriver).Push(message)
}

func (q *Queue) Pop() *types.QueueItem {
	return (*q.storageDriver).Pop()
}
