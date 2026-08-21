package paymentcompleted

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	event "common/events/paymentcompleted"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const queueName = "notification.payment-completed.v1"

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

// ApplyOnce makes transport redelivery harmless: Inbox identity and a durable
// notification intent are committed atomically. Actual external delivery is a
// separate Notification-owned responsibility.
func (s *Store) ApplyOnce(ctx context.Context, e event.V1) (bool, error) {
	if s == nil || s.db == nil || strings.TrimSpace(e.EventID) == "" || e.OrderID == 0 {
		return false, errors.New("invalid payment completed notification event")
	}
	replay := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		inbox := map[string]any{"event_id": e.EventID, "event_type": event.EventTypeV1, "order_id": e.OrderID}
		result := tx.Table("payment_event_inbox").Clauses(clause.OnConflict{DoNothing: true}).Create(inbox)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			replay = true
			return nil
		}
		intent := map[string]any{
			"event_id": e.EventID, "order_id": e.OrderID, "subject_kind": e.SubjectKind, "subject_id": e.SubjectID,
			"template_code": "payment_completed", "status": "pending", "payload": fmt.Sprintf("%s|%s|%d", e.PlanCode, e.Currency, e.AmountMinor),
		}
		return tx.Table("payment_notification_intents").Create(intent).Error
	})
	return replay, err
}

type Consumer struct {
	channel *amqp.Channel
	store   *Store
}

func NewConsumer(conn *amqp.Connection, exchange string, store *Store) (*Consumer, error) {
	if conn == nil || store == nil || strings.TrimSpace(exchange) == "" {
		return nil, errors.New("notification payment consumer missing dependency")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	// The consumer owns its queue and dead-letter queue. Platform owns broker
	// availability, but the consumer team owns poison-message triage because only
	// that team can decide whether replaying the business reaction is safe.
	dlx := exchange + ".dlx"
	if err := ch.ExchangeDeclare(dlx, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	deadQueue := queueName + ".dlq"
	if _, err := ch.QueueDeclare(deadQueue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.QueueBind(deadQueue, event.EventTypeV1, dlx, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	q, err := ch.QueueDeclare(queueName, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    dlx,
		"x-dead-letter-routing-key": event.EventTypeV1,
	})
	if err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.QueueBind(q.Name, event.EventTypeV1, exchange, false, nil); err != nil {
		_ = ch.Close()
		return nil, err
	}
	if err := ch.Qos(16, 0, false); err != nil {
		_ = ch.Close()
		return nil, err
	}
	return &Consumer{channel: ch, store: store}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	deliveries, err := c.channel.ConsumeWithContext(ctx, queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		var e event.V1
		if err := json.Unmarshal(delivery.Body, &e); err != nil || e.EventID == "" {
			_ = delivery.Nack(false, false)
			continue
		}
		if _, err := c.store.ApplyOnce(ctx, e); err != nil {
			_ = delivery.Nack(false, true)
			continue
		}
		if err := delivery.Ack(false); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (c *Consumer) Close() error {
	if c == nil || c.channel == nil {
		return nil
	}
	return c.channel.Close()
}
