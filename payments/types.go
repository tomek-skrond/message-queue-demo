package main

import (
	"github.com/google/uuid"
)

type PaymentRequest struct {
	ID        uuid.UUID
	Price     int
	PricePaid int
	Status    string
}
