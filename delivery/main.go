package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID
	Name        string
	FoodOrdered string
	Price       int
}

func startDelivery(body []byte) {
	// Simulate delivery process
	log.Println("Delivery process started for item: ", string(body))
}

func main() {

	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)

	config := &RabbitMQConfig{
		ConnStr:           connstr,
		QueueName:         "order_processing",
		QueueExchange:     "payment_x_delivery",
		QueueRoutingKey:   "payment_to_delivery",
		QueueConsumerName: "payment_consumer",
		QueueExchangeType: "direct",
	}

	client, err := NewRabbitMQClient(*config)
	if err != nil {
		log.Fatalf("Error creating RabbitMQ client: %s", err)
	}
	defer client.Close()

	err = client.Consume(config.QueueName, config.QueueConsumerName, startDelivery)
	if err != nil {
		log.Fatalf("Error consuming delivery queue: %s", err)
	}

	select {}
}
