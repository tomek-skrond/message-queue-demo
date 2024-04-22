package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	DB         *Storage
	mqsession  *MQSession
	listenPort string
}

func NewAPIServer(lp string, db *Storage) (*APIServer, error) {
	session, err := NewMQSession()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &APIServer{
		DB:         db,
		mqsession:  session,
		listenPort: lp,
	}, nil
}

func (s *APIServer) Start() {
	go s.checkForNewMessages()

	r := mux.NewRouter()
	r.HandleFunc("/deliveries", s.handleGetDeliveries).Methods("GET")
	log.Fatalln(http.ListenAndServe(s.listenPort, r))
}

func (s *APIServer) handleGetDeliveries(w http.ResponseWriter, r *http.Request) {

}

func (s *APIServer) checkForNewMessages() {
	// messages := make(chan []byte)

	q, err := s.mqsession.channel.QueueDeclare(
		"payment_delivery_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("consuming from q:", q)
	deliveries, err := s.mqsession.channel.Consume(
		"payment_delivery_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	var forever chan struct{}

	go func() {
		for msg := range deliveries {
			log.Println("received msg:", func(b []byte) Delivery {
				var d Delivery
				_ = json.Unmarshal(b, &d)
				return d
			}(msg.Body))

			msg := msg.Body
			if err := s.insertMessagesIntoDB(msg); err != nil {
				log.Println(err)
			}
		}
	}()

	log.Println("[*] Waiting for messages.")
	<-forever

}

func (s *APIServer) insertMessagesIntoDB(msg []byte) error {
	var newDelivery Delivery
	if err := json.Unmarshal(msg, &newDelivery); err != nil {
		return err
	}

	if err := s.DB.CreateDelivery(&newDelivery); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) CreateDelivery(delivery *Delivery) error {
	if err := s.db.Create(delivery).Error; err != nil {
		return err
	}
	return nil
}
