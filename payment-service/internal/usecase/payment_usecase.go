package usecase

import (
	"context"
	"fmt"
	paymentpb "pb/types/payment"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/domain/repository"
)

type PaymentUsecase interface {
	GetWalletDashboard(ctx context.Context, req *paymentpb.Wallet) (*paymentpb.WalletDashboard, error)
	GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error)
	Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error)
	Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error)
	MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error)
	GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error)
	CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, req *paymentpb.DeletePaymentMethodRequest) (*paymentpb.Empty, error)
	GetWalletReport(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ReportResponse, error)
	ExportWalletData(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ExportResponse, error)
	GetDealTransactionsDashboard(ctx context.Context, req *paymentpb.Empty) (*paymentpb.DealTransactionsDashboard, error)
	GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) (*paymentpb.DealTransactionList, error)
	GetDealTransaction(ctx context.Context, req *paymentpb.GetDealTransactionRequest) (*paymentpb.DealTransaction, error)
	CreateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error)
	UpdateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error)
	ApproveDealTransaction(ctx context.Context, req *paymentpb.ApproveDealTransactionRequest) (*paymentpb.DealTransaction, error)
	GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error)
	CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error)
	UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error)
	DeleteTransactionType(ctx context.Context, req *paymentpb.DeleteTransactionTypeRequest) (*paymentpb.Empty, error)
}

type paymentUsecase struct {
	paymentRepository           repository.PaymentMethodRepository
	walletRepository            repository.WalletsRepository
	walletTransactionRepository repository.WalletTransactionsRepository
	walletAuditLogRepository    repository.WalletAuditLogRepository
	withdrawalRequestRepository repository.WithdrawalRequestsRepository
}

func NewPaymentUsecase(
	paymentRepository repository.PaymentMethodRepository,
	walletRepository repository.WalletsRepository,
	walletTransactionRepository repository.WalletTransactionsRepository,
	walletAuditLogRepository repository.WalletAuditLogRepository,
	withdrawalRequestRepository repository.WithdrawalRequestsRepository,
) PaymentUsecase {
	return &paymentUsecase{
		paymentRepository:           paymentRepository,
		walletRepository:            walletRepository,
		walletTransactionRepository: walletTransactionRepository,
		walletAuditLogRepository:    walletAuditLogRepository,
		withdrawalRequestRepository: withdrawalRequestRepository,
	}
}

func (p *paymentUsecase) GetWalletDashboard(ctx context.Context, req *paymentpb.Wallet) (*paymentpb.WalletDashboard, error) {
	// Get wallet details
	wallet, err := p.walletRepository.GetWalletById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Get total income and expense
	totalIncome, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, wallet.Id, "INCOME")
	if err != nil {
		return nil, err
	}

	totalExpense, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, wallet.Id, "EXPENSE")
	if err != nil {
		return nil, err
	}

	// Get pending transactions count
	pendingCount, err := p.walletTransactionRepository.GetTransactionCountByStatus(ctx, wallet.Id, "PENDING")
	if err != nil {
		return nil, err
	}

	// Build alerts
	var alerts []string
	if pendingCount > 0 {
		alerts = append(alerts, fmt.Sprintf("You have %d pending transactions", pendingCount))
	}

	return &paymentpb.WalletDashboard{
		Balance:             wallet.Balance,
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		PendingTransactions: int32(pendingCount),
		Alerts:             alerts,
	}, nil
}

func (p *paymentUsecase) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	transactions, total, err := p.walletTransactionRepository.GetTransactions(ctx, req.WalletId, req.Filter, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &paymentpb.WalletTransactionList{
		Transactions: transactions,
		Pagination: &paymentpb.PaginationResponse{
			Total: int32(total),
			Page:  req.Pagination.Page,
			Size:  req.Pagination.Size,
		},
	}, nil
}

func (p *paymentUsecase) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &entity.WalletTransaction{
		WalletId:          req.WalletId,
		Type:             "DEPOSIT",
		Amount:           req.Amount,
		Status:           "COMPLETED",
		PaymentMethod:    req.PaymentMethod,
		ExternalPaymentId: req.ExternalPaymentId,
	}

	err = p.walletTransactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Update wallet balance
	oldBalance := wallet.Balance
	wallet.Balance += req.Amount
	err = p.walletRepository.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &entity.WalletAuditLog{
		WalletId:   req.WalletId,
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Deposit via %s", req.PaymentMethod),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(ctx, auditLog)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction.ToProto(), nil
}

