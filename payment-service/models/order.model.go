package models

import (
	cmodels "common/models"

	"github.com/payOSHQ/payos-lib-golang"
)

// Order đại diện cho thông tin đơn hàng
// @swagger:model
type Order struct {
	cmodels.BaseEntity
	ID                     uint64
	Amount                 int
	Description            string
	Status                 string
	Reason                 string
	Bin                    string // Mã BIN của ngân hàng
	AccountNumber          string // Số tài khoản thanh toán
	AccountName            string `gorm:"size:40"` // Tên tài khoản thanh toán
	OrderCode              uint   `gorm:"unique"`  // Mã đơn hàng
	Currency               string `gorm:"size:10"` // Đơn vị tiền tệ
	PaymentLinkID          string // ID của link thanh toán
	CheckoutURL            string // URL để khách hàng thanh toán
	QRCode                 string // Mã QR để thanh toán
	AmountPaid             int
	AccountNumberPaid      string `gorm:"size:30"`
	Reference              string
	TransactionDateTime    string
	Code                   string
	Message                string
	CounterAccountBankID   string `gorm:"size:20"`
	CounterAccountBankName string
	CounterAccountName     string
	CounterAccountNumber   string
	VirtualAccountName     string
	VirtualAccountNumber   string
}

type OrderRequest struct {
	Amount      int
	Description string
}

func (Order) TableName() string {
	return "orders"
}

func (o *Order) MapWebhookToOrder(webhookData payos.WebhookType) {
	data := webhookData.Data

	o.AmountPaid = data.Amount
	o.Description = data.Description
	o.AccountNumberPaid = data.AccountNumber
	o.Reference = data.Reference
	o.TransactionDateTime = data.TransactionDateTime
	o.Currency = data.Currency
	o.PaymentLinkID = data.PaymentLinkId
	o.Code = data.Code
	o.Message = data.Desc
	o.CounterAccountBankID = *data.CounterAccountBankId
	o.CounterAccountBankName = *data.CounterAccountBankName
	o.CounterAccountName = *data.CounterAccountName
	o.CounterAccountNumber = *data.CounterAccountNumber
	o.VirtualAccountName = *data.VirtualAccountName
	o.VirtualAccountNumber = *data.VirtualAccountNumber

}
