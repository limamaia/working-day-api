package messaging

import (
	"fmt"
	"log"
	"sync"
	"working-day-api/config"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	mu         sync.Mutex
}

var (
	Rabbit *RabbitMQ
	once   sync.Once
)

type Messenger interface {
	Publish(topic string, message []byte) error
}

type RabbitMessenger struct{}

func (r *RabbitMessenger) Publish(topic string, message []byte) error {
	return Rabbit.Publish(topic, message)
}

func Connection(config *config.AppVars) {
	once.Do(func() {
		dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/",
			config.RabbitMQ.Username,
			config.RabbitMQ.Password,
			config.RabbitMQ.Host,
			config.RabbitMQ.Port,
		)

		conn, err := amqp.Dial(dsn)
		if err != nil {
			log.Fatalf("Error connecting to RabbitMQ: %v", err)
		}

		Rabbit = &RabbitMQ{
			Connection: conn,
		}
		log.Println("Connection to RabbitMQ successfully established.")
	})
}

func (r *RabbitMQ) Publish(queue string, body []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	channel, err := r.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message successfully published to queue %s.", queue)
	return nil
}

func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Connection != nil {
		err := r.Connection.Close()
		if err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		} else {
			log.Println("RabbitMQ connection successfully closed.")
		}
	}
}
