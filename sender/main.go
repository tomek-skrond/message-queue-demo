package main

import (
	"fmt"
	"log"
	"os"
)

type MealResponse struct {
	Meals []struct {
		StrMeal string `json:"strMeal"`
	} `json:"meals"`
}

func main() {

	queueHost := os.Getenv("QUEUE_HOSTNAME")
	connStr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHost)
	sender, err := NewSender(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	sender.sendOrders()

}
