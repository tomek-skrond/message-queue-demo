package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type APIServer struct {
	DB         *Storage
	listenPort string
}

func NewAPIServer(lp string, db *Storage) (*APIServer, error) {
	return &APIServer{
		DB:         db,
		listenPort: lp,
	}, nil
}

func (s *APIServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/order", s.handleOrder).Methods("POST")

	log.Fatalln(http.ListenAndServe(s.listenPort, r))
}

func toJSON(obj interface{}) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func (s *APIServer) handleOrder(w http.ResponseWriter, r *http.Request) {
	order := fetchOrderFromExternalSource()
	fmt.Println(order)

	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
	fmt.Println(connstr)

	if err := s.DB.CreateOrder(&order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := toJSON(order)

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
