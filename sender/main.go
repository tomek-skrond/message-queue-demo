package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type MealResponse struct {
	Meals []struct {
		StrMeal string `json:"strMeal"`
	} `json:"meals"`
}

func GenerateOrder() string {
	resp, err := http.Get("https://www.themealdb.com/api/json/v1/1/random.php")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	var mealResp MealResponse
	if err := json.NewDecoder(resp.Body).Decode(&mealResp); err != nil {
		log.Fatalln("Failed to decode JSON:", err)
	}

	// Extract the value of 'strMeal' if available
	if len(mealResp.Meals) > 0 {
		mealName := mealResp.Meals[0].StrMeal
		fmt.Println("Generated Meal Name:", mealName)
		// generatedOrder, err := NewOrder(fmt.Sprint("Order ", mealName), mealName)
		if err != nil {
			log.Fatalln(err)
		}
		return mealName
	}
	return "error"
}

func main() {

	connStr := "amqp://guest:guest@localhost:5672/"
	sender, err := NewSender(connStr)
	if err != nil {
		log.Fatalln(err)
	}

	var name string
	//for testing: i := 0; i <= 15; i++
	for {
		name = GenerateOrder()
		if err != nil {
			log.Fatalln(err)
		}

		if err := sender.Send(Order{
			uuid.New(),
			fmt.Sprint("Order ", name),
			name,
		}); err != nil {
			log.Println(err)
		}

		time.Sleep(15 * time.Second)
	}

}
