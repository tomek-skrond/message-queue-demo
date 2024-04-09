package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)

	config := &RabbitMQConfig{
		ConnStr:           connstr,
		QueueName:         "order_processing",
		QueueRoutingKey:   "order_to_payment",
		QueueExchange:     "order_x_payment",
		QueueExchangeType: "direct",
	}

	client, err := NewRabbitMQClient(*config)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	err = client.Consume(config.QueueName, config.QueueConsumerName, processPayment)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// r := mux.NewRouter()
	// r.HandleFunc("/pay", handlePayment).Methods("POST")
	// log.Fatalln(http.ListenAndServe(":7777", r))

	select {}
}

// func handlePayment(w http.ResponseWriter, r *http.Request) {
// 	// orderUUID := r.URL.Query().Get("id")
// 	// paymentPrice := r.URL.Query().Get("price")
// }

func processPayment(body []byte) {
	paymentReq := &PaymentRequest{}
	err := json.Unmarshal(body, paymentReq)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("lelelelelele")
	fmt.Println("price:", paymentReq.Price)

	if paymentReq.Price >= 10 {
		queueHostName := os.Getenv("QUEUE_HOSTNAME")
		connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
		fmt.Println("lelelelelele")

		deliveryConfig := &RabbitMQConfig{
			ConnStr:           connstr,
			QueueName:         "order_processing",
			QueueExchange:     "payment_x_delivery",
			QueueRoutingKey:   "payment_to_delivery",
			QueueExchangeType: "direct",
		}

		deliveryClient, err := NewRabbitMQClient(*deliveryConfig)
		if err != nil {
			log.Fatalln(err)
		}
		defer deliveryClient.Close()

		err = deliveryClient.Publish(context.Background(),
			deliveryClient.QueueExchange,
			deliveryClient.QueueRoutingKey,
			body,
		)
		if err != nil {
			log.Fatalln(err)
			return
		}
	}
}
