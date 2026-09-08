package handler

import (
	"context"
	"fmt"
	"time"

	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"

	pb_marketing "pb/types/marketing"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MarketingGrpcHandler struct {
	pb_marketing.UnimplementedMarketingServiceServer
	campaignUsecase            usecase.CampaignUsecase
	packageUsecase             usecase.PackageUsecase
	paymentUsecase             usecase.PaymentUsecase
	campaignDetailUsecase      usecase.CampaignDetailUsecase
	budgetUsecase              usecase.BudgetUsecase
	advertisingSettingsUsecase usecase.AdvertisingSettingsUsecase
	notificationClient         *client.NotificationClient
}

func NewMarketingGrpcHandler(campaignUsecase usecase.CampaignUsecase, packageUsecase usecase.PackageUsecase, paymentUsecase usecase.PaymentUsecase, campaignDetailUsecase usecase.CampaignDetailUsecase, budgetUsecase usecase.BudgetUsecase, advertisingSettingsUsecase usecase.AdvertisingSettingsUsecase, notificationClient *client.NotificationClient) *MarketingGrpcHandler {
	return &MarketingGrpcHandler{
		campaignUsecase:            campaignUsecase,
		packageUsecase:             packageUsecase,
		paymentUsecase:             paymentUsecase,
		campaignDetailUsecase:      campaignDetailUsecase,
		budgetUsecase:              budgetUsecase,
		advertisingSettingsUsecase: advertisingSettingsUsecase,
		notificationClient:         notificationClient,
	}
}

func (h *MarketingGrpcHandler) CreateCampaign(ctx context.Context, req *pb_marketing.CreateCampaignRequest) (*pb_marketing.CampaignResponse, error) {
	createDTO := &dto.CampaignCreateDTO{
		Name:           req.Name,
		Description:    req.Description,
		Type:           domain.CampaignType(req.Type),
		ProductID:      &req.ProductId,
		ProductName:    req.ProductName,
		ProductType:    req.ProductType,
		PackageID:      &req.PackageId,
		PackageName:    req.PackageName,
		Budget:         req.Budget,
		Duration:       int(req.Duration),
		TargetAudience: req.TargetAudience,
		TargetLocation: req.TargetLocation,
		Priority:       int(req.Priority),
		InternalNote:   req.InternalNote,
		TermsAccepted:  req.TermsAccepted,
	}

	if req.StartDate != "" {
		if startDate, err := time.Parse(time.RFC3339, req.StartDate); err == nil {
			createDTO.StartDate = &startDate
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse(time.RFC3339, req.EndDate); err == nil {
			createDTO.EndDate = &endDate
		}
	}

	campaign, err := h.campaignUsecase.CreateCampaign(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create campaign: %v", err)
	}

	return h.campaignToGrpcResponse(campaign), nil
}

func (h *MarketingGrpcHandler) UpdateCampaign(ctx context.Context, req *pb_marketing.UpdateCampaignRequest) (*pb_marketing.CampaignResponse, error) {
	updateDTO := &dto.CampaignUpdateDTO{
		Name:           req.Name,
		Description:    req.Description,
		Type:           domain.CampaignType(req.Type),
		ProductID:      &req.ProductId,
		ProductName:    req.ProductName,
		ProductType:    req.ProductType,
		PackageID:      &req.PackageId,
		PackageName:    req.PackageName,
		Budget:         req.Budget,
		Duration:       int(req.Duration),
		TargetAudience: req.TargetAudience,
		TargetLocation: req.TargetLocation,
		Status:         domain.CampaignStatus(req.Status),
		Priority:       int(req.Priority),
		InternalNote:   req.InternalNote,
	}

	if req.StartDate != "" {
		if startDate, err := time.Parse(time.RFC3339, req.StartDate); err == nil {
			updateDTO.StartDate = &startDate
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse(time.RFC3339, req.EndDate); err == nil {
			updateDTO.EndDate = &endDate
		}
	}

	campaign, err := h.campaignUsecase.UpdateCampaign(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update campaign: %v", err)
	}

	return h.campaignToGrpcResponse(campaign), nil
}

func (h *MarketingGrpcHandler) DeleteCampaign(ctx context.Context, req *pb_marketing.DeleteCampaignRequest) (*pb_marketing.DeleteCampaignResponse, error) {
	err := h.campaignUsecase.DeleteCampaign(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete campaign: %v", err)
	}

	return &pb_marketing.DeleteCampaignResponse{
		Success: true,
		Message: "Campaign deleted successfully",
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaign(ctx context.Context, req *pb_marketing.GetCampaignRequest) (*pb_marketing.CampaignResponse, error) {
	campaign, err := h.campaignUsecase.GetCampaign(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "campaign not found: %v", err)
	}

	return h.campaignToGrpcResponse(campaign), nil
}

func (h *MarketingGrpcHandler) ListCampaigns(ctx context.Context, req *pb_marketing.ListCampaignsRequest) (*pb_marketing.ListCampaignsResponse, error) {
	searchDTO := dto.CampaignSearchDTO{
		Status: domain.CampaignStatus(req.Status),
		Type:   domain.CampaignType(req.Type),
		Page:   int(req.Page),
		Limit:  int(req.Size),
	}

	campaigns, total, err := h.campaignUsecase.SearchCampaigns(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list campaigns: %v", err)
	}

	grpcCampaigns := make([]*pb_marketing.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		grpcCampaigns[i] = h.campaignToGrpcCampaign(&campaign)
	}

	return &pb_marketing.ListCampaignsResponse{
		Data:  grpcCampaigns,
		Total: uint32(total),
	}, nil
}

func (h *MarketingGrpcHandler) SearchCampaigns(ctx context.Context, req *pb_marketing.SearchCampaignsRequest) (*pb_marketing.SearchCampaignsResponse, error) {
	searchDTO := dto.CampaignSearchDTO{
		Name:      req.Name,
		Type:      domain.CampaignType(req.Type),
		Status:    domain.CampaignStatus(req.Status),
		ProductID: &req.ProductId,
		PackageID: &req.PackageId,
		Page:      int(req.Page),
		Limit:     int(req.Size),
	}

	if req.FromDate != "" {
		if fromDate, err := time.Parse(time.RFC3339, req.FromDate); err == nil {
			searchDTO.FromDate = &fromDate
		}
	}
	if req.ToDate != "" {
		if toDate, err := time.Parse(time.RFC3339, req.ToDate); err == nil {
			searchDTO.ToDate = &toDate
		}
	}

	campaigns, total, err := h.campaignUsecase.SearchCampaigns(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search campaigns: %v", err)
	}

	grpcCampaigns := make([]*pb_marketing.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		grpcCampaigns[i] = h.campaignToGrpcCampaign(&campaign)
	}

	return &pb_marketing.SearchCampaignsResponse{
		Data:  grpcCampaigns,
		Total: uint32(total),
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaignsByProduct(ctx context.Context, req *pb_marketing.GetCampaignsByProductRequest) (*pb_marketing.SearchCampaignsResponse, error) {
	campaigns, err := h.campaignUsecase.GetCampaignsByProduct(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaigns by product: %v", err)
	}

	grpcCampaigns := make([]*pb_marketing.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		grpcCampaigns[i] = h.campaignToGrpcCampaign(&campaign)
	}

	return &pb_marketing.SearchCampaignsResponse{
		Data:  grpcCampaigns,
		Total: uint32(len(campaigns)),
	}, nil
}

func (h *MarketingGrpcHandler) UpdateCampaignStatus(ctx context.Context, req *pb_marketing.UpdateCampaignStatusRequest) (*sharepb.SubmitResponse, error) {
	err := h.campaignUsecase.UpdateCampaignStatus(ctx, req.Id, domain.CampaignStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update campaign status: %v", err)
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Campaign status updated successfully",
	}, nil
}

func (h *MarketingGrpcHandler) TopupCampaignBudget(ctx context.Context, req *pb_marketing.TopupCampaignBudgetRequest) (*pb_marketing.CampaignBudgetTopupResponse, error) {
	topupDTO := &dto.CampaignBudgetTopupDTO{
		CampaignID:        req.Id,
		Amount:            req.Amount,
		PaymentMethod:     req.PaymentMethod,
		ExternalPaymentID: req.ExternalPaymentId,
		Description:       req.Description,
	}

	result, err := h.campaignUsecase.TopupCampaignBudget(ctx, topupDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to topup campaign budget: %v", err)
	}

	return &pb_marketing.CampaignBudgetTopupResponse{
		CampaignId:        result.CampaignID,
		TransactionId:     result.TransactionID,
		Amount:            result.Amount,
		PreviousBalance:   result.PreviousBalance,
		NewBalance:        result.NewBalance,
		Currency:          result.Currency,
		Status:            result.Status,
		TransactionCode:   result.TransactionCode,
		PaymentMethod:     result.PaymentMethod,
		ExternalPaymentId: result.ExternalPaymentID,
		CreatedAt:         result.CreatedAt,
		Message:           result.Message,
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaignWalletInfo(ctx context.Context, req *pb_marketing.GetCampaignWalletInfoRequest) (*pb_marketing.CampaignWalletInfoResponse, error) {
	result, err := h.campaignUsecase.GetCampaignWalletInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign wallet info: %v", err)
	}

	return &pb_marketing.CampaignWalletInfoResponse{
		CampaignId:      result.CampaignID,
		WalletId:        result.WalletID,
		Balance:         result.Balance,
		Currency:        result.Currency,
		TotalDeposited:  result.TotalDeposited,
		TotalSpent:      result.TotalSpent,
		AvailableBudget: result.AvailableBudget,
		CampaignBudget:  result.CampaignBudget,
		RemainingBudget: result.RemainingBudget,
		LastUpdated:     result.LastUpdated,
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaignTransactions(ctx context.Context, req *pb_marketing.GetCampaignTransactionsRequest) (*pb_marketing.CampaignTransactionsResponse, error) {
	transactions, total, err := h.campaignUsecase.GetCampaignTransactions(ctx, req.Id, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign transactions: %v", err)
	}

	grpcTransactions := make([]*pb_marketing.CampaignBudgetTopupResponse, len(transactions))
	for i, t := range transactions {
		grpcTransactions[i] = &pb_marketing.CampaignBudgetTopupResponse{
			CampaignId:        t.CampaignID,
			TransactionId:     t.TransactionID,
			Amount:            t.Amount,
			Currency:          t.Currency,
			Status:            t.Status,
			TransactionCode:   t.TransactionCode,
			ExternalPaymentId: t.ExternalPaymentID,
			CreatedAt:         t.CreatedAt,
			Message:           t.Message,
		}
	}

	return &pb_marketing.CampaignTransactionsResponse{
		Data:  grpcTransactions,
		Total: uint32(total),
	}, nil
}

func (h *MarketingGrpcHandler) CreatePackage(ctx context.Context, req *pb_marketing.CreatePackageRequest) (*pb_marketing.PackageResponse, error) {
	createDTO := &dto.PackageCreateDTO{
		Name:                 req.Name,
		Description:          req.Description,
		Type:                 domain.PackageType(req.Type),
		Price:                req.Price,
		Currency:             req.Currency,
		Duration:             int(req.Duration),
		Features:             req.Features,
		MaxImpressions:       getIntPtr(req.MaxImpressions),
		MaxClicks:            getIntPtr(req.MaxClicks),
		MinBudget:            getFloat64Ptr(req.MinBudget),
		MaxBudget:            getFloat64Ptr(req.MaxBudget),
		EstimatedReach:       getIntPtr(req.EstimatedReach),
		EstimatedClicks:      getIntPtr(req.EstimatedClicks),
		IsPopular:            req.IsPopular,
		IsRecommended:        req.IsRecommended,
		Icon:                 req.Icon,
		Color:                req.Color,
		RequiresVerification: req.RequiresVerification,
		RequiresPayment:      req.RequiresPayment,
	}

	package_, err := h.packageUsecase.CreatePackage(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create package: %v", err)
	}

	return h.packageToGrpcResponse(package_), nil
}

func (h *MarketingGrpcHandler) GetPackage(ctx context.Context, req *pb_marketing.GetPackageRequest) (*pb_marketing.PackageResponse, error) {
	package_, err := h.packageUsecase.GetPackage(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "package not found: %v", err)
	}

	return h.packageToGrpcResponse(package_), nil
}

func (h *MarketingGrpcHandler) UpdatePackage(ctx context.Context, req *pb_marketing.UpdatePackageRequest) (*pb_marketing.PackageResponse, error) {
	updateDTO := &dto.PackageUpdateDTO{
		Name:                 req.Name,
		Description:          req.Description,
		Type:                 domain.PackageType(req.Type),
		Status:               domain.PackageStatus(req.Status),
		Price:                req.Price,
		Currency:             req.Currency,
		Duration:             int(req.Duration),
		Features:             req.Features,
		MaxImpressions:       getIntPtr(req.MaxImpressions),
		MaxClicks:            getIntPtr(req.MaxClicks),
		MinBudget:            getFloat64Ptr(req.MinBudget),
		MaxBudget:            getFloat64Ptr(req.MaxBudget),
		EstimatedReach:       getIntPtr(req.EstimatedReach),
		EstimatedClicks:      getIntPtr(req.EstimatedClicks),
		IsPopular:            req.IsPopular,
		IsRecommended:        req.IsRecommended,
		Icon:                 req.Icon,
		Color:                req.Color,
		RequiresVerification: req.RequiresVerification,
		RequiresPayment:      req.RequiresPayment,
	}

	package_, err := h.packageUsecase.UpdatePackage(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update package: %v", err)
	}

	return h.packageToGrpcResponse(package_), nil
}

func (h *MarketingGrpcHandler) DeletePackage(ctx context.Context, req *pb_marketing.DeletePackageRequest) (*pb_marketing.DeletePackageResponse, error) {
	err := h.packageUsecase.DeletePackage(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete package: %v", err)
	}

	return &pb_marketing.DeletePackageResponse{
		Success: true,
		Message: "Package deleted successfully",
	}, nil
}

func (h *MarketingGrpcHandler) ListPackages(ctx context.Context, req *pb_marketing.ListPackagesRequest) (*pb_marketing.ListPackagesResponse, error) {
	searchDTO := dto.PackageSearchDTO{
		Type:   domain.PackageType(req.Type),
		Status: domain.PackageStatus(req.Status),
		Page:   int(req.Page),
		Limit:  int(req.Size),
	}

	packages, total, err := h.packageUsecase.SearchPackages(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list packages: %v", err)
	}

	grpcPackages := make([]*pb_marketing.Package, len(packages))
	for i, package_ := range packages {
		grpcPackages[i] = h.packageToGrpcPackage(package_)
	}

	return &pb_marketing.ListPackagesResponse{
		Data:  grpcPackages,
		Total: uint32(total),
	}, nil
}

func (h *MarketingGrpcHandler) CreatePayment(ctx context.Context, req *pb_marketing.CreatePaymentRequest) (*pb_marketing.PaymentResponse, error) {
	// CreatePayment is similar to ProcessPayment but for initial payment creation
	paymentDTO := &dto.PaymentRequestDTO{
		CampaignID:     req.CampaignId,
		PackageID:      req.PackageId,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		WalletID:       req.WalletId,
		RequireInvoice: req.RequireInvoice,
		InvoiceEmail:   req.InvoiceEmail,
		TermsAccepted:  req.TermsAccepted,
	}

	result, err := h.paymentUsecase.ProcessPayment(ctx, paymentDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return &pb_marketing.PaymentResponse{
		Success:         result.Success,
		TransactionId:   result.TransactionID,
		CampaignId:      result.CampaignID,
		Amount:          result.Amount,
		Currency:        result.Currency,
		PaymentMethod:   result.PaymentMethod,
		Status:          result.Status,
		TransactionCode: result.TransactionCode,
		WalletBalance:   result.WalletBalance,
		PreviousBalance: result.PreviousBalance,
		CreatedAt:       result.CreatedAt,
		Message:         result.Message,
		NextStep:        result.NextStep,
	}, nil
}

func (h *MarketingGrpcHandler) GetPayment(ctx context.Context, req *pb_marketing.GetPaymentRequest) (*pb_marketing.PaymentResponse, error) {
	// GetPayment retrieves payment information by payment ID using usecase
	result, err := h.paymentUsecase.GetPaymentByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &pb_marketing.PaymentResponse{
		Success:         result.Success,
		TransactionId:   result.TransactionID,
		CampaignId:      result.CampaignID,
		Amount:          result.Amount,
		Currency:        result.Currency,
		PaymentMethod:   result.PaymentMethod,
		Status:          result.Status,
		TransactionCode: result.TransactionCode,
		WalletBalance:   result.WalletBalance,
		PreviousBalance: result.PreviousBalance,
		CreatedAt:       result.CreatedAt,
		Message:         result.Message,
		NextStep:        result.NextStep,
	}, nil
}

func (h *MarketingGrpcHandler) ListPayments(ctx context.Context, req *pb_marketing.ListPaymentsRequest) (*pb_marketing.ListPaymentsResponse, error) {
	page := 0
	size := 10

	if req.Page > 0 {
		page = int(req.Page)
	}

	if req.Size > 0 {
		size = int(req.Size)
	}
	history, total, err := h.paymentUsecase.GetPaymentHistory(ctx, req.CampaignId, page, size)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list payments: %v", err)
	}

	response := &pb_marketing.ListPaymentsResponse{
		Data:  make([]*pb_marketing.PaymentHistoryItemResponse, len(history)),
		Total: uint32(total),
	}

	for i, item := range history {
		response.Data[i] = &pb_marketing.PaymentHistoryItemResponse{
			Id:              item.ID,
			CampaignId:      item.CampaignID,
			PackageId:       item.PackageID,
			Amount:          item.Amount,
			Currency:        item.Currency,
			PaymentMethod:   item.PaymentMethod,
			Status:          item.Status,
			TransactionCode: item.TransactionCode,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
			CampaignName:    item.CampaignName,
			PackageName:     item.PackageName,
		}
	}

	return response, nil
}

func (h *MarketingGrpcHandler) GetPaymentSummary(ctx context.Context, req *pb_marketing.GetPaymentSummaryRequest) (*pb_marketing.PaymentSummaryResponse, error) {
	var customBudget *float64
	if req.CustomBudget > 0 {
		customBudget = &req.CustomBudget
	}

	var customDuration *int
	if req.CustomDuration > 0 {
		duration := int(req.CustomDuration)
		customDuration = &duration
	}

	result, err := h.paymentUsecase.GetPaymentSummary(ctx, req.PackageId, req.ProductId, customBudget, customDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment summary: %v", err)
	}

	return &pb_marketing.PaymentSummaryResponse{
		PackageId:       result.PackageID,
		PackageName:     result.PackageName,
		PackageType:     result.PackageType,
		ProductId:       result.ProductID,
		ProductName:     result.ProductName,
		ProductImage:    result.ProductImage,
		Duration:        uint32(result.Duration),
		TotalCost:       result.TotalCost,
		Currency:        result.Currency,
		CustomBudget:    getFloat64Value(result.CustomBudget),
		CustomDuration:  getUint32Value(result.CustomDuration),
		EstimatedReach:  getUint32Value(result.EstimatedReach),
		EstimatedClicks: getUint32Value(result.EstimatedClicks),
	}, nil
}

func (h *MarketingGrpcHandler) ValidatePayment(ctx context.Context, req *pb_marketing.ValidatePaymentRequest) (*pb_marketing.PaymentValidationResponse, error) {
	paymentDTO := &dto.PaymentRequestDTO{
		CampaignID:     req.CampaignId,
		PackageID:      req.PackageId,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		WalletID:       req.WalletId,
		ProductID:      req.ProductId,
		ProductName:    req.ProductName,
		ProductType:    req.ProductType,
		CustomBudget:   getFloat64Ptr(req.CustomBudget),
		CustomDuration: getIntPtr(req.CustomDuration),
		RequireInvoice: req.RequireInvoice,
		InvoiceEmail:   req.InvoiceEmail,
		TermsAccepted:  req.TermsAccepted,
	}

	result, err := h.paymentUsecase.ValidatePayment(ctx, paymentDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate payment: %v", err)
	}

	response := &pb_marketing.PaymentValidationResponse{
		IsValid:       result.IsValid,
		Errors:        result.Errors,
		Warnings:      result.Warnings,
		CanProceed:    result.CanProceed,
		RequiresTopup: result.RequiresTopup,
		TopupAmount:   result.TopupAmount,
	}

	if result.WalletBalance != nil {
		response.WalletBalance = &pb_marketing.WalletBalanceResponse{
			WalletId:       result.WalletBalance.WalletID,
			Balance:        result.WalletBalance.Balance,
			Currency:       result.WalletBalance.Currency,
			IsSufficient:   result.WalletBalance.IsSufficient,
			RequiredAmount: result.WalletBalance.RequiredAmount,
			Shortfall:      result.WalletBalance.Shortfall,
			LastUpdated:    result.WalletBalance.LastUpdated,
		}
	}

	response.PaymentMethods = make([]*pb_marketing.PaymentMethodResponse, len(result.PaymentMethods))
	for i, method := range result.PaymentMethods {
		response.PaymentMethods[i] = &pb_marketing.PaymentMethodResponse{
			Id:          method.ID,
			Name:        method.Name,
			Code:        method.Code,
			Description: method.Description,
			IsActive:    method.IsActive,
			IsDefault:   method.IsDefault,
			Icon:        method.Icon,
		}
	}

	return response, nil
}

func (h *MarketingGrpcHandler) ProcessPayment(ctx context.Context, req *pb_marketing.ProcessPaymentRequest) (*pb_marketing.PaymentResponse, error) {
	paymentDTO := &dto.PaymentRequestDTO{
		CampaignID:     req.CampaignId,
		PackageID:      req.PackageId,
		Amount:         req.Amount,
		PaymentMethod:  req.PaymentMethod,
		WalletID:       req.WalletId,
		ProductID:      req.ProductId,
		ProductName:    req.ProductName,
		ProductType:    req.ProductType,
		CustomBudget:   getFloat64Ptr(req.CustomBudget),
		CustomDuration: getIntPtr(req.CustomDuration),
		RequireInvoice: req.RequireInvoice,
		InvoiceEmail:   req.InvoiceEmail,
		TermsAccepted:  req.TermsAccepted,
	}

	result, err := h.paymentUsecase.ProcessPayment(ctx, paymentDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &pb_marketing.PaymentResponse{
		Success:         result.Success,
		TransactionId:   result.TransactionID,
		CampaignId:      result.CampaignID,
		Amount:          result.Amount,
		Currency:        result.Currency,
		PaymentMethod:   result.PaymentMethod,
		Status:          result.Status,
		TransactionCode: result.TransactionCode,
		WalletBalance:   result.WalletBalance,
		PreviousBalance: result.PreviousBalance,
		CreatedAt:       result.CreatedAt,
		Message:         result.Message,
		NextStep:        result.NextStep,
	}, nil
}

func (h *MarketingGrpcHandler) GetWalletBalance(ctx context.Context, req *pb_marketing.GetWalletBalanceRequest) (*pb_marketing.WalletBalanceResponse, error) {
	result, err := h.paymentUsecase.GetWalletBalance(ctx, req.WalletId, req.RequiredAmount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get wallet balance: %v", err)
	}

	return &pb_marketing.WalletBalanceResponse{
		WalletId:       result.WalletID,
		Balance:        result.Balance,
		Currency:       result.Currency,
		IsSufficient:   result.IsSufficient,
		RequiredAmount: result.RequiredAmount,
		Shortfall:      result.Shortfall,
		LastUpdated:    result.LastUpdated,
	}, nil
}

func (h *MarketingGrpcHandler) GetPaymentMethods(ctx context.Context, req *pb_marketing.Empty) (*pb_marketing.PaymentMethodsResponse, error) {
	methods, err := h.paymentUsecase.GetPaymentMethods(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment methods: %v", err)
	}

	response := &pb_marketing.PaymentMethodsResponse{
		PaymentMethods: make([]*pb_marketing.PaymentMethodResponse, len(methods)),
	}

	for i, method := range methods {
		response.PaymentMethods[i] = &pb_marketing.PaymentMethodResponse{
			Id:          method.ID,
			Name:        method.Name,
			Code:        method.Code,
			Description: method.Description,
			IsActive:    method.IsActive,
			IsDefault:   method.IsDefault,
			Icon:        method.Icon,
		}
	}

	return response, nil
}

func (h *MarketingGrpcHandler) GetPaymentHistory(ctx context.Context, req *pb_marketing.GetPaymentHistoryRequest) (*pb_marketing.PaymentHistoryResponse, error) {
	history, total, err := h.paymentUsecase.GetPaymentHistory(ctx, req.CampaignId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment history: %v", err)
	}

	response := &pb_marketing.PaymentHistoryResponse{
		Data:  make([]*pb_marketing.PaymentHistoryItemResponse, len(history)),
		Total: uint32(total),
	}

	for i, item := range history {
		response.Data[i] = &pb_marketing.PaymentHistoryItemResponse{
			Id:              item.ID,
			CampaignId:      item.CampaignID,
			PackageId:       item.PackageID,
			Amount:          item.Amount,
			Currency:        item.Currency,
			PaymentMethod:   item.PaymentMethod,
			Status:          item.Status,
			TransactionCode: item.TransactionCode,
			CreatedAt:       item.CreatedAt,
			UpdatedAt:       item.UpdatedAt,
			CampaignName:    item.CampaignName,
			PackageName:     item.PackageName,
		}
	}

	return response, nil
}

func (h *MarketingGrpcHandler) GetCampaignPayment(ctx context.Context, req *pb_marketing.GetCampaignPaymentRequest) (*pb_marketing.CampaignPaymentResponse, error) {
	result, err := h.paymentUsecase.GetCampaignPayment(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign payment: %v", err)
	}

	response := &pb_marketing.CampaignPaymentResponse{
		CampaignId:     result.CampaignID,
		PackageId:      result.PackageID,
		Amount:         result.Amount,
		Currency:       result.Currency,
		PaymentStatus:  result.PaymentStatus,
		PaymentMethod:  result.PaymentMethod,
		RequireInvoice: result.RequireInvoice,
	}

	if result.TransactionID != nil {
		response.TransactionId = *result.TransactionID
	}
	if result.PaidAt != nil {
		response.PaidAt = *result.PaidAt
	}
	if result.InvoiceNumber != nil {
		response.InvoiceNumber = *result.InvoiceNumber
	}
	if result.InvoiceEmail != nil {
		response.InvoiceEmail = *result.InvoiceEmail
	}

	return response, nil
}

func (h *MarketingGrpcHandler) GetCampaignPerformance(ctx context.Context, req *pb_marketing.GetCampaignPerformanceRequest) (*pb_marketing.CampaignPerformanceResponse, error) {
	result, err := h.campaignDetailUsecase.GetCampaignPerformance(ctx, req.Id, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign performance: %v", err)
	}

	dailyData := make([]*pb_marketing.CampaignDailyData, len(result.DailyData))
	for i, data := range result.DailyData {
		dailyData[i] = &pb_marketing.CampaignDailyData{
			Date:        data.Date,
			Impressions: uint32(data.Impressions),
			Clicks:      uint32(data.Clicks),
			Ctr:         data.CTR,
			Cpc:         data.CPC,
			Spent:       data.Spent,
		}
	}

	return &pb_marketing.CampaignPerformanceResponse{
		Data: &pb_marketing.CampaignPerformance{
			CampaignId:  result.CampaignID,
			Date:        result.Date,
			Impressions: uint32(result.Impressions),
			Clicks:      uint32(result.Clicks),
			Ctr:         result.CTR,
			Cpc:         result.CPC,
			Spent:       result.Spent,
			DailyData:   dailyData,
			Summary: &pb_marketing.CampaignPerformanceSummary{
				TotalImpressions: uint32(result.Summary.TotalImpressions),
				TotalClicks:      uint32(result.Summary.TotalClicks),
				AverageCtr:       result.Summary.AverageCTR,
				AverageCpc:       result.Summary.AverageCPC,
				TotalSpent:       result.Summary.TotalSpent,
				RemainingBudget:  result.Summary.RemainingBudget,
				Roi:              result.Summary.ROI,
			},
		},
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaignHistory(ctx context.Context, req *pb_marketing.GetCampaignHistoryRequest) (*pb_marketing.CampaignHistoryResponse, error) {
	history, total, err := h.campaignDetailUsecase.GetCampaignHistory(ctx, req.Id, int(req.Page), int(req.Size))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign history: %v", err)
	}

	response := &pb_marketing.CampaignHistoryResponse{
		Data:  make([]*pb_marketing.CampaignHistoryItem, len(history)),
		Total: uint32(total),
	}

	for i, item := range history {
		response.Data[i] = &pb_marketing.CampaignHistoryItem{
			Id:              item.ID,
			CampaignId:      item.CampaignID,
			Action:          item.Action,
			Description:     item.Description,
			PerformedBy:     item.PerformedBy,
			PerformedByUser: item.PerformedByUser,
			PerformedAt:     item.PerformedAt,
			IpAddress:       item.IPAddress,
			UserAgent:       item.UserAgent,
		}
	}

	return response, nil
}

func (h *MarketingGrpcHandler) ExtendCampaign(ctx context.Context, req *pb_marketing.ExtendCampaignRequest) (*pb_marketing.PaymentResponse, error) {
	extendDTO := &dto.CampaignExtendDTO{
		CampaignID:       req.Id,
		ExtensionDays:    int(req.ExtensionDays),
		AdditionalBudget: req.AdditionalBudget,
		PaymentMethod:    req.PaymentMethod,
		TermsAccepted:    req.TermsAccepted,
	}

	result, err := h.campaignDetailUsecase.ExtendCampaign(ctx, extendDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to extend campaign: %v", err)
	}

	return &pb_marketing.PaymentResponse{
		Success:         result.Success,
		TransactionId:   result.TransactionID,
		CampaignId:      result.CampaignID,
		Amount:          result.Amount,
		Currency:        result.Currency,
		PaymentMethod:   result.PaymentMethod,
		Status:          result.Status,
		TransactionCode: result.TransactionCode,
		WalletBalance:   result.WalletBalance,
		PreviousBalance: result.PreviousBalance,
		CreatedAt:       result.CreatedAt,
		Message:         result.Message,
		NextStep:        result.NextStep,
	}, nil
}

func (h *MarketingGrpcHandler) ExportCampaignData(ctx context.Context, req *pb_marketing.ExportCampaignDataRequest) (*pb_marketing.CampaignExportResponse, error) {
	result, err := h.campaignDetailUsecase.ExportCampaignData(ctx, req.Id, req.ExportType, req.DateRange)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to export campaign data: %v", err)
	}

	performanceData := make([]*pb_marketing.CampaignDailyData, len(result.PerformanceData))
	for i, data := range result.PerformanceData {
		performanceData[i] = &pb_marketing.CampaignDailyData{
			Date:        data.Date,
			Impressions: uint32(data.Impressions),
			Clicks:      uint32(data.Clicks),
			Ctr:         data.CTR,
			Cpc:         data.CPC,
			Spent:       data.Spent,
		}
	}

	return &pb_marketing.CampaignExportResponse{
		Data: &pb_marketing.CampaignExport{
			CampaignId:      result.CampaignID,
			CampaignName:    result.CampaignName,
			ExportType:      result.ExportType,
			DateRange:       result.DateRange,
			PerformanceData: performanceData,
			Summary: &pb_marketing.CampaignPerformanceSummary{
				TotalImpressions: uint32(result.Summary.TotalImpressions),
				TotalClicks:      uint32(result.Summary.TotalClicks),
				AverageCtr:       result.Summary.AverageCTR,
				AverageCpc:       result.Summary.AverageCPC,
				TotalSpent:       result.Summary.TotalSpent,
				RemainingBudget:  result.Summary.RemainingBudget,
				Roi:              result.Summary.ROI,
			},
			GeneratedAt: result.GeneratedAt,
			GeneratedBy: result.GeneratedBy,
		},
	}, nil
}

func (h *MarketingGrpcHandler) GetCampaignOptimization(ctx context.Context, req *pb_marketing.GetCampaignOptimizationRequest) (*pb_marketing.CampaignOptimizationResponse, error) {
	result, err := h.campaignDetailUsecase.GetCampaignOptimization(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign optimization: %v", err)
	}

	return &pb_marketing.CampaignOptimizationResponse{
		Data: &pb_marketing.CampaignOptimization{
			CampaignId:      result.CampaignID,
			Suggestions:     result.Suggestions,
			Priority:        result.Priority,
			EstimatedImpact: result.EstimatedImpact,
			GeneratedAt:     result.GeneratedAt,
		},
	}, nil
}

func (h *MarketingGrpcHandler) CreateBudget(ctx context.Context, req *pb_marketing.CreateBudgetRequest) (*pb_marketing.BudgetResponse, error) {
	// Map protobuf request to DTO structure
	createDTO := &dto.BudgetConfigCreateDTO{
		CampaignID:     req.CampaignId,
		MonthlyLimit:   req.Amount,
		AlertThreshold: 80.0, // Default alert threshold
		AlertEnabled:   true,
		AutoTopup:      false,
	}

	// Create budget config
	budgetConfig, err := h.budgetUsecase.CreateBudgetConfig(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create budget: %v", err)
	}

	// Map DTO response to protobuf structure
	budget := &pb_marketing.Budget{
		Id:          budgetConfig.ID,
		CampaignId:  budgetConfig.CampaignID,
		Amount:      budgetConfig.MonthlyLimit,
		Currency:    "VND", // Default currency
		Description: "Monthly budget limit",
		Status:      "ACTIVE",
		CreatedAt:   budgetConfig.CreatedAt,
		UpdatedAt:   budgetConfig.UpdatedAt,
	}

	return &pb_marketing.BudgetResponse{
		Data: budget,
	}, nil
}

func (h *MarketingGrpcHandler) GetBudget(ctx context.Context, req *pb_marketing.GetBudgetRequest) (*pb_marketing.BudgetResponse, error) {
	// Get budget config by campaign ID (assuming the ID is campaign ID)
	budgetConfig, err := h.budgetUsecase.GetBudgetConfig(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "budget not found: %v", err)
	}

	// Map DTO response to protobuf structure
	budget := &pb_marketing.Budget{
		Id:          budgetConfig.ID,
		CampaignId:  budgetConfig.CampaignID,
		Amount:      budgetConfig.MonthlyLimit,
		Currency:    "VND", // Default currency
		Description: "Monthly budget limit",
		Status:      "ACTIVE",
		CreatedAt:   budgetConfig.CreatedAt,
		UpdatedAt:   budgetConfig.UpdatedAt,
	}

	return &pb_marketing.BudgetResponse{
		Data: budget,
	}, nil
}

func (h *MarketingGrpcHandler) UpdateBudget(ctx context.Context, req *pb_marketing.UpdateBudgetRequest) (*pb_marketing.BudgetResponse, error) {
	// Map protobuf request to DTO structure
	updateDTO := &dto.BudgetConfigUpdateDTO{
		MonthlyLimit: &req.Amount,
	}

	// Update budget config
	budgetConfig, err := h.budgetUsecase.UpdateBudgetConfig(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update budget: %v", err)
	}

	// Map DTO response to protobuf structure
	budget := &pb_marketing.Budget{
		Id:          budgetConfig.ID,
		CampaignId:  budgetConfig.CampaignID,
		Amount:      budgetConfig.MonthlyLimit,
		Currency:    "VND", // Default currency
		Description: "Monthly budget limit",
		Status:      "ACTIVE",
		CreatedAt:   budgetConfig.CreatedAt,
		UpdatedAt:   budgetConfig.UpdatedAt,
	}

	return &pb_marketing.BudgetResponse{
		Data: budget,
	}, nil
}

func (h *MarketingGrpcHandler) DeleteBudget(ctx context.Context, req *pb_marketing.DeleteBudgetRequest) (*pb_marketing.DeleteBudgetResponse, error) {
	err := h.budgetUsecase.DeleteBudgetConfig(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete budget: %v", err)
	}

	return &pb_marketing.DeleteBudgetResponse{
		Success: true,
		Message: "Budget deleted successfully",
	}, nil
}

func (h *MarketingGrpcHandler) ListBudgets(ctx context.Context, req *pb_marketing.ListBudgetsRequest) (*pb_marketing.ListBudgetsResponse, error) {
	// Create search DTO

	searchDTO := dto.BudgetSearchDTO{
		Page: 0,
		Size: 10,
	}

	if req.Page > 0 {
		searchDTO.Page = int(req.Page)
	}

	if req.Size > 0 {
		searchDTO.Size = int(req.Size)
	}

	// Get budget list
	budgets, total, err := h.budgetUsecase.GetBudgetList(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list budgets: %v", err)
	}

	// Map DTO responses to protobuf structure
	grpcBudgets := make([]*pb_marketing.Budget, len(budgets))
	for i, budget := range budgets {
		grpcBudgets[i] = &pb_marketing.Budget{
			Id:          budget.ID,
			CampaignId:  budget.CampaignID,
			Amount:      budget.MonthlyLimit,
			Currency:    "VND", // Default currency
			Description: "Monthly budget limit",
			Status:      "ACTIVE",
			CreatedAt:   budget.CreatedAt,
			UpdatedAt:   budget.UpdatedAt,
		}
	}

	return &pb_marketing.ListBudgetsResponse{
		Data:  grpcBudgets,
		Total: uint32(total),
	}, nil
}

func (h *MarketingGrpcHandler) CreateAdvertisingSettings(ctx context.Context, req *pb_marketing.CreateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Map protobuf request to DTO structure
	createDTO := &dto.AdvertisingSettingsCreateDTO{
		OrganizationID:      1,       // TODO: Get from context
		UserID:              1,       // TODO: Get from context
		DefaultCampaignType: "VIP",   // Default campaign type
		DefaultBudget:       100000,  // Default budget
		DefaultDuration:     7,       // Default duration
		MonthlyBudgetLimit:  1000000, // Default monthly limit
		AutoSuggestEnabled:  false,
		AutoRunEnabled:      false,
		AlertOnBudgetExceed: false,
	}

	// Create advertising settings
	settings, err := h.advertisingSettingsUsecase.CreateAdvertisingSettings(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create advertising settings: %v", err)
	}

	// Map DTO response to protobuf structure
	grpcSettings := h.advertisingSettingsToGrpc(settings)

	return &pb_marketing.AdvertisingSettingsResponse{
		Data: grpcSettings,
	}, nil
}

func (h *MarketingGrpcHandler) GetAdvertisingSettings(ctx context.Context, req *pb_marketing.GetAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Get advertising settings by user ID (assuming the ID is user ID)
	userID := uint64(1) // TODO: Get from context or use req.Id if it represents user ID
	settings, err := h.advertisingSettingsUsecase.GetAdvertisingSettings(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "advertising settings not found: %v", err)
	}

	// Map DTO response to protobuf structure
	grpcSettings := h.advertisingSettingsToGrpc(settings)

	return &pb_marketing.AdvertisingSettingsResponse{
		Data: grpcSettings,
	}, nil
}

func (h *MarketingGrpcHandler) UpdateAdvertisingSettings(ctx context.Context, req *pb_marketing.UpdateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Map protobuf request to DTO structure
	updateDTO := &dto.AdvertisingSettingsUpdateDTO{
		// Map fields based on protobuf structure
		// For now, using minimal mapping
	}

	// Update advertising settings
	settings, err := h.advertisingSettingsUsecase.UpdateAdvertisingSettings(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update advertising settings: %v", err)
	}

	// Map DTO response to protobuf structure
	grpcSettings := h.advertisingSettingsToGrpc(settings)

	return &pb_marketing.AdvertisingSettingsResponse{
		Data: grpcSettings,
	}, nil
}

func (h *MarketingGrpcHandler) DeleteAdvertisingSettings(ctx context.Context, req *pb_marketing.DeleteAdvertisingSettingsRequest) (*pb_marketing.DeleteAdvertisingSettingsResponse, error) {
	err := h.advertisingSettingsUsecase.DeleteAdvertisingSettings(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete advertising settings: %v", err)
	}

	return &pb_marketing.DeleteAdvertisingSettingsResponse{
		Success: true,
		Message: "Advertising settings deleted successfully",
	}, nil
}

func (h *MarketingGrpcHandler) ListAdvertisingSettings(ctx context.Context, req *pb_marketing.ListAdvertisingSettingsRequest) (*pb_marketing.ListAdvertisingSettingsResponse, error) {
	searchDTO := dto.AdvertisingSettingsSearchDTO{
		Page:  int(req.Page),
		Limit: int(req.Size),
	}

	settings, total, err := h.advertisingSettingsUsecase.SearchAdvertisingSettings(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list advertising settings: %v", err)
	}

	grpcSettings := make([]*pb_marketing.AdvertisingSettings, len(settings))
	for i, setting := range settings {
		grpcSettings[i] = h.advertisingSettingsToGrpc(setting)
	}

	return &pb_marketing.ListAdvertisingSettingsResponse{
		Data:  grpcSettings,
		Total: uint32(total),
	}, nil
}

func getFloat64Value(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func getUint32Value(ptr *int) uint32 {
	if ptr == nil {
		return 0
	}
	return uint32(*ptr)
}

func getFloat64Ptr(value float64) *float64 {
	if value == 0 {
		return nil
	}
	return &value
}

func getIntPtr(value uint32) *int {
	if value == 0 {
		return nil
	}
	intValue := int(value)
	return &intValue
}

func getUint64Ptr(value uint64) *uint64 {
	if value == 0 {
		return nil
	}
	return &value
}

func (h *MarketingGrpcHandler) campaignToGrpcResponse(campaign *dto.CampaignResponseDTO) *pb_marketing.CampaignResponse {
	return &pb_marketing.CampaignResponse{
		Data: h.campaignToGrpcCampaign(campaign),
	}
}

func (h *MarketingGrpcHandler) campaignToGrpcCampaign(campaign *dto.CampaignResponseDTO) *pb_marketing.Campaign {
	grpcCampaign := &pb_marketing.Campaign{
		Id:             campaign.ID,
		Name:           campaign.Name,
		Description:    campaign.Description,
		Type:           string(campaign.Type),
		Status:         string(campaign.Status),
		ProductName:    campaign.ProductName,
		ProductType:    campaign.ProductType,
		PackageName:    campaign.PackageName,
		Budget:         campaign.Budget,
		Duration:       uint32(campaign.Duration),
		TargetAudience: campaign.TargetAudience,
		TargetLocation: campaign.TargetLocation,
		IsActive:       campaign.IsActive,
		Priority:       uint32(campaign.Priority),
		InternalNote:   campaign.InternalNote,
		CreatedBy:      campaign.CreatedBy,
		OrganizationId: campaign.OrganizationID,
		WalletBalance:  campaign.WalletBalance,
		Currency:       campaign.Currency,
		CreatedAt:      campaign.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      campaign.UpdatedAt.Format(time.RFC3339),
	}

	if campaign.ProductID != nil {
		grpcCampaign.ProductId = *campaign.ProductID
	}
	if campaign.PackageID != nil {
		grpcCampaign.PackageId = *campaign.PackageID
	}
	if campaign.WalletID != nil {
		grpcCampaign.WalletId = *campaign.WalletID
	}

	if campaign.StartDate != nil {
		grpcCampaign.StartDate = campaign.StartDate.Format(time.RFC3339)
	}
	if campaign.EndDate != nil {
		grpcCampaign.EndDate = campaign.EndDate.Format(time.RFC3339)
	}

	return grpcCampaign
}

func (h *MarketingGrpcHandler) packageToGrpcResponse(package_ *dto.PackageResponseDTO) *pb_marketing.PackageResponse {
	return &pb_marketing.PackageResponse{
		Data: h.packageToGrpcPackage(package_),
	}
}

func (h *MarketingGrpcHandler) packageToGrpcPackage(package_ *dto.PackageResponseDTO) *pb_marketing.Package {
	grpcPackage := &pb_marketing.Package{
		Id:                   package_.ID,
		Name:                 package_.Name,
		Description:          package_.Description,
		Type:                 string(package_.Type),
		Status:               string(package_.Status),
		Price:                package_.Price,
		Currency:             package_.Currency,
		Duration:             uint32(package_.Duration),
		Features:             package_.Features,
		IsPopular:            package_.IsPopular,
		IsRecommended:        package_.IsRecommended,
		Icon:                 package_.Icon,
		Color:                package_.Color,
		RequiresVerification: package_.RequiresVerification,
		RequiresPayment:      package_.RequiresPayment,
		CreatedBy:            package_.CreatedBy,
		OrganizationId:       package_.OrganizationID,
		CreatedAt:            package_.CreatedAt,
		UpdatedAt:            package_.UpdatedAt,
	}

	// Handle optional fields
	if package_.MaxImpressions != nil {
		grpcPackage.MaxImpressions = uint32(*package_.MaxImpressions)
	}
	if package_.MaxClicks != nil {
		grpcPackage.MaxClicks = uint32(*package_.MaxClicks)
	}
	if package_.MinBudget != nil {
		grpcPackage.MinBudget = *package_.MinBudget
	}
	if package_.MaxBudget != nil {
		grpcPackage.MaxBudget = *package_.MaxBudget
	}
	if package_.EstimatedReach != nil {
		grpcPackage.EstimatedReach = uint32(*package_.EstimatedReach)
	}
	if package_.EstimatedClicks != nil {
		grpcPackage.EstimatedClicks = uint32(*package_.EstimatedClicks)
	}

	return grpcPackage
}

func (h *MarketingGrpcHandler) advertisingSettingsToGrpc(setting *dto.AdvertisingSettingsDTO) *pb_marketing.AdvertisingSettings {
	// Map DTO to protobuf structure
	grpcSettings := &pb_marketing.AdvertisingSettings{
		Id:             setting.ID,
		CampaignId:     0, // Not available in DTO, using 0 as default
		Platform:       setting.DefaultCampaignType,
		TargetAudience: "DEFAULT", // Not available in DTO, using default
		TargetLocation: "DEFAULT", // Not available in DTO, using default
		AdFormat:       "DEFAULT", // Not available in DTO, using default
		Schedule:       "DEFAULT", // Not available in DTO, using default
		Status:         "ACTIVE",
		CreatedAt:      setting.CreatedAt,
		UpdatedAt:      setting.UpdatedAt,
	}

	// Convert settings map
	if setting.ProductTypeSettings != nil {
		grpcSettings.Settings = make(map[string]string)
		for key, value := range setting.ProductTypeSettings {
			// Convert ProductTypeSettingDTO to string representation
			grpcSettings.Settings[key] = fmt.Sprintf("budget:%.2f,duration:%d,priority:%d",
				value.DefaultBudget, value.DefaultDuration, value.Priority)
		}
	}

	return grpcSettings
}