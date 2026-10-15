package controller

import (
	"context"
	"net/http"

	userpb "pb/types/user"
	"userv4/internal/httpapi"
	"userv4/modules/user/dto"
	"userv4/modules/user/services"
)

// KYCController owns KYC HTTP routes under /v2/user.
type KYCController struct {
	svc services.KYCService
}

func NewKYCController(svc services.KYCService) *KYCController {
	return &KYCController{svc: svc}
}

func (c *KYCController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/user/kyc/submit", c.submitKYC)
	httpapi.GET(mux, "/v2/user/admin/kycs", c.listKYCs)
	httpapi.POST(mux, "/v2/user/admin/kyc/approve/{id}", c.approveKYC)
	httpapi.POST(mux, "/v2/user/admin/kyc/reject/{id}", c.rejectKYC)
}

func (c *KYCController) submitKYC(ctx context.Context, req userpb.SubmitKYCRequest) (*httpapi.Result, error) {
	result, err := c.svc.SubmitKYC(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *KYCController) listKYCs(ctx context.Context, req dto.ListKYCsQuery) (*httpapi.Result, error) {
	result, err := c.svc.ListKYCs(ctx, &userpb.ListKYCsRequest{
		Page:      req.Page,
		Size:      req.Size,
		Search:    req.Search,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *KYCController) approveKYC(ctx context.Context, req dto.ApproveKYCRequest) (*httpapi.Result, error) {
	if err := c.svc.ApproveKYC(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *KYCController) rejectKYC(ctx context.Context, req dto.RejectKYCRequest) (*httpapi.Result, error) {
	if err := c.svc.RejectKYC(ctx, req.ID, req.Reason); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}
