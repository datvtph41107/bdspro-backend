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

	"payment/internal/usecase/outbox"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	mu       sync.Mutex
	channel  *amqp.Channel
	exchange string
	returns  <-chan amqp.Return
	closed   <-chan *amqp.Error
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
	returns := ch.NotifyReturn(make(chan amqp.Return, 1))
	closed := ch.NotifyClose(make(chan *amqp.Error, 1))
	return &Publisher{channel: ch, exchange: strings.TrimSpace(exchange), returns: returns, closed: closed}, nil
}

func (p *Publisher) Publish(ctx context.Context, message outbox.Message) error {
	if p == nil || p.channel == nil || message.ID == 0 || strings.TrimSpace(message.EventType) == "" ||
		strings.TrimSpace(message.RoutingKey) == "" || strings.TrimSpace(message.EventID) == "" || len(message.Payload) == 0 {
		return errors.New("invalid payment outbox message")
	}
	if p.channel.IsClosed() {
		return outbox.ErrPublisherUnavailable
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := p.unavailable(); err != nil {
		return err
	}

	// Rabbit channels are serialized to preserve one deterministic
	// publish/return/confirm outcome at a time. A lost confirm remains UNKNOWN
	// and the durable Outbox retries; consumer Inbox absorbs duplicates.
	p.mu.Lock()
	defer p.mu.Unlock()

	// There should be no stale return because only one publish is in flight.
	// Drain defensively so an earlier caller cancellation cannot poison the
	// evidence for a later event.
	select {
	case <-p.returns:
	default:
	}

	deferred, err := p.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		p.exchange,
		message.RoutingKey,
		true, // mandatory: unroutable publication must be observable
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    message.EventID,
			Type:         message.EventType,
			Timestamp:    time.Now().UTC(),
			Headers: amqp.Table{
				"schema_version": int32(message.SchemaVersion),
				"outbox_id":      int64(message.ID),
			},
			Body: message.Payload,
		},
	)
	if err != nil {
		if errors.Is(err, amqp.ErrClosed) {
			return fmt.Errorf("%w: publish channel closed: %v", outbox.ErrPublisherUnavailable, err)
		}
		return fmt.Errorf("publish payment event: %w", err)
	}

	acked, err := deferred.WaitContext(ctx)
	if err != nil {
		if unavailable := p.unavailable(); unavailable != nil {
			return unavailable
		}
		return fmt.Errorf("payment event confirmation unknown: %w", err)
	}
	if !acked {
		return errors.New("rabbitmq negatively acknowledged payment event")
	}

	// For mandatory unroutable messages RabbitMQ sends basic.return before the
	// corresponding publisher confirm. Since the confirm has arrived and the
	// channel permits one in-flight publish, a queued return is authoritative.
	select {
	case returned := <-p.returns:
		if returned.MessageId == "" || returned.MessageId == message.EventID {
			return fmt.Errorf("payment event unroutable: reply=%d %s", returned.ReplyCode, returned.ReplyText)
		}
		return fmt.Errorf("unexpected RabbitMQ return for message %q", returned.MessageId)
	default:
		return nil
	}
}

func (p *Publisher) unavailable() error {
	if p == nil || p.channel == nil {
		return outbox.ErrPublisherUnavailable
	}
	select {
	case closeErr, ok := <-p.closed:
		if !ok || closeErr == nil {
			return outbox.ErrPublisherUnavailable
		}
		return fmt.Errorf("%w: %v", outbox.ErrPublisherUnavailable, closeErr)
	default:
		return nil
	}
}

func (p *Publisher) Close() error {
	if p == nil || p.channel == nil {
		return nil
	}
	return p.channel.Close()
}

var _ outbox.Publisher = (*Publisher)(nil)
