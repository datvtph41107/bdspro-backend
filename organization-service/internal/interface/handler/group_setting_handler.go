package handler

import (
	"context"
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type GroupSettingHandler struct {
	grouppb.UnimplementedGroupSettingServiceServer
	GroupSettingUsecase     usecase.GroupSettingUsecase
	GroupSettingTransformer transformer.GroupSettingTransformer
	GroupSettingValidator   validator.GroupSettingValidator
}

func NewGroupSettingHandler(
	groupSettingUsecase usecase.GroupSettingUsecase,
	groupSettingTransformer transformer.GroupSettingTransformer,
	groupSettingValidator validator.GroupSettingValidator,
) *GroupSettingHandler {
	return &GroupSettingHandler{
		GroupSettingUsecase:     groupSettingUsecase,
		GroupSettingTransformer: groupSettingTransformer,
		GroupSettingValidator:   groupSettingValidator,
	}
}

func (h *GroupSettingHandler) CreateGroupSetting(ctx context.Context, req *grouppb.CreateGroupSettingRequest) (*grouppb.CreateGroupSettingResponse, error) {
	if err := h.GroupSettingValidator.ValidateCreateGroupSettingRequest(req); err != nil {
		return nil, err
	}

	setting := h.GroupSettingTransformer.CreateGroupSettingRequestToEntity(req)

	setting, err := h.GroupSettingUsecase.CreateGroupSetting(ctx, setting)
	if err != nil {
		return nil, err
	}

	return h.GroupSettingTransformer.EntityToCreateGroupSettingResponse(setting), nil
}

func (h *GroupSettingHandler) UpdateGroupSetting(ctx context.Context, req *grouppb.UpdateGroupSettingRequest) (*grouppb.UpdateGroupSettingResponse, error) {
	if err := h.GroupSettingValidator.ValidateUpdateGroupSettingRequest(req); err != nil {
		return nil, err
	}

	setting := h.GroupSettingTransformer.UpdateGroupSettingRequestToEntity(req)

	setting, err := h.GroupSettingUsecase.UpdateGroupSetting(ctx, setting)
	if err != nil {
		return nil, err
	}

	return h.GroupSettingTransformer.EntityToUpdateGroupSettingResponse(setting), nil
}

func (h *GroupSettingHandler) DeleteGroupSetting(ctx context.Context, req *grouppb.DeleteGroupSettingRequest) (*grouppb.DeleteGroupSettingResponse, error) {
	if err := h.GroupSettingValidator.ValidateDeleteGroupSettingRequest(req); err != nil {
		return nil, err
	}

	err := h.GroupSettingUsecase.DeleteGroupSetting(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &grouppb.DeleteGroupSettingResponse{Id: req.Id}, nil
}

func (h *GroupSettingHandler) GetGroupSetting(ctx context.Context, req *grouppb.GetGroupSettingRequest) (*grouppb.GetGroupSettingResponse, error) {
	if err := h.GroupSettingValidator.ValidateGetGroupSettingRequest(req); err != nil {
		return nil, err
	}

	setting, err := h.GroupSettingUsecase.GetGroupSettingByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.GroupSettingTransformer.EntityToGetGroupSettingResponse(setting), nil
}

func (h *GroupSettingHandler) GetGroupSettings(ctx context.Context, req *grouppb.GetGroupSettingsRequest) (*grouppb.GetGroupSettingsResponse, error) {
	if err := h.GroupSettingValidator.ValidateGetGroupSettingsRequest(req); err != nil {
		return nil, err
	}

	page :=0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	settings, total, err := h.GroupSettingUsecase.GetGroupSettingsByGroupID(ctx, req.GroupId, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupSettingTransformer.EntityToGetGroupSettingsResponse(settings, total), nil
} 