package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
)

func fetchOrderFromExternalSource() Order {
	// Simulate fetching order data from an external source
	ordName, err := fetchOrders()
	if err != nil {
		return Order{Name: "err", FoodOrdered: "err"}
	}
	ord, err := NewOrder(
		fmt.Sprintf("Order %s", ordName),
		ordName,
	)
	if err != nil {
		return Order{Name: "err", FoodOrdered: "err"}
	}
	return ord
}

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

	fmt.Println(mealResp.Meals[0].StrMeal)
	// Extract the value of 'strMeal' if available
	if len(mealResp.Meals) > 0 {
		mealName := mealResp.Meals[0].StrMeal
		fmt.Println("Generated Meal Name:", mealName)
		// generatedOrder, err := NewOrder(fmt.Sprint("Order ", mealName), mealName)
		if mealName == "" {
			return "", fmt.Errorf("meal name empty! %s", mealName)
		}
		return mealName, nil
	}
	return "", fmt.Errorf("error")
}

func fetchPrice() int {
	return rand.IntN(60) + 1
}

func MarshalBody[T any](b T) ([]byte, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return []byte("conversion err"), err
	}
	return body, err
}
