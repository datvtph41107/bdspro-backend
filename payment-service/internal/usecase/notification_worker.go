package usecase

import (
	"context"
	"fmt"

	_utils "common/utils"
	sharepb "pb/types/shared"

	"payment/internal/domain/entity"
	"payment/internal/enums"
	"payment/internal/infrastructure/client"
)

type NotificationDTO struct {
	Transaction *entity.WalletTransaction
	Message     string
}

type notificationWorker struct {
	notificationClient client.NotificationClient
	channel            chan *NotificationDTO
}

func NewNotificationWorker(notificationClient client.NotificationClient) *notificationWorker {
	return &notificationWorker{channel: make(chan *NotificationDTO), notificationClient: notificationClient}
}

func (w *notificationWorker) Start() {
	go func() {
		for notification := range w.channel {
			w.sendToNotificationService(context.Background(), notification.Transaction, notification.Message)
		}
	}()
}

func (w *notificationWorker) PushToQueue(ctx context.Context, transaction *entity.WalletTransaction, message string) error {
	w.channel <- &NotificationDTO{Transaction: transaction, Message: message}
	return nil
}

func (w *notificationWorker) sendToNotificationService(ctx context.Context, transaction *entity.WalletTransaction, message string) error {
	var notificationType int32
	if transaction.Type == entity.TransactionTypeDeposit {
		notificationType = int32(enums.NotificationTypeDeposit)
	} else if transaction.Type == entity.TransactionTypeWithdrawal {
		notificationType = int32(enums.NotificationTypeWithdraw)
	} else if transaction.Type == entity.TransactionTypePayment {
		notificationType = int32(enums.NotificationTypePayment)
	} else {
		return fmt.Errorf("invalid transaction type")
	}
	notification := &sharepb.NotificationRequest{
		AttachData: []string{_utils.FormatVND(transaction.Amount)},
		Type:       notificationType,
		Title:      "Thông báo thanh toán",
		Message:    []string{message},
		OwnerId:    uint64(transaction.WalletId),
		TargetId:   uint64(transaction.Id),
	}
	response, err := w.notificationClient.Create(ctx, notification)
	if err != nil {
		return err
	}
	fmt.Println("response", response)
	return nil
}
