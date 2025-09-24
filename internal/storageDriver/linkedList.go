package storageDriver

import "github.com/tiago123456789/tqueue/pkg/types"

type Node struct {
	Value interface{}
	Next  *Node
}

type LinkedList struct {
	Head *Node
	Tail *Node
	Size int
}

func (l *LinkedList) Append(value interface{}) {
	if l.Head == nil {
		head := Node{
			Value: value,
			Next:  nil,
		}
		l.Head = &head
		l.Tail = &head
		l.Size++
		return
	}

	if l.Head.Next == nil {
		item := Node{
			Value: value,
			Next:  nil,
		}
		l.Head.Next = &item
		l.Tail = &item
		l.Size++
		return
	}

	oldTail := l.Tail
	l.Tail = &Node{
		Value: value,
		Next:  nil,
	}
	oldTail.Next = l.Tail
	l.Size++
}

func (l *LinkedList) GetById(queueItem types.QueueItem) *Node {
	if l.Head == nil {
		return nil
	}

	start := l.Head
	previous := l.Head
	for start != nil && start.Value != nil {
		if item, ok := start.Value.(types.QueueItem); ok && item.Id == queueItem.Id {
			if start == l.Head {
				l.Head = start.Next
				l.Size--
				return start
			}

			if start == l.Tail {
				l.Tail = previous
				previous.Next = nil
				l.Size--
				return start
			}
			previous.Next = start.Next
			l.Size--
			return start
		}
		previous = start
		start = start.Next
	}

	return nil
}

func (l *LinkedList) Get() *Node {
	if l.Head == nil {
		return nil
	}

	current := l.Head

	if current.Value == nil {
		return nil
	}

	newHead := current.Next
	l.Head = newHead
	l.Size--
	return current
}
