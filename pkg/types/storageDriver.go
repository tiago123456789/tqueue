package types

type IStorageDriver interface {
	GetByIdFromUnavaible(id string) *QueueItem
	PushToUnavaible(message string)
	Push(message string)
	Pop() *QueueItem
	TotalMessages() int
	RequeueUnavailableMessages()
}
