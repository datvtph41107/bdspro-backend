package client

import (
	"context"
	"fmt"

	"pb/clients"
	pb_payment "pb/types/payment"
)

type PaymentClient struct {
	client pb_payment.PaymentServiceClient
}

func NewPaymentClient(rpcClient *clients.PaymentClient) *PaymentClient {
	return &PaymentClient{
		client: rpcClient.ServiceClient,
	}
}

// CreateWallet creates a new wallet using the payment service
func (c *PaymentClient) CreateWallet(ctx context.Context, req *pb_payment.CreateWalletRequest) (*pb_payment.CreateWalletResponse, error) {
	// Convert to gRPC request using existing payment proto

	grpcResp, err := c.client.CreateWallet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet via gRPC: %w", err)
	}

	// Convert response
	return grpcResp, nil
}

// GetWallet retrieves wallet information
func (c *PaymentClient) GetWalletByWalletId(ctx context.Context, walletID uint64) (*pb_payment.Wallet, error) {

	// Call gRPC service
	grpcReq := &pb_payment.GetWalletByWalletIdRequest{
		WalletId: uint32(walletID),
	}

	grpcResp, err := c.client.GetWalletByWalletId(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet via gRPC: %w", err)
	}

	// Convert response
	return grpcResp, nil
}

func (c *PaymentClient) GetWalletByUserId(ctx context.Context, userId uint64) (*pb_payment.Wallet, error) {

	grpcReq := &pb_payment.GetWalletByUserIdRequest{
		UserId: uint32(userId),
	}

	grpcResp, err := c.client.GetWalletByUserId(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet via gRPC: %w", err)
	}

	return grpcResp, nil
}

// GetWalletDashboard gets wallet dashboard information
func (c *PaymentClient) GetWalletDashboard(ctx context.Context) (*WalletDashboardResponse, error) {
	grpcReq := &pb_payment.Empty{}

	grpcResp, err := c.client.GetWalletDashboard(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet dashboard via gRPC: %w", err)
	}

	return &WalletDashboardResponse{
		Balance:             grpcResp.Balance,
		TotalIncome:         grpcResp.TotalIncome,
		TotalExpense:        grpcResp.TotalExpense,
		PendingTransactions: grpcResp.PendingTransactions,
		Alerts:              grpcResp.Alerts,
	}, nil
}

// GetWalletTransactions gets wallet transactions
func (c *PaymentClient) GetWalletTransactions(ctx context.Context, req *GetWalletTransactionsRequest) (*WalletTransactionListResponse, error) {
	grpcReq := &pb_payment.GetWalletTransactionsRequest{
		Type:     req.Type,
		Status:   req.Status,
		FromDate: req.FromDate,
		ToDate:   req.ToDate,
		Page:     uint32(req.Page),
		Size:     uint32(req.Size),
	}

	grpcResp, err := c.client.GetWalletTransactions(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions via gRPC: %w", err)
	}

	// Convert transactions
	transactions := make([]WalletTransaction, len(grpcResp.Data))
	for i, t := range grpcResp.Data {
		transactions[i] = WalletTransaction{
			ID:                uint64(t.Id),
			WalletID:          uint64(t.WalletId),
			Type:              t.Type,
			Amount:            t.Amount,
			Status:            t.Status,
			RelatedService:    t.RelatedService,
			RelatedID:         t.RelatedId,
			ExternalPaymentID: t.ExternalPaymentId,
			TransactionCode:   t.TransactionCode,
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		}
	}

	return &WalletTransactionListResponse{
		Transactions: transactions,
		Total:        uint64(grpcResp.Total),
	}, nil
}

// GetWalletTransactionsByWalletId gets wallet transactions for a specific wallet ID
func (c *PaymentClient) GetWalletTransactionsByWalletId(ctx context.Context, walletID uint64, req *GetWalletTransactionsRequest) (*WalletTransactionListResponse, error) {
	grpcReq := &pb_payment.GetWalletTransactionsByWalletIdRequest{
		WalletId: uint32(walletID),
		Type:     req.Type,
		Status:   req.Status,
		FromDate: req.FromDate,
		ToDate:   req.ToDate,
		Page:     uint32(req.Page),
		Size:     uint32(req.Size),
	}

	grpcResp, err := c.client.GetWalletTransactionsByWalletId(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions by wallet ID via gRPC: %w", err)
	}

	// Convert transactions
	transactions := make([]WalletTransaction, len(grpcResp.Data))
	for i, t := range grpcResp.Data {
		transactions[i] = WalletTransaction{
			ID:                uint64(t.Id),
			WalletID:          uint64(t.WalletId),
			Type:              t.Type,
			Amount:            t.Amount,
			Status:            t.Status,
			RelatedService:    t.RelatedService,
			RelatedID:         t.RelatedId,
			ExternalPaymentID: t.ExternalPaymentId,
			TransactionCode:   t.TransactionCode,
			CreatedAt:         t.CreatedAt,
			UpdatedAt:         t.UpdatedAt,
		}
	}

	return &WalletTransactionListResponse{
		Transactions: transactions,
		Total:        uint64(grpcResp.Total),
	}, nil
}