func (p *paymentUsecase) Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}

	// Check balance
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create withdrawal request
	withdrawal := &entity.WithdrawalRequest{
		WalletId:        req.WalletId,
		Amount:          req.Amount,
		Status:          "PENDING",
		PaymentMethod:   req.PaymentMethod,
		BankAccountInfo: req.BankAccountInfo,
		RequestedBy:     0, // TODO: Get from context
	}

	err = p.withdrawalRequestRepository.CreateWithdrawalRequest(ctx, withdrawal)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &entity.WalletTransaction{
		WalletId:       req.WalletId,
		Type:          "WITHDRAWAL",
		Amount:        -req.Amount,
		Status:        "PENDING",
		PaymentMethod: req.PaymentMethod,
	}

	err = p.walletTransactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return withdrawal.ToProto(), nil
}

func (p *paymentUsecase) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}

	// Check balance
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create transaction record
	transaction := &entity.WalletTransaction{
		WalletId:       req.WalletId,
		Type:          "PAYMENT",
		Amount:        -req.Amount,
		Status:        "COMPLETED",
		RelatedService: req.RelatedService,
		RelatedId:     req.RelatedId,
	}

	err = p.walletTransactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Update wallet balance
	oldBalance := wallet.Balance
	wallet.Balance -= req.Amount
	err = p.walletRepository.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &entity.WalletAuditLog{
		WalletId:   req.WalletId,
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Payment for %s:%s", req.RelatedService, req.RelatedId),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(ctx, auditLog)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction.ToProto(), nil
}

func (p *paymentUsecase) GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error) {
	methods, total, err := p.paymentRepository.GetPaymentMethods(ctx, req.OrganizationId, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &paymentpb.PaymentMethodList{
		PaymentMethods: methods,
		Pagination: &paymentpb.PaginationResponse{
			Total: int32(total),
			Page:  req.Pagination.Page,
			Size:  req.Pagination.Size,
		},
	}, nil
}

func (p *paymentUsecase) CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	method := &entity.PaymentMethod{
		OrganizationId: req.OrganizationId,
		Name:          req.Name,
		Code:          req.Code,
		IsActive:      req.IsActive,
	}

	_, err := p.paymentRepository.CreatePaymentMethod(ctx, method)
	if err != nil {
		return nil, err
	}

	return method.ToProto(), nil
}

func (p *paymentUsecase) UpdatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	// Get existing payment method
	method, err := p.paymentRepository.GetPaymentMethodById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Update fields
	method.Name = req.Name
	method.Code = req.Code
	method.IsActive = req.IsActive

	_, err = p.paymentRepository.UpdatePaymentMethod(ctx, method)
	if err != nil {
		return nil, err
	}

	return method.ToProto(), nil
}

func (p *paymentUsecase) DeletePaymentMethod(ctx context.Context, req *paymentpb.DeletePaymentMethodRequest) (*paymentpb.Empty, error) {
	err := p.paymentRepository.DeletePaymentMethod(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &paymentpb.Empty{}, nil
}

func (p *paymentUsecase) GetWalletReport(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ReportResponse, error) {
	transactions, total, totalAmount, err := p.walletTransactionRepository.GetTransactionsReport(ctx, req)
	if err != nil {
		return nil, err
	}

	return &paymentpb.ReportResponse{
		TotalAmount:       totalAmount,
		TotalTransactions: int32(total),
		Transactions:      transactions,
	}, nil
}

func (p *paymentUsecase) ExportWalletData(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ExportResponse, error) {
	data, err := p.walletTransactionRepository.ExportTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	return &paymentpb.ExportResponse{
		Data:   data,
		Format: "CSV",
	}, nil
}

func (p *paymentUsecase) GetDealTransactionsDashboard(ctx context.Context, req *paymentpb.Empty) (*paymentpb.DealTransactionsDashboard, error) {
	totalIncome, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, 0, "DEAL_INCOME")
	if err != nil {
		return nil, err
	}

	totalExpense, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, 0, "DEAL_EXPENSE")
	if err != nil {
		return nil, err
	}

	pendingCount, err := p.walletTransactionRepository.GetTransactionCountByStatus(ctx, 0, "DEAL_PENDING")
	if err != nil {
		return nil, err
	}

	var alerts []string
	if pendingCount > 0 {
		alerts = append(alerts, fmt.Sprintf("You have %d pending deal transactions", pendingCount))
	}

	return &paymentpb.DealTransactionsDashboard{
		TotalIncome:        totalIncome,
		TotalExpense:       totalExpense,
		PendingTransactions: int32(pendingCount),
		Alerts:             alerts,
	}, nil
}

