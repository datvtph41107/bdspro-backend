package usecase

import (
	"context"
	"fmt"

	"marketing/infrastructure/client"
	"marketing/internal/domain/model"
	"marketing/internal/domain/repository"
	"marketing/internal/dto"

	_utils "common/utils"
	shared_enum "pb/enums"
	notificationpb "pb/types/notification"

	"github.com/hyperledger/fabric/common/flogging"
)

type CampaignUsecase interface {
	CreateCampaign(ctx context.Context, createDTO *dto.CampaignCreateDTO) (*dto.CampaignResponseDTO, error)
	UpdateCampaign(ctx context.Context, id uint64, updateDTO *dto.CampaignUpdateDTO) (*dto.CampaignResponseDTO, error)
	DeleteCampaign(ctx context.Context, id uint64) error
	GetCampaign(ctx context.Context, id uint64) (*dto.CampaignResponseDTO, error)
	SearchCampaigns(ctx context.Context, searchDTO dto.CampaignSearchDTO) ([]dto.CampaignResponseDTO, int64, error)
	GetCampaignsByProduct(ctx context.Context, productID uint64) ([]dto.CampaignResponseDTO, error)
	UpdateCampaignStatus(ctx context.Context, id uint64, status model.CampaignStatus) error

	// Budget management
	TopupCampaignBudget(ctx context.Context, topupDTO *dto.CampaignBudgetTopupDTO) (*dto.CampaignBudgetTopupResponseDTO, error)
	GetCampaignWalletInfo(ctx context.Context, campaignID uint64) (*dto.CampaignWalletInfoDTO, error)
	GetCampaignTransactions(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.CampaignBudgetTopupResponseDTO, int64, error)
}

type campaignUsecase struct {
	campaignRepo       repository.CampaignRepository
	paymentClient      *client.PaymentClient
	notificationClient *client.NotificationClient
	logger             *flogging.FabricLogger
}

func NewCampaignUsecase(campaignRepo repository.CampaignRepository, paymentClient *client.PaymentClient, notificationClient *client.NotificationClient) CampaignUsecase {
	return &campaignUsecase{
		campaignRepo:       campaignRepo,
		paymentClient:      paymentClient,
		notificationClient: notificationClient,
		logger:             flogging.MustGetLogger("campaignUsecase"),
	}
}

// createHistoryRecord helper function to create history records
func (u *campaignUsecase) createHistoryRecord(ctx context.Context, actionType shared_enum.EHistory, targetId uint64, targetType shared_enum.ETargetHistory, title string, notes []string) {
	if u.notificationClient == nil {
		return
	}

	historyReq := &notificationpb.HistoryDTO{
		ActionType: int32(actionType),
		TargetId:   targetId,
		TargetType: int32(targetType),
		Title:      title,
		Note:       notes,
		OwnerType:  int32(shared_enum.EOwnerTypeOrgnization),
	}

	_, err := u.notificationClient.CreateHistory(ctx, historyReq)
	if err != nil {
		// Log error but don't fail the main operation
		fmt.Printf("Failed to create history record: %v\n", err)
	}
}

