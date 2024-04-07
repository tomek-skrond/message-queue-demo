package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func fetchOrders() (string, error) {
	resp, err := http.Get("https://www.themealdb.com/api/json/v1/1/random.php")
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	defer resp.Body.Close()

	var mealResp MealResponse
	if err := json.NewDecoder(resp.Body).Decode(&mealResp); err != nil {
		log.Fatalln("Failed to decode JSON:", err)
		return "", err
	}

	// Extract the value of 'strMeal' if available
	if len(mealResp.Meals) > 0 {
		mealName := mealResp.Meals[0].StrMeal
		fmt.Println("Generated Meal Name:", mealName)
		// generatedOrder, err := NewOrder(fmt.Sprint("Order ", mealName), mealName)
		if err != nil {
			log.Fatalln(err)
		}
		return mealName, nil
	}
	return "", fmt.Errorf("error")
}

func MarshalBody[T any](b T) ([]byte, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return []byte("conversion err"), err
	}
	return body, err
}
