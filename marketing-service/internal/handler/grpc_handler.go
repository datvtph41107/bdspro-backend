package handler

import (
	"context"
	"marketing/infrastructure/client"
	"marketing/internal/domain/model"
	"marketing/internal/dto"
	"marketing/internal/usecase"
	"time"

	pb_marketing "pb/types/marketing"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MarketingGrpcHandler struct {
	pb_marketing.UnimplementedMarketingServiceServer
	campaignUsecase     usecase.CampaignUsecase
	packageUsecase      usecase.PackageUsecase
	paymentUsecase      usecase.PaymentUsecase
	campaignDetailUsecase usecase.CampaignDetailUsecase
	budgetUsecase       usecase.BudgetUsecase
	advertisingSettingsUsecase usecase.AdvertisingSettingsUsecase
	notificationClient  *client.NotificationClient
}

func NewMarketingGrpcHandler(campaignUsecase usecase.CampaignUsecase, packageUsecase usecase.PackageUsecase, paymentUsecase usecase.PaymentUsecase, campaignDetailUsecase usecase.CampaignDetailUsecase, budgetUsecase usecase.BudgetUsecase, advertisingSettingsUsecase usecase.AdvertisingSettingsUsecase, notificationClient *client.NotificationClient) *MarketingGrpcHandler {
	return &MarketingGrpcHandler{
		campaignUsecase:      campaignUsecase,
		packageUsecase:       packageUsecase,
		paymentUsecase:       paymentUsecase,
		campaignDetailUsecase: campaignDetailUsecase,
		budgetUsecase:        budgetUsecase,
		advertisingSettingsUsecase: advertisingSettingsUsecase,
		notificationClient:   notificationClient,
	}
}

// CreateCampaign implements the CreateCampaign gRPC method
func (h *MarketingGrpcHandler) CreateCampaign(ctx context.Context, req *pb_marketing.CreateCampaignRequest) (*pb_marketing.CampaignResponse, error) {
	// Convert gRPC request to DTO
	createDTO := &dto.CampaignCreateDTO{
		Name:            req.Name,
		Description:     req.Description,
		Type:            model.CampaignType(req.Type),
		ProductID:       &req.ProductId,
		ProductName:     req.ProductName,
		ProductType:     req.ProductType,
		PackageID:       &req.PackageId,
		PackageName:     req.PackageName,
		Budget:          req.Budget,
		Duration:        int(req.Duration),
		TargetAudience:  req.TargetAudience,
		TargetLocation:  req.TargetLocation,
		Priority:        int(req.Priority),
		InternalNote:    req.InternalNote,
		TermsAccepted:   req.TermsAccepted,
	}

	// Parse dates if provided
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

	// Call usecase
	campaign, err := h.campaignUsecase.CreateCampaign(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create campaign: %v", err)
	}

	// Convert to gRPC response
	return h.campaignToGrpcResponse(campaign), nil
}

// UpdateCampaign implements the UpdateCampaign gRPC method
func (h *MarketingGrpcHandler) UpdateCampaign(ctx context.Context, req *pb_marketing.UpdateCampaignRequest) (*pb_marketing.CampaignResponse, error) {
	// Convert gRPC request to DTO
	updateDTO := &dto.CampaignUpdateDTO{
		Name:            req.Name,
		Description:     req.Description,
		Type:            model.CampaignType(req.Type),
		ProductID:       &req.ProductId,
		ProductName:     req.ProductName,
		ProductType:     req.ProductType,
		PackageID:       &req.PackageId,
		PackageName:     req.PackageName,
		Budget:          req.Budget,
		Duration:        int(req.Duration),
		TargetAudience:  req.TargetAudience,
		TargetLocation:  req.TargetLocation,
		Status:          model.CampaignStatus(req.Status),
		Priority:        int(req.Priority),
		InternalNote:    req.InternalNote,
	}

	// Parse dates if provided
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

	// Call usecase
	campaign, err := h.campaignUsecase.UpdateCampaign(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update campaign: %v", err)
	}

	// Convert to gRPC response
	return h.campaignToGrpcResponse(campaign), nil
}

// DeleteCampaign implements the DeleteCampaign gRPC method
func (h *MarketingGrpcHandler) DeleteCampaign(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	err := h.campaignUsecase.DeleteCampaign(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete campaign: %v", err)
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Campaign deleted successfully",
	}, nil
}

// GetCampaign implements the GetCampaign gRPC method
func (h *MarketingGrpcHandler) GetCampaign(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignResponse, error) {
	campaign, err := h.campaignUsecase.GetCampaign(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "campaign not found: %v", err)
	}

	// Convert to gRPC response
	return h.campaignToGrpcResponse(campaign), nil
}

// SearchCampaigns implements the SearchCampaigns gRPC method
func (h *MarketingGrpcHandler) SearchCampaigns(ctx context.Context, req *pb_marketing.SearchCampaignsRequest) (*pb_marketing.SearchCampaignsResponse, error) {
	// Convert gRPC request to DTO
	searchDTO := dto.CampaignSearchDTO{
		Name:        req.Name,
		Type:        model.CampaignType(req.Type),
		Status:      model.CampaignStatus(req.Status),
		ProductID:   &req.ProductId,
		PackageID:   &req.PackageId,
		Page:        int(req.Page),
		Limit:       int(req.Size),
	}

	// Parse dates if provided
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

	// Convert to gRPC response
	grpcCampaigns := make([]*pb_marketing.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		grpcCampaigns[i] = h.campaignToGrpcCampaign(&campaign)
	}

	return &pb_marketing.SearchCampaignsResponse{
		Data:  grpcCampaigns,
		Total: uint32(total),
	}, nil
}

