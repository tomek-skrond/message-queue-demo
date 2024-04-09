package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ENV_VAR_TLS = "TLS_PATH"
)

type RabbitMQConfigFunc func(RabbitMQConfig)

type RabbitMQConfig struct {
	ConnStr           string
	TlsConfig         *tls.Config
	TLS               bool
	QueueName         string
	QueueExchange     string
	QueueExchangeType string
	QueueRoutingKey   string
	QueueConsumerName string
}

type RabbitMQClient struct {
	RabbitMQConfig
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient(config RabbitMQConfig, configFunc ...RabbitMQConfigFunc) (*RabbitMQClient, error) {
	for _, decorator := range configFunc {
		decorator(config)
	}

	fmt.Println(config.QueueExchangeType)

	var connErr error
	var conn *amqp.Connection
	if config.TLS {
		conn, connErr = amqp.DialTLS(config.ConnStr, config.TlsConfig)
		if connErr != nil {
			return nil, connErr
		}
	} else {
		conn, connErr = amqp.Dial(config.ConnStr)
		if connErr != nil {
			return nil, connErr
		}
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if config.QueueName != "" {
		_, err := channel.QueueDeclare(
			config.QueueName, // name
			true,             // durable
			false,            // delete when unused
			false,            // exclusive
			false,            // no-wait
			nil,              // arguments
		)
		if err != nil {
			conn.Close()
			return nil, err
		}
		// fmt.Println("Queue Declared")
		// fmt.Println(q.Name)
	}
	if config.QueueExchange != "" {
		fmt.Println("queue exchange type:", config.QueueExchangeType)
		if config.QueueExchangeType == "" {
			config.QueueConsumerName = "fanout"
		}
		err = channel.ExchangeDeclare(
			config.QueueExchange,
			config.QueueExchangeType,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	if config.QueueExchange != "" && config.QueueName != "" {
		err = channel.QueueBind(
			config.QueueName,
			config.QueueRoutingKey,
			config.QueueExchange,
			false,
			nil,
		)
		if err != nil {
			conn.Close()
			log.Fatalf("Failed to bind queue to exchange: %s", err)
			return nil, err
		}
	}
	return &RabbitMQClient{
		RabbitMQConfig: config,
		Conn:           conn,
		Channel:        channel,
	}, nil
}

func (c *RabbitMQClient) Close() error {
	if err := c.Channel.Close(); err != nil {
		return err
	}
	return c.Conn.Close()
}

func (c *RabbitMQClient) Publish(ctx context.Context, exchange, routing string, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	log.Printf(" [x] Sent %s to %s\n", body, c.RabbitMQConfig.QueueName)
	return c.Channel.PublishWithContext(ctx,
		exchange,
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (c *RabbitMQClient) Consume(queueName string, consumer string, handler func([]byte)) error {
	msgs, err := c.Channel.Consume(
		queueName, // queue
		consumer,  // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()
	return nil
}

func WithTLS(config RabbitMQConfig) {
	tlsPath := os.Getenv(ENV_VAR_TLS)
	if tlsPath == "" {
		log.Fatalln("no tls path specified")
	}
	tlsConfig, err := ConfigureTLS(tlsPath)
	if err != nil {
		log.Fatalln(err)
	}
	config.TLS = true
	config.TlsConfig = tlsConfig
}

func ConfigureTLS(tlsPath string) (*tls.Config, error) {
	caCertPath := fmt.Sprintf("%s/ca_certificate.pem", tlsPath)
	clientCertPath := fmt.Sprintf("%s/client_certificate.pem", tlsPath)
	clientKeyPath := fmt.Sprintf("%s/client_key.pem", tlsPath)

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, err
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCert)

	tlsConf := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{cert},
		// ServerName:   "localhost", // Optional
	}

	return tlsConf, nil
}
