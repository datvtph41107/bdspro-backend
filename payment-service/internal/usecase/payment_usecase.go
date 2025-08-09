package usecase

import (
	"context"
	"fmt"
	paymentpb "pb/types/payment"
	"time"

	"payment/appconfig"
	"payment/internal/domain/entity"
	"payment/internal/domain/repository"
	"payment/internal/infrastructure/client"
	"payment/pkg/utils"

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

	GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error)
	CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error)
	UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error)
	DeleteTransactionType(ctx context.Context, req *paymentpb.DeleteTransactionTypeRequest) (*paymentpb.Empty, error)
	CreateWallet(ctx context.Context, req *paymentpb.CreateWalletRequest) (*paymentpb.CreateWalletResponse, error)
	GetWalletByUserId(ctx context.Context, req *paymentpb.GetWalletByUserIdRequest) (*paymentpb.Wallet, error)
	GetWalletByWalletId(ctx context.Context, req *paymentpb.GetWalletByWalletIdRequest) (*paymentpb.Wallet, error)
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
		Data:  transactions,
		Total: uint32(total),
	}, nil
}

func (p *paymentUsecase) GetWalletTransactionsByWalletId(ctx context.Context, req *paymentpb.GetWalletTransactionsByWalletIdRequest) (*paymentpb.WalletTransactionList, error) {
	// Get current user ID for security validation
	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	wallet, err := p.walletRepository.GetWalletById(ctx, req.WalletId)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Security check: Ensure user owns the wallet (for non-admin users)
	// Note: In a real system, you might want to add admin role check here
	if wallet.UserId != userID {
		return nil, fmt.Errorf("unauthorized access to wallet")
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
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return &paymentpb.WalletTransactionList{
		Data:  transactions,
		Total: uint32(total),
	}, nil
}

func (p *paymentUsecase) TransferBetweenWallets(ctx context.Context, req *paymentpb.TransferBetweenWalletsRequest) (*paymentpb.TransferBetweenWalletsResponse, error) {
	// Validate request
	if req.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: must be greater than 0")
	}

	if req.FromWalletId == req.ToWalletId {
		return nil, fmt.Errorf("cannot transfer to the same wallet")
	}

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := context.WithValue(ctx, "tx", tx)

	// Get sender wallet with proper locking
	senderWallet, err := p.walletRepository.GetWalletByIdForUpdate(txCtx, req.FromWalletId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get sender wallet: %w", err)
	}

	// Get receiver wallet with proper locking
	receiverWallet, err := p.walletRepository.GetWalletByIdForUpdate(txCtx, req.ToWalletId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get receiver wallet: %w", err)
	}

	// Check sender balance
	if senderWallet.Balance < req.Amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance in sender wallet")
	}

	// Generate transaction code with different timestamps to avoid conflicts
	now := time.Now()
	senderTransactionCode := fmt.Sprintf("TRF%d%d", now.UnixNano(), senderWallet.Id)
	time.Sleep(1 * time.Nanosecond) // Ensure different timestamp
	receiverTransactionCode := fmt.Sprintf("TRF%d%d", time.Now().UnixNano(), receiverWallet.Id)

	// Create sender transaction within transaction context
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
	err = p.walletTransactionRepository.CreateTransaction(txCtx, senderTransaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create sender transaction: %w", err)
	}

	// Create receiver transaction within transaction context
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
	err = p.walletTransactionRepository.CreateTransaction(txCtx, receiverTransaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create receiver transaction: %w", err)
	}

	// Update sender balance within transaction context
	oldSenderBalance := senderWallet.Balance
	senderWallet.Balance -= req.Amount
	err = p.walletRepository.UpdateWallet(txCtx, senderWallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update sender wallet: %w", err)
	}

	// Update receiver balance within transaction context
	oldReceiverBalance := receiverWallet.Balance
	receiverWallet.Balance += req.Amount
	err = p.walletRepository.UpdateWallet(txCtx, receiverWallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update receiver wallet: %w", err)
	}

	// Create audit logs within transaction context
	senderAuditLog := &entity.WalletAuditLog{
		WalletId:   senderWallet.Id,
		OldBalance: oldSenderBalance,
		NewBalance: senderWallet.Balance,
		Reason:     fmt.Sprintf("Transfer to wallet %d: %s", req.ToWalletId, req.Description),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(txCtx, senderAuditLog)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create sender audit log: %w", err)
	}

	receiverAuditLog := &entity.WalletAuditLog{
		WalletId:   receiverWallet.Id,
		OldBalance: oldReceiverBalance,
		NewBalance: receiverWallet.Balance,
		Reason:     fmt.Sprintf("Transfer from wallet %d: %s", req.FromWalletId, req.Description),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(txCtx, receiverAuditLog)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create receiver audit log: %w", err)
	}

	// Commit transaction - if this fails, everything is rolled back
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notifications (outside transaction - fire and forget)
	go func() {
		p.notificationWorker.PushToQueue(ctx, senderTransaction, "Chuyển tiền thành công")
		p.notificationWorker.PushToQueue(ctx, receiverTransaction, "Nhận tiền thành công")
	}()

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
	// Validate amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: must be greater than 0")
	}

	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := context.WithValue(ctx, "tx", tx)

	// Get wallet within transaction context
	wallet, err := p.walletRepository.GetWalletByUserId(txCtx, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Generate random transaction code
	transactionCode := fmt.Sprintf("DEP%d%d", time.Now().UnixNano(), wallet.Id)

	// Create transaction record within transaction context
	transaction := &entity.WalletTransaction{
		WalletId:          wallet.Id,
		Type:              entity.TransactionTypeDeposit,
		Amount:            req.Amount,
		Status:            entity.TransactionStatusPending, // Set initial status as pending
		PaymentMethod:     req.PaymentMethod,
		ExternalPaymentId: req.ExternalPaymentId,
		TransactionCode:   transactionCode,
	}

	err = p.walletTransactionRepository.CreateTransaction(txCtx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notification (outside transaction - fire and forget)
	go func() {
		p.notificationWorker.PushToQueue(ctx, transaction, "Tạo giao dịch nạp tiền thành công")
	}()

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
	// Validate amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: must be greater than 0")
	}

	// Get current user ID for security validation
	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := context.WithValue(ctx, "tx", tx)

	// Get wallet with proper locking within transaction context
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(txCtx, req.WalletId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Security check: Ensure user owns the wallet
	if wallet.UserId != userID {
		tx.Rollback()
		return nil, fmt.Errorf("unauthorized access to wallet")
	}

	// Check balance
	if wallet.Balance < req.Amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create withdrawal request within transaction context
	withdrawal := &entity.WithdrawalRequest{
		WalletId:        req.WalletId,
		Amount:          req.Amount,
		Status:          string(entity.TransactionStatusPending),
		PaymentMethod:   req.PaymentMethod,
		BankAccountInfo: req.BankAccountInfo,
		RequestedBy:     userID, // Use actual user ID instead of 0
	}

	err = p.withdrawalRequestRepository.CreateWithdrawalRequest(txCtx, withdrawal)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create withdrawal request: %w", err)
	}

	// Create transaction record within transaction context
	transaction := &entity.WalletTransaction{
		WalletId:      req.WalletId,
		Type:          entity.TransactionTypeWithdrawal,
		Amount:        -req.Amount,
		Status:        entity.TransactionStatusPending,
		PaymentMethod: req.PaymentMethod,
	}

	err = p.walletTransactionRepository.CreateTransaction(txCtx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notification (outside transaction - fire and forget)
	go func() {
		p.notificationWorker.PushToQueue(ctx, transaction, "Rút tiền thành công")
	}()

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
	// Validate amount
	if req.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: must be greater than 0")
	}

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := context.WithValue(ctx, "tx", tx)

	userID := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	// Get wallet with proper locking - use correct method
	wallet, err := p.walletRepository.GetWalletByUserId(txCtx, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check balance
	if wallet.Balance < req.Amount {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Generate transaction code with correct prefix
	transactionCode := fmt.Sprintf("PAY%d%d", time.Now().UnixNano(), wallet.Id)

	transaction := &entity.WalletTransaction{
		WalletId:        wallet.Id, // ✅ Fixed: Use wallet.Id instead of req.WalletId
		Type:            "PAYMENT",
		Amount:          -req.Amount,
		Status:          "COMPLETED",
		RelatedService:  req.RelatedService,
		RelatedId:       req.RelatedId,
		TransactionCode: transactionCode,
	}

	err = p.walletTransactionRepository.CreateTransaction(txCtx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update wallet balance within transaction context
	oldBalance := wallet.Balance
	wallet.Balance -= req.Amount
	err = p.walletRepository.UpdateWallet(txCtx, wallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create audit log within transaction context
	auditLog := &entity.WalletAuditLog{
		WalletId:   wallet.Id, // ✅ Fixed: Use wallet.Id instead of req.WalletId
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Payment for %s:%s", req.RelatedService, req.RelatedId),
		ChangedBy:  userID, // Use actual user ID instead of 0
	}
	err = p.walletAuditLogRepository.CreateAuditLog(txCtx, auditLog)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create audit log: %w", err)
	}

	// Commit transaction - if this fails, everything is rolled back
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notification (outside transaction - fire and forget)
	go func() {
		p.notificationWorker.PushToQueue(ctx, transaction, "Thanh toán thành công")
	}()

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
		Data:  methods,
		Total: uint32(total),
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

func (p *paymentUsecase) GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error) {
	types, err := p.walletTransactionRepository.GetTransactionTypes(ctx)
	if err != nil {
		return nil, err
	}

	return &paymentpb.TransactionTypeList{
		Data: types,
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
	// Validate request
	if req.UserId == 0 {
		return nil, fmt.Errorf("user ID cannot be zero")
	}
	if req.Balance < 0 {
		return nil, fmt.Errorf("initial balance cannot be negative")
	}
	if req.Currency == "" {
		return nil, fmt.Errorf("currency cannot be empty")
	}

	// Check if wallet already exists for this user
	existingWallet, err := p.walletRepository.GetWalletByUserId(ctx, req.UserId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing wallet: %w", err)
	}
	if existingWallet != nil {
		return nil, fmt.Errorf("wallet already exists for user %d", req.UserId)
	}

	wallet := &entity.Wallet{
		UserId:   req.UserId,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	wallet, err = p.walletRepository.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &paymentpb.CreateWalletResponse{
		Id: wallet.Id,
	}, nil
}

func (p *paymentUsecase) GetWalletByUserId(ctx context.Context, req *paymentpb.GetWalletByUserIdRequest) (*paymentpb.Wallet, error) {
	wallet, err := p.walletRepository.GetWalletByUserId(ctx, req.UserId)
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

func (p *paymentUsecase) GetWalletByWalletId(ctx context.Context, req *paymentpb.GetWalletByWalletIdRequest) (*paymentpb.Wallet, error) {
	wallet, err := p.walletRepository.GetWalletById(ctx, req.WalletId)
	if err != nil {
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
	// Validate inputs
	if transactionCode == "" {
		return nil, fmt.Errorf("transaction code cannot be empty")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than 0")
	}

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := context.WithValue(ctx, "tx", tx)

	// Get transaction by code within transaction context
	transaction, err := p.walletTransactionRepository.GetTransactionByCode(txCtx, transactionCode)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get transaction by code: %w", err)
	}

	// Validate amount
	if transaction.Amount != amount {
		tx.Rollback()
		return nil, fmt.Errorf("transaction amount mismatch: expected %f, got %f", transaction.Amount, amount)
	}

	// Check if transaction is already completed
	if transaction.Status == entity.TransactionStatusCompleted {
		tx.Rollback()
		return nil, fmt.Errorf("transaction already completed")
	}

	// Get wallet with proper locking within transaction context
	wallet, err := p.walletRepository.GetWalletByIdForUpdate(txCtx, transaction.WalletId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Update wallet balance within transaction context
	oldBalance := wallet.Balance
	wallet.Balance += transaction.Amount
	err = p.walletRepository.UpdateWallet(txCtx, wallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Update transaction status within transaction context
	transaction.Status = entity.TransactionStatusCompleted
	err = p.walletTransactionRepository.UpdateTransaction(txCtx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Create audit log within transaction context
	auditLog := &entity.WalletAuditLog{
		WalletId:   transaction.WalletId,
		OldBalance: oldBalance,
		NewBalance: wallet.Balance,
		Reason:     fmt.Sprintf("Deposit completed via %s", transaction.PaymentMethod),
		ChangedBy:  0, // System
	}
	err = p.walletAuditLogRepository.CreateAuditLog(txCtx, auditLog)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create audit log: %w", err)
	}

	// Commit transaction - if this fails, everything is rolled back
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notification (outside transaction - fire and forget)
	go func() {
		p.notificationWorker.PushToQueue(ctx, transaction, "Nạp tiền thành công")
	}()

	return &paymentpb.WalletTransaction{
		Id:       transaction.Id,
		WalletId: transaction.WalletId,
		Type:     string(transaction.Type),
		Amount:   transaction.Amount,
		Status:   string(transaction.Status),
	}, nil
}
