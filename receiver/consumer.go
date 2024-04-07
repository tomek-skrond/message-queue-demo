package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	MQConn    *amqp.Connection
	MQChannel *amqp.Channel
	MQQueue   *amqp.Queue
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

func NewConsumer(connStr string) (*Consumer, error) {
	tlsConf, err := ConfigureTLS()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	conn, err := amqp.DialTLS(connStr, tlsConf)
	// conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close()

	q, err := ch.QueueDeclare(
		"orders", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		MQConn:    conn,
		MQChannel: ch,
		MQQueue:   &q,
	}, nil
}

func (c *Consumer) WaitForMessages() {
	defer c.MQConn.Close()
	defer c.MQChannel.Close()
	msgs, err := c.MQChannel.Consume(
		c.MQQueue.Name, // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalln(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
