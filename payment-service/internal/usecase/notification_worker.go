package usecase

import (
	"context"
	"fmt"
	pb_noti "pb/types/notification"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/infrastructure/client"
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
		notificationType = 110
	} else if transaction.Type == entity.TransactionTypeWithdrawal {
		notificationType = 111
	} else if transaction.Type == entity.TransactionTypePayment {
		notificationType = 112
	} else {
		return fmt.Errorf("invalid transaction type")
	}
	notification := &pb_noti.NotiNewRequest{
		AttachData: []string{fmt.Sprintf("%f", transaction.Amount)},
		Type:       notificationType,
		Title:      "Thông báo thanh toán",
		Message:    message,
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
