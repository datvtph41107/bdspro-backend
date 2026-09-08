package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type BusinessDomainHandler struct {
	organizationpb.UnimplementedBusinessDomainServiceServer
	BusinessDomainUsecase     usecase.BusinessDomainUsecase
	BusinessDomainTransformer transformer.BusinessDomainTransformer
	BusinessDomainValidator   validator.BusinessDomainValidator
}

func NewBusinessDomainHandler(
	businessDomainUsecase usecase.BusinessDomainUsecase,
	businessDomainTransformer transformer.BusinessDomainTransformer,
	businessDomainValidator validator.BusinessDomainValidator,
) *BusinessDomainHandler {
	return &BusinessDomainHandler{
		BusinessDomainUsecase:     businessDomainUsecase,
		BusinessDomainTransformer: businessDomainTransformer,
		BusinessDomainValidator:   businessDomainValidator,
	}
}

// @Summary Tạo lĩnh vực hoạt động
// @Description Tạo mới một lĩnh vực hoạt động
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateBusinessDomainRequest true "Thông tin lĩnh vực hoạt động"
// @Security BearerAuth
// @Router /business-domain [post]
func (h *BusinessDomainHandler) CreateBusinessDomain(ctx context.Context, req *organizationpb.CreateBusinessDomainRequest) (*organizationpb.CreateBusinessDomainResponse, error) {
	if err := h.BusinessDomainValidator.ValidateCreateBusinessDomainRequest(req); err != nil {
		return nil, err
	}

	businessDomain := h.BusinessDomainTransformer.CreateBusinessDomainRequestToEntity(req)
	businessDomain, err := h.BusinessDomainUsecase.CreateBusinessDomain(ctx, businessDomain)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntityToCreateBusinessDomainResponse(businessDomain), nil
}

// @Summary Cập nhật lĩnh vực hoạt động
// @Description Cập nhật thông tin lĩnh vực hoạt động
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Param request body organizationpb.UpdateBusinessDomainRequest true "Thông tin lĩnh vực hoạt động"
// @Security BearerAuth
// @Param id path int true "ID của lĩnh vực hoạt động"
// @Router /business-domain/{id} [put]
func (h *BusinessDomainHandler) UpdateBusinessDomain(ctx context.Context, req *organizationpb.UpdateBusinessDomainRequest) (*organizationpb.UpdateBusinessDomainResponse, error) {
	if err := h.BusinessDomainValidator.ValidateUpdateBusinessDomainRequest(req); err != nil {
		return nil, err
	}

	businessDomain := h.BusinessDomainTransformer.UpdateBusinessDomainRequestToEntity(req)
	businessDomain, err := h.BusinessDomainUsecase.UpdateBusinessDomain(ctx, businessDomain)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntityToUpdateBusinessDomainResponse(businessDomain), nil
}

// @Summary Xóa lĩnh vực hoạt động
// @Description Xóa lĩnh vực hoạt động (soft delete)
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Param id path int true "ID của lĩnh vực hoạt động"
// @Security BearerAuth
// @Router /business-domain/{id} [delete]
func (h *BusinessDomainHandler) DeleteBusinessDomain(ctx context.Context, req *organizationpb.DeleteBusinessDomainRequest) (*organizationpb.DeleteBusinessDomainResponse, error) {
	if err := h.BusinessDomainValidator.ValidateDeleteBusinessDomainRequest(req); err != nil {
		return nil, err
	}

	businessDomain, err := h.BusinessDomainUsecase.GetBusinessDomain(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	err = h.BusinessDomainUsecase.DeleteBusinessDomain(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntityToDeleteBusinessDomainResponse(businessDomain), nil
}

// @Summary Lấy chi tiết lĩnh vực hoạt động
// @Description Lấy thông tin chi tiết lĩnh vực hoạt động theo ID
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Param id path int true "ID của lĩnh vực hoạt động"
// @Security BearerAuth
// @Router /business-domain/{id} [get]
func (h *BusinessDomainHandler) GetBusinessDomain(ctx context.Context, req *organizationpb.GetBusinessDomainRequest) (*organizationpb.GetBusinessDomainResponse, error) {
	if err := h.BusinessDomainValidator.ValidateGetBusinessDomainRequest(req); err != nil {
		return nil, err
	}

	businessDomain, err := h.BusinessDomainUsecase.GetBusinessDomain(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntityToGetBusinessDomainResponse(businessDomain), nil
}

// @Summary Lấy tất cả lĩnh vực hoạt động
// @Description Lấy danh sách tất cả lĩnh vực hoạt động
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /business-domain [get]
func (h *BusinessDomainHandler) GetAllBusinessDomains(ctx context.Context, req *organizationpb.GetAllBusinessDomainsRequest) (*organizationpb.GetAllBusinessDomainsResponse, error) {
	businessDomains, err := h.BusinessDomainUsecase.GetAllBusinessDomains(ctx)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntitiesToGetAllBusinessDomainsResponse(businessDomains), nil
}

// @Summary Lấy lĩnh vực hoạt động đang active
// @Description Lấy danh sách các lĩnh vực hoạt động đang active
// @Tags Lĩnh vực hoạt động
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /business-domain/active [get]
func (h *BusinessDomainHandler) GetActiveBusinessDomains(ctx context.Context, req *organizationpb.GetActiveBusinessDomainsRequest) (*organizationpb.GetActiveBusinessDomainsResponse, error) {
	businessDomains, err := h.BusinessDomainUsecase.GetActiveBusinessDomains(ctx)
	if err != nil {
		return nil, err
	}

	return h.BusinessDomainTransformer.EntitiesToGetActiveBusinessDomainsResponse(businessDomains), nil
}