// TransferBetweenWallets transfers money between two wallets
func (c *PaymentClient) TransferBetweenWallets(ctx context.Context, fromWalletID, toWalletID uint64, amount float64, description, relatedService, relatedID string) (*TransferBetweenWalletsResponse, error) {
	grpcReq := &pb_payment.TransferBetweenWalletsRequest{
		FromWalletId:   uint32(fromWalletID),
		ToWalletId:     uint32(toWalletID),
		Amount:         amount,
		Description:    description,
		RelatedService: relatedService,
		RelatedId:      relatedID,
	}

	grpcResp, err := c.client.TransferBetweenWallets(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer between wallets via gRPC: %w", err)
	}

	return &TransferBetweenWalletsResponse{
		FromTransactionID: uint64(grpcResp.FromTransactionId),
		ToTransactionID:   uint64(grpcResp.ToTransactionId),
		Amount:            grpcResp.Amount,
		Status:            grpcResp.Status,
		TransactionCode:   grpcResp.TransactionCode,
		CreatedAt:         grpcResp.CreatedAt,
	}, nil
}

// MakePayment makes a payment from wallet
func (c *PaymentClient) MakePayment(ctx context.Context, req *MakePaymentRequest) (*WalletTransaction, error) {
	grpcReq := &pb_payment.PaymentRequest{
		WalletId:       uint32(req.WalletID),
		Amount:         req.Amount,
		RelatedService: req.RelatedService,
		RelatedId:      req.RelatedID,
	}

	grpcResp, err := c.client.MakePayment(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make payment via gRPC: %w", err)
	}

	return &WalletTransaction{
		ID:                uint64(grpcResp.Id),
		WalletID:          uint64(grpcResp.WalletId),
		Type:              grpcResp.Type,
		Amount:            grpcResp.Amount,
		Status:            grpcResp.Status,
		RelatedService:    grpcResp.RelatedService,
		RelatedID:         grpcResp.RelatedId,
		ExternalPaymentID: grpcResp.ExternalPaymentId,
		TransactionCode:   grpcResp.TransactionCode,
		CreatedAt:         grpcResp.CreatedAt,
		UpdatedAt:         grpcResp.UpdatedAt,
	}, nil
}

// Deposit money to wallet
func (c *PaymentClient) Deposit(ctx context.Context, req *DepositRequest) (*WalletTransaction, error) {
	grpcReq := &pb_payment.DepositRequest{
		WalletId:          uint32(req.WalletID),
		Amount:            req.Amount,
		PaymentMethod:     req.PaymentMethod,
		ExternalPaymentId: req.ExternalPaymentID,
	}

	grpcResp, err := c.client.Deposit(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit via gRPC: %w", err)
	}

	return &WalletTransaction{
		ID:                uint64(grpcResp.Id),
		WalletID:          uint64(grpcResp.WalletId),
		Type:              grpcResp.Type,
		Amount:            grpcResp.Amount,
		Status:            grpcResp.Status,
		RelatedService:    grpcResp.RelatedService,
		RelatedID:         grpcResp.RelatedId,
		ExternalPaymentID: grpcResp.ExternalPaymentId,
		TransactionCode:   grpcResp.TransactionCode,
		CreatedAt:         grpcResp.CreatedAt,
		UpdatedAt:         grpcResp.UpdatedAt,
	}, nil
}

// GetPaymentMethods gets available payment methods
func (c *PaymentClient) GetPaymentMethods(ctx context.Context, req *GetPaymentMethodsRequest) (*PaymentMethodListResponse, error) {
	grpcReq := &pb_payment.GetPaymentMethodsRequest{
		OrganizationId: uint32(req.OrganizationID),
		Page:           uint32(req.Page),
		Size:           uint32(req.Size),
	}

	grpcResp, err := c.client.GetPaymentMethods(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods via gRPC: %w", err)
	}

	// Convert payment methods
	methods := make([]PaymentMethod, len(grpcResp.Data))
	for i, m := range grpcResp.Data {
		methods[i] = PaymentMethod{
			ID:             uint64(m.Id),
			OrganizationID: uint64(m.OrganizationId),
			Name:           m.Name,
			Code:           m.Code,
			IsActive:       m.IsActive,
			CreatedAt:      m.CreatedAt,
			UpdatedAt:      m.UpdatedAt,
		}
	}

	return &PaymentMethodListResponse{
		PaymentMethods: methods,
		Total:          uint64(grpcResp.Total),
	}, nil
}

// Request/Response structs for internal use
type CreateWalletRequest struct {
	OwnerID       uint64  `json:"owner_id"`
	OwnerType     string  `json:"owner_type"`
	Currency      string  `json:"currency"`
	Description   string  `json:"description"`
	CampaignID    uint64  `json:"campaign_id"`
	InitialAmount float64 `json:"initial_amount"`
}

type CreateWalletResponse struct {
	Success bool `json:"success"`
	Data    struct {
		WalletID   string  `json:"wallet_id"`
		Balance    float64 `json:"balance"`
		Currency   string  `json:"currency"`
		OwnerID    uint64  `json:"owner_id"`
		OwnerType  string  `json:"owner_type"`
		CampaignID uint64  `json:"campaign_id"`
		CreatedAt  string  `json:"created_at"`
		UpdatedAt  string  `json:"updated_at"`
	} `json:"data"`
	Message string `json:"message"`
}

type WalletDashboardResponse struct {
	Balance             float64  `json:"balance"`
	TotalIncome         float64  `json:"total_income"`
	TotalExpense        float64  `json:"total_expense"`
	PendingTransactions int32    `json:"pending_transactions"`
	Alerts              []string `json:"alerts"`
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

type DepositRequest struct {
	WalletID          uint64  `json:"wallet_id"`
	Amount            float64 `json:"amount"`
	PaymentMethod     string  `json:"payment_method"`
	ExternalPaymentID string  `json:"external_payment_id"`
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
