package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	lp := ":9999"
	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
	fmt.Println(connstr)

	r := mux.NewRouter()

	log.Fatalln(http.ListenAndServe(lp, r))
}
