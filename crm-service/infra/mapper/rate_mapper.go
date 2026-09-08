package mapper

import (
	_utils "common/utils"
	"context"
	"crm/infra/client"
	"crm/internal/domain"
	"crm/internal/dto"
	crmpb "pb/types/crm"
)

type RateMapper struct {
	userClient *client.UserClient
}

// @bind: crm/infra/mapper.RateMapper
func NewRateMapper(userClient *client.UserClient) *RateMapper {
	return &RateMapper{
		userClient: userClient,
	}
}

func (m *RateMapper) RateToRateResponse(rate *domain.Rate) *dto.RateResponse {
	if rate == nil {
		return nil
	}

	response := &dto.RateResponse{
		ID:        rate.ID,
		Anonymous: rate.Anonymous,
		OwnerID:   rate.OwnerID,
		OwnerOf:   rate.OwnerOf,
		Score:     rate.Score,
		Comment:   rate.Comment,
		ParentID:  rate.ParentID,
		CreatedBy: rate.CreatedBy,
		CreatedAt: rate.CreatedAt,
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, dto.AttachEntity{
				ID:       attach.ID,
				FileURL:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}

func (m *RateMapper) MapRateToPbs(ctx context.Context, rates []domain.Rate) []*crmpb.RateResponse {
	ratesPb := make([]*crmpb.RateResponse, len(rates))
	for i, rate := range rates {
		ratesPb[i] = m.MapRateItemToPb(ctx, &rate)
	}

	m.MapUserToPb(ctx, ratesPb)
	return ratesPb
}

func (m *RateMapper) MapRateToPb(ctx context.Context, rate *domain.Rate) *crmpb.RateResponse {
	if rate == nil {
		return nil
	}

	response := &crmpb.RateResponse{
		Id:        rate.ID,
		Anonymous: rate.Anonymous,
		OwnerId:   rate.OwnerID,
		OwnerOf:   rate.OwnerOf,
		Score:     uint32(rate.Score),
		Comment:   rate.Comment,
		CreatedAt: _utils.FormatTimeToString(rate.CreatedAt),
	}

	if rate.ParentID != nil {
		response.ParentId = rate.ParentID
	}

	if rate.CreatedBy != nil {
		response.CreatedBy = rate.CreatedBy
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, &crmpb.AttachEntity{
				Id:       attach.ID,
				FileUrl:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}

func (m *RateMapper) MapRateItemToPb(ctx context.Context, rate *domain.Rate) *crmpb.RateResponse {
	if rate == nil {
		return nil
	}

	response := &crmpb.RateResponse{
		Id:        rate.ID,
		Anonymous: rate.Anonymous,
		Score:     uint32(rate.Score),
		Comment:   rate.Comment,
		CreatedAt: _utils.FormatTimeToString(rate.CreatedAt),
	}

	if rate.ParentID != nil {
		response.ParentId = rate.ParentID
	}

	if rate.CreatedBy != nil {
		response.CreatedBy = rate.CreatedBy
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, &crmpb.AttachEntity{
				Id:       attach.ID,
				FileUrl:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}

func (m *RateMapper) MapUserToPb(ctx context.Context, rates []*crmpb.RateResponse) {
	profileIDSet := make(map[uint64]struct{})
	for _, rate := range rates {
		// Chỉ lấy thông tin user khi không anonymous
		if rate.CreatedBy != nil && !rate.Anonymous {
			profileIDSet[*rate.CreatedBy] = struct{}{}
		}
	}

	profileMap, err := m.userClient.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, rate := range rates {
		// Chỉ map user khi không anonymous
		if rate.CreatedBy != nil && !rate.Anonymous {
			if profile, ok := profileMap[*rate.CreatedBy]; ok {
				rate.CreatedUser = profile
			}
		}
	}
}

// MapUserToPbSingle map user profile cho single rate response
func (m *RateMapper) MapUserToPbSingle(ctx context.Context, rate *crmpb.RateResponse) {
	if rate == nil {
		return
	}

	// Chỉ lấy thông tin user khi không anonymous
	if rate.CreatedBy != nil && !rate.Anonymous {
		profileIDSet := map[uint64]struct{}{
			*rate.CreatedBy: {},
		}

		profileMap, err := m.userClient.GetMapByIDs(ctx, profileIDSet)
		if err != nil {
			return
		}

		if profile, ok := profileMap[*rate.CreatedBy]; ok {
			rate.CreatedUser = profile
		}
	}
}