package main

import "log"

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	consumer, err := NewConsumer(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	consumer.WaitForMessages()
}
