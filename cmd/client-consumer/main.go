package main

import (
	"fmt"
	"log"
	"os"

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
		Handler: func(message string) error {
			fmt.Println("Handler => Message received:", message)
			return nil
		},
	})
	_, err = consumer.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Disconnect()

	consumer.Start()
}
