package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"
	"marketing/infrastructure/client"
	"marketing/internal/domain/model"
	"marketing/internal/domain/repository"
	"marketing/internal/dto"
)

type PaymentUsecase interface {
	// D.10.2.3 specific methods
	GetPaymentSummary(ctx context.Context, packageID uint64, productID uint64, customBudget *float64, customDuration *int) (*dto.PaymentSummaryDTO, error)
	ValidatePayment(ctx context.Context, paymentDTO *dto.PaymentRequestDTO) (*dto.PaymentValidationDTO, error)
	ProcessPayment(ctx context.Context, paymentDTO *dto.PaymentRequestDTO) (*dto.PaymentResponseDTO, error)
	GetWalletBalance(ctx context.Context, walletID string, requiredAmount float64) (*dto.WalletBalanceDTO, error)
	GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodDTO, error)
	
	// Payment history
	GetPaymentHistory(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.PaymentHistoryDTO, int64, error)
	GetCampaignPayment(ctx context.Context, campaignID uint64) (*dto.CampaignPaymentDTO, error)
}

type paymentUsecase struct {
	campaignRepo repository.CampaignRepository
	packageRepo  repository.PackageRepository
	paymentClient *client.PaymentClient
}

func NewPaymentUsecase(campaignRepo repository.CampaignRepository, packageRepo repository.PackageRepository, paymentClient *client.PaymentClient) PaymentUsecase {
	return &paymentUsecase{
		campaignRepo:  campaignRepo,
		packageRepo:   packageRepo,
		paymentClient: paymentClient,
	}
}

// GetPaymentSummary gets payment summary for D.10.2.3 popup
func (u *paymentUsecase) GetPaymentSummary(ctx context.Context, packageID uint64, productID uint64, customBudget *float64, customDuration *int) (*dto.PaymentSummaryDTO, error) {
	// Get package details
	pkg, err := u.packageRepo.GetByID(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("package not found: %w", err)
	}

	// Calculate total cost
	totalCost := pkg.Price
	var estimatedReach, estimatedClicks *int

	// For VIP packages with custom duration
	if pkg.Type == model.PackageTypeVIP && customDuration != nil {
		durationRatio := float64(*customDuration) / float64(pkg.Duration)
		totalCost = pkg.Price * durationRatio
		if pkg.EstimatedReach != nil {
			reach := int(float64(*pkg.EstimatedReach) * durationRatio)
			estimatedReach = &reach
		}
		if pkg.EstimatedClicks != nil {
			clicks := int(float64(*pkg.EstimatedClicks) * durationRatio)
			estimatedClicks = &clicks
		}
	}

	// For Ads packages with custom budget
	if (pkg.Type == model.PackageTypeMetaAds || pkg.Type == model.PackageTypeGoogleAds) && customBudget != nil {
		totalCost = *customBudget
		if pkg.EstimatedReach != nil {
			budgetRatio := *customBudget / pkg.Price
			reach := int(float64(*pkg.EstimatedReach) * budgetRatio)
			estimatedReach = &reach
		}
		if pkg.EstimatedClicks != nil {
			budgetRatio := *customBudget / pkg.Price
			clicks := int(float64(*pkg.EstimatedClicks) * budgetRatio)
			estimatedClicks = &clicks
		}
	}

	// TODO: Get product details from product service
	productName := "Product Name" // Placeholder
	productImage := "product_image.jpg" // Placeholder

	return &dto.PaymentSummaryDTO{
		PackageID:       packageID,
		PackageName:     pkg.Name,
		PackageType:     string(pkg.Type),
		ProductID:       productID,
		ProductName:     productName,
		ProductImage:    productImage,
		Duration:        pkg.Duration,
		TotalCost:       totalCost,
		Currency:        pkg.Currency,
		CustomBudget:    customBudget,
		CustomDuration:  customDuration,
		EstimatedReach:  estimatedReach,
		EstimatedClicks: estimatedClicks,
	}, nil
}