// CreateCampaign creates a new advertising campaign and creates a wallet for it
func (u *campaignUsecase) CreateCampaign(ctx context.Context, createDTO *dto.CampaignCreateDTO) (*dto.CampaignResponseDTO, error) {
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	userID := _utils.GetUserIDFromContext(ctx)

	// Validate terms acceptance
	if !createDTO.TermsAccepted {
		return nil, fmt.Errorf("terms must be accepted to create campaign")
	}

	// Check time conflict for same product and campaign type
	if createDTO.StartDate != nil && createDTO.EndDate != nil {
		hasConflict, err := u.campaignRepo.CheckTimeConflict(ctx, *createDTO.ProductID, createDTO.Type, createDTO.StartDate, createDTO.EndDate, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check time conflict: %w", err)
		}
		if hasConflict {
			return nil, fmt.Errorf("campaign time conflicts with existing campaign for this product")
		}
	}

	// Create campaign entity
	campaign := &model.Campaign{
		Name:           createDTO.Name,
		Description:    createDTO.Description,
		Type:           createDTO.Type,
		Status:         model.CampaignStatusDraft,
		ProductID:      createDTO.ProductID,
		ProductName:    createDTO.ProductName,
		ProductType:    createDTO.ProductType,
		PackageID:      createDTO.PackageID,
		PackageName:    createDTO.PackageName,
		Budget:         createDTO.Budget,
		Duration:       createDTO.Duration,
		StartDate:      createDTO.StartDate,
		EndDate:        createDTO.EndDate,
		TargetAudience: createDTO.TargetAudience,
		TargetLocation: createDTO.TargetLocation,
		IsActive:       true,
		Priority:       createDTO.Priority,
		InternalNote:   createDTO.InternalNote,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		OrganizationID: organizationID,
		Currency:       "VND", // Default currency
	}

	// Save campaign to database first
	createdCampaign, err := u.campaignRepo.Create(ctx, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	// Create wallet for the campaign
	walletReq := &client.CreateWalletRequest{
		OwnerID:       createdCampaign.ID,
		OwnerType:     "CAMPAIGN",
		Currency:      createdCampaign.Currency,
		Description:   fmt.Sprintf("Wallet for campaign: %s", createdCampaign.Name),
		CampaignID:    createdCampaign.ID,
		InitialAmount: 0, // Start with 0 balance, user will deposit later
	}

	walletResp, err := u.paymentClient.CreateWallet(ctx, walletReq)
	if err != nil {
		// Log error but don't fail the campaign creation
		// TODO: Add proper logging
		fmt.Printf("Warning: Failed to create wallet for campaign %d: %v\n", createdCampaign.ID, err)
	} else {
		// Update campaign with wallet information
		createdCampaign.WalletID = &walletResp.Data.WalletID
		createdCampaign.WalletBalance = walletResp.Data.Balance

		// Update campaign in database with wallet info
		_, err = u.campaignRepo.Update(ctx, createdCampaign.ID, createdCampaign)
		if err != nil {
			// Log error but don't fail the response
			fmt.Printf("Warning: Failed to update campaign with wallet info: %v\n", err)
		}
	}

	// Convert to response DTO
	response := u.toResponseDTO(createdCampaign)

	// Create history record
	u.createHistoryRecord(ctx,
		shared_enum.HistoryCampaignCreate,
		createdCampaign.ID,
		shared_enum.TargetHistoryCampaign,
		"Tạo chiến dịch quảng cáo",
		[]string{createdCampaign.Name, fmt.Sprintf("Ngân sách: %.0f %s", createdCampaign.Budget, createdCampaign.Currency)},
	)

	return response, nil
}

// UpdateCampaign updates an existing advertising campaign
func (u *campaignUsecase) UpdateCampaign(ctx context.Context, id uint64, updateDTO *dto.CampaignUpdateDTO) (*dto.CampaignResponseDTO, error) {
	userID := _utils.GetUserIDFromContext(ctx)

	// Get existing campaign
	existingCampaign, err := u.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check time conflict if dates are being updated
	if updateDTO.StartDate != nil && updateDTO.EndDate != nil {
		hasConflict, err := u.campaignRepo.CheckTimeConflict(ctx, *existingCampaign.ProductID, existingCampaign.Type, updateDTO.StartDate, updateDTO.EndDate, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check time conflict: %w", err)
		}
		if hasConflict {
			return nil, fmt.Errorf("campaign time conflicts with existing campaign for this product")
		}
	}

	// Update fields
	if updateDTO.Name != "" {
		existingCampaign.Name = updateDTO.Name
	}
	if updateDTO.Description != "" {
		existingCampaign.Description = updateDTO.Description
	}
	if updateDTO.Type != "" {
		existingCampaign.Type = updateDTO.Type
	}
	if updateDTO.Status != "" {
		existingCampaign.Status = updateDTO.Status
	}
	if updateDTO.ProductID != nil {
		existingCampaign.ProductID = updateDTO.ProductID
	}
	if updateDTO.ProductName != "" {
		existingCampaign.ProductName = updateDTO.ProductName
	}
	if updateDTO.ProductType != "" {
		existingCampaign.ProductType = updateDTO.ProductType
	}
	if updateDTO.PackageID != nil {
		existingCampaign.PackageID = updateDTO.PackageID
	}
	if updateDTO.PackageName != "" {
		existingCampaign.PackageName = updateDTO.PackageName
	}
	if updateDTO.Budget > 0 {
		existingCampaign.Budget = updateDTO.Budget
	}
	if updateDTO.Duration > 0 {
		existingCampaign.Duration = updateDTO.Duration
	}
	if updateDTO.StartDate != nil {
		existingCampaign.StartDate = updateDTO.StartDate
	}
	if updateDTO.EndDate != nil {
		existingCampaign.EndDate = updateDTO.EndDate
	}
	if updateDTO.TargetAudience != "" {
		existingCampaign.TargetAudience = updateDTO.TargetAudience
	}
	if updateDTO.TargetLocation != "" {
		existingCampaign.TargetLocation = updateDTO.TargetLocation
	}
	if updateDTO.IsActive != nil {
		existingCampaign.IsActive = *updateDTO.IsActive
	}
	existingCampaign.Priority = updateDTO.Priority
	existingCampaign.InternalNote = updateDTO.InternalNote
	existingCampaign.UpdatedBy = userID

	// Save to database
	updatedCampaign, err := u.campaignRepo.Update(ctx, id, existingCampaign)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	response := u.toResponseDTO(updatedCampaign)

	// Create history record
	u.createHistoryRecord(ctx,
		shared_enum.HistoryCampaignUpdate,
		updatedCampaign.ID,
		shared_enum.TargetHistoryCampaign,
		"Cập nhật chiến dịch quảng cáo",
		[]string{updatedCampaign.Name, fmt.Sprintf("Ngân sách: %.0f %s", updatedCampaign.Budget, updatedCampaign.Currency)},
	)

	return response, nil
}

// DeleteCampaign deletes an advertising campaign (soft delete)
func (u *campaignUsecase) DeleteCampaign(ctx context.Context, id uint64) error {
	// Get campaign info before deletion for history
	campaign, err := u.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete campaign
	err = u.campaignRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Create history record
	u.createHistoryRecord(ctx,
		shared_enum.HistoryCampaignDelete,
		id,
		shared_enum.TargetHistoryCampaign,
		"Xóa chiến dịch quảng cáo",
		[]string{campaign.Name},
	)

	return nil
}

// GetCampaign retrieves a campaign by ID
func (u *campaignUsecase) GetCampaign(ctx context.Context, id uint64) (*dto.CampaignResponseDTO, error) {
	campaign, err := u.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.toResponseDTO(campaign), nil
}

// SearchCampaigns searches for campaigns with filters
func (u *campaignUsecase) SearchCampaigns(ctx context.Context, searchDTO dto.CampaignSearchDTO) ([]dto.CampaignResponseDTO, int64, error) {
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	campaigns, total, err := u.campaignRepo.Search(ctx, organizationID, searchDTO)
	if err != nil {
		return nil, 0, err
	}

	// Convert to response DTOs
	responseDTOs := make([]dto.CampaignResponseDTO, len(campaigns))
	for i, campaign := range campaigns {
		responseDTOs[i] = *u.toResponseDTO(campaign)
	}

	return responseDTOs, total, nil
}

// GetCampaignsByProduct retrieves campaigns for a specific product
func (u *campaignUsecase) GetCampaignsByProduct(ctx context.Context, productID uint64) ([]dto.CampaignResponseDTO, error) {
	campaigns, err := u.campaignRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responseDTOs := make([]dto.CampaignResponseDTO, len(campaigns))
	for i, campaign := range campaigns {
		responseDTOs[i] = *u.toResponseDTO(campaign)
	}

	return responseDTOs, nil
}

// UpdateCampaignStatus updates the status of a campaign
func (u *campaignUsecase) UpdateCampaignStatus(ctx context.Context, id uint64, status model.CampaignStatus) error {
	// Get campaign info before update for history
	campaign, err := u.campaignRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	oldStatus := campaign.Status

	// Update status
	err = u.campaignRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	// Create history record based on status change
	var actionType shared_enum.EHistory
	var title string
	var notes []string

	switch status {
	case model.CampaignStatusActive:
		actionType = shared_enum.HistoryCampaignActivate
		title = "Kích hoạt chiến dịch quảng cáo"
		notes = []string{campaign.Name, "Trạng thái: ACTIVE"}
	case model.CampaignStatusPaused:
		actionType = shared_enum.HistoryCampaignPause
		title = "Tạm dừng chiến dịch quảng cáo"
		notes = []string{campaign.Name, "Trạng thái: PAUSED"}
	case model.CampaignStatusCompleted:
		actionType = shared_enum.HistoryCampaignUpdate
		title = "Hoàn thành chiến dịch quảng cáo"
		notes = []string{campaign.Name, "Trạng thái: COMPLETED"}
	default:
		actionType = shared_enum.HistoryCampaignUpdate
		title = "Cập nhật trạng thái chiến dịch"
		notes = []string{campaign.Name, fmt.Sprintf("Từ %s → %s", oldStatus, status)}
	}

	u.createHistoryRecord(ctx, actionType, id, shared_enum.TargetHistoryCampaign, title, notes)

	return nil
}

// TopupCampaignBudget adds budget to campaign wallet
func (u *campaignUsecase) TopupCampaignBudget(ctx context.Context, topupDTO *dto.CampaignBudgetTopupDTO) (*dto.CampaignBudgetTopupResponseDTO, error) {
	// Get campaign to validate and get wallet info
	campaign, err := u.campaignRepo.GetByID(ctx, topupDTO.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate campaign ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Check if campaign has wallet
	if campaign.WalletID == nil {
		return nil, fmt.Errorf("campaign does not have a wallet")
	}

	// Parse wallet ID to get numeric ID
	var walletIDNum uint32
	_, err = fmt.Sscanf(*campaign.WalletID, "wallet_%d", &walletIDNum)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID format: %s", *campaign.WalletID)
	}

	// Call payment service to deposit
	depositReq := &client.DepositRequest{
		WalletID:          uint64(walletIDNum),
		Amount:            topupDTO.Amount,
		PaymentMethod:     topupDTO.PaymentMethod,
		ExternalPaymentID: topupDTO.ExternalPaymentID,
	}

	transaction, err := u.paymentClient.Deposit(ctx, depositReq)
	if err != nil {
		return nil, fmt.Errorf("failed to deposit to wallet: %w", err)
	}

	// Get updated wallet info
	walletInfo, err := u.paymentClient.GetWallet(ctx, *campaign.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated wallet info: %w", err)
	}

	// Update campaign with new wallet balance
	campaign.WalletBalance = walletInfo.Data.Balance
	_, err = u.campaignRepo.Update(ctx, campaign.ID, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign wallet balance: %w", err)
	}

	response := &dto.CampaignBudgetTopupResponseDTO{
		CampaignID:        topupDTO.CampaignID,
		TransactionID:     transaction.ID,
		Amount:            topupDTO.Amount,
		PreviousBalance:   walletInfo.Data.Balance - topupDTO.Amount,
		NewBalance:        walletInfo.Data.Balance,
		Currency:          campaign.Currency,
		Status:            transaction.Status,
		TransactionCode:   transaction.TransactionCode,
		PaymentMethod:     topupDTO.PaymentMethod,
		ExternalPaymentID: topupDTO.ExternalPaymentID,
		CreatedAt:         transaction.CreatedAt,
		Message:           "Budget topup successful",
	}

	// Create history record
	u.createHistoryRecord(ctx,
		shared_enum.HistoryCampaignTopup,
		topupDTO.CampaignID,
		shared_enum.TargetHistoryCampaign,
		"Nạp tiền chiến dịch quảng cáo",
		[]string{campaign.Name, fmt.Sprintf("Số tiền: %.0f %s", topupDTO.Amount, campaign.Currency), transaction.TransactionCode},
	)

	return response, nil
}

// GetCampaignWalletInfo gets campaign wallet information
func (u *campaignUsecase) GetCampaignWalletInfo(ctx context.Context, campaignID uint64) (*dto.CampaignWalletInfoDTO, error) {
	// Get campaign
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate campaign ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Check if campaign has wallet
	if campaign.WalletID == nil {
		return nil, fmt.Errorf("campaign does not have a wallet")
	}

	// Get wallet info from payment service
	walletInfo, err := u.paymentClient.GetWallet(ctx, *campaign.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet info: %w", err)
	}

	// Get wallet transactions to calculate totals
	transactionsReq := &client.GetWalletTransactionsRequest{
		Page: 1,
		Size: 100, // Get recent transactions
	}
	// todo: use internal transaction
	transactions, err := u.paymentClient.GetWalletTransactions(ctx, transactionsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet transactions: %w", err)
	}

	// Calculate totals
	var totalDeposited, totalSpent float64
	for _, t := range transactions.Transactions {
		if t.Type == "DEPOSIT" {
			totalDeposited += t.Amount
		} else if t.Type == "PAYMENT" {
			totalSpent += t.Amount
		}
	}

	// Calculate remaining budget
	remainingBudget := campaign.Budget - totalSpent
	if remainingBudget < 0 {
		remainingBudget = 0
	}

	return &dto.CampaignWalletInfoDTO{
		CampaignID:      campaignID,
		WalletID:        *campaign.WalletID,
		Balance:         walletInfo.Data.Balance,
		Currency:        campaign.Currency,
		TotalDeposited:  totalDeposited,
		TotalSpent:      totalSpent,
		AvailableBudget: walletInfo.Data.Balance,
		CampaignBudget:  campaign.Budget,
		RemainingBudget: remainingBudget,
		LastUpdated:     walletInfo.Data.UpdatedAt,
	}, nil
}

// GetCampaignTransactions gets campaign wallet transactions
func (u *campaignUsecase) GetCampaignTransactions(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.CampaignBudgetTopupResponseDTO, int64, error) {
	// Get campaign
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, 0, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate campaign ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, 0, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Check if campaign has wallet
	if campaign.WalletID == nil {
		return nil, 0, fmt.Errorf("campaign does not have a wallet")
	}

	// Parse wallet ID
	var walletIDNum uint32
	_, err = fmt.Sscanf(*campaign.WalletID, "wallet_%d", &walletIDNum)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid wallet ID format: %s", *campaign.WalletID)
	}

	// Get transactions from payment service
	transactionsReq := &client.GetWalletTransactionsRequest{
		Page: page,
		Size: limit,
	}

	transactions, err := u.paymentClient.GetWalletTransactions(ctx, transactionsReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wallet transactions: %w", err)
	}

	// Convert to response DTOs
	responseDTOs := make([]*dto.CampaignBudgetTopupResponseDTO, len(transactions.Transactions))
	for i, t := range transactions.Transactions {
		responseDTOs[i] = &dto.CampaignBudgetTopupResponseDTO{
			CampaignID:        campaignID,
			TransactionID:     t.ID,
			Amount:            t.Amount,
			Currency:          campaign.Currency,
			Status:            t.Status,
			TransactionCode:   t.TransactionCode,
			ExternalPaymentID: t.ExternalPaymentID,
			CreatedAt:         t.CreatedAt,
			Message:           fmt.Sprintf("%s transaction", t.Type),
		}
	}

	return responseDTOs, int64(transactions.Total), nil
}

