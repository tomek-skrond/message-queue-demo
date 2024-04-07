package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	queueHost := os.Getenv("QUEUE_HOSTNAME")
	connStr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHost)
	consumer, err := NewConsumer(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	consumer.WaitForMessages()
}
