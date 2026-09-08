package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sharedevent "common/events/paymentcompleted"
	eventingdomain "notification/internal/domain/eventing"
	paymentcompleted "notification/internal/usecase/eventing/paymentcompleted"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PaymentCompletedQueue = "notification.payment-completed.v1"

type PaymentCompletedConsumer struct {
	channel *amqp.Channel
	service *paymentcompleted.Service
}

func NewPaymentCompletedConsumer(conn *amqp.Connection, exchange string, service *paymentcompleted.Service) (*PaymentCompletedConsumer, error) {
	if conn == nil || service == nil || strings.TrimSpace(exchange) == "" {
		return nil, errors.New("payment completed consumer missing dependency")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	dlx := exchange + ".dlx"
	if err := ch.ExchangeDeclare(dlx, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	deadQueue := PaymentCompletedQueue + ".dlq"
	if _, err := ch.QueueDeclare(deadQueue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.QueueBind(deadQueue, sharedevent.EventTypeV1, dlx, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	q, err := ch.QueueDeclare(PaymentCompletedQueue, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange": dlx, "x-dead-letter-routing-key": sharedevent.EventTypeV1,
	})
	if err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.QueueBind(q.Name, sharedevent.EventTypeV1, exchange, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.Qos(16, 0, false); err != nil {
		_ = ch.Close()
		return nil, err
	}
	return &PaymentCompletedConsumer{channel: ch, service: service}, nil
}

func (c *PaymentCompletedConsumer) Run(ctx context.Context) error {
	deliveries, err := c.channel.ConsumeWithContext(ctx, PaymentCompletedQueue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		var e sharedevent.V1
		if err := json.Unmarshal(delivery.Body, &e); err != nil || !e.IsValid() ||
			delivery.Type != sharedevent.EventTypeV1 || delivery.MessageId == "" || delivery.MessageId != e.EventID {
			if nackErr := delivery.Nack(false, false); nackErr != nil {
				return fmt.Errorf("nack malformed payment completed event: %w", nackErr)
			}
			continue
		}
		if _, err := c.service.Handle(ctx, e); err != nil {
			requeue := !errors.Is(err, eventingdomain.ErrEventIdentityConflict)
			if nackErr := delivery.Nack(false, requeue); nackErr != nil {
				return fmt.Errorf("nack payment completed event: %w", nackErr)
			}
			continue
		}
		if err := delivery.Ack(false); err != nil {
			return fmt.Errorf("ack payment completed: %w", err)
		}
	}
	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return errors.New("RabbitMQ payment-completed delivery channel closed")
}

func (c *PaymentCompletedConsumer) Close() error {
	if c == nil || c.channel == nil {
		return nil
	}
	return c.channel.Close()
}
