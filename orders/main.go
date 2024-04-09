package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/order", handleOrder).Methods("POST")

	log.Fatalln(http.ListenAndServe(":8000", r))

}

func toJSON(obj interface{}) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	order := fetchOrderFromExternalSource()

	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
	config := &RabbitMQConfig{
		ConnStr:           connstr,
		QueueName:         "order_processing",
		QueueRoutingKey:   "order_to_payment",
		QueueExchange:     "order_x_payment",
		QueueExchangeType: "direct",
	}
	// client, err := WithTLS(config, NewRabbitMQClient)
	client, err := NewRabbitMQClient(*config)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	err = client.Publish(context.Background(), config.QueueExchange, config.QueueRoutingKey, toJSON(order))
	if err != nil {
		http.Error(w, "Error publishing message", http.StatusInternalServerError)
		return
	}

}
