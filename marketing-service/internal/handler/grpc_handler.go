package handler

import (
	"context"
	"time"

	"marketing/infrastructure/client"
	"marketing/internal/domain/model"
	"marketing/internal/dto"
	"marketing/internal/usecase"

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
		Type:           model.CampaignType(req.Type),
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
		Type:           model.CampaignType(req.Type),
		ProductID:      &req.ProductId,
		ProductName:    req.ProductName,
		ProductType:    req.ProductType,
		PackageID:      &req.PackageId,
		PackageName:    req.PackageName,
		Budget:         req.Budget,
		Duration:       int(req.Duration),
		TargetAudience: req.TargetAudience,
		TargetLocation: req.TargetLocation,
		Status:         model.CampaignStatus(req.Status),
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
		Status: model.CampaignStatus(req.Status),
		Type:   model.CampaignType(req.Type),
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
		Type:      model.CampaignType(req.Type),
		Status:    model.CampaignStatus(req.Status),
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
	err := h.campaignUsecase.UpdateCampaignStatus(ctx, req.Id, model.CampaignStatus(req.Status))
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
	return nil, status.Errorf(codes.Unimplemented, "method CreatePackage not implemented")
}

func (h *MarketingGrpcHandler) GetPackage(ctx context.Context, req *pb_marketing.GetPackageRequest) (*pb_marketing.PackageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackage not implemented")
}

func (h *MarketingGrpcHandler) UpdatePackage(ctx context.Context, req *pb_marketing.UpdatePackageRequest) (*pb_marketing.PackageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePackage not implemented")
}

func (h *MarketingGrpcHandler) DeletePackage(ctx context.Context, req *pb_marketing.DeletePackageRequest) (*pb_marketing.DeletePackageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePackage not implemented")
}

func (h *MarketingGrpcHandler) ListPackages(ctx context.Context, req *pb_marketing.ListPackagesRequest) (*pb_marketing.ListPackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPackages not implemented")
}

func (h *MarketingGrpcHandler) CreatePayment(ctx context.Context, req *pb_marketing.CreatePaymentRequest) (*pb_marketing.PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePayment not implemented")
}

func (h *MarketingGrpcHandler) GetPayment(ctx context.Context, req *pb_marketing.GetPaymentRequest) (*pb_marketing.PaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPayment not implemented")
}

func (h *MarketingGrpcHandler) ListPayments(ctx context.Context, req *pb_marketing.ListPaymentsRequest) (*pb_marketing.ListPaymentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPayments not implemented")
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

func (h *MarketingGrpcHandler) CreateCampaignDetail(ctx context.Context, req *pb_marketing.CreateCampaignDetailRequest) (*pb_marketing.CampaignDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCampaignDetail not implemented")
}

func (h *MarketingGrpcHandler) GetCampaignDetail(ctx context.Context, req *pb_marketing.GetCampaignDetailRequest) (*pb_marketing.CampaignDetailResponse, error) {
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

func (h *MarketingGrpcHandler) UpdateCampaignDetail(ctx context.Context, req *pb_marketing.UpdateCampaignDetailRequest) (*pb_marketing.CampaignDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCampaignDetail not implemented")
}

func (h *MarketingGrpcHandler) DeleteCampaignDetail(ctx context.Context, req *pb_marketing.DeleteCampaignDetailRequest) (*pb_marketing.DeleteCampaignDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCampaignDetail not implemented")
}

func (h *MarketingGrpcHandler) ListCampaignDetails(ctx context.Context, req *pb_marketing.ListCampaignDetailsRequest) (*pb_marketing.ListCampaignDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCampaignDetails not implemented")
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
	return nil, status.Errorf(codes.Unimplemented, "method CreateBudget not implemented")
}

func (h *MarketingGrpcHandler) GetBudget(ctx context.Context, req *pb_marketing.GetBudgetRequest) (*pb_marketing.BudgetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBudget not implemented")
}

func (h *MarketingGrpcHandler) UpdateBudget(ctx context.Context, req *pb_marketing.UpdateBudgetRequest) (*pb_marketing.BudgetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBudget not implemented")
}

func (h *MarketingGrpcHandler) DeleteBudget(ctx context.Context, req *pb_marketing.DeleteBudgetRequest) (*pb_marketing.DeleteBudgetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBudget not implemented")
}

func (h *MarketingGrpcHandler) ListBudgets(ctx context.Context, req *pb_marketing.ListBudgetsRequest) (*pb_marketing.ListBudgetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBudgets not implemented")
}

func (h *MarketingGrpcHandler) CreateAdvertisingSettings(ctx context.Context, req *pb_marketing.CreateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAdvertisingSettings not implemented")
}

func (h *MarketingGrpcHandler) GetAdvertisingSettings(ctx context.Context, req *pb_marketing.GetAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdvertisingSettings not implemented")
}

func (h *MarketingGrpcHandler) UpdateAdvertisingSettings(ctx context.Context, req *pb_marketing.UpdateAdvertisingSettingsRequest) (*pb_marketing.AdvertisingSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAdvertisingSettings not implemented")
}

func (h *MarketingGrpcHandler) DeleteAdvertisingSettings(ctx context.Context, req *pb_marketing.DeleteAdvertisingSettingsRequest) (*pb_marketing.DeleteAdvertisingSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAdvertisingSettings not implemented")
}

func (h *MarketingGrpcHandler) ListAdvertisingSettings(ctx context.Context, req *pb_marketing.ListAdvertisingSettingsRequest) (*pb_marketing.ListAdvertisingSettingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAdvertisingSettings not implemented")
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
