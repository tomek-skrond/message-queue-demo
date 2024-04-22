package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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
			var d *Delivery
			if err := json.Unmarshal(msg.Body, &d); err != nil {
				log.Println(err)
				return
			}

			msg := msg.Body
			if err := s.insertMessagesIntoDB(msg); err != nil {
				log.Println(err)
				return
			}

			s.startDelivery(d)
		}
	}()

	log.Println("[*] Waiting for messages.")
	<-forever

}

func (s *APIServer) startDelivery(d *Delivery) {
	log.Println("[DELIVERY] Delivery started for order id:", d.ID)
	log.Println("[DELIVERY] Status: waiting for the delivery guy...")
	time.Sleep(time.Duration(time.Duration(rand.Intn(5)).Seconds())) //what the f
	log.Println("[DELIVERY] Status: successful")
	d.Delivered = true
	if err := s.DB.UpdateDelivery(d); err != nil {
		log.Println("Delivery error", err)
		return
	}
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
