package usecase

import (
	"context"
	"fmt"
	paymentpb "pb/types/payment"
	"time"

	"github.com/tuandq2112/payment-service/appconfig"
	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/domain/repository"
	"github.com/tuandq2112/payment-service/internal/infrastructure/client"
	"github.com/tuandq2112/payment-service/pkg/utils"
	"gorm.io/gorm"
)

type PaymentUsecase interface {
	GetWalletDashboard(ctx context.Context) (*paymentpb.WalletDashboard, error)
	GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error)
	GetWalletTransactionsByWalletId(ctx context.Context, req *paymentpb.GetWalletTransactionsByWalletIdRequest) (*paymentpb.WalletTransactionList, error)
	TransferBetweenWallets(ctx context.Context, req *paymentpb.TransferBetweenWalletsRequest) (*paymentpb.TransferBetweenWalletsResponse, error)
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
	CreateWallet(ctx context.Context, req *paymentpb.CreateWalletRequest) (*paymentpb.CreateWalletResponse, error)
	GetWallet(ctx context.Context, req *paymentpb.GetWalletRequest) (*paymentpb.Wallet, error)
	HandleDepositWebhook(ctx context.Context, transactionCode string, amount float64) (*paymentpb.WalletTransaction, error)
}

type paymentUsecase struct {
	paymentRepository           repository.PaymentMethodRepository
	walletRepository            repository.WalletsRepository
	walletTransactionRepository repository.WalletTransactionsRepository
	walletAuditLogRepository    repository.WalletAuditLogRepository
	withdrawalRequestRepository repository.WithdrawalRequestsRepository
	notificationWorker          *notificationWorker
}

func NewPaymentUsecase(
	paymentRepository repository.PaymentMethodRepository,
	walletRepository repository.WalletsRepository,
	walletTransactionRepository repository.WalletTransactionsRepository,
	walletAuditLogRepository repository.WalletAuditLogRepository,
	withdrawalRequestRepository repository.WithdrawalRequestsRepository,
	notificationClient client.NotificationClient,
) PaymentUsecase {
	notificationWorker := NewNotificationWorker(notificationClient)
	notificationWorker.Start()
	return &paymentUsecase{
		paymentRepository:           paymentRepository,
		walletRepository:            walletRepository,
		walletTransactionRepository: walletTransactionRepository,
		walletAuditLogRepository:    walletAuditLogRepository,
		withdrawalRequestRepository: withdrawalRequestRepository,
		notificationWorker:          notificationWorker,
	}
}

func (p *paymentUsecase) GetWalletDashboard(ctx context.Context) (*paymentpb.WalletDashboard, error) {
	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	// Get wallet details
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, userID)
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
		TotalIncome:         totalIncome,
		TotalExpense:        totalExpense,
		PendingTransactions: int32(pendingCount),
		Alerts:              alerts,
	}, nil
}

func (p *paymentUsecase) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	transactions, total, err := p.walletTransactionRepository.GetTransactions(ctx, wallet.Id, req)
	if err != nil {
		return nil, err
	}
	return &paymentpb.WalletTransactionList{
		Transactions: transactions,
		Total:        uint32(total),
	}, nil
}

func (p *paymentUsecase) GetWalletTransactionsByWalletId(ctx context.Context, req *paymentpb.GetWalletTransactionsByWalletIdRequest) (*paymentpb.WalletTransactionList, error) {
	wallet, err := p.walletRepository.GetWalletById(ctx, req.WalletId)
	if err != nil {
		return nil, err
	}
	transactions, total, err := p.walletTransactionRepository.GetTransactions(ctx, wallet.Id, &paymentpb.GetWalletTransactionsRequest{
		Page:     req.Page,
		Size:     req.Size,
		Type:     req.Type,
		Status:   req.Status,
		FromDate: req.FromDate,
		ToDate:   req.ToDate,
	})
	if err != nil {
		return nil, err
	}
	return &paymentpb.WalletTransactionList{
		Transactions: transactions,
		Total:        uint32(total),
	}, nil
}

