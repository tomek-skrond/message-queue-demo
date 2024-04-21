package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	lp := ":7777"
	queueHostName := os.Getenv("QUEUE_HOSTNAME")
	connstr := fmt.Sprintf("amqp://guest:guest@%s:5672/", queueHostName)
	fmt.Println(connstr)

	db, err := NewStorage()
	if err != nil {
		log.Fatalln(err)
	}
	s, err := NewAPIServer(lp, db)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
}

func (s *APIServer) handleAddPayment(w http.ResponseWriter, r *http.Request) {
	var paymentReq *PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentReq)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(toJSON(""))
}

func toJSON(obj interface{}) []byte {
	data, _ := json.Marshal(obj)
	return data
}

func (s *APIServer) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/pay", s.handleProcessPayment).Methods("POST")
	log.Fatalln(http.ListenAndServe(s.listenPort, r))
}

func (s *APIServer) handleProcessPayment(w http.ResponseWriter, r *http.Request) {
	// decode payment from user
	var paymentComing *PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentComing); err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	isPayed, err := s.isAlreadyPaid(paymentComing)
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if !isPayed {
		realPrice, err := s.checkPrice(paymentComing)
		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal server error"))
			return
		}
		if paymentComing.Price >= realPrice {
			if err := s.DB.updatePaymentStatus(paymentComing, "paid"); err != nil {
				log.Fatalln(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal server error"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("not enough money"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order is already paid for"))
}

func (s *APIServer) isAlreadyPaid(payment *PaymentRequest) (bool, error) {
	payments, err := s.DB.GetPayments()
	if err != nil {
		return false, err
	}

	for _, p := range payments {
		if payment.ID == p.ID {
			if p.Status == "paid" {
				return true, nil
			}
			if p.Status == "pending" {
				return false, nil
			}
			return false, fmt.Errorf("invalid status")
		}
	}
	return false, fmt.Errorf("payment not found")
}
func (s *APIServer) checkPrice(payment *PaymentRequest) (int, error) {
	payments, err := s.DB.GetPayments()
	if err != nil {
		return -1, err
	}
	for _, p := range payments {
		if p.ID == payment.ID {
			//returns real price for order
			return p.Price, nil
		}
	}
	return -1, fmt.Errorf("object not found")
}
