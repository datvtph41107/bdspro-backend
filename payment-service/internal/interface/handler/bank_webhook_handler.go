package handler

import (
	"context"
	"fmt"
	paymentpb "pb/types/payment"
	"time"
)

type BankWebhookData struct {
	Gateway         string  `json:"gateway"`
	TransactionDate string  `json:"transactionDate"`
	AccountNumber   string  `json:"accountNumber"`
	SubAccount      string  `json:"subAccount"`
	Code            string  `json:"code"`
	Content         string  `json:"content"`
	TransferType    string  `json:"transferType"`
	Description     string  `json:"description"`
	TransferAmount  float64 `json:"transferAmount"`
	ReferenceCode   string  `json:"referenceCode"`
	Accumulated     float64 `json:"accumulated"`
	ID             uint32  `json:"id"`
}

type BankWebhookHandler struct {
	paymentHandler *PaymentHandler
}

func NewBankWebhookHandler(paymentHandler *PaymentHandler) *BankWebhookHandler {
	return &BankWebhookHandler{
		paymentHandler: paymentHandler,
	}
}

func (h *BankWebhookHandler) HandleBankWebhook(ctx context.Context, data *BankWebhookData) error {
	// Only process incoming transfers
	if data.TransferType != "in" {
		return nil
	}

	// Parse transaction date
	_, err := time.Parse("2006-01-02 15:04:05", data.TransactionDate)
	if err != nil {
		return fmt.Errorf("invalid transaction date format: %v", err)
	}

	// Create deposit request
	depositReq := &paymentpb.DepositRequest{
		WalletId:          0, // TODO: Need to map subAccount to walletId
		Amount:            data.TransferAmount,
		PaymentMethod:     data.Gateway,
		ExternalPaymentId: data.ReferenceCode,
	}

	// Process deposit
	_, err = h.paymentHandler.Deposit(ctx, depositReq)
	if err != nil {
		return fmt.Errorf("failed to process deposit: %v", err)
	}

	return nil
} 

// {
// 	"gateway": "BIDV",
// 	"transactionDate": "2025-06-19 22:56:39",
// 	"accountNumber": "1600966010",
// 	"subAccount": "96247199",
// 	"code": null,
// 	"content": "MBVCB.9914229506.183091.DO QUOC TUAN chuyen tien.CT tu 1016242286 DO QUOC TUAN toi 96247199 DO QUOC TUAN tai BIDV",
// 	"transferType": "in",
// 	"description": "BankAPINotify MBVCB.9914229506.183091.DO QUOC TUAN chuyen tien.CT tu 1016242286 DO QUOC TUAN toi 96247199 DO QUOC TUAN tai BIDV",
// 	"transferAmount": 2000,
// 	"referenceCode": "a42779fe-4420-435a-823e-39854b66a7f1",
// 	"accumulated": 0,
// 	"id": 15339673
//   }

// Apikey APi cuc xin
