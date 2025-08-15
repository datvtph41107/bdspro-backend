package usecase

import (
	"context"
	"fmt"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"
)

type DealTransactionUsecase interface {
	CreateTransaction(ctx context.Context, req *dto.CreateDealTransactionRequest, userID uint64) (*dto.DealTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *dto.UpdateDealTransactionRequest, userID uint64) (*dto.DealTransactionResponse, error)
	DeleteTransaction(ctx context.Context, transactionID uint64, userID uint64) error
	GetContractByID(ctx context.Context, contractId uint64) (*entity.DealTransaction, error)
	GetTransactions(ctx context.Context, req *dto.GetDealTransactionsRequest, userID uint64) (*dto.DealTransactionListResponse, error)
	ApproveTransaction(ctx context.Context, req *dto.ApproveDealTransactionRequest, userID uint64) error
	GetFinancialSummary(ctx context.Context, dealID uint64, userID uint64) (*dto.DealFinancialSummary, error)
	SignContract(ctx context.Context, req *dto.SignContractRequest, userID uint64) (*dto.DealContractResponse, error)
	TransferContract(ctx context.Context, req *dto.TransferContractRequest, userID uint64) (*dto.DealContractResponse, error)
}

type dealTransactionUsecase struct {
	dealTransactionRepo    repository.DealTransactionRepository
	dealRepo               repository.GroupDealRepository
	organizationMemberRepo repository.OrganizationMemberRepository
	notiClient             iusecase.INotiClient
	groupAuthUsecase       GroupAuthUsecase
}

func NewDealTransactionUsecase(
	dealTransactionRepo repository.DealTransactionRepository,
	dealRepo repository.GroupDealRepository,
	organizationMemberRepo repository.OrganizationMemberRepository,
	notiClient iusecase.INotiClient,
	groupAuthUsecase GroupAuthUsecase,
) DealTransactionUsecase {
	return &dealTransactionUsecase{
		dealTransactionRepo:    dealTransactionRepo,
		dealRepo:               dealRepo,
		organizationMemberRepo: organizationMemberRepo,
		notiClient:             notiClient,
		groupAuthUsecase:       groupAuthUsecase,
	}
}

func (u *dealTransactionUsecase) CreateTransaction(ctx context.Context, req *dto.CreateDealTransactionRequest, userID uint64) (*dto.DealTransactionResponse, error) {
	// Todo
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionCreateDeal) {
		return nil, fmt.Errorf("permission denied")
	}
	deal, err := u.dealRepo.GetByID(ctx, req.DealID)
	if err != nil {
		return nil, fmt.Errorf("deal not found: %w", err)
	}

	if !deal.IsManual {
		return nil, fmt.Errorf("deal is not manual")
	}

	// Tạo transaction
	transaction := &entity.DealTransaction{
		DealID:          req.DealID,
		Type:            req.Type,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		Category:        req.Category,
		Note:            req.Note,
		IsManual:        req.IsManual,
		CreatedBy:       userID,
		Status:          enums.ApprovedStatusPending,
	}

	if err := u.dealTransactionRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Gửi thông báo cho admin/phụ trách
	u.sendTransactionNotification(ctx, req.DealID, "NEW_TRANSACTION", transaction.ID)

	return u.toResponse(transaction), nil
}

func (u *dealTransactionUsecase) UpdateTransaction(ctx context.Context, req *dto.UpdateDealTransactionRequest, userID uint64) (*dto.DealTransactionResponse, error) {
	// Lấy transaction hiện tại
	transaction, err := u.dealTransactionRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// Kiểm tra quyền
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionUpdateDeal) {
		return nil, fmt.Errorf("permission denied")
	}

	// Chỉ cho phép sửa khi chưa được approve
	if transaction.Status != enums.ApprovedStatusPending {
		return nil, fmt.Errorf("cannot update approved/rejected transaction")
	}

	// Cập nhật thông tin
	transaction.Amount = req.Amount
	transaction.Description = req.Description
	transaction.TransactionDate = req.TransactionDate
	transaction.Category = req.Category
	transaction.Note = req.Note
	transaction.IsManual = req.IsManual

	if err := u.dealTransactionRepo.Update(ctx, transaction); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return u.toResponse(transaction), nil
}

func (u *dealTransactionUsecase) DeleteTransaction(ctx context.Context, transactionID uint64, userID uint64) error {
	// Lấy transaction
	transaction, err := u.dealTransactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Kiểm tra quyền
	if !u.groupAuthUsecase.HasPermission(ctx, enums.PermissionDeleteDeal) {
		return fmt.Errorf("permission denied")
	}

	// Chỉ cho phép xóa khi chưa được approve
	if transaction.Status != enums.ApprovedStatusPending {
		return fmt.Errorf("cannot delete approved/rejected transaction")
	}

	if err := u.dealTransactionRepo.Delete(ctx, transactionID); err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}

func (u *dealTransactionUsecase) GetContractByID(ctx context.Context, contractId uint64) (*entity.DealTransaction, error) {
	transaction, err := u.dealTransactionRepo.GetByID(ctx, contractId)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	return transaction, nil
}

