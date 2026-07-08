package handler

import (
	"context"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type BlockService struct {
	UC *usecase.BlockUsecase
	crmpb.UnimplementedBlockServiceServer
}

func NewBlockService(uc *usecase.BlockUsecase) *BlockService {
	return &BlockService{
		UC: uc,
	}
}

// @Summary Lấy danh sách người dùng bị chặn
// @Description Lấy danh sách ID của những người mà người dùng hiện tại đã chặn
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /block/list [get]
func (r *BlockService) ListBlockedUsers(c context.Context, req *sharepb.IdRequest) (*crmpb.ListBlockResponse, error) {
	result, err := r.UC.ListBlockedUsers(c)
	return &crmpb.ListBlockResponse{
		Data: result,
	}, err
}

// @Summary Chặn một người dùng
// @Description Chặn một người dùng với profileId cụ thể
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID người dùng cần chặn"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /block/{profileId} [post]
func (r *BlockService) BlockUser(c context.Context, req *crmpb.BlockUserRequest) (*sharepb.SubmitResponse, error) {
	result, err := r.UC.BlockUser(c, req.ProfileId)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Message: "success",
		Id:      result.BlockedId,
	}, err
}

// @Summary Bỏ chặn một người dùng
// @Description Bỏ chặn một người dùng với profileId cụ thể
// @Tags Chặn
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID người dùng cần bỏ chặn"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /block/{profileId} [delete]
func (r *BlockService) UnblockUser(c context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	result, err := r.UC.UnblockUser(c, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Message: "success",
		Id:      result,
	}, err
}