package mapper

import (
	_utils "common/utils"
	"context"
	"feedback/infra/client"
	"feedback/internal/domain"
	"feedback/internal/dto"
	feedbackpb "pb/types/feedback"
)

type RateMapper struct {
	userClient *client.UserClient
}

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

	// if rate.CreatedUser != nil {
	// 	response.CreatedUser = &dto.UserItem{
	// 		ProfileID: rate.CreatedUser.ProfileID,
	// 		FullName:  rate.CreatedUser.FullName,
	// 		Avatar:    rate.CreatedUser.Avatar,
	// 	}
	// }

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

func (m *RateMapper) RateItemToRateResponse(rate *domain.Rate) *dto.RateResponse {
	if rate == nil {
		return nil
	}

	response := &dto.RateResponse{
		ID:        rate.ID,
		Anonymous: rate.Anonymous,
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

func (m *RateMapper) RateCreateRequestToRate(req *dto.RateCreateRequest) *domain.Rate {
	if req == nil {
		return nil
	}

	rate := &domain.Rate{
		Anonymous: req.Anonymous,
		OwnerID:   req.OwnerID,
		OwnerOf:   req.OwnerOf,
		Score:     req.Score,
		Comment:   req.Comment,
		ParentID:  req.ParentID,
	}

	if len(req.Attachs) > 0 {
		for _, attachURL := range req.Attachs {
			rate.Attachs = append(rate.Attachs, domain.AttachEntity{
				FileURL:  attachURL,
				FileType: "image", // Default type
				FileName: "attachment",
			})
		}
	}

	return rate
}

func (m *RateMapper) RateUpdateRequestToRate(req *dto.RateUpdateRequest) *domain.Rate {
	if req == nil {
		return nil
	}

	rate := &domain.Rate{
		Score:   req.Score,
		Comment: req.Comment,
	}

	if len(req.Attachs) > 0 {
		for _, attachURL := range req.Attachs {
			rate.Attachs = append(rate.Attachs, domain.AttachEntity{
				FileURL:  attachURL,
				FileType: "image", // Default type
				FileName: "attachment",
			})
		}
	}

	return rate
}

func (m *RateMapper) RateStatsToRateStatsResponse(stats domain.RateStats) dto.RateStatsResponse {
	return dto.RateStatsResponse{
		TotalReviews: stats.TotalReviews,
		AverageScore: stats.AverageScore,
		Star1Count:   stats.Star1Count,
		Star2Count:   stats.Star2Count,
		Star3Count:   stats.Star3Count,
		Star4Count:   stats.Star4Count,
		Star5Count:   stats.Star5Count,
	}
}

func (m *RateMapper) MapRateToPbs(ctx context.Context, rates []domain.Rate) []*feedbackpb.RateResponse {
	ratesPb := make([]*feedbackpb.RateResponse, len(rates))
	for i, rate := range rates {
		ratesPb[i] = m.MapRateItemToPb(ctx, &rate)
	}

	m.MapUserToPb(ctx, ratesPb)
	return ratesPb
}

func (m *RateMapper) MapRateToPb(ctx context.Context, rate *domain.Rate) *feedbackpb.RateResponse {
	if rate == nil {
		return nil
	}

	response := &feedbackpb.RateResponse{
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
		// Không cần gọi user service ở đây vì sẽ được xử lý trong MapUserToPb
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, &feedbackpb.AttachEntity{
				Id:       attach.ID,
				FileUrl:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}

func (m *RateMapper) MapRateItemToPb(ctx context.Context, rate *domain.Rate) *feedbackpb.RateResponse {
	if rate == nil {
		return nil
	}

	response := &feedbackpb.RateResponse{
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
		// Không cần gọi user service ở đây vì sẽ được xử lý trong MapUserToPb
	}

	if len(rate.Attachs) > 0 {
		for _, attach := range rate.Attachs {
			response.Attachs = append(response.Attachs, &feedbackpb.AttachEntity{
				Id:       attach.ID,
				FileUrl:  attach.FileURL,
				FileName: attach.FileName,
				FileType: attach.FileType,
			})
		}
	}

	return response
}

func (m *RateMapper) MapUserToPb(ctx context.Context, rates []*feedbackpb.RateResponse) {
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
func (m *RateMapper) MapUserToPbSingle(ctx context.Context, rate *feedbackpb.RateResponse) {
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
