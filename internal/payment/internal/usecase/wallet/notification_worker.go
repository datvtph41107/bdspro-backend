package wallet

import (
	"context"
	"fmt"

	_utils "common/utils"
	sharepb "pb/types/shared"

	notificationclient "payment/infra/client/notification"
	"payment/internal/domain/wallet"
	"payment/internal/enums"
)

type NotificationDTO struct {
	Transaction *wallet.WalletTransaction
	Message     string
}

type NotificationWorker struct {
	notificationClient notificationclient.NotificationClient
	channel            chan *NotificationDTO
}

func NewNotificationWorker(notificationClient notificationclient.NotificationClient) *NotificationWorker {
	return &NotificationWorker{channel: make(chan *NotificationDTO, 128), notificationClient: notificationClient}
}

func (w *NotificationWorker) Run(ctx context.Context) error {
	if w == nil || w.notificationClient == nil {
		return fmt.Errorf("notification worker is not configured")
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case notification := <-w.channel:
			if notification == nil || notification.Transaction == nil {
				continue
			}
			// Notification is post-commit background responsibility. Process context,
			// not the completed request context, owns its cancellation.
			_ = w.sendToNotificationService(ctx, notification.Transaction, notification.Message)
		}
	}
}

func (w *NotificationWorker) PushToQueue(ctx context.Context, transaction *wallet.WalletTransaction, message string) error {
	if w == nil || transaction == nil {
		return fmt.Errorf("notification worker is not configured")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.channel <- &NotificationDTO{Transaction: transaction, Message: message}:
		return nil
	default:
		return fmt.Errorf("notification queue is full")
	}
}

func (w *NotificationWorker) sendToNotificationService(ctx context.Context, transaction *wallet.WalletTransaction, message string) error {
	var notificationType int32
	if transaction.Type == wallet.TransactionTypeDeposit {
		notificationType = int32(enums.NotificationTypeDeposit)
	} else if transaction.Type == wallet.TransactionTypeWithdrawal {
		notificationType = int32(enums.NotificationTypeWithdraw)
	} else if transaction.Type == wallet.TransactionTypePayment {
		notificationType = int32(enums.NotificationTypePayment)
	} else {
		return fmt.Errorf("invalid transaction type")
	}
	notification := &sharepb.NotificationRequest{
		AttachData: []string{_utils.FormatVND(wallet.MinorToWire(transaction.Amount))},
		Type:       notificationType,
		Title:      "Thông báo thanh toán",
		Message:    []string{message},
		OwnerId:    uint64(transaction.WalletId),
		TargetId:   uint64(transaction.Id),
	}
	_, err := w.notificationClient.Create(ctx, notification)
	return err
}
