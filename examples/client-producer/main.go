package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tiago123456789/tqueue/pkg/client"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	producer := client.NewProducer(&client.ProducerOptions{
		Address:  "localhost:8080",
		User:     os.Getenv("USER_ADMIN"),
		Password: os.Getenv("PASSWORD"),
		Queue:    "test",
	})

	_, err = producer.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Disconnect()

	for i := 0; i < 10; i++ {
		if err := producer.Send(fmt.Sprintf("Message %d", i)); err != nil {
			log.Println(err)
			return
		}
	}

	time.Sleep(5 * time.Second)
}
