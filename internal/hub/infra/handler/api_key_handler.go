package handler

import (
	"context"

	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
	sharepb "pb/types/shared"
)

// ApiKeyHandler triển khai ApiKeyService được định nghĩa trong protobuf.
type ApiKeyHandler struct {
	hubpb.UnimplementedApiKeyServiceServer
	usecase _usecase.IApiKeyUsecase
	mapper  *mapper.ApiKeyMapper
}

func NewApiKeyHandler(usecase _usecase.IApiKeyUsecase, mapper *mapper.ApiKeyMapper) *ApiKeyHandler {
	return &ApiKeyHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// @Summary Tạo API key
// @Description Tạo mới một API key sử dụng cho các dịch vụ nội bộ
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param body body hubpb.ApiKeyProto true "Thông tin API key"
// @Success 200 {object} hubpb.ApiKeyProto
// @Router /v2/hub/admin/api-keys [post]
func (h *ApiKeyHandler) CreateApiKey(ctx context.Context, req *hubpb.ApiKeyProto) (*hubpb.ApiKeyProto, error) {
	entity := h.mapper.ProtoToEntity(req)
	if entity != nil {
		entity.ID = 0
	}

	created, err := h.usecase.CreateApiKey(ctx, entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(created), nil
}

// @Summary Danh sách API key
// @Description Lấy danh sách tất cả API key
// @Tags ApiKey
// @Accept json
// @Produce json
// @Success 200 {object} hubpb.ApiKeyResponse
// @Router /v2/hub/admin/api-keys [get]
func (h *ApiKeyHandler) GetApiKey(ctx context.Context, req *hubpb.ApiKeyProto) (*hubpb.ApiKeyResponse, error) {
	apiKeys, err := h.usecase.GetApiKeys(ctx)
	if err != nil {
		return nil, err
	}

	data := h.mapper.EntitiesToProtos(apiKeys)

	return &hubpb.ApiKeyResponse{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}

// @Summary Xóa API key
// @Description Xóa (soft delete) một API key theo ID
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param id path uint64 true "ID của API key"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/hub/admin/api-keys/{id} [delete]
func (h *ApiKeyHandler) DeleteApiKey(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.usecase.DeleteApiKey(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.GetId(),
		Message: "Xóa API key thành công",
	}, nil
}

// @Summary Kiểm tra API key
// @Description Kiểm tra API key có tồn tại và còn hiệu lực không
// @Tags ApiKey
// @Accept json
// @Produce json
// @Param body body hubpb.VerifyApiKeyRequest true "Thông tin API key cần kiểm tra"
// @Success 200 {object} hubpb.VerifyApiKeyResponse
// @Router /v2/hub/admin/api-keys/verify [post]
func (h *ApiKeyHandler) VerifyApiKey(ctx context.Context, req *hubpb.VerifyApiKeyRequest) (*hubpb.VerifyApiKeyResponse, error) {
	entity, errDTO := h.usecase.VerifyApiKey(ctx, req.GetApiKey())
	if errDTO != nil {
		if errDTO.Code >= 500 {
			return nil, errDTO
		}
		return &hubpb.VerifyApiKeyResponse{
			Valid:   false,
			Message: errDTO.Message,
		}, nil
	}

	return &hubpb.VerifyApiKeyResponse{
		Valid:   true,
		Data:    h.mapper.EntityToProto(entity),
		Message: "API key hợp lệ",
	}, nil
}