// GetCampaignsByProduct implements the GetCampaignsByProduct gRPC method
func (h *MarketingGrpcHandler) GetCampaignsByProduct(ctx context.Context, req *pb_marketing.GetCampaignsByProductRequest) (*pb_marketing.SearchCampaignsResponse, error) {
	campaigns, err := h.campaignUsecase.GetCampaignsByProduct(ctx, req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaigns by product: %v", err)
	}

	// Convert to gRPC response
	grpcCampaigns := make([]*pb_marketing.Campaign, len(campaigns))
	for i, campaign := range campaigns {
		grpcCampaigns[i] = h.campaignToGrpcCampaign(&campaign)
	}

	return &pb_marketing.SearchCampaignsResponse{
		Data:  grpcCampaigns,
		Total: uint32(len(campaigns)),
	}, nil
}

// UpdateCampaignStatus implements the UpdateCampaignStatus gRPC method
func (h *MarketingGrpcHandler) UpdateCampaignStatus(ctx context.Context, req *pb_marketing.UpdateCampaignStatusRequest) (*sharepb.SubmitResponse, error) {
	err := h.campaignUsecase.UpdateCampaignStatus(ctx, req.Id, model.CampaignStatus(req.Status))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update campaign status: %v", err)
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Campaign status updated successfully",
	}, nil
}

// TopupCampaignBudget implements the TopupCampaignBudget gRPC method
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

// GetCampaignWalletInfo implements the GetCampaignWalletInfo gRPC method
func (h *MarketingGrpcHandler) GetCampaignWalletInfo(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignWalletInfoResponse, error) {
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

// GetCampaignTransactions implements the GetCampaignTransactions gRPC method
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

// Payment management - D.10.2.3

// GetPaymentSummary implements the GetPaymentSummary gRPC method
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

// ValidatePayment implements the ValidatePayment gRPC method
func (h *MarketingGrpcHandler) ValidatePayment(ctx context.Context, req *pb_marketing.ValidatePaymentRequest) (*pb_marketing.PaymentValidationResponse, error) {
	paymentDTO := &dto.PaymentRequestDTO{
		CampaignID:      req.CampaignId,
		PackageID:       req.PackageId,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		WalletID:        req.WalletId,
		ProductID:       req.ProductId,
		ProductName:     req.ProductName,
		ProductType:     req.ProductType,
		CustomBudget:    getFloat64Ptr(req.CustomBudget),
		CustomDuration:  getIntPtr(req.CustomDuration),
		RequireInvoice:  req.RequireInvoice,
		InvoiceEmail:    req.InvoiceEmail,
		TermsAccepted:   req.TermsAccepted,
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

// ProcessPayment implements the ProcessPayment gRPC method
func (h *MarketingGrpcHandler) ProcessPayment(ctx context.Context, req *pb_marketing.ProcessPaymentRequest) (*pb_marketing.PaymentResponse, error) {
	paymentDTO := &dto.PaymentRequestDTO{
		CampaignID:      req.CampaignId,
		PackageID:       req.PackageId,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		WalletID:        req.WalletId,
		ProductID:       req.ProductId,
		ProductName:     req.ProductName,
		ProductType:     req.ProductType,
		CustomBudget:    getFloat64Ptr(req.CustomBudget),
		CustomDuration:  getIntPtr(req.CustomDuration),
		RequireInvoice:  req.RequireInvoice,
		InvoiceEmail:    req.InvoiceEmail,
		TermsAccepted:   req.TermsAccepted,
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

// GetWalletBalance implements the GetWalletBalance gRPC method
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

// GetPaymentMethods implements the GetPaymentMethods gRPC method
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

// GetPaymentHistory implements the GetPaymentHistory gRPC method
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

// GetCampaignPayment implements the GetCampaignPayment gRPC method
func (h *MarketingGrpcHandler) GetCampaignPayment(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignPaymentResponse, error) {
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

// Campaign detail management - D.10.2.4

// GetCampaignDetail implements the GetCampaignDetail gRPC method
func (h *MarketingGrpcHandler) GetCampaignDetail(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignDetailResponse, error) {
	result, err := h.campaignDetailUsecase.GetCampaignDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign detail: %v", err)
	}

	return &pb_marketing.CampaignDetailResponse{
		Data: &pb_marketing.CampaignDetail{
			Id:              result.ID,
			Name:            result.Name,
			Description:     result.Description,
			Type:            result.Type,
			Status:          result.Status,
			ProductId:       result.ProductID,
			ProductName:     result.ProductName,
			ProductImage:    result.ProductImage,
			ProductType:     result.ProductType,
			PackageId:       result.PackageID,
			PackageName:     result.PackageName,
			PackageType:     result.PackageType,
			Budget:          result.Budget,
			Duration:        uint32(result.Duration),
			StartDate:       result.StartDate,
			EndDate:         result.EndDate,
			Currency:        result.Currency,
			Impressions:     uint32(result.Impressions),
			Clicks:          uint32(result.Clicks),
			Ctr:             result.CTR,
			Cpc:             result.CPC,
			SpentAmount:     result.SpentAmount,
			RemainingBudget: result.RemainingBudget,
			WalletId:        result.WalletID,
			WalletBalance:   result.WalletBalance,
			TransactionId:   result.TransactionID,
			CanEdit:         result.CanEdit,
			CanPause:        result.CanPause,
			CanResume:       result.CanResume,
			CanExtend:       result.CanExtend,
			CanTopup:        result.CanTopup,
			CreatedBy:       result.CreatedBy,
			OrganizationId:  result.OrganizationID,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		},
	}, nil
}

// GetCampaignPerformance implements the GetCampaignPerformance gRPC method
func (h *MarketingGrpcHandler) GetCampaignPerformance(ctx context.Context, req *pb_marketing.GetCampaignPerformanceRequest) (*pb_marketing.CampaignPerformanceResponse, error) {
	result, err := h.campaignDetailUsecase.GetCampaignPerformance(ctx, req.Id, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign performance: %v", err)
	}

	// Convert daily data
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
			CampaignId: result.CampaignID,
			Date:       result.Date,
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

// GetCampaignHistory implements the GetCampaignHistory gRPC method
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

// ExtendCampaign implements the ExtendCampaign gRPC method
func (h *MarketingGrpcHandler) ExtendCampaign(ctx context.Context, req *pb_marketing.ExtendCampaignRequest) (*pb_marketing.PaymentResponse, error) {
	extendDTO := &dto.CampaignExtendDTO{
		CampaignID:      req.Id,
		ExtensionDays:   int(req.ExtensionDays),
		AdditionalBudget: req.AdditionalBudget,
		PaymentMethod:   req.PaymentMethod,
		TermsAccepted:   req.TermsAccepted,
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

// ExportCampaignData implements the ExportCampaignData gRPC method
func (h *MarketingGrpcHandler) ExportCampaignData(ctx context.Context, req *pb_marketing.ExportCampaignDataRequest) (*pb_marketing.CampaignExportResponse, error) {
	result, err := h.campaignDetailUsecase.ExportCampaignData(ctx, req.Id, req.ExportType, req.DateRange)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to export campaign data: %v", err)
	}

	// Convert performance data
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
			GeneratedAt:  result.GeneratedAt,
			GeneratedBy:  result.GeneratedBy,
		},
	}, nil
}

// GetCampaignOptimization implements the GetCampaignOptimization gRPC method
func (h *MarketingGrpcHandler) GetCampaignOptimization(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignOptimizationResponse, error) {
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

// GetCampaignTransactionHistory method đã được xóa
// Sử dụng trực tiếp API từ notification service:
// GET /v2/notification/history/campaign/{id}
// GET /v2/notification/history/package/{id}

// Budget Management APIs - D.10.2.5
// TODO: Uncomment when protobuf definitions are added

/*
// CreateBudgetConfig implements the CreateBudgetConfig gRPC method
func (h *MarketingGrpcHandler) CreateBudgetConfig(ctx context.Context, req *pb_marketing.CreateBudgetConfigRequest) (*pb_marketing.BudgetConfigResponse, error) {
	// Convert gRPC request to DTO
	createDTO := &dto.BudgetConfigCreateDTO{
		CampaignID:         req.CampaignId,
		MonthlyLimit:       req.MonthlyLimit,
		DailyLimit:         getFloat64Ptr(req.DailyLimit),
		TotalLimit:         getFloat64Ptr(req.TotalLimit),
		AlertThreshold:     req.AlertThreshold,
		AlertEnabled:       req.AlertEnabled,
		AlertEmail:         req.AlertEmail,
		AlertPhone:         req.AlertPhone,
		AutoTopup:          req.AutoTopup,
		AutoTopupAmount:    req.AutoTopupAmount,
		AutoTopupThreshold: req.AutoTopupThreshold,
	}

	// Call usecase
	config, err := h.budgetUsecase.CreateBudgetConfig(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create budget config: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetConfigResponse{
		Data: &pb_marketing.BudgetConfig{
			Id:                 config.ID,
			CampaignId:         config.CampaignID,
			OrganizationId:     config.OrganizationID,
			MonthlyLimit:       config.MonthlyLimit,
			DailyLimit:         getFloat64Value(config.DailyLimit),
			TotalLimit:         getFloat64Value(config.TotalLimit),
			AlertThreshold:     config.AlertThreshold,
			AlertEnabled:       config.AlertEnabled,
			AlertEmail:         config.AlertEmail,
			AlertPhone:         config.AlertPhone,
			AutoTopup:          config.AutoTopup,
			AutoTopupAmount:    config.AutoTopupAmount,
			AutoTopupThreshold: config.AutoTopupThreshold,
			CreatedBy:          config.CreatedBy,
			UpdatedBy:          config.UpdatedBy,
			CreatedAt:          config.CreatedAt,
			UpdatedAt:          config.UpdatedAt,
		},
	}, nil
}

// UpdateBudgetConfig implements the UpdateBudgetConfig gRPC method
func (h *MarketingGrpcHandler) UpdateBudgetConfig(ctx context.Context, req *pb_marketing.UpdateBudgetConfigRequest) (*pb_marketing.BudgetConfigResponse, error) {
	// Convert gRPC request to DTO
	updateDTO := &dto.BudgetConfigUpdateDTO{
		MonthlyLimit:       getFloat64Ptr(req.MonthlyLimit),
		DailyLimit:         getFloat64Ptr(req.DailyLimit),
		TotalLimit:         getFloat64Ptr(req.TotalLimit),
		AlertThreshold:     getFloat64Ptr(req.AlertThreshold),
		AlertEnabled:       &req.AlertEnabled,
		AlertEmail:         &req.AlertEmail,
		AlertPhone:         &req.AlertPhone,
		AutoTopup:          &req.AutoTopup,
		AutoTopupAmount:    getFloat64Ptr(req.AutoTopupAmount),
		AutoTopupThreshold: getFloat64Ptr(req.AutoTopupThreshold),
	}

	// Call usecase
	config, err := h.budgetUsecase.UpdateBudgetConfig(ctx, req.CampaignId, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update budget config: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetConfigResponse{
		Data: &pb_marketing.BudgetConfig{
			Id:                 config.ID,
			CampaignId:         config.CampaignID,
			OrganizationId:     config.OrganizationID,
			MonthlyLimit:       config.MonthlyLimit,
			DailyLimit:         getFloat64Value(config.DailyLimit),
			TotalLimit:         getFloat64Value(config.TotalLimit),
			AlertThreshold:     config.AlertThreshold,
			AlertEnabled:       config.AlertEnabled,
			AlertEmail:         config.AlertEmail,
			AlertPhone:         config.AlertPhone,
			AutoTopup:          config.AutoTopup,
			AutoTopupAmount:    config.AutoTopupAmount,
			AutoTopupThreshold: config.AutoTopupThreshold,
			CreatedBy:          config.CreatedBy,
			UpdatedBy:          config.UpdatedBy,
			CreatedAt:          config.CreatedAt,
			UpdatedAt:          config.UpdatedAt,
		},
	}, nil
}

// GetBudgetConfig implements the GetBudgetConfig gRPC method
func (h *MarketingGrpcHandler) GetBudgetConfig(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.BudgetConfigResponse, error) {
	// Call usecase
	config, err := h.budgetUsecase.GetBudgetConfig(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budget config: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetConfigResponse{
		Data: &pb_marketing.BudgetConfig{
			Id:                 config.ID,
			CampaignId:         config.CampaignID,
			OrganizationId:     config.OrganizationID,
			MonthlyLimit:       config.MonthlyLimit,
			DailyLimit:         getFloat64Value(config.DailyLimit),
			TotalLimit:         getFloat64Value(config.TotalLimit),
			AlertThreshold:     config.AlertThreshold,
			AlertEnabled:       config.AlertEnabled,
			AlertEmail:         config.AlertEmail,
			AlertPhone:         config.AlertPhone,
			AutoTopup:          config.AutoTopup,
			AutoTopupAmount:    config.AutoTopupAmount,
			AutoTopupThreshold: config.AutoTopupThreshold,
			CreatedBy:          config.CreatedBy,
			UpdatedBy:          config.UpdatedBy,
			CreatedAt:          config.CreatedAt,
			UpdatedAt:          config.UpdatedAt,
		},
	}, nil
}

// DeleteBudgetConfig implements the DeleteBudgetConfig gRPC method
func (h *MarketingGrpcHandler) DeleteBudgetConfig(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Call usecase
	err := h.budgetUsecase.DeleteBudgetConfig(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete budget config: %v", err)
	}

	return &sharepb.SubmitResponse{
		Success: true,
		Message: "Budget config deleted successfully",
	}, nil
}

// GetBudgetList implements the GetBudgetList gRPC method
func (h *MarketingGrpcHandler) GetBudgetList(ctx context.Context, req *pb_marketing.GetBudgetListRequest) (*pb_marketing.BudgetListResponse, error) {
	// Convert gRPC request to DTO
	searchDTO := dto.BudgetSearchDTO{
		CampaignID:      getUint64Ptr(req.CampaignId),
		CampaignStatus:  req.CampaignStatus,
		CampaignType:    req.CampaignType,
		AlertEnabled:    &req.AlertEnabled,
		IsAlertTriggered: &req.IsAlertTriggered,
		UsageMin:        getFloat64Ptr(req.UsageMin),
		UsageMax:        getFloat64Ptr(req.UsageMax),
		Page:            int(req.Page),
		Limit:           int(req.Limit),
	}

	// Call usecase
	budgets, total, err := h.budgetUsecase.GetBudgetList(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budget list: %v", err)
	}

	// Convert to gRPC response
	var budgetList []*pb_marketing.BudgetList
	for _, budget := range budgets {
		budgetList = append(budgetList, &pb_marketing.BudgetList{
			Id:                budget.ID,
			CampaignId:        budget.CampaignID,
			CampaignName:      budget.CampaignName,
			CampaignStatus:    budget.CampaignStatus,
			CampaignType:      budget.CampaignType,
			TotalBudget:       budget.TotalBudget,
			SpentAmount:       budget.SpentAmount,
			RemainingBudget:   budget.RemainingBudget,
			UsagePercentage:   budget.UsagePercentage,
			MonthlyLimit:      budget.MonthlyLimit,
			MonthlySpent:      budget.MonthlySpent,
			MonthlyUsage:      budget.MonthlyUsage,
			AlertEnabled:      budget.AlertEnabled,
			IsAlertTriggered:  budget.IsAlertTriggered,
			Currency:          budget.Currency,
			CreatedAt:         budget.CreatedAt,
			UpdatedAt:         budget.UpdatedAt,
		})
	}

	return &pb_marketing.BudgetListResponse{
		Data:  budgetList,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetBudgetStatus implements the GetBudgetStatus gRPC method
func (h *MarketingGrpcHandler) GetBudgetStatus(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.BudgetStatusResponse, error) {
	// Call usecase
	status, err := h.budgetUsecase.GetBudgetStatus(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budget status: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetStatusResponse{
		Data: &pb_marketing.BudgetStatus{
			CampaignId:        status.CampaignID,
			CampaignName:      status.CampaignName,
			CampaignStatus:    status.CampaignStatus,
			TotalBudget:       status.TotalBudget,
			SpentAmount:       status.SpentAmount,
			RemainingBudget:   status.RemainingBudget,
			UsagePercentage:   status.UsagePercentage,
			MonthlyLimit:      status.MonthlyLimit,
			MonthlySpent:      status.MonthlySpent,
			MonthlyRemaining:  status.MonthlyRemaining,
			MonthlyUsage:      status.MonthlyUsage,
			DailyLimit:        getFloat64Value(status.DailyLimit),
			DailySpent:        status.DailySpent,
			DailyRemaining:    getFloat64Value(status.DailyRemaining),
			DailyUsage:        getFloat64Value(status.DailyUsage),
			AlertThreshold:    status.AlertThreshold,
			AlertEnabled:      status.AlertEnabled,
			IsAlertTriggered:  status.IsAlertTriggered,
			AlertMessage:      status.AlertMessage,
			Currency:          status.Currency,
			LastUpdated:       status.LastUpdated,
		},
	}, nil
}

// TopupBudget implements the TopupBudget gRPC method
func (h *MarketingGrpcHandler) TopupBudget(ctx context.Context, req *pb_marketing.TopupBudgetRequest) (*pb_marketing.BudgetTopupResponse, error) {
	// Convert gRPC request to DTO
	topupDTO := &dto.BudgetTopupDTO{
		CampaignID:    req.CampaignId,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Description:   req.Description,
		TermsAccepted: req.TermsAccepted,
	}

	// Call usecase
	result, err := h.budgetUsecase.TopupBudget(ctx, topupDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to topup budget: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetTopupResponse{
		Success:        result.Success,
		TransactionId:  result.TransactionID,
		CampaignId:     result.CampaignID,
		Amount:         result.Amount,
		PreviousBudget: result.PreviousBudget,
		NewBudget:      result.NewBudget,
		Currency:       result.Currency,
		Status:         result.Status,
		Message:        result.Message,
		CreatedAt:      result.CreatedAt,
	}, nil
}

// GetBudgetAlerts implements the GetBudgetAlerts gRPC method
func (h *MarketingGrpcHandler) GetBudgetAlerts(ctx context.Context, req *pb_marketing.GetBudgetAlertsRequest) (*pb_marketing.BudgetAlertsResponse, error) {
	// Call usecase
	alerts, total, err := h.budgetUsecase.GetBudgetAlerts(ctx, req.CampaignId, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budget alerts: %v", err)
	}

	// Convert to gRPC response
	var alertList []*pb_marketing.BudgetAlert
	for _, alert := range alerts {
		alertList = append(alertList, &pb_marketing.BudgetAlert{
			Id:           alert.ID,
			CampaignId:   alert.CampaignID,
			CampaignName: alert.CampaignName,
			AlertType:    alert.AlertType,
			AlertLevel:   alert.AlertLevel,
			Message:      alert.Message,
			CurrentUsage: alert.CurrentUsage,
			Threshold:    alert.Threshold,
			IsRead:       alert.IsRead,
			CreatedAt:    alert.CreatedAt,
		})
	}

	return &pb_marketing.BudgetAlertsResponse{
		Data:  alertList,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// MarkAlertAsRead implements the MarkAlertAsRead gRPC method
func (h *MarketingGrpcHandler) MarkAlertAsRead(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Call usecase
	err := h.budgetUsecase.MarkAlertAsRead(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to mark alert as read: %v", err)
	}

	return &sharepb.SubmitResponse{
		Success: true,
		Message: "Alert marked as read successfully",
	}, nil
}

// GetBudgetUsageSummary implements the GetBudgetUsageSummary gRPC method
func (h *MarketingGrpcHandler) GetBudgetUsageSummary(ctx context.Context, req *pb_marketing.GetBudgetUsageSummaryRequest) (*pb_marketing.BudgetUsageSummaryResponse, error) {
	// Call usecase
	summary, err := h.budgetUsecase.GetBudgetUsageSummary(ctx, req.CampaignId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get budget usage summary: %v", err)
	}

	// Convert usage data
	var usageData []*pb_marketing.BudgetUsage
	for _, data := range summary.UsageData {
		usageData = append(usageData, &pb_marketing.BudgetUsage{
			CampaignId:   data.CampaignID,
			Date:         data.Date,
			Impressions:  uint32(data.Impressions),
			Clicks:       uint32(data.Clicks),
			SpentAmount:  data.SpentAmount,
			Ctr:          data.CTR,
			Cpc:          data.CPC,
			DailyLimit:   getFloat64Value(data.DailyLimit),
			DailyUsage:   getFloat64Value(data.DailyUsage),
		})
	}

	// Convert to gRPC response
	return &pb_marketing.BudgetUsageSummaryResponse{
		Data: &pb_marketing.BudgetUsageSummary{
			CampaignId:        summary.CampaignID,
			Period:            summary.Period,
			TotalImpressions:  uint32(summary.TotalImpressions),
			TotalClicks:       uint32(summary.TotalClicks),
			TotalSpent:        summary.TotalSpent,
			AverageCtr:        summary.AverageCTR,
			AverageCpc:        summary.AverageCPC,
			UsageData:         usageData,
		},
	}, nil
}
*/

// Advertising Settings APIs - D.10.2.8
// TODO: Uncomment when protobuf definitions are added

/*
// CreateAdvertisingSettings implements the CreateAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) CreateAdvertisingSettings(ctx context.Context, req *pb_marketing.CreateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Convert gRPC request to DTO
	createDTO := &dto.AdvertisingSettingsCreateDTO{
		OrganizationID:         req.OrganizationId,
		UserID:                 req.UserId,
		DefaultCampaignType:    req.DefaultCampaignType,
		DefaultPackageID:       getUint64Ptr(req.DefaultPackageId),
		DefaultBudget:          req.DefaultBudget,
		DefaultDuration:        int(req.DefaultDuration),
		AutoSuggestEnabled:     req.AutoSuggestEnabled,
		AutoRunEnabled:         req.AutoRunEnabled,
		AutoRunBudget:          req.AutoRunBudget,
		AutoRunPackageID:       getUint64Ptr(req.AutoRunPackageId),
		MonthlyBudgetLimit:     req.MonthlyBudgetLimit,
		AlertOnBudgetExceed:    req.AlertOnBudgetExceed,
		AlertEmail:             req.AlertEmail,
		AlertPhone:             req.AlertPhone,
		ApplyToAllOrganization: req.ApplyToAllOrganization,
	}

	// Call usecase
	settings, err := h.advertisingSettingsUsecase.CreateAdvertisingSettings(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create advertising settings: %v", err)
	}

	// Convert to gRPC response
	return h.advertisingSettingsToGrpcResponse(settings), nil
}

// UpdateAdvertisingSettings implements the UpdateAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) UpdateAdvertisingSettings(ctx context.Context, req *pb_marketing.UpdateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Convert gRPC request to DTO
	updateDTO := &dto.AdvertisingSettingsUpdateDTO{
		DefaultCampaignType:    &req.DefaultCampaignType,
		DefaultPackageID:       getUint64Ptr(req.DefaultPackageId),
		DefaultBudget:          getFloat64Ptr(req.DefaultBudget),
		DefaultDuration:        getIntPtr(req.DefaultDuration),
		AutoSuggestEnabled:     &req.AutoSuggestEnabled,
		AutoRunEnabled:         &req.AutoRunEnabled,
		AutoRunBudget:          getFloat64Ptr(req.AutoRunBudget),
		AutoRunPackageID:       getUint64Ptr(req.AutoRunPackageId),
		MonthlyBudgetLimit:     getFloat64Ptr(req.MonthlyBudgetLimit),
		AlertOnBudgetExceed:    &req.AlertOnBudgetExceed,
		AlertEmail:             &req.AlertEmail,
		AlertPhone:             &req.AlertPhone,
		ApplyToAllOrganization: &req.ApplyToAllOrganization,
	}

	// Call usecase
	settings, err := h.advertisingSettingsUsecase.UpdateAdvertisingSettings(ctx, req.Id, updateDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update advertising settings: %v", err)
	}

	// Convert to gRPC response
	return h.advertisingSettingsToGrpcResponse(settings), nil
}

// GetAdvertisingSettings implements the GetAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) GetAdvertisingSettings(ctx context.Context, req *pb_marketing.GetAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	// Call usecase
	settings, err := h.advertisingSettingsUsecase.GetAdvertisingSettings(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get advertising settings: %v", err)
	}

	// Convert to gRPC response
	return h.advertisingSettingsToGrpcResponse(settings), nil
}

// DeleteAdvertisingSettings implements the DeleteAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) DeleteAdvertisingSettings(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Call usecase
	err := h.advertisingSettingsUsecase.DeleteAdvertisingSettings(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete advertising settings: %v", err)
	}

	return &sharepb.SubmitResponse{
		Success: true,
		Message: "Advertising settings deleted successfully",
	}, nil
}

// SearchAdvertisingSettings implements the SearchAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) SearchAdvertisingSettings(ctx context.Context, req *pb_marketing.SearchAdvertisingSettingsRequest) (*pb_marketing.SearchAdvertisingSettingsResponse, error) {
	// Convert gRPC request to DTO
	searchDTO := dto.AdvertisingSettingsSearchDTO{
		OrganizationID:      getUint64Ptr(req.OrganizationId),
		UserID:              getUint64Ptr(req.UserId),
		DefaultCampaignType: req.DefaultCampaignType,
		AutoSuggestEnabled:  &req.AutoSuggestEnabled,
		AutoRunEnabled:      &req.AutoRunEnabled,
		IsActive:            &req.IsActive,
		Page:                int(req.Page),
		Limit:               int(req.Limit),
	}

	// Call usecase
	settings, total, err := h.advertisingSettingsUsecase.SearchAdvertisingSettings(ctx, searchDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search advertising settings: %v", err)
	}

	// Convert to gRPC response
	grpcSettings := make([]*pb_marketing.AdvertisingSettings, len(settings))
	for i, setting := range settings {
		grpcSettings[i] = h.advertisingSettingsToGrpcSettings(setting)
	}

	return &pb_marketing.SearchAdvertisingSettingsResponse{
		Data: &pb_marketing.AdvertisingSettingsList{
			Settings: grpcSettings,
			Total:    total,
			Page:     req.Page,
			Limit:    req.Limit,
		},
	}, nil
}

// GetCampaignDefaults implements the GetCampaignDefaults gRPC method
func (h *MarketingGrpcHandler) GetCampaignDefaults(ctx context.Context, req *pb_marketing.GetCampaignDefaultsRequest) (*pb_marketing.CampaignDefaultsResponse, error) {
	// Call usecase
	defaults, err := h.advertisingSettingsUsecase.GetCampaignDefaults(ctx, req.ProductType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get campaign defaults: %v", err)
	}

	// Convert to gRPC response
	grpcDefaults := &pb_marketing.CampaignDefaults{
		CampaignType:   defaults.CampaignType,
		PackageId:      getUint64Value(defaults.PackageID),
		PackageName:    defaults.PackageName,
		Budget:         defaults.Budget,
		Duration:       uint32(defaults.Duration),
		TargetAudience: defaults.TargetAudience,
		TargetLocation: defaults.TargetLocation,
		Priority:       uint32(defaults.Priority),
		AutoRunEnabled: defaults.AutoRunEnabled,
		SettingsSource: defaults.SettingsSource,
	}

	return &pb_marketing.CampaignDefaultsResponse{
		Data: grpcDefaults,
	}, nil
}

// GetAutoCampaignSuggestions implements the GetAutoCampaignSuggestions gRPC method
func (h *MarketingGrpcHandler) GetAutoCampaignSuggestions(ctx context.Context, req *pb_marketing.GetAutoCampaignSuggestionsRequest) (*pb_marketing.AutoCampaignSuggestionsResponse, error) {
	// Call usecase
	suggestions, total, err := h.advertisingSettingsUsecase.GetAutoCampaignSuggestions(ctx, int(req.Page), int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get auto campaign suggestions: %v", err)
	}

	// Convert to gRPC response
	grpcSuggestions := make([]*pb_marketing.AutoCampaignSuggestion, len(suggestions))
	for i, suggestion := range suggestions {
		grpcSuggestions[i] = &pb_marketing.AutoCampaignSuggestion{
			ProductId:              suggestion.ProductID,
			ProductName:            suggestion.ProductName,
			ProductType:            suggestion.ProductType,
			SuggestedCampaignType: suggestion.SuggestedCampaignType,
			SuggestedPackageId:    suggestion.SuggestedPackageID,
			SuggestedPackageName:  suggestion.SuggestedPackageName,
			SuggestedBudget:       suggestion.SuggestedBudget,
			SuggestedDuration:     uint32(suggestion.SuggestedDuration),
			Reason:                suggestion.Reason,
			Confidence:            suggestion.Confidence,
			EstimatedReach:        getUint32Value(suggestion.EstimatedReach),
			EstimatedClicks:       getUint32Value(suggestion.EstimatedClicks),
		}
	}

	return &pb_marketing.AutoCampaignSuggestionsResponse{
		Data: &pb_marketing.AutoCampaignSuggestionsList{
			Suggestions: grpcSuggestions,
			Total:       total,
			Page:        req.Page,
			Limit:       req.Limit,
		},
	}, nil
}

// AcceptAutoCampaignSuggestion implements the AcceptAutoCampaignSuggestion gRPC method
func (h *MarketingGrpcHandler) AcceptAutoCampaignSuggestion(ctx context.Context, req *sharepb.IdRequest) (*pb_marketing.CampaignResponse, error) {
	// Call usecase
	campaignDTO, err := h.advertisingSettingsUsecase.AcceptAutoCampaignSuggestion(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept auto campaign suggestion: %v", err)
	}

	// Create campaign from suggestion
	campaign, err := h.campaignUsecase.CreateCampaign(ctx, campaignDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create campaign from suggestion: %v", err)
	}

	// Convert to gRPC response
	return h.campaignToGrpcResponse(campaign), nil
}

// RejectAutoCampaignSuggestion implements the RejectAutoCampaignSuggestion gRPC method
func (h *MarketingGrpcHandler) RejectAutoCampaignSuggestion(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	// Call usecase
	err := h.advertisingSettingsUsecase.RejectAutoCampaignSuggestion(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reject auto campaign suggestion: %v", err)
	}

	return &sharepb.SubmitResponse{
		Success: true,
		Message: "Auto campaign suggestion rejected successfully",
	}, nil
}

// GenerateAutoCampaignSuggestions implements the GenerateAutoCampaignSuggestions gRPC method
func (h *MarketingGrpcHandler) GenerateAutoCampaignSuggestions(ctx context.Context, req *pb_marketing.GenerateAutoCampaignSuggestionsRequest) (*pb_marketing.AutoCampaignSuggestionResponse, error) {
	// Convert gRPC request to DTO
	requestDTO := &dto.AutoCampaignSuggestionRequestDTO{
		ProductID:         req.ProductId,
		ProductName:       req.ProductName,
		ProductType:       req.ProductType,
		ProductPrice:      getFloat64Ptr(req.ProductPrice),
		ProductLocation:   req.ProductLocation,
		UserBudget:        getFloat64Ptr(req.UserBudget),
		PreferredDuration: getIntPtr(req.PreferredDuration),
	}

	// Call usecase
	response, err := h.advertisingSettingsUsecase.GenerateAutoCampaignSuggestions(ctx, requestDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate auto campaign suggestions: %v", err)
	}

	// Convert to gRPC response
	grpcSuggestions := make([]*pb_marketing.AutoCampaignSuggestion, len(response.Suggestions))
	for i, suggestion := range response.Suggestions {
		grpcSuggestions[i] = &pb_marketing.AutoCampaignSuggestion{
			ProductId:              suggestion.ProductID,
			ProductName:            suggestion.ProductName,
			ProductType:            suggestion.ProductType,
			SuggestedCampaignType: suggestion.SuggestedCampaignType,
			SuggestedPackageId:    suggestion.SuggestedPackageID,
			SuggestedPackageName:  suggestion.SuggestedPackageName,
			SuggestedBudget:       suggestion.SuggestedBudget,
			SuggestedDuration:     uint32(suggestion.SuggestedDuration),
			Reason:                suggestion.Reason,
			Confidence:            suggestion.Confidence,
			EstimatedReach:        getUint32Value(suggestion.EstimatedReach),
			EstimatedClicks:       getUint32Value(suggestion.EstimatedClicks),
		}
	}

	grpcBestSuggestion := &pb_marketing.AutoCampaignSuggestion{}
	if response.BestSuggestion != nil {
		grpcBestSuggestion = &pb_marketing.AutoCampaignSuggestion{
			ProductId:              response.BestSuggestion.ProductID,
			ProductName:            response.BestSuggestion.ProductName,
			ProductType:            response.BestSuggestion.ProductType,
			SuggestedCampaignType: response.BestSuggestion.SuggestedCampaignType,
			SuggestedPackageId:    response.BestSuggestion.SuggestedPackageID,
			SuggestedPackageName:  response.BestSuggestion.SuggestedPackageName,
			SuggestedBudget:       response.BestSuggestion.SuggestedBudget,
			SuggestedDuration:     uint32(response.BestSuggestion.SuggestedDuration),
			Reason:                response.BestSuggestion.Reason,
			Confidence:            response.BestSuggestion.Confidence,
			EstimatedReach:        getUint32Value(response.BestSuggestion.EstimatedReach),
			EstimatedClicks:       getUint32Value(response.BestSuggestion.EstimatedClicks),
		}
	}

	grpcSettingsUsed := &pb_marketing.AdvertisingSettings{
		DefaultCampaignType: response.SettingsUsed.DefaultCampaignType,
		DefaultBudget:       response.SettingsUsed.DefaultBudget,
		DefaultDuration:     uint32(response.SettingsUsed.DefaultDuration),
	}

	return &pb_marketing.AutoCampaignSuggestionResponse{
		Data: &pb_marketing.AutoCampaignSuggestionData{
			ProductId:        response.ProductID,
			Suggestions:      grpcSuggestions,
			TotalSuggestions: uint32(response.TotalSuggestions),
			BestSuggestion:   grpcBestSuggestion,
			SettingsUsed:     grpcSettingsUsed,
		},
	}, nil
}

// ValidateAdvertisingSettings implements the ValidateAdvertisingSettings gRPC method
func (h *MarketingGrpcHandler) ValidateAdvertisingSettings(ctx context.Context, req *pb_marketing.ValidateAdvertisingSettingsRequest) (*pb_marketing.SettingsValidationResponse, error) {
	// Convert gRPC request to DTO
	createDTO := &dto.AdvertisingSettingsCreateDTO{
		OrganizationID:         req.OrganizationId,
		UserID:                 req.UserId,
		DefaultCampaignType:    req.DefaultCampaignType,
		DefaultPackageID:       getUint64Ptr(req.DefaultPackageId),
		DefaultBudget:          req.DefaultBudget,
		DefaultDuration:        int(req.DefaultDuration),
		AutoSuggestEnabled:     req.AutoSuggestEnabled,
		AutoRunEnabled:         req.AutoRunEnabled,
		AutoRunBudget:          req.AutoRunBudget,
		AutoRunPackageID:       getUint64Ptr(req.AutoRunPackageId),
		MonthlyBudgetLimit:     req.MonthlyBudgetLimit,
		AlertOnBudgetExceed:    req.AlertOnBudgetExceed,
		AlertEmail:             req.AlertEmail,
		AlertPhone:             req.AlertPhone,
		ApplyToAllOrganization: req.ApplyToAllOrganization,
	}

	// Call usecase
	validation, err := h.advertisingSettingsUsecase.ValidateSettings(ctx, createDTO)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate advertising settings: %v", err)
	}

	// Convert to gRPC response
	return &pb_marketing.SettingsValidationResponse{
		Data: &pb_marketing.SettingsValidation{
			IsValid:            validation.IsValid,
			Errors:             validation.Errors,
			Warnings:           validation.Warnings,
			Suggestions:        validation.Suggestions,
			CanApplyToAll:      validation.CanApplyToAll,
			RequiresPermission: validation.RequiresPermission,
			PermissionLevel:    validation.PermissionLevel,
		},
	}, nil
}

// Helper methods for advertising settings

func (h *MarketingGrpcHandler) advertisingSettingsToGrpcResponse(settings *dto.AdvertisingSettingsDTO) *pb_marketing.AdvertisingSettingsResponse {
	return &pb_marketing.AdvertisingSettingsResponse{
		Data: h.advertisingSettingsToGrpcSettings(settings),
	}
}

func (h *MarketingGrpcHandler) advertisingSettingsToGrpcSettings(settings *dto.AdvertisingSettingsDTO) *pb_marketing.AdvertisingSettings {
	return &pb_marketing.AdvertisingSettings{
		Id:                     settings.ID,
		OrganizationId:         settings.OrganizationID,
		UserId:                 settings.UserID,
		DefaultCampaignType:    settings.DefaultCampaignType,
		DefaultPackageId:       getUint64Value(settings.DefaultPackageID),
		DefaultPackageName:     settings.DefaultPackageName,
		DefaultBudget:          settings.DefaultBudget,
		DefaultDuration:        uint32(settings.DefaultDuration),
		AutoSuggestEnabled:     settings.AutoSuggestEnabled,
		AutoRunEnabled:         settings.AutoRunEnabled,
		AutoRunBudget:          settings.AutoRunBudget,
		AutoRunPackageId:       getUint64Value(settings.AutoRunPackageID),
		AutoRunPackageName:     settings.AutoRunPackageName,
		MonthlyBudgetLimit:     settings.MonthlyBudgetLimit,
		AlertOnBudgetExceed:    settings.AlertOnBudgetExceed,
		AlertEmail:             settings.AlertEmail,
		AlertPhone:             settings.AlertPhone,
		ApplyToAllOrganization: settings.ApplyToAllOrganization,
		IsActive:               settings.IsActive,
		CreatedBy:              settings.CreatedBy,
		UpdatedBy:              settings.UpdatedBy,
		CreatedAt:              settings.CreatedAt,
		UpdatedAt:              settings.UpdatedAt,
	}
}
*/

// Helper methods to convert between DTOs and gRPC messages

// Helper functions for type conversion
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

	// Handle optional fields
	if campaign.ProductID != nil {
		grpcCampaign.ProductId = *campaign.ProductID
	}
	if campaign.PackageID != nil {
		grpcCampaign.PackageId = *campaign.PackageID
	}
	if campaign.WalletID != nil {
		grpcCampaign.WalletId = *campaign.WalletID
	}

	// Add dates if they exist
	if campaign.StartDate != nil {
		grpcCampaign.StartDate = campaign.StartDate.Format(time.RFC3339)
	}
	if campaign.EndDate != nil {
		grpcCampaign.EndDate = campaign.EndDate.Format(time.RFC3339)
	}

	return grpcCampaign
} 