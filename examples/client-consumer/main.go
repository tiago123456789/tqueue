package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		Queue:    "test",
		Handler: func(message client.Message) error {
			body, err := json.Marshal(message)
			if err != nil {
				return err
			}

			fmt.Println("Message received:", string(body))

			resp, err := http.Post("https://boring-cricket-53.webhook.cool", "application/json", bytes.NewBuffer(body))
			defer resp.Body.Close()
			if err != nil {
				return err
			}
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
