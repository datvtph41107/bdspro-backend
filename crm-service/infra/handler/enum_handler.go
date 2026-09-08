package handler

import (
	"context"
	"crm/internal/usecase"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EnumService struct {
	crmpb.UnimplementedEnumServiceServer
	UC *usecase.EnumUsecase
}

func NewEnumService(uc *usecase.EnumUsecase) *EnumService {
	return &EnumService{
		UC: uc,
	}
}

// GetEnumByName trả về danh sách enum theo tên
// @Summary Lấy danh sách enum theo tên
// @Description Trả về danh sách ItemV3Proto với id và name tương ứng với enum
// @Tags Enum
// @Accept json
// @Produce json
// @Param name path string true "Tên enum (tagContact, statusContact, visibility, priority, ownerOf, step, sourceLead, friendStatus, condition, ruleThen, targetType, reportStatus)"
// @Success 200 {object} sharepb.ItemV3PageProto
// @Router /v2/enums/crm/{name} [get]
func (s *EnumService) GetEnumByName(ctx context.Context, req *crmpb.GetEnumRequest) (*sharepb.ItemV3PageProto, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Tên enum không được để trống")
	}

	items, err := s.UC.GetEnumByName(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.ItemV3PageProto{
		Data:  items,
		Total: uint32(len(items)),
	}, nil
}