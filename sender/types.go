package main

import uuid "github.com/google/uuid"

type Order struct {
	ID          uuid.UUID
	Name        string
	FoodOrdered string
}

func NewOrder(name string, food string) (*Order, error) {
	return &Order{
		ID:          uuid.New(),
		Name:        name,
		FoodOrdered: food,
	}, nil
}
