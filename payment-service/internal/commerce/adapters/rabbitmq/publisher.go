// Package rabbitmq implements only the Payment integration-event publisher.
// The caller/process root owns the AMQP connection lifetime. This adapter owns
// its channel and publisher-confirm protocol; business retry remains in the
// durable Payment outbox owner.
package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"payment/internal/commerce/application"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	mu       sync.Mutex
	channel  *amqp.Channel
	exchange string
}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	if conn == nil || strings.TrimSpace(exchange) == "" {
		return nil, errors.New("rabbitmq publisher requires connection and exchange")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq publisher channel: %w", err)
	}
	if err := ch.ExchangeDeclare(strings.TrimSpace(exchange), "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("declare payment event exchange: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("enable rabbitmq publisher confirms: %w", err)
	}
	return &Publisher{channel: ch, exchange: strings.TrimSpace(exchange)}, nil
}

func (p *Publisher) Publish(ctx context.Context, message application.OutboxMessage) error {
	if p == nil || p.channel == nil || message.ID == 0 || strings.TrimSpace(message.EventType) == "" ||
		strings.TrimSpace(message.RoutingKey) == "" || len(message.Payload) == 0 {
		return errors.New("invalid payment outbox message")
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	// Channel publish sequence + deferred confirmation are serialized here so
	// one Outbox claim maps to one known broker confirmation outcome. If the
	// context expires while waiting, outcome is UNKNOWN and the durable Outbox
	// remains retryable; duplicate delivery is handled by consumer Inbox.
	p.mu.Lock()
	deferred, err := p.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		p.exchange,
		message.RoutingKey,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    message.EventKey,
			Type:         message.EventType,
			Timestamp:    time.Now().UTC(),
			Headers: amqp.Table{
				"schema_version": int32(message.SchemaVersion),
				"outbox_id":      int64(message.ID),
			},
			Body: message.Payload,
		},
	)
	p.mu.Unlock()
	if err != nil {
		return fmt.Errorf("publish payment event: %w", err)
	}
	acked, err := deferred.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("payment event confirmation unknown: %w", err)
	}
	if !acked {
		return errors.New("rabbitmq negatively acknowledged payment event")
	}
	return nil
}

func (p *Publisher) Close() error {
	if p == nil || p.channel == nil {
		return nil
	}
	return p.channel.Close()
}

var _ application.EventPublisher = (*Publisher)(nil)