func (u *dealTransactionUsecase) GetTransactions(ctx context.Context, req *dto.GetDealTransactionsRequest, userID uint64) (*dto.DealTransactionListResponse, error) {
	transactions, total, err := u.dealTransactionRepo.GetByDealID(ctx, req.DealID, req.Type, req.Page, req.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	responses := make([]dto.DealTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *u.toResponse(transaction)
	}

	return &dto.DealTransactionListResponse{
		Transactions: responses,
		Total:        total,
	}, nil
}

func (u *dealTransactionUsecase) ApproveTransaction(ctx context.Context, req *dto.ApproveDealTransactionRequest, userID uint64) error {
	// Lấy transaction
	transaction, err := u.dealTransactionRepo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Xác định status dựa trên rejectReason
	status := enums.ApprovedStatusApproved
	if req.RejectReason != "" {
		status = enums.ApprovedStatusRejected
	}

	if err := u.dealTransactionRepo.Approve(ctx, req.ID, userID, status, req.RejectReason); err != nil {
		return fmt.Errorf("failed to approve transaction: %w", err)
	}

	// Gửi thông báo cho người tạo
	u.sendTransactionNotification(ctx, transaction.DealID, "TRANSACTION_APPROVED", transaction.ID)

	return nil
}

func (u *dealTransactionUsecase) GetFinancialSummary(ctx context.Context, dealID uint64, userID uint64) (*dto.DealFinancialSummary, error) {
	summary, err := u.dealTransactionRepo.GetFinancialSummary(ctx, dealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial summary: %w", err)
	}

	// Thêm cảnh báo
	summary.WarningMessages = u.generateWarnings(summary)

	return summary, nil
}

func (u *dealTransactionUsecase) toResponse(transaction *entity.DealTransaction) *dto.DealTransactionResponse {
	response := &dto.DealTransactionResponse{
		ID:              transaction.ID,
		DealID:          transaction.DealID,
		Type:            transaction.Type,
		TypeName:        enums.TransactionTypeMap[transaction.Type],
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		TransactionDate: transaction.TransactionDate,
		Category:        transaction.Category,
		Note:            transaction.Note,
		IsManual:        transaction.IsManual,
		CreatedBy:       transaction.CreatedBy,
		Status:          transaction.Status,
		StatusName:      enums.ApprovedStatusMap[transaction.Status],
		ApprovedBy:      transaction.ApprovedBy,
		ApprovedAt:      transaction.ApprovedAt,
		RejectReason:    transaction.RejectReason,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
		CustomerID:      transaction.CustomerID,
		TransactionID:   &transaction.ID,
	}

	// Lấy document IDs
	if transaction.AttachDocuments != nil {
		response.DocumentIds = make([]uint64, len(transaction.AttachDocuments))
		for i, doc := range transaction.AttachDocuments {
			response.DocumentIds[i] = doc.ID
		}
	}

	return response
}

func (u *dealTransactionUsecase) generateWarnings(summary *dto.DealFinancialSummary) []string {
	var warnings []string

	// Cảnh báo khi lợi nhuận thấp
	if summary.NetProfit < 0 {
		warnings = append(warnings, "⚠️ Deal đang thua lỗ")
	}

	// Cảnh báo khi lợi nhuận dưới mục tiêu
	if summary.TargetProfit > 0 && summary.NetProfit < summary.TargetProfit*0.8 {
		warnings = append(warnings, "⚠️ Lợi nhuận thấp hơn 80% mục tiêu")
	}

	// Cảnh báo khi chi phí quá cao so với doanh thu
	if summary.TotalRevenue > 0 && summary.TotalExpense > summary.TotalRevenue*0.9 {
		warnings = append(warnings, "⚠️ Chi phí chiếm hơn 90% doanh thu")
	}

	// Cảnh báo khi chưa có doanh thu
	if summary.TotalRevenue == 0 {
		warnings = append(warnings, "ℹ️ Chưa có doanh thu được ghi nhận")
	}

	return warnings
}

func (u *dealTransactionUsecase) sendTransactionNotification(ctx context.Context, dealID uint64, notificationType string, transactionID uint64) {
	// Gửi thông báo cho admin/phụ trách
	// Implementation sẽ được thêm sau
}

func (u *dealTransactionUsecase) SignContract(ctx context.Context, req *dto.SignContractRequest, userID uint64) (*dto.DealContractResponse, error) {
	// Lấy transaction
	// transaction, err := u.dealTransactionRepo.GetByID(ctx, req.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("transaction not found: %w", err)
	// }

	return nil, nil
}

func (u *dealTransactionUsecase) TransferContract(ctx context.Context, req *dto.TransferContractRequest, userID uint64) (*dto.DealContractResponse, error) {
	// Lấy transaction
	// transaction, err := u.dealTransactionRepo.GetByID(ctx, req.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("transaction not found: %w", err)
	// }

	return nil, nil
}
