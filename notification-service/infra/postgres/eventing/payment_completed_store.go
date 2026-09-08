package eventing

import (
	commonenum "common/domain/enum"
	"context"
	"fmt"
	"notification/internal/domain/eventing"
	paymentcompleted "notification/internal/usecase/eventing/paymentcompleted"
	"strconv"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentCompletedStore struct{ db *gorm.DB }

func NewPaymentCompletedStore(db *gorm.DB) *PaymentCompletedStore {
	return &PaymentCompletedStore{db: db}
}

// AcceptOnce atomically transfers responsibility from Rabbit transport to
// Notification durable state. Inbox + logical reaction + optional delivery
// intent are one transaction, so ACK may safely happen after this returns.
func (s *PaymentCompletedStore) AcceptOnce(ctx context.Context, f eventing.PaymentCompleted) (bool, error) {
	if s == nil || s.db == nil {
		return false, fmt.Errorf("payment completed store database is nil")
	}
	replay := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		payloadHash := f.PayloadHash()
		inbox := map[string]any{
			"event_id": f.EventID, "event_type": "payment.completed.v1",
			"schema_version": 1, "order_id": f.OrderID, "payload_hash": payloadHash,
		}
		result := tx.Table("payment_event_inbox").Clauses(clause.OnConflict{DoNothing: true}).Create(inbox)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			var existing struct {
				PayloadHash string `gorm:"column:payload_hash"`
			}
			if err := tx.Table("payment_event_inbox").Select("payload_hash").Where("event_id = ?", f.EventID).Take(&existing).Error; err != nil {
				return err
			}
			if existing.PayloadHash != payloadHash {
				return fmt.Errorf("%w: event_id=%s", eventing.ErrEventIdentityConflict, f.EventID)
			}
			replay = true
			return nil
		}

		reaction := map[string]any{
			"source_event_id": f.EventID, "order_id": f.OrderID,
			"subject_kind": f.SubjectKind, "subject_id": f.SubjectID,
			"kind": "payment_completed", "title": "Thanh toán thành công",
			"body": paymentcompleted.NotificationPayload(f),
		}
		if err := tx.Table("event_notifications").Create(reaction).Error; err != nil {
			return err
		}

		if paymentcompleted.ShouldCreatePushIntent(f) {
			profileID, err := strconv.ParseUint(f.SubjectID, 10, 64)
			if err != nil || profileID == 0 {
				return fmt.Errorf("payment notification profile identity is invalid")
			}
			// The Inbox-owned reaction is also projected into the existing
			// customer notification surface in the same transaction. The
			// unique Inbox event remains the deduplication authority, so broker
			// redelivery cannot create a second visible notification.
			visible := map[string]any{
				"owner_of":          commonenum.EOwnerOfMember,
				"owner_id":          profileID,
				"created_at":        f.CompletedAt,
				"updated_at":        f.CompletedAt,
				"visible_at":        f.CompletedAt,
				"title":             "Thanh toán thành công",
				"message":           pq.Array([]string{paymentcompleted.NotificationPayload(f)}),
				"attach_data":       pq.Array([]string{}),
				"target_id":         f.OrderID,
				"notification_type": commonenum.NotificationSystem,
				"is_read":           false,
			}
			if err := tx.Table("notification").Create(visible).Error; err != nil {
				return err
			}

			intent := map[string]any{
				"event_id": f.EventID, "subject_kind": f.SubjectKind, "subject_id": f.SubjectID,
				"channel": "push", "template_code": "payment_completed",
				"status": "pending", "payload": paymentcompleted.NotificationPayload(f),
			}
			if err := tx.Table("notification_delivery_intents").Create(intent).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return replay, err
}
