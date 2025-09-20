package queue

import "sync"

type IQueue interface {
	Push(message string)
	Pop() string
	TotalMessages() int
}

type Queue struct {
	mu       sync.Mutex
	messages []string
}

func (q *Queue) TotalMessages() int {
	return len(q.messages)
}

func (q *Queue) Push(message string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.messages = append(q.messages, message)
}

func (q *Queue) Pop() string {
	if q.TotalMessages() == 0 {
		return ""
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	message := q.messages[0]
	q.messages = q.messages[1:]
	return message
}