func (p *paymentUsecase) TransferBetweenWallets(ctx context.Context, req *paymentpb.TransferBetweenWalletsRequest) (*paymentpb.TransferBetweenWalletsResponse, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get sender wallet
	senderWallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, req.FromWalletId)
	if err != nil {
		return nil, err
	}

	// Get receiver wallet
	receiverWallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, req.ToWalletId)
	if err != nil {
		return nil, err
	}

	// Check sender balance
	if senderWallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance in sender wallet")
	}

	// Generate transaction code
	senderTransactionCode := fmt.Sprintf("TRF%d%d", time.Now().UnixNano(), senderWallet.Id)
	receiverTransactionCode := fmt.Sprintf("TRF%d%d", time.Now().UnixNano(), receiverWallet.Id)
	// Create sender transaction
	senderTransaction := &entity.WalletTransaction{
		WalletId:        req.FromWalletId,
		Type:            entity.TransactionTypeTransferOut,
		Amount:          -req.Amount,
		Status:          entity.TransactionStatusCompleted,
		RelatedService:  req.RelatedService,
		RelatedId:       req.RelatedId,
		TransactionCode: senderTransactionCode,
		Description:     req.Description,
	}
	err = p.walletTransactionRepository.CreateTransaction(ctx, senderTransaction)
	if err != nil {
		return nil, err
	}

	// Create receiver transaction
	receiverTransaction := &entity.WalletTransaction{
		WalletId:        req.ToWalletId,
		Type:            entity.TransactionTypeTransferIn,
		Amount:          req.Amount,
		Status:          entity.TransactionStatusCompleted,
		RelatedService:  req.RelatedService,
		RelatedId:       req.RelatedId,
		TransactionCode: receiverTransactionCode,
		Description:     req.Description,
	}
	err = p.walletTransactionRepository.CreateTransaction(ctx, receiverTransaction)
	if err != nil {
		return nil, err
	}

	// Update sender balance
	oldSenderBalance := senderWallet.Balance
	senderWallet.Balance -= req.Amount
	err = p.walletRepository.UpdateWallet(ctx, senderWallet)
	if err != nil {
		return nil, err
	}

	// Update receiver balance
	oldReceiverBalance := receiverWallet.Balance
	receiverWallet.Balance += req.Amount
	err = p.walletRepository.UpdateWallet(ctx, receiverWallet)
	if err != nil {
		return nil, err
	}

	// Create audit logs
	senderAuditLog := &entity.WalletAuditLog{
		WalletId:   senderWallet.Id,
		OldBalance: oldSenderBalance,
		NewBalance: senderWallet.Balance,
		Reason:     fmt.Sprintf("Transfer to wallet %d: %s", req.ToWalletId, req.Description),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(ctx, senderAuditLog)
	if err != nil {
		return nil, err
	}

	receiverAuditLog := &entity.WalletAuditLog{
		WalletId:   receiverWallet.Id,
		OldBalance: oldReceiverBalance,
		NewBalance: receiverWallet.Balance,
		Reason:     fmt.Sprintf("Transfer from wallet %d: %s", req.FromWalletId, req.Description),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(ctx, receiverAuditLog)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Send notifications
	go p.notificationWorker.PushToQueue(ctx, senderTransaction, "Chuyển tiền thành công")
	go p.notificationWorker.PushToQueue(ctx, receiverTransaction, "Nhận tiền thành công")

	return &paymentpb.TransferBetweenWalletsResponse{
		FromTransactionId: senderTransaction.Id,
		ToTransactionId:   receiverTransaction.Id,
		Amount:            req.Amount,
		Status:            "COMPLETED",
		TransactionCode:   senderTransactionCode,
		CreatedAt:         senderTransaction.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (p *paymentUsecase) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Generate random transaction code
	transactionCode := fmt.Sprintf("DEP%d%d", time.Now().UnixNano(), wallet.Id)

	// Create transaction record
	transaction := &entity.WalletTransaction{
		WalletId:          wallet.Id,
		Type:              entity.TransactionTypeDeposit,
		Amount:            req.Amount,
		Status:            entity.TransactionStatusPending, // Set initial status as pending
		PaymentMethod:     req.PaymentMethod,
		ExternalPaymentId: req.ExternalPaymentId,
		TransactionCode:   transactionCode,
	}

	err = p.walletTransactionRepository.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	go p.notificationWorker.PushToQueue(ctx, transaction, "Tạo giao dịch nạp tiền thành công")
	return &paymentpb.WalletTransaction{
		Id:              transaction.Id,
		WalletId:        transaction.WalletId,
		Type:            string(transaction.Type),
		Amount:          transaction.Amount,
		Status:          string(transaction.Status),
		TransactionCode: transaction.TransactionCode,
	}, nil
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
		Status:          string(entity.TransactionStatusPending),
		PaymentMethod:   req.PaymentMethod,
		BankAccountInfo: req.BankAccountInfo,
		RequestedBy:     0,
	}

	err = p.withdrawalRequestRepository.CreateWithdrawalRequest(ctx, withdrawal)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &entity.WalletTransaction{
		WalletId:      req.WalletId,
		Type:          entity.TransactionTypeWithdrawal,
		Amount:        -req.Amount,
		Status:        entity.TransactionStatusPending,
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
	go p.notificationWorker.PushToQueue(ctx, transaction, "Rút tiền thành công")
	return &paymentpb.WithdrawalRequest{
		Id:              withdrawal.Id,
		WalletId:        withdrawal.WalletId,
		Amount:          withdrawal.Amount,
		Status:          withdrawal.Status,
		PaymentMethod:   withdrawal.PaymentMethod,
		BankAccountInfo: withdrawal.BankAccountInfo,
	}, nil
}

func (p *paymentUsecase) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	// Get wallet
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check balance
	if wallet.Balance < req.Amount {
		return nil, fmt.Errorf("insufficient balance")
	}
	transactionCode := fmt.Sprintf("DEP%d%d", time.Now().UnixNano(), wallet.Id)

	transaction := &entity.WalletTransaction{
		WalletId:        req.WalletId,
		Type:            "PAYMENT",
		Amount:          -req.Amount,
		Status:          "COMPLETED",
		RelatedService:  req.RelatedService,
		RelatedId:       req.RelatedId,
		TransactionCode: transactionCode,
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
	go p.notificationWorker.PushToQueue(ctx, transaction, "Thanh toán thành công")
	return &paymentpb.WalletTransaction{
		Id:       transaction.Id,
		WalletId: transaction.WalletId,
		Type:     string(transaction.Type),
		Amount:   transaction.Amount,
		Status:   string(transaction.Status),
	}, nil
}

func (p *paymentUsecase) GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error) {
	page := 0
	size := 10
	if req.Page != 0 {
		page = int(req.Page)
	}
	if req.Size != 0 {
		size = int(req.Size)
	}
	methods, total, err := p.paymentRepository.GetPaymentMethods(ctx, req.OrganizationId, page, size)
	if err != nil {
		return nil, err
	}

	return &paymentpb.PaymentMethodList{
		PaymentMethods: methods,
		Total:          uint32(total),
	}, nil
}

func (p *paymentUsecase) CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	method := &entity.PaymentMethod{
		OrganizationId: req.OrganizationId,
		Name:           req.Name,
		Code:           req.Code,
		IsActive:       req.IsActive,
	}

	method, err := p.paymentRepository.CreatePaymentMethod(ctx, method)
	if err != nil {
		return nil, err
	}

	return &paymentpb.PaymentMethod{
		Id:             method.Id,
		OrganizationId: method.OrganizationId,
		Name:           method.Name,
		Code:           method.Code,
		IsActive:       method.IsActive,
		CreatedAt:      method.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      method.UpdatedAt.Format(time.RFC3339),
	}, nil
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

	method, err = p.paymentRepository.UpdatePaymentMethod(ctx, method)
	if err != nil {
		return nil, err
	}

	return &paymentpb.PaymentMethod{
		Id:             method.Id,
		OrganizationId: method.OrganizationId,
		Name:           method.Name,
		Code:           method.Code,
		IsActive:       method.IsActive,
		CreatedAt:      method.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      method.UpdatedAt.Format(time.RFC3339),
	}, nil
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
	totalIncome, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, 0, string(entity.TransactionTypeDealIncome))
	if err != nil {
		return nil, err
	}

	totalExpense, err := p.walletTransactionRepository.GetTotalAmountByType(ctx, 0, string(entity.TransactionTypeDealExpense))
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
		TotalIncome:         totalIncome,
		TotalExpense:        totalExpense,
		PendingTransactions: int32(pendingCount),
		Alerts:              alerts,
	}, nil
}

