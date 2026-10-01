package client

import (
	"context"
	"errors"

	pb_payment "pb/types/payment"
)

// PaymentClient is a thin adapter used by CRM payment/campaign/budget services.
type PaymentClient struct{}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{}
}

func (c *PaymentClient) CreateWallet(ctx context.Context, req *pb_payment.CreateWalletRequest) (*pb_payment.CreateWalletResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) GetWalletByWalletId(ctx context.Context, walletID uint64) (*pb_payment.Wallet, error) {
	_ = ctx
	_ = walletID
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) GetWalletByUserId(ctx context.Context, userId uint64) (*pb_payment.Wallet, error) {
	_ = ctx
	_ = userId
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) GetWalletTransactionsByWalletId(ctx context.Context, walletID uint64, req *GetWalletTransactionsRequest) (*WalletTransactionListResponse, error) {
	_ = ctx
	_ = walletID
	_ = req
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) TransferBetweenWallets(ctx context.Context, fromWalletID, toWalletID uint64, amount float64, description, relatedService, relatedID string) (*TransferBetweenWalletsResponse, error) {
	_ = ctx
	_ = fromWalletID
	_ = toWalletID
	_ = amount
	_ = description
	_ = relatedService
	_ = relatedID
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) MakePayment(ctx context.Context, req *MakePaymentRequest) (*WalletTransaction, error) {
	_ = ctx
	_ = req
	return nil, errors.New("payment client not configured")
}

func (c *PaymentClient) GetPaymentMethods(ctx context.Context, req *GetPaymentMethodsRequest) (*PaymentMethodListResponse, error) {
	_ = ctx
	_ = req
	return nil, errors.New("payment client not configured")
}

type GetWalletTransactionsRequest struct {
	Type     string `json:"type"`
	Status   string `json:"status"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
}

type WalletTransactionListResponse struct {
	Transactions []WalletTransaction `json:"transactions"`
	Total        uint64              `json:"total"`
}

type WalletTransaction struct {
	ID                uint64  `json:"id"`
	WalletID          uint64  `json:"wallet_id"`
	Type              string  `json:"type"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	RelatedService    string  `json:"related_service"`
	RelatedID         string  `json:"related_id"`
	ExternalPaymentID string  `json:"external_payment_id"`
	TransactionCode   string  `json:"transaction_code"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type MakePaymentRequest struct {
	WalletID       uint64  `json:"wallet_id"`
	Amount         float64 `json:"amount"`
	RelatedService string  `json:"related_service"`
	RelatedID      string  `json:"related_id"`
}

type GetPaymentMethodsRequest struct {
	OrganizationID uint64 `json:"organization_id"`
	Page           int    `json:"page"`
	Size           int    `json:"size"`
}

type PaymentMethodListResponse struct {
	PaymentMethods []PaymentMethod `json:"payment_methods"`
	Total          uint64          `json:"total"`
}

type TransferBetweenWalletsResponse struct {
	FromTransactionID uint64  `json:"from_transaction_id"`
	ToTransactionID   uint64  `json:"to_transaction_id"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	TransactionCode   string  `json:"transaction_code"`
	CreatedAt         string  `json:"created_at"`
}

type PaymentMethod struct {
	ID             uint64 `json:"id"`
	OrganizationID uint64 `json:"organization_id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
