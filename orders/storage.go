package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage() (*Storage, error) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	sslmode := os.Getenv("SSLMODE")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s TimeZone=Europe/Warsaw",
		host,
		user,
		pass,
		db_name,
		sslmode)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	if db.AutoMigrate(&Order{}) != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) CreateOrder(o *Order) error {
	if err := s.db.Create(o).Error; err != nil {
		return err
	}
	return nil
}