func (p *paymentUsecase) GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) (*paymentpb.DealTransactionList, error) {
	transactions, total, err := p.walletTransactionRepository.GetDealTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	return &paymentpb.DealTransactionList{
		Transactions: transactions,
		Total:        uint32(total),
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
		Type:           entity.TransactionTypeDealIncome,
		Amount:         req.Amount,
		Status:         entity.TransactionStatusPending,
		RelatedService: "DEAL",
		RelatedId:      fmt.Sprintf("%d", req.Id),
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

func (p *paymentUsecase) CreateWallet(ctx context.Context, req *paymentpb.CreateWalletRequest) (*paymentpb.CreateWalletResponse, error) {
	wallet := &entity.Wallet{
		UserId:   req.UserId,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	wallet, err := p.walletRepository.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return &paymentpb.CreateWalletResponse{
		Id: wallet.Id,
	}, nil
}

func (p *paymentUsecase) GetWallet(ctx context.Context, req *paymentpb.GetWalletRequest) (*paymentpb.Wallet, error) {
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &paymentpb.Wallet{
		Id:       wallet.Id,
		UserId:   wallet.UserId,
		Balance:  wallet.Balance,
		Currency: wallet.Currency,
	}, nil
}

func (p *paymentUsecase) HandleDepositWebhook(ctx context.Context, transactionCode string, amount float64) (*paymentpb.WalletTransaction, error) {
	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer tx.Rollback()

	// Get transaction by code
	transaction, err := p.walletTransactionRepository.GetTransactionByCode(ctx, transactionCode)
	if err != nil {
		return nil, err
	}
	if transaction.Amount != amount {
		return nil, fmt.Errorf("transaction amount mismatch")
	}

	// Check if transaction is already completed
	if transaction.Status == entity.TransactionStatusCompleted {
		return nil, fmt.Errorf("transaction already completed")
	}

	// Get wallet
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(ctx, transaction.WalletId)
	if err != nil {
		return nil, err
	}

	// Update wallet balance
	oldBalance := wallet.Balance
	wallet.Balance += transaction.Amount
	err = p.walletRepository.UpdateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	// Update transaction status
	transaction.Status = entity.TransactionStatusCompleted
	err = p.walletTransactionRepository.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Create audit log
	auditLog := &entity.WalletAuditLog{
		WalletId:   transaction.WalletId,
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Deposit completed via %s", transaction.PaymentMethod),
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
	go p.notificationWorker.PushToQueue(ctx, transaction, "Nạp tiền thành công")
	return &paymentpb.WalletTransaction{
		Id:       transaction.Id,
		WalletId: transaction.WalletId,
		Type:     string(transaction.Type),
		Amount:   transaction.Amount,
		Status:   string(transaction.Status),
	}, nil
}
