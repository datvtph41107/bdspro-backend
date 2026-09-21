package admin_handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal"
	"bdspro/internal/dto"
	admin_usecases "bdspro/internal/usecases/admin"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminPropertyTypeHandler struct {
	bdspropb.UnimplementedAdminPropertyTypeServiceServer
	AdminPropertyTypeUsecase *admin_usecases.AdminPropertyTypeUsecase
	PropertyTypeMapper       *mapper.PropertyTypeMapper
}

func NewAdminPropertyTypeHandler(
	uc *admin_usecases.AdminPropertyTypeUsecase,
	propertyTypeMapper *mapper.PropertyTypeMapper,
) *AdminPropertyTypeHandler {
	return &AdminPropertyTypeHandler{
		AdminPropertyTypeUsecase: uc,
		PropertyTypeMapper:       propertyTypeMapper,
	}
}

func (h *AdminPropertyTypeHandler) GetList(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.PropertyTypeListResponse, error) {
	result, total, err := h.AdminPropertyTypeUsecase.GetList(ctx,
		&dto.TextSearchRequest{
			Pagable: _dto.Pagable{
				Page: req.Page,
				Size: req.Size,
			},
			Text: req.Text,
		})
	if err != nil {
		return nil, err
	}

	propertyTypes := h.PropertyTypeMapper.MapPropertyTypePbList(result)
	return &bdspropb.PropertyTypeListResponse{
		Data:          propertyTypes,
		TotalElements: total,
	}, nil
}

func (h *AdminPropertyTypeHandler) Create(ctx context.Context, req *bdspropb.PropertyType) (*bdspropb.PropertyType, error) {
	propertyType := h.PropertyTypeMapper.PropertyTypePbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(service.PropertyNameRequired)
	}
	propertyType, err := h.AdminPropertyTypeUsecase.Create(ctx, propertyType)
	if err != nil {
		return nil, err
	}
	return h.PropertyTypeMapper.MapPropertyTypePb(propertyType), nil
}

func (h *AdminPropertyTypeHandler) Update(ctx context.Context, req *bdspropb.PropertyType) (*bdspropb.PropertyType, error) {
	propertyType := h.PropertyTypeMapper.PropertyTypePbToDomain(req)
	if req.Name == "" {
		return nil, _errors.ReturnError(service.PropertyTypeNameRequired)
	}
	propertyType, err := h.AdminPropertyTypeUsecase.Update(ctx, req.Id, propertyType)
	if err != nil {
		return nil, err
	}
	return h.PropertyTypeMapper.MapPropertyTypePb(propertyType), nil
}

func (h *AdminPropertyTypeHandler) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	_, err := h.AdminPropertyTypeUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "success",
	}, nil
}
