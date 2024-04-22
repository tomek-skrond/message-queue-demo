package main

import "github.com/google/uuid"

type Delivery struct {
	ID        uuid.UUID
	Delivered bool `gorm:"defult=false"`
}