func (p *paymentUsecase) GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) (*paymentpb.DealTransactionList, error) {
	transactions, total, err := p.walletTransactionRepository.GetDealTransactions(ctx, req.Filter, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &paymentpb.DealTransactionList{
		Transactions: transactions,
		Pagination: &paymentpb.PaginationResponse{
			Total: int32(total),
			Page:  req.Pagination.Page,
			Size:  req.Pagination.Size,
		},
	}, nil
}

func (p *paymentUsecase) GetDealTransaction(ctx context.Context, req *paymentpb.GetDealTransactionRequest) (*paymentpb.DealTransaction, error) {
	transaction, err := p.walletTransactionRepository.GetDealTransactionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (p *paymentUsecase) CreateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Create deal transaction record
	transaction := &entity.WalletTransaction{
		Type:          "DEAL",
		Amount:        req.Amount,
		Status:        "PENDING",
		RelatedService: "DEAL",
		RelatedId:     fmt.Sprintf("%d", req.Id),
	}

	err := p.walletTransactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return req, nil
}

func (p *paymentUsecase) UpdateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error) {
	// Get existing transaction
	transaction, err := p.walletTransactionRepository.GetDealTransactionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Update fields
	transaction.Amount = req.Amount
	transaction.Status = req.Status
	transaction.Description = req.Description
	transaction.Type = req.Type

	err = p.walletTransactionRepository.UpdateDealTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (p *paymentUsecase) ApproveDealTransaction(ctx context.Context, req *paymentpb.ApproveDealTransactionRequest) (*paymentpb.DealTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get deal transaction
	transaction, err := p.walletTransactionRepository.GetDealTransactionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Get associated wallet transaction
	walletTx, err := p.walletTransactionRepository.GetTransactionByDealId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, walletTx.WalletId)
	if err != nil {
		return nil, err
	}

	// Check balance for expense transactions
	if walletTx.Amount < 0 && wallet.Balance < -walletTx.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}

	// Update wallet balance
	oldBalance := wallet.Balance
	wallet.Balance += walletTx.Amount
	err = p.walletRepository.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	// Update transaction status
	transaction.Status = "COMPLETED"
	transaction.ApprovedBy = req.ApprovedBy
	err = p.walletTransactionRepository.UpdateDealTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Update wallet transaction status
	walletTx.Status = "COMPLETED"
	err = p.walletTransactionRepository.UpdateTransaction(ctx, walletTx)
	if err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &entity.WalletAuditLog{
		WalletId:   wallet.Id,
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Deal transaction %d approved", req.Id),
		ChangedBy:  req.ApprovedBy,
	}
	err = p.walletAuditLogRepository.CreateAuditLog(ctx, auditLog)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (p *paymentUsecase) GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error) {
	types, err := p.walletTransactionRepository.GetTransactionTypes(ctx)
	if err != nil {
		return nil, err
	}

	return &paymentpb.TransactionTypeList{
		Types: types,
	}, nil
}

func (p *paymentUsecase) CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	err := p.walletTransactionRepository.CreateTransactionType(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (p *paymentUsecase) UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	err := p.walletTransactionRepository.UpdateTransactionType(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (p *paymentUsecase) DeleteTransactionType(ctx context.Context, req *paymentpb.DeleteTransactionTypeRequest) (*paymentpb.Empty, error) {
	err := p.walletTransactionRepository.DeleteTransactionType(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &paymentpb.Empty{}, nil
}

	
