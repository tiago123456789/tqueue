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
}
