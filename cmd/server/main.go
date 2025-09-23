package main

import (
	"flag"
	"log"

	"github.com/tiago123456789/tqueue/internal/publishEngine"
	"github.com/tiago123456789/tqueue/internal/queue"
	"github.com/tiago123456789/tqueue/internal/tcp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	storageMode := flag.String("storage", "inmemory", "Storage mode (inmemory or sqlite)")

	flag.Parse()

	if *storageMode != "inmemory" && *storageMode != "sqlite" {
		log.Fatal("Invalid storage mode")
	}

	var queueManager = queue.NewQueueManager(*storageMode)

	tcpManager := tcp.NewTcpManager(
		queueManager, publishEngine.PublishEngine)
	tcpManager.StartServer()

}