// ValidatePayment validates payment request
func (u *paymentUsecase) ValidatePayment(ctx context.Context, paymentDTO *dto.PaymentRequestDTO) (*dto.PaymentValidationDTO, error) {
	validation := &dto.PaymentValidationDTO{
		IsValid: true,
		Errors:  []string{},
		Warnings: []string{},
	}

	// Validate terms acceptance
	if !paymentDTO.TermsAccepted {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Terms must be accepted")
	}

	// Validate amount
	if paymentDTO.Amount <= 0 {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "Amount must be greater than 0")
	}

	// Get wallet balance
	if paymentDTO.WalletID != "" {
		walletBalance, err := u.GetWalletBalance(ctx, paymentDTO.WalletID, paymentDTO.Amount)
		if err != nil {
			validation.Warnings = append(validation.Warnings, "Unable to check wallet balance")
		} else {
			validation.WalletBalance = walletBalance
			if !walletBalance.IsSufficient {
				validation.RequiresTopup = true
				validation.TopupAmount = walletBalance.Shortfall
				validation.Warnings = append(validation.Warnings, "Insufficient wallet balance")
			}
		}
	}

	// Get payment methods
	paymentMethods, err := u.GetPaymentMethods(ctx)
	if err != nil {
		validation.Warnings = append(validation.Warnings, "Unable to get payment methods")
	} else {
		validation.PaymentMethods = paymentMethods
	}

	// Determine if can proceed
	validation.CanProceed = validation.IsValid && !validation.RequiresTopup

	return validation, nil
}

// ProcessPayment processes the payment and creates campaign
func (u *paymentUsecase) ProcessPayment(ctx context.Context, paymentDTO *dto.PaymentRequestDTO) (*dto.PaymentResponseDTO, error) {
	// Validate payment first
	validation, err := u.ValidatePayment(ctx, paymentDTO)
	if err != nil {
		return nil, fmt.Errorf("payment validation failed: %w", err)
	}

	if !validation.CanProceed {
		return nil, fmt.Errorf("payment validation failed: %s", validation.Errors)
	}

	// Parse wallet ID to get numeric ID
	var walletIDNum uint32
	_, err = fmt.Sscanf(paymentDTO.WalletID, "wallet_%d", &walletIDNum)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID format: %s", paymentDTO.WalletID)
	}

	// Process payment via payment service
	paymentReq := &client.MakePaymentRequest{
		WalletID:       uint64(walletIDNum), // Fix linter error
		Amount:         paymentDTO.Amount,
		RelatedService: "MARKETING",
		RelatedID:      fmt.Sprintf("campaign_%d", paymentDTO.CampaignID),
	}

	transaction, err := u.paymentClient.MakePayment(ctx, paymentReq)
	if err != nil {
		return nil, fmt.Errorf("payment processing failed: %w", err)
	}

	// Get updated wallet balance
	walletInfo, err := u.paymentClient.GetWallet(ctx, paymentDTO.WalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated wallet info: %w", err)
	}

	// Update campaign with payment info and transaction ID
	campaign, err := u.campaignRepo.GetByID(ctx, paymentDTO.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Update campaign status to ACTIVE and assign transaction ID
	campaign.Status = model.CampaignStatusActive
	campaign.WalletBalance = walletInfo.Data.Balance
	campaign.TransactionID = &transaction.ID // Gán transaction ID vào campaign
	_, err = u.campaignRepo.Update(ctx, campaign.ID, campaign)
	if err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	// Generate invoice if required
	if paymentDTO.RequireInvoice {
		// TODO: Generate invoice logic here
		// invoice := fmt.Sprintf("INV-%d-%d", campaign.ID, transaction.ID)
		// invoiceNumber = &invoice
	}

	return &dto.PaymentResponseDTO{
		Success:         true,
		TransactionID:   transaction.ID,
		CampaignID:      paymentDTO.CampaignID,
		Amount:          paymentDTO.Amount,
		Currency:        campaign.Currency,
		PaymentMethod:   paymentDTO.PaymentMethod,
		Status:          transaction.Status,
		TransactionCode: transaction.TransactionCode,
		WalletBalance:   walletInfo.Data.Balance,
		PreviousBalance: walletInfo.Data.Balance + paymentDTO.Amount,
		CreatedAt:       transaction.CreatedAt,
		Message:         "Payment processed successfully",
		NextStep:        "campaign_details",
	}, nil
}

