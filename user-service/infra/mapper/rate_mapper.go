package mapper

import (
	_dto "common/dto"
	_models "common/models"
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"

	"user/internal/dto"
	user_model "user/internal/models/user"
)

type RateMapper struct{}

func NewRateMapper() *RateMapper {
	return &RateMapper{}
}

func (m *RateMapper) PbToEntityRate(rate *userpb.Rate) *user_model.Rate {
	result := &user_model.Rate{
		BaseEntity: _models.BaseEntity{
			ID: rate.Id,
		},
		TargetID:   rate.TargetId,
		TargetType: uint8(rate.TargetType),
		ParentID:   rate.ParentId,
		Anonymous:  rate.Anonymous,
		Score:      uint8(rate.Score),
		Comment:    rate.Comment,
	}

	for _, attach := range rate.Attachs {
		result.Attachs = append(result.Attachs, m.PbToEntityAttach(attach))
	}

	return result
}

func (m *RateMapper) PbToEntityAttach(attach *userpb.Attach) user_model.AttachEntity {
	return user_model.AttachEntity{
		ID:        attach.Id,
		ImageURL:  attach.ImageUrl,
		ImageType: attach.ImageType,
	}
}

func (m *RateMapper) EntityToPbRateItem(rate *user_model.RateItem) *userpb.Rate {
	result := &userpb.Rate{
		Id:        rate.ID,
		Anonymous: rate.Anonymous,
		Score:     uint64(rate.Score),
		Comment:   rate.Comment,
		ParentId:  rate.ParentID,
		CreatedBy: rate.CreatedBy,
		CreatedAt: _utils.FormatTimeToString(rate.CreatedAt),
	}

	if rate.CreatedUser != nil && !rate.Anonymous {
		result.CreatedUser = &sharepb.ProfileItem{
			Id:           rate.CreatedUser.ProfileID,
			FullName:     rate.CreatedUser.FullName,
			Avatar:       rate.CreatedUser.Avatar,
			TickVerified: rate.CreatedUser.TickVerified,
		}
	}

	for _, attach := range rate.Attachs {
		result.Attachs = append(result.Attachs, &userpb.Attach{
			Id:        attach.ID,
			ImageUrl:  attach.ImageURL,
			ImageType: attach.ImageType,
		})
	}

	return result
}

func (m *RateMapper) EntityToPbRate(rate *user_model.Rate) *userpb.Rate {
	result := &userpb.Rate{
		Id:        rate.ID,
		ParentId:  rate.ParentID,
		Anonymous: rate.Anonymous,
		Score:     uint64(rate.Score),
		Comment:   rate.Comment,
		CreatedBy: rate.CreatedBy,
		CreatedAt: _utils.FormatTimeToString(rate.CreatedAt),
	}

	for _, attach := range rate.Attachs {
		result.Attachs = append(result.Attachs, &userpb.Attach{
			Id:        attach.ID,
			ImageUrl:  attach.ImageURL,
			ImageType: attach.ImageType,
		})
	}

	if rate.CreatedUser != nil && !rate.Anonymous {
		result.CreatedUser = &sharepb.ProfileItem{
			Id:           rate.CreatedUser.ProfileID,
			FullName:     rate.CreatedUser.FullName,
			Avatar:       rate.CreatedUser.Avatar,
			TickVerified: rate.CreatedUser.TickVerified,
		}
	}

	return result
}

func (m *RateMapper) PbToEntityRateSearch(req *userpb.GetRateListRequest) *dto.RateSearchDTO {
	return &dto.RateSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		TargetId:   req.TargetId,
		TargetType: req.TargetType,
		ParentId:   req.ParentId,
		CreatedBy:  req.CreatedBy,
		UpdatedBy:  req.UpdatedBy,
	}
}

func (m *RateMapper) EntityToPbRateStats(stats *user_model.RateStats) *userpb.RateStats {
	return &userpb.RateStats{
		TotalReviews: uint64(stats.TotalReviews),
		AverageScore: float32(stats.AverageScore),
	}
}
