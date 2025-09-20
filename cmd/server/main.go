package main

import (
	"log"

	"github.com/tiago123456789/tqueue/internal/publishEngine"
	"github.com/tiago123456789/tqueue/internal/queue"
	"github.com/tiago123456789/tqueue/internal/tcp"

	"github.com/joho/godotenv"
)

var queueManager = queue.NewQueueManager()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	tcpManager := tcp.NewTcpManager(queueManager, publishEngine.PublishEngine)
	tcpManager.StartServer()

	// ll := queue.LinkedList{}

	// ll.Append(types.QueueItem{Id: "a", Message: "a", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "b", Message: "b", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "c", Message: "c", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "d", Message: "d", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "e", Message: "e", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "f", Message: "f", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "g", Message: "g", AvailableAt: time.Now()})
	// ll.Append(types.QueueItem{Id: "h", Message: "h", AvailableAt: time.Now()})

	// log.Println(ll.GetById(types.QueueItem{Id: "g", Message: "g", AvailableAt: time.Now()}))
	// log.Println(ll.GetById(types.QueueItem{Id: "i", Message: "i", AvailableAt: time.Now()}))

	// log.Println("@@@@@@@@@@@@@@@@@@@@@")
	// log.Println(ll.Head)
	// log.Println(ll.Tail)
	// log.Println(ll.Tail.Next)
	// log.Println(ll.Size)
}
