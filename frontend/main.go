package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Order represents an order item
type Order struct {
	ID    string `json:"id"`
	Food  string `json:"food"`
	Price int    `json:"price"`
}

// Delivery represents delivery status
type Delivery struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Details string `json:"details"`
}

var orders []Order

func main() {
	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Define routes
	http.HandleFunc("/", homeHandler)
	// http.HandleFunc("/order", orderHandler)
	// http.HandleFunc("/pay", payHandler)
	// http.HandleFunc("/delivery", deliveryHandler)

	// Start server
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Render home page
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, orders)
}

// func orderHandler(w http.ResponseWriter, r *http.Request) {
// 	// Process order
// 	// Simulating ordering random food
// 	order := Order{
// 		ID:    "random-id", // Generate UUID or use any other method to generate unique ID
// 		Food:  "Random Food",
// 		Price: 10, // Change as per your requirement
// 	}
// 	orders = append(orders, order)

// 	// Redirect to home page
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

func payHandler(w http.ResponseWriter, r *http.Request) {
	// Process payment
	var paymentData struct {
		ID    string `json:"id"`
		Price int    `json:"price"`
	}
	err := json.NewDecoder(r.Body).Decode(&paymentData)
	if err != nil {
		http.Error(w, "Failed to decode payment data", http.StatusBadRequest)
		return
	}

	// Do something with payment data
	fmt.Printf("Payment received for order ID %s with price %d\n", paymentData.ID, paymentData.Price)

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

func deliveryHandler(w http.ResponseWriter, r *http.Request) {
	// Simulating delivery status
	delivery := Delivery{
		ID:      r.URL.Query().Get("id"),
		Status:  "In Progress",
		Details: "Your order is being delivered",
	}

	// Respond with delivery status
	json.NewEncoder(w).Encode(delivery)
}
