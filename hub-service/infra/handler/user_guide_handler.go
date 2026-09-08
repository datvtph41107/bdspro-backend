package handler

import (
	"context"
	"time"

	_dto "common/domain/dto"
	_utils "common/utils"
	"hub/infra/mapper"
	"hub/internal/domain"
	"hub/internal/repo"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type UserGuideHandler struct {
	hubpb.UnimplementedUserGuideServiceServer
	userGuideUsecase _usecase.IUserGuideUsecase
	userGuideMapper  *mapper.UserGuideMapper
	stepRepo         repo.IUserGuideStepRepo
	syncProvider     *_utils.SyncUtil
}

func NewUserGuideHandler(
	userGuideUsecase _usecase.IUserGuideUsecase,
	userGuideMapper *mapper.UserGuideMapper,
	stepRepo repo.IUserGuideStepRepo,
	syncProvider *_utils.SyncUtil,
) *UserGuideHandler {
	return &UserGuideHandler{
		userGuideUsecase: userGuideUsecase,
		userGuideMapper:  userGuideMapper,
		stepRepo:         stepRepo,
		syncProvider:     syncProvider,
	}
}

// @Summary Tạo user guide mới
// @Description Tạo một user guide mới với title, description, groupKey và steps
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param body body hubpb.CreateUserGuideRequest true "Thông tin user guide"
// @Success 200 {object} hubpb.CreateUserGuideResponse
// @Router /v2/hub/user-guides [post]
func (h *UserGuideHandler) CreateUserGuide(ctx context.Context, req *hubpb.CreateUserGuideRequest) (*hubpb.CreateUserGuideResponse, error) {
	entity := h.userGuideMapper.CreateRequestToEntity(req)

	result, err := h.userGuideUsecase.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	// Create steps if provided
	if len(req.Steps) > 0 {
		steps := make([]*domain.UserGuideStepEntity, 0, len(req.Steps))
		for _, stepInput := range req.Steps {
			step := h.userGuideMapper.StepInputToEntity(stepInput, result.ID)
			steps = append(steps, step)
		}

		if err := h.stepRepo.CreateBatch(ctx, steps); err != nil {
			return nil, err
		}
	}

	// Get steps for response
	steps, _ := h.stepRepo.GetByUserGuideID(ctx, result.ID)

	// Sync: khi tạo user guide có key thì cập nhật timestamp chung trong Redis
	if result.Key != "" && h.syncProvider != nil {
		ts := int64(0)
		if result.UpdatedAt != nil {
			ts = result.UpdatedAt.UnixMilli()
		} else {
			ts = time.Now().UnixMilli()
		}
		_ = h.syncProvider.PutTimestamp(ctx, _utils.SyncKeyUserGuideByKey, ts)
	}

	return &hubpb.CreateUserGuideResponse{
		Data: h.userGuideMapper.EntityToDetailProto(result, steps),
	}, nil
}

// @Summary Cập nhật user guide
// @Description Cập nhật thông tin user guide theo ID và steps
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param id path uint64 true "User Guide ID"
// @Param body body hubpb.UpdateUserGuideRequest true "Thông tin user guide cần cập nhật"
// @Success 200 {object} hubpb.UpdateUserGuideResponse
// @Router /v2/hub/user-guides/{id} [put]
func (h *UserGuideHandler) UpdateUserGuide(ctx context.Context, req *hubpb.UpdateUserGuideRequest) (*hubpb.UpdateUserGuideResponse, error) {
	// Get existing entity
	existing, err := h.userGuideUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Update entity with request data
	h.userGuideMapper.UpdateRequestToEntity(existing, req)

	// Save updated entity
	result, err := h.userGuideUsecase.Update(ctx, req.GetId(), existing)
	if err != nil {
		return nil, err
	}

	// Update steps if provided
	if len(req.Steps) > 0 {
		// Delete old steps
		h.stepRepo.DeleteByUserGuideID(ctx, result.ID)

		// Create new steps
		steps := make([]*domain.UserGuideStepEntity, 0, len(req.Steps))
		for _, stepInput := range req.Steps {
			step := h.userGuideMapper.StepInputToEntity(stepInput, result.ID)
			steps = append(steps, step)
		}

		if err := h.stepRepo.CreateBatch(ctx, steps); err != nil {
			return nil, err
		}
	}

	// Get steps for response
	steps, _ := h.stepRepo.GetByUserGuideID(ctx, result.ID)

	// Sync: khi cập nhật user guide có key thì cập nhật timestamp chung trong Redis
	if result.Key != "" && h.syncProvider != nil {
		ts := int64(0)
		if result.UpdatedAt != nil {
			ts = result.UpdatedAt.UnixMilli()
		} else {
			ts = time.Now().UnixMilli()
		}
		_ = h.syncProvider.PutTimestamp(ctx, _utils.SyncKeyUserGuideByKey, ts)
	}

	return &hubpb.UpdateUserGuideResponse{
		Data: h.userGuideMapper.EntityToDetailProto(result, steps),
	}, nil
}

// @Summary Xóa user guide
// @Description Xóa user guide theo ID (soft delete)
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param id path uint64 true "User Guide ID"
// @Success 200 {object} hubpb.DeleteUserGuideResponse
// @Router /v2/hub/user-guides/{id} [delete]
func (h *UserGuideHandler) DeleteUserGuide(ctx context.Context, req *hubpb.DeleteUserGuideRequest) (*hubpb.DeleteUserGuideResponse, error) {
	_, err := h.userGuideUsecase.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &hubpb.DeleteUserGuideResponse{
		Success: true,
	}, nil
}

// @Summary Lấy chi tiết user guide
// @Description Lấy thông tin chi tiết user guide theo ID với steps
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param id path uint64 true "User Guide ID"
// @Success 200 {object} hubpb.GetUserGuideResponse
// @Router /v2/hub/user-guides/{id} [get]
func (h *UserGuideHandler) GetUserGuide(ctx context.Context, req *hubpb.GetUserGuideRequest) (*hubpb.GetUserGuideResponse, error) {
	result, err := h.userGuideUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Get steps
	steps, _ := h.stepRepo.GetByUserGuideID(ctx, result.ID)

	return &hubpb.GetUserGuideResponse{
		Data: h.userGuideMapper.EntityToDetailProto(result, steps),
	}, nil
}

// @Summary Lấy danh sách user guides
// @Description Lấy danh sách user guides với tìm kiếm theo title, filter theo groupKey và phân trang
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param title query string false "Tìm kiếm theo title"
// @Param groupKey query string false "Lọc theo groupKey"
// @Param page query int false "Số trang (default: 1)"
// @Param size query int false "Số bản ghi trên trang (default: 20, max: 100)"
// @Success 200 {object} hubpb.GetUserGuideListResponse
// @Router /v2/hub/user-guides [get]
func (h *UserGuideHandler) GetUserGuideList(ctx context.Context, req *hubpb.GetUserGuideListRequest) (*hubpb.GetUserGuideListResponse, error) {
	// Create pagable from request
	pagable := &_dto.Pagable{
		Page: uint32(req.GetPage()),
		Size: uint32(req.GetSize()),
	}

	// Get filtered list
	data, total, err := h.userGuideUsecase.GetListWithFilter(
		ctx,
		req.Title,
		req.GetGroupKey(),
		pagable,
	)
	if err != nil {
		return nil, err
	}

	return &hubpb.GetUserGuideListResponse{
		Data:  h.userGuideMapper.EntitiesToProto(data),
		Total: total,
	}, nil
}

// @Summary Lấy danh sách nhóm user guides
// @Description Lấy danh sách các groupKey và số lượng user guide trong mỗi nhóm
// @Tags UserGuide
// @Accept json
// @Produce json
// @Success 200 {object} hubpb.GetUserGuideGroupsResponse
// @Router /v2/hub/user-guides/groups [get]
func (h *UserGuideHandler) GetUserGuideGroups(ctx context.Context, req *hubpb.GetUserGuideGroupsRequest) (*hubpb.GetUserGuideGroupsResponse, error) {
	groups, err := h.userGuideUsecase.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	return &hubpb.GetUserGuideGroupsResponse{
		Data: h.userGuideMapper.GroupMapsToProto(groups),
	}, nil
}

// @Summary Lấy danh sách user guides đơn giản cho user
// @Description API đơn giản cho user lấy danh sách user guides, chỉ trả về id, title và description (tối đa 80 ký tự). Hỗ trợ search theo text (lowercase, khoảng trống thay bằng %)
// @Tags UserGuide/Simple
// @Accept json
// @Produce json
// @Param text query string false "Search theo title"
// @Param groupKey query string false "Lọc theo groupKey"
// @Param page query int false "Số trang (default: 1)"
// @Param size query int false "Số bản ghi trên trang (default: 20, max: 100)"
// @Success 200 {object} hubpb.GetUserGuideSimpleResponse
// @Router /v2/hub/user-guides/simple [get]
func (h *UserGuideHandler) GetUserGuideSimple(ctx context.Context, req *hubpb.GetUserGuideSimpleRequest) (*hubpb.GetUserGuideSimpleResponse, error) {
	// Create pagable from request
	pagable := &_dto.Pagable{
		Page: uint32(req.GetPage()),
		Size: uint32(req.GetSize()),
	}

	// Get simple list with text search
	data, total, err := h.userGuideUsecase.GetSimpleList(
		ctx,
		req.Text,
		req.GetGroupKey(),
		pagable,
	)
	if err != nil {
		return nil, err
	}

	// Convert to simple proto format
	var simpleItems []*hubpb.UserGuideSimple
	for _, item := range data {
		description := item.Description
		// Truncate description to max 50 characters
		if len(description) > 50 {
			description = description[:50] + "..."
		}

		simpleItems = append(simpleItems, &hubpb.UserGuideSimple{
			Id:    item.ID,
			Title: item.Title,
			// Description: description,
		})
	}

	return &hubpb.GetUserGuideSimpleResponse{
		Data:  simpleItems,
		Total: total,
	}, nil
}

// @Summary Lấy chi tiết user guide theo key
// @Description Lấy thông tin chi tiết user guide theo key với steps
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param key path string true "User Guide Key"
// @Success 200 {object} hubpb.GetUserGuideByKeyResponse
// @Router /v2/hub/user-guides/by-key/{key} [get]
func (h *UserGuideHandler) GetUserGuideByKey(ctx context.Context, req *hubpb.GetUserGuideByKeyRequest) (*hubpb.GetUserGuideByKeyResponse, error) {
	result, err := h.userGuideUsecase.GetByKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	// Get steps
	// steps, _ := h.stepRepo.GetByUserGuideID(ctx, result.ID)

	return &hubpb.GetUserGuideByKeyResponse{
		Data: &hubpb.UserGuideDetail{
			Id:    result.ID,
			Title: result.Title,
			Mode:  hubpb.UserGuideMode(result.Mode),
		},
	}, nil
}

// @Summary Lấy toàn bộ user guide có key (sync theo timestamp)
// @Description Một key Redis chung; client gửi timestamp lần trước, nếu bằng với Redis thì trả unchanged, không trả full data
// @Tags UserGuide
// @Accept json
// @Produce json
// @Param timestamp query int64 false "Timestamp lần sync trước (0 = luôn lấy full)"
// @Success 200 {object} hubpb.GetAllUserGuidesByKeyResponse
// @Router /v2/hub/user-guides/by-key [get]
func (h *UserGuideHandler) GetAllUserGuidesByKey(ctx context.Context, req *hubpb.GetAllUserGuidesByKeyRequest) (*hubpb.GetAllUserGuidesByKeyResponse, error) {
	key := _utils.SyncKeyUserGuideByKey
	clientTs := req.GetTimestamp()

	if h.syncProvider != nil && clientTs > 0 {
		h.syncProvider.PutTimeRequest(ctx, key, clientTs)
		updated := h.syncProvider.HasUpdated(ctx, key, clientTs)
		if !updated {
			// Client đã có bản mới, không cần trả full data
			ts := h.syncProvider.GetTimestamp(ctx, key)
			return &hubpb.GetAllUserGuidesByKeyResponse{
				Data:      nil,
				Timestamp: ts,
				Unchanged: true,
			}, nil
		}
	}

	list, err := h.userGuideUsecase.GetAllWithKey(ctx)
	if err != nil {
		return nil, err
	}

	var maxTs int64
	for _, e := range list {
		if e.UpdatedAt != nil {
			ms := e.UpdatedAt.UnixMilli()
			if ms > maxTs {
				maxTs = ms
			}
		}
	}
	if maxTs == 0 {
		maxTs = time.Now().UnixMilli()
	}

	if h.syncProvider != nil {
		_ = h.syncProvider.PutTimestamp(ctx, key, maxTs)
	}

	return &hubpb.GetAllUserGuidesByKeyResponse{
		Data:      h.userGuideMapper.EntitiesToProto(list),
		Timestamp: maxTs,
		Unchanged: false,
	}, nil
}
