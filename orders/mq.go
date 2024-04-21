package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQSession struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewMQSession() (*MQSession, error) {

	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
	fmt.Println(connstr)

	conn, err := amqp.Dial(connstr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"order_payment_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println("queue used:", q.Name, q)

	return &MQSession{conn, ch}, nil
}

func (mqs *MQSession) Close() error {
	if mqs.connection == nil {
		return nil
	}
	return mqs.connection.Close()
}
