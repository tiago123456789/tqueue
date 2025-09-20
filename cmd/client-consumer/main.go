package main

import (
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

	consumer := client.NewConsumer(&client.ConsumerOptions{
		Address:  "localhost:8080",
		User:     os.Getenv("USER_ADMIN"),
		Password: os.Getenv("PASSWORD"),
		Queue:    "t",
		Handler: func(message client.Message) error {
			log.Println("message: ", message.Message)
			log.Println("id: ", message.Id)
			time.Sleep(5 * time.Second)
			return nil
		},
	})
	_, err = consumer.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// defer consumer.Disconnect()

	consumer.Start()
}