// toResponseDTO converts domain entity to response DTO
func (u *campaignUsecase) toResponseDTO(campaign *model.Campaign) *dto.CampaignResponseDTO {
	return &dto.CampaignResponseDTO{
		ID:             campaign.ID,
		Name:           campaign.Name,
		Description:    campaign.Description,
		Type:           campaign.Type,
		Status:         campaign.Status,
		ProductID:      campaign.ProductID,
		ProductName:    campaign.ProductName,
		ProductType:    campaign.ProductType,
		PackageID:      campaign.PackageID,
		PackageName:    campaign.PackageName,
		Budget:         campaign.Budget,
		Duration:       campaign.Duration,
		StartDate:      campaign.StartDate,
		EndDate:        campaign.EndDate,
		TargetAudience: campaign.TargetAudience,
		TargetLocation: campaign.TargetLocation,
		IsActive:       campaign.IsActive,
		Priority:       campaign.Priority,
		InternalNote:   campaign.InternalNote,
		CreatedBy:      campaign.CreatedBy,
		OrganizationID: campaign.OrganizationID,
		WalletID:       campaign.WalletID,
		WalletBalance:  campaign.WalletBalance,
		Currency:       campaign.Currency,
		CreatedAt:      campaign.CreatedAt,
		UpdatedAt:      campaign.UpdatedAt,
	}
}
