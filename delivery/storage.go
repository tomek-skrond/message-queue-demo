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

	fmt.Println("envs: ", user, pass, host, db_name, sslmode)
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

	if db.AutoMigrate(&Delivery{}) != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) CreateDelivery(delivery *Delivery) error {
	if err := s.db.Create(delivery).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdateDelivery(delivery *Delivery) error {
	if err := s.db.Where("id = ?", delivery.ID).Updates(delivery).Error; err != nil {
		return err
	}
	return nil
}
