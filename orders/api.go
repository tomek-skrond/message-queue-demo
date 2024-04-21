package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
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

	r := mux.NewRouter()
	r.HandleFunc("/order", s.handleOrder).Methods("POST")
	r.HandleFunc("/order", s.handleGetOrder).Methods("GET")

	log.Fatalln(http.ListenAndServe(s.listenPort, r))
}

func toJSON(obj interface{}) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func (s *APIServer) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	orders, err := s.DB.GetAllOrders()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// defer s.mqsession.Close()

	resp := toJSON(orders)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}
func (s *APIServer) pushToOrderPaymentsQueue(newOrder Order) error {

	jsonOrder := toJSON(newOrder)

	if err := s.mqsession.channel.Publish(
		"",
		"order_payment_queue",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonOrder,
		},
	); err != nil {
		log.Println(err)
		return err
	}
	log.Println("published to queue:", newOrder)

	return nil
}
func (s *APIServer) handleOrder(w http.ResponseWriter, r *http.Request) {
	order := fetchOrderFromExternalSource()
	fmt.Println(order)

	if err := s.pushToOrderPaymentsQueue(order); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// defer s.mqsession.Close()

	if err := s.DB.CreateOrder(&order); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := toJSON(order)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
