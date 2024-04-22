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

	if db.AutoMigrate(&PaymentRequest{}) != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil

}

func (s *Storage) CreatePaymentRequest(p *PaymentRequest) error {
	if err := s.db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetPayments() ([]*PaymentRequest, error) {
	payments := []*PaymentRequest{}
	result := s.db.Find(&payments)
	if result.Error != nil {
		return nil, result.Error
	}

	return payments, nil
}

func (s *Storage) UpdatePaymentByID(p *PaymentRequest) error {
	if err := s.db.Where("id = ?", p.ID).Updates(p).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) updatePaymentStatus(p *PaymentRequest, status string) error {
	p.Status = status
	p.PricePaid = p.Price
	// if err := s.db.Where("id = ?", p.ID).Update("status", status).Error; err != nil {
	// 	return err
	// }
	if err := s.db.Model(&p).Omit("price").Updates(p).Error; err != nil {
		return err
	}
	return nil
}

// func (s *Storage) ReadOrders() {}