// GetWalletBalance gets wallet balance and checks sufficiency
func (u *paymentUsecase) GetWalletBalance(ctx context.Context, walletID string, requiredAmount float64) (*dto.WalletBalanceDTO, error) {
	walletInfo, err := u.paymentClient.GetWallet(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet info: %w", err)
	}

	isSufficient := walletInfo.Data.Balance >= requiredAmount
	shortfall := requiredAmount - walletInfo.Data.Balance
	if shortfall < 0 {
		shortfall = 0
	}

	return &dto.WalletBalanceDTO{
		WalletID:       walletID,
		Balance:        walletInfo.Data.Balance,
		Currency:       walletInfo.Data.Currency,
		IsSufficient:   isSufficient,
		RequiredAmount: requiredAmount,
		Shortfall:      shortfall,
		LastUpdated:    walletInfo.Data.UpdatedAt,
	}, nil
}

// GetPaymentMethods gets available payment methods
func (u *paymentUsecase) GetPaymentMethods(ctx context.Context) ([]dto.PaymentMethodDTO, error) {
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	
	req := &client.GetPaymentMethodsRequest{
		OrganizationID: organizationID,
		Page:           1,
		Size:           10,
	}

	response, err := u.paymentClient.GetPaymentMethods(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	methods := make([]dto.PaymentMethodDTO, len(response.PaymentMethods))
	for i, method := range response.PaymentMethods {
		methods[i] = dto.PaymentMethodDTO{
			ID:          method.ID,
			Name:        method.Name,
			Code:        method.Code,
			Description: "", // Not available in current response
			IsActive:    method.IsActive,
			IsDefault:   false, // TODO: Add to payment service
			Icon:        "",    // TODO: Add to payment service
		}
	}

	return methods, nil
}

// GetPaymentHistory gets payment history for campaign
func (u *paymentUsecase) GetPaymentHistory(ctx context.Context, campaignID uint64, page, limit int) ([]*dto.PaymentHistoryDTO, int64, error) {
	// Get campaign
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, 0, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, 0, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// Get wallet transactions
	if campaign.WalletID == nil {
		return []*dto.PaymentHistoryDTO{}, 0, nil
	}

	transactionsReq := &client.GetWalletTransactionsRequest{
		Page: page,
		Size: limit,
	}

	transactions, err := u.paymentClient.GetWalletTransactions(ctx, transactionsReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wallet transactions: %w", err)
	}

	// Convert to payment history DTOs
	history := make([]*dto.PaymentHistoryDTO, len(transactions.Transactions))
	for i, t := range transactions.Transactions {
		history[i] = &dto.PaymentHistoryDTO{
			ID:              t.ID,
			CampaignID:      campaignID,
			PackageID:       0, // TODO: Extract from RelatedID
			Amount:          t.Amount,
			Currency:        campaign.Currency,
			PaymentMethod:   "", // TODO: Add to transaction
			Status:          t.Status,
			TransactionCode: t.TransactionCode,
			CreatedAt:       t.CreatedAt,
			UpdatedAt:       t.UpdatedAt,
			CampaignName:    campaign.Name,
			PackageName:     campaign.PackageName,
		}
	}

	return history, int64(transactions.Total), nil
}

// GetCampaignPayment gets campaign payment information
func (u *paymentUsecase) GetCampaignPayment(ctx context.Context, campaignID uint64) (*dto.CampaignPaymentDTO, error) {
	campaign, err := u.campaignRepo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("campaign not found: %w", err)
	}

	// Validate ownership
	organizationID := _utils.GetOrganizationIdFromContext(ctx)
	if campaign.OrganizationID != organizationID {
		return nil, fmt.Errorf("unauthorized: campaign does not belong to your organization")
	}

	// TODO: Get payment details from payment service
	// For now, return basic info
	return &dto.CampaignPaymentDTO{
		CampaignID:    campaignID,
		PackageID:     0, // TODO: Get from campaign
		Amount:        campaign.Budget,
		Currency:      campaign.Currency,
		PaymentStatus: "PAID", // TODO: Get actual status
		PaymentMethod: "WALLET",
		RequireInvoice: false,
	}, nil
} 