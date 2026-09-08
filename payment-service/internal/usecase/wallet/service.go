package wallet

import (
	"context"
	"fmt"
	paymentpb "pb/types/payment"
	"time"

	walletdomain "payment/internal/domain/wallet"
	"payment/internal/requestactor"
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
	CountProcessingTransactions(ctx context.Context) (int64, error)
	CountCompletedTransactions(ctx context.Context) (int64, error)
	GetDashboardMetrics(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*paymentpb.DashboardMetric, error)
	CreateOrUpdateDashboardMetric(ctx context.Context, metric *walletdomain.DashboardMetric) error
}

type paymentUsecase struct {
	paymentRepository           PaymentMethodRepository
	walletRepository            WalletRepository
	walletTransactionRepository WalletTransactionRepository
	walletAuditLogRepository    WalletAuditRepository
	withdrawalRequestRepository WithdrawalRepository
	notificationWorker          *NotificationWorker
	dashboardMetricRepository   DashboardMetricRepository
}

func NewPaymentUsecase(
	paymentRepository PaymentMethodRepository,
	walletRepository WalletRepository,
	walletTransactionRepository WalletTransactionRepository,
	walletAuditLogRepository WalletAuditRepository,
	withdrawalRequestRepository WithdrawalRepository,
	notificationWorker *NotificationWorker,
	dashboardMetricRepository DashboardMetricRepository,
) PaymentUsecase {
	return &paymentUsecase{
		paymentRepository:           paymentRepository,
		walletRepository:            walletRepository,
		walletTransactionRepository: walletTransactionRepository,
		walletAuditLogRepository:    walletAuditLogRepository,
		withdrawalRequestRepository: withdrawalRequestRepository,
		notificationWorker:          notificationWorker,
		dashboardMetricRepository:   dashboardMetricRepository,
	}
}

func (p *paymentUsecase) GetWalletDashboard(ctx context.Context) (*paymentpb.WalletDashboard, error) {
	userID := utils.GetUserID(ctx, requestactor.UserContextKey)

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
		Balance:             walletdomain.MinorToWire(wallet.Balance),
		TotalIncome:         walletdomain.MinorToWire(totalIncome),
		TotalExpense:        walletdomain.MinorToWire(totalExpense),
		PendingTransactions: int32(pendingCount),
		Alerts:              alerts,
	}, nil
}

func (p *paymentUsecase) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	userID := utils.GetUserID(ctx, requestactor.UserContextKey)
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
	userID := utils.GetUserID(ctx, requestactor.UserContextKey)

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
	amountMinor, err := walletdomain.MinorFromWire(req.Amount)
	if err != nil || amountMinor <= 0 {
		return nil, fmt.Errorf("invalid amount: must be positive whole minor units")
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
	txCtx := tx.Context(ctx)

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
	if senderWallet.Balance < amountMinor {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance in sender wallet")
	}

	// Generate transaction code with different timestamps to avoid conflicts
	now := time.Now()
	senderTransactionCode := fmt.Sprintf("TRF%d%d", now.UnixNano(), senderWallet.Id)
	time.Sleep(1 * time.Nanosecond) // Ensure different timestamp
	receiverTransactionCode := fmt.Sprintf("TRF%d%d", time.Now().UnixNano(), receiverWallet.Id)

	// Create sender transaction within transaction context
	senderTransaction := &walletdomain.WalletTransaction{
		WalletId:        req.FromWalletId,
		Type:            walletdomain.TransactionTypeTransferOut,
		Amount:          -amountMinor,
		Status:          walletdomain.TransactionStatusCompleted,
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
	receiverTransaction := &walletdomain.WalletTransaction{
		WalletId:        req.ToWalletId,
		Type:            walletdomain.TransactionTypeTransferIn,
		Amount:          amountMinor,
		Status:          walletdomain.TransactionStatusCompleted,
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
	senderWallet.Balance -= amountMinor
	err = p.walletRepository.UpdateWallet(txCtx, senderWallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update sender wallet: %w", err)
	}

	// Update receiver balance within transaction context
	oldReceiverBalance := receiverWallet.Balance
	receiverWallet.Balance += amountMinor
	err = p.walletRepository.UpdateWallet(txCtx, receiverWallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update receiver wallet: %w", err)
	}

	// Create audit logs within transaction context
	senderAuditLog := &walletdomain.WalletAuditLog{
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

	receiverAuditLog := &walletdomain.WalletAuditLog{
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

	// Legacy wallet notifications are explicitly best-effort post-commit.
	_ = p.notificationWorker.PushToQueue(ctx, senderTransaction, "Chuyển tiền thành công")
	_ = p.notificationWorker.PushToQueue(ctx, receiverTransaction, "Nhận tiền thành công")

	return &paymentpb.TransferBetweenWalletsResponse{
		FromTransactionId: senderTransaction.Id,
		ToTransactionId:   receiverTransaction.Id,
		Amount:            walletdomain.MinorToWire(amountMinor),
		Status:            "COMPLETED",
		TransactionCode:   senderTransactionCode,
		CreatedAt:         senderTransaction.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (p *paymentUsecase) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	amountMinor, err := walletdomain.MinorFromWire(req.Amount)
	if err != nil || amountMinor <= 0 {
		return nil, fmt.Errorf("invalid amount: must be positive whole minor units")
	}

	userID := utils.GetUserID(ctx, requestactor.UserContextKey)

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := tx.Context(ctx)

	// Get wallet within transaction context
	wallet, err := p.walletRepository.GetWalletByUserId(txCtx, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Generate random transaction code
	transactionCode := fmt.Sprintf("DEP%d%d", time.Now().UnixNano(), wallet.Id)

	// Create transaction record within transaction context
	transaction := &walletdomain.WalletTransaction{
		WalletId:          wallet.Id,
		Type:              walletdomain.TransactionTypeDeposit,
		Amount:            amountMinor,
		Status:            walletdomain.TransactionStatusPending, // Set initial status as pending
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

	// Legacy wallet notification is explicitly best-effort post-commit.
	_ = p.notificationWorker.PushToQueue(ctx, transaction, "Tạo giao dịch nạp tiền thành công")

	return &paymentpb.WalletTransaction{
		Id:              transaction.Id,
		WalletId:        transaction.WalletId,
		Type:            string(transaction.Type),
		Amount:          walletdomain.MinorToWire(transaction.Amount),
		Status:          string(transaction.Status),
		TransactionCode: transaction.TransactionCode,
	}, nil
}

func (p *paymentUsecase) Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error) {
	amountMinor, err := walletdomain.MinorFromWire(req.Amount)
	if err != nil || amountMinor <= 0 {
		return nil, fmt.Errorf("invalid amount: must be positive whole minor units")
	}

	// Get current user ID for security validation
	userID := utils.GetUserID(ctx, requestactor.UserContextKey)

	// Start transaction
	tx := p.walletRepository.BeginTx(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Create transaction context
	txCtx := tx.Context(ctx)

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
	if wallet.Balance < amountMinor {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Create withdrawal request within transaction context
	withdrawal := &walletdomain.WithdrawalRequest{
		WalletId:        req.WalletId,
		Amount:          amountMinor,
		Status:          string(walletdomain.TransactionStatusPending),
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
	transaction := &walletdomain.WalletTransaction{
		WalletId:      req.WalletId,
		Type:          walletdomain.TransactionTypeWithdrawal,
		Amount:        -amountMinor,
		Status:        walletdomain.TransactionStatusPending,
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

	// Legacy wallet notification is explicitly best-effort post-commit.
	_ = p.notificationWorker.PushToQueue(ctx, transaction, "Rút tiền thành công")

	return &paymentpb.WithdrawalRequest{
		Id:              withdrawal.Id,
		WalletId:        withdrawal.WalletId,
		Amount:          walletdomain.MinorToWire(withdrawal.Amount),
		Status:          withdrawal.Status,
		PaymentMethod:   withdrawal.PaymentMethod,
		BankAccountInfo: withdrawal.BankAccountInfo,
	}, nil
}

func (p *paymentUsecase) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	amountMinor, err := walletdomain.MinorFromWire(req.Amount)
	if err != nil || amountMinor <= 0 {
		return nil, fmt.Errorf("invalid amount: must be positive whole minor units")
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
	txCtx := tx.Context(ctx)

	userID := utils.GetUserID(ctx, requestactor.UserContextKey)
	// Get wallet with proper locking - use correct method
	wallet, err := p.walletRepository.GetWalletByUserId(txCtx, userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check balance
	if wallet.Balance < amountMinor {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance")
	}

	// Generate transaction code with correct prefix
	transactionCode := fmt.Sprintf("PAY%d%d", time.Now().UnixNano(), wallet.Id)

	transaction := &walletdomain.WalletTransaction{
		WalletId:        wallet.Id, // ✅ Fixed: Use wallet.Id instead of req.WalletId
		Type:            "PAYMENT",
		Amount:          -amountMinor,
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
	wallet.Balance -= amountMinor
	err = p.walletRepository.UpdateWallet(txCtx, wallet)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create audit log within transaction context
	auditLog := &walletdomain.WalletAuditLog{
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

	// Legacy wallet notification is explicitly best-effort post-commit.
	_ = p.notificationWorker.PushToQueue(ctx, transaction, "Thanh toán thành công")

	return &paymentpb.WalletTransaction{
		Id:       transaction.Id,
		WalletId: transaction.WalletId,
		Type:     string(transaction.Type),
		Amount:   walletdomain.MinorToWire(transaction.Amount),
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
	method := &walletdomain.PaymentMethod{
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
		TotalAmount:       walletdomain.MinorToWire(totalAmount),
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
	initialBalance, err := walletdomain.MinorFromWire(req.Balance)
	if err != nil || initialBalance < 0 {
		return nil, fmt.Errorf("initial balance must be non-negative whole minor units")
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

	wallet := &walletdomain.Wallet{
		UserId:   req.UserId,
		Balance:  initialBalance,
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
		Balance:  walletdomain.MinorToWire(wallet.Balance),
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
		Balance:  walletdomain.MinorToWire(wallet.Balance),
		Currency: wallet.Currency,
	}, nil
}

func (p *paymentUsecase) HandleDepositWebhook(ctx context.Context, transactionCode string, amount float64) (*paymentpb.WalletTransaction, error) {
	// Validate inputs
	if transactionCode == "" {
		return nil, fmt.Errorf("transaction code cannot be empty")
	}
	amountMinor, err := walletdomain.MinorFromWire(amount)
	if err != nil || amountMinor <= 0 {
		return nil, fmt.Errorf("amount must be positive whole minor units")
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
	txCtx := tx.Context(ctx)

	// Get transaction by code within transaction context
	transaction, err := p.walletTransactionRepository.GetTransactionByCode(txCtx, transactionCode)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get transaction by code: %w", err)
	}

	// Validate amount
	if transaction.Amount != amountMinor {
		tx.Rollback()
		return nil, fmt.Errorf("transaction amount mismatch: expected %d minor units, got %d", transaction.Amount, amountMinor)
	}

	// Check if transaction is already completed
	if transaction.Status == walletdomain.TransactionStatusCompleted {
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
	transaction.Status = walletdomain.TransactionStatusCompleted
	err = p.walletTransactionRepository.UpdateTransaction(txCtx, transaction)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Create audit log within transaction context
	auditLog := &walletdomain.WalletAuditLog{
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

	// Legacy wallet notification is explicitly best-effort post-commit.
	_ = p.notificationWorker.PushToQueue(ctx, transaction, "Nạp tiền thành công")

	return &paymentpb.WalletTransaction{
		Id:       transaction.Id,
		WalletId: transaction.WalletId,
		Type:     string(transaction.Type),
		Amount:   walletdomain.MinorToWire(transaction.Amount),
		Status:   string(transaction.Status),
	}, nil
}

func (uc *paymentUsecase) CountProcessingTransactions(ctx context.Context) (int64, error) {
	// Đếm số giao dịch có status PENDING (đang xử lý)
	count, err := uc.walletTransactionRepository.CountTransactionsByStatus(ctx, "PENDING")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (uc *paymentUsecase) CountCompletedTransactions(ctx context.Context) (int64, error) {
	// Đếm số giao dịch có status COMPLETED
	count, err := uc.walletTransactionRepository.CountTransactionsByStatus(ctx, "COMPLETED")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *paymentUsecase) GetDashboardMetrics(ctx context.Context, fromDate time.Time, toDate time.Time) ([]*paymentpb.DashboardMetric, error) {
	metrics, err := p.dashboardMetricRepository.GetDashboardMetrics(ctx, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	metricsResponse := make([]*paymentpb.DashboardMetric, len(metrics))
	for i, metric := range metrics {
		metricsResponse[i] = &paymentpb.DashboardMetric{
			Time:            metric.Time.Format(time.RFC3339),
			ProcessingCount: metric.ProcessingCount,
			CompletedCount:  metric.CompletedCount,
		}
	}
	return metricsResponse, nil
}

func (p *paymentUsecase) CreateOrUpdateDashboardMetric(ctx context.Context, metric *walletdomain.DashboardMetric) error {
	_, err := p.dashboardMetricRepository.CreateOrUpdate(ctx, metric)
	if err != nil {
		return err
	}
	return nil
}
