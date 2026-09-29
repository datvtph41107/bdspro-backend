package handler

import (
	_dto "common/domain/dto"
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strings"
	"user/infra/mapper"
	"user/internal/models"
	"user/internal/usecases"
)

type MainAreaHandler struct {
	userpb.UnimplementedMainAreaServiceServer
	areaUsecase usecases.IAdminMainAreaUsecase
	mapper      *mapper.MainAreaMapper
}

func NewMainAreaHandler(
	areaUsecase usecases.IAdminMainAreaUsecase,
	mapper *mapper.MainAreaMapper,
) *MainAreaHandler {
	return &MainAreaHandler{
		areaUsecase: areaUsecase,
		mapper:      mapper,
	}
}

// GetMainAreas lấy danh sách nhãn khu vực hoạt động
// @Summary Lấy danh sách nhãn khu vực hoạt động
// @Description Lấy danh sách nhãn khu vực hoạt động (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.IdRequest true "Thông tin nhãn"
// @Success 200 {object} userpb.MainAreaPage
// @Router /v2/user/admin/main-area [get]
func (h *MainAreaHandler) GetMainAreas(ctx context.Context, req *sharepb.IdRequest) (*userpb.MainAreaPage, error) {
	entities, err := h.areaUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	data := h.mapper.EntitiesToProtos(entities)

	return &userpb.MainAreaPage{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}

// CreateMainArea tạo nhãn khu vực hoạt động mới
// @Summary Tạo nhãn khu vực hoạt động
// @Description Tạo nhãn khu vực hoạt động mới (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.MainArea true "Thông tin nhãn"
// @Success 200 {object} userpb.MainArea
// @Router /v2/user/admin/main-area [post]
func (h *MainAreaHandler) CreateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error) {
	entity := h.mapper.ProtoToEntity(req)

	result, err := h.areaUsecase.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(result), nil
}

// UpdateMainArea cập nhật nhãn khu vực hoạt động
// @Summary Cập nhật nhãn khu vực hoạt động
// @Description Cập nhật nhãn khu vực hoạt động (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID nhãn"
// @Param request body userpb.MainArea true "Thông tin nhãn"
// @Success 200 {object} userpb.MainArea
// @Router /v2/user/admin/main-area/{id} [put]
func (h *MainAreaHandler) UpdateMainArea(ctx context.Context, req *userpb.MainArea) (*userpb.MainArea, error) {
	entity, err := h.areaUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	entity.Name = req.GetName()
	if req.Description != nil {
		entity.Description = req.GetDescription()
	}

	updated, err := h.areaUsecase.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(updated), nil
}

// DeleteMainArea xóa nhãn khu vực hoạt động
// @Summary Xóa nhãn khu vực hoạt động
// @Description Xóa mềm nhãn khu vực hoạt động (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID nhãn"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/user/admin/main-area/{id} [delete]
func (h *MainAreaHandler) DeleteMainArea(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if _, err := h.areaUsecase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.GetId(),
		Message: "Xóa nhãn khu vực hoạt động thành công",
	}, nil
}

// GetMainAreaItems lấy danh sách nhãn khu vực hoạt động
// @Summary Lấy danh sách nhãn khu vực hoạt động
// @Description Lấy danh sách nhãn khu vực hoạt động (User)
// @Tags UserDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request query string false "Từ khóa tìm kiếm"
// @Success 200 {object} sharepb.ItemV3PageProto
// @Router /v2/user/main-area [get]
func (h *MainAreaHandler) GetMainAreaItems(ctx context.Context, req *sharepb.RequestV3Proto) (*sharepb.ItemV3PageProto, error) {
	entities, total, err := h.areaUsecase.Search(ctx, req.GetText(), &_dto.Pagable{
		Page: req.GetPage(),
		Size: req.GetSize(),
	})
	if err != nil {
		return nil, err
	}

	filtered := entities
	searchText := strings.TrimSpace(req.GetText())
	if searchText != "" {
		searchTextLower := strings.ToLower(searchText)
		filtered = make([]models.MainAreaEntity, 0, len(entities))
		for _, entity := range entities {
			nameMatch := strings.Contains(strings.ToLower(entity.Name), searchTextLower)
			descMatch := entity.Description != "" && strings.Contains(strings.ToLower(entity.Description), searchTextLower)
			if nameMatch || descMatch {
				filtered = append(filtered, entity)
			}
		}
	}

	data := h.mapper.EntitiesToItemProtos(filtered)

	return &sharepb.ItemV3PageProto{
		Data:  data,
		Total: total,
	}, nil
}
