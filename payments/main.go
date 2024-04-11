package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	PaymentDb = []*PaymentRequest{}
)

func main() {

	go func() {
		r := mux.NewRouter()
		r.HandleFunc("/pay", handleAddNewPayment).Methods("POST")
		log.Fatalln(http.ListenAndServe(":7777", r))
	}()

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

	select {}
}

func handleAddNewPayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq PaymentRequest
	r.Header.Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println("xdfxdfxdf", paymentReq.ID)
	payedPrice, err := strconv.Atoi(mux.Vars(r)["price"])
	if err != nil {
		log.Fatalln(err)
		return
	}
	var paymentUuid uuid.UUID
	paymentId := []byte(mux.Vars(r)["id"])
	paymentUuid = uuid.UUID(paymentId)

	paymentReq.Price = payedPrice
	paymentReq.ID = paymentUuid
	fmt.Println("new payment: ")
	PaymentDb = append(PaymentDb, &paymentReq)
	fmt.Println("payment DB: ", PaymentDb)
}

func checkPaymentIssued(payment *PaymentRequest) bool {
	for _, r := range PaymentDb {
		if r.ID == payment.ID {
			return true
		}
	}
	return false
}

func getPricePayed(pr *PaymentRequest) int {
	for _, record := range PaymentDb {
		if record.ID == pr.ID {
			return record.Price
		}
	}
	return 0
}

func processPayment(body []byte) {
	// Take the price from orders coming from order service
	orderPaymentReq := &PaymentRequest{}
	err := json.Unmarshal(body, orderPaymentReq)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if checkPaymentIssued(orderPaymentReq) {
		// fetch recent paid amount of money for order's ID (orderPaymentReq.ID)
		payedPrice := getPricePayed(orderPaymentReq)
		// fmt.Println("lelelelelele")
		fmt.Println("price:", orderPaymentReq.Price)
		realPrice := orderPaymentReq.Price

		if payedPrice >= realPrice {
			queueHostName := os.Getenv("QUEUE_HOSTNAME")
			connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
			// fmt.Println("lelelelelele")

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
		} else {
			fmt.Println("Not enough money")
		}
	}
}
