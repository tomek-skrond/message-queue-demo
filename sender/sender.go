package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	MQ_CONN_ERR          = "MQ Connection Error"
	MQ_CHANNEL_INIT_ERR  = "MQ Channel Init Error"
	MQ_QUEUE_DECLARE_ERR = "MQ Queue Declare Error"
)

type Sender struct {
	connString string
}

func ConfigureTLS() (*tls.Config, error) {
	caCert, err := os.ReadFile("../cert/ca_certificate.pem")
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair("../cert/client_certificate.pem", "../cert/client_key.pem")
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost", // Optional
	}

	return tlsConf, nil
}

func NewSender(connStr string) (*Sender, error) {
	return &Sender{
		connString: connStr,
	}, nil
}

func (s *Sender) ConnectToMQ() (*amqp.Connection, error) {
	tlsConf, err := ConfigureTLS()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	conn, err := amqp.DialTLS(s.connString, tlsConf)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", MQ_CONN_ERR, err)
	}

	return conn, nil
}

func (s *Sender) CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", MQ_CHANNEL_INIT_ERR, err)
	}
	return ch, nil
}

func (s *Sender) DeclareQueue(ch *amqp.Channel) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		"orders",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", MQ_QUEUE_DECLARE_ERR, err)
	}
	return &q, err
}

func (s *Sender) Send(order Order) error {
	conn, err := s.ConnectToMQ()
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := s.CreateChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()
	q, err := s.DeclareQueue(ch)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonBody, err := MarshalBody(order)
	if err != nil {
		return err
	}

	if err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	); err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", jsonBody)
	return nil
}

func MarshalBody[T any](b T) ([]byte, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return []byte("conversion err"), err
	}
	return body, err
}
