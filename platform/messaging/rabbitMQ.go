package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	PaymentQueue = "payments"
)

type MessagingClient interface {
	PublishPayment(ctx context.Context, paymentID string) error
	ConsumePayments(ctx context.Context) (<-chan amqp.Delivery, error)
	Close() error
}

type rabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQClient(url string) (MessagingClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		PaymentQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &rabbitMQClient{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *rabbitMQClient) PublishPayment(ctx context.Context, paymentID string) error {
	body, err := json.Marshal(map[string]string{"payment_id": paymentID})
	if err != nil {
		return fmt.Errorf("failed to marshal payment ID: %w", err)
	}

	err = r.channel.PublishWithContext(ctx,
		"",
		PaymentQueue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published payment message: %s", paymentID)
	return nil
}

func (r *rabbitMQClient) ConsumePayments(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		PaymentQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %w", err)
	}
	return msgs, nil
}

func (r *rabbitMQClient) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
