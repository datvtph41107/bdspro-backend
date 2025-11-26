package mapper

import (
	_utils "common/utils"
	sharepb "pb/types/shared"
	"user/internal/dto"
	"user/internal/models"
)

type KYCMapper struct{}

func NewKYCMapper() *KYCMapper {
	return &KYCMapper{}
}

func (m *KYCMapper) ToProto(entity *models.KYCEntity) *sharepb.KYCV3Proto {
	return m.ToProtoWithOwnerCheck(entity, false)
}

// ToProtoWithOwnerCheck trả về KYC proto, chỉ trả về 4 trường nhạy cảm nếu isOwner = true
func (m *KYCMapper) ToProtoWithOwnerCheck(entity *models.KYCEntity, isOwner bool) *sharepb.KYCV3Proto {
	if entity == nil {
		return nil
	}

	resp := &sharepb.KYCV3Proto{
		Id:           entity.ID,
		ProfileId:    entity.ProfileID,
		FullName:     entity.FullName,
		IdentityCard: entity.IdentityCard,
		FrontImage:   entity.FrontImage,
		BackImage:    entity.BackImage,
		SelfieImage:  entity.SelfieImage,
		Status:       uint32(entity.Status),
		RejectReason: entity.RejectReason,
	}

	// Chỉ trả về 4 trường nhạy cảm nếu là owner request
	if isOwner {
		if entity.IDNumber != "" {
			resp.IdNumber = &entity.IDNumber
		}
		if entity.DateOfBirth != "" {
			resp.DateOfBirth = &entity.DateOfBirth
		}
		if entity.ExpiryDate != "" {
			resp.ExpiryDate = &entity.ExpiryDate
		}
	}

	if entity.ReviewedBy != nil {
		resp.ReviewedBy = *entity.ReviewedBy
	}

	if entity.CreatedAt != nil {
		resp.CreatedAt = _utils.FormatTimeToString(entity.CreatedAt)
	}
	if entity.UpdatedAt != nil {
		resp.UpdatedAt = _utils.FormatTimeToString(entity.UpdatedAt)
	}
	if entity.ReviewedAt != nil {
		resp.ReviewedAt = _utils.FormatTimeToString(entity.ReviewedAt)
	}

	return resp
}

func (m *KYCMapper) FromSubmitRequest(req *dto.SubmitKYCRequest, profileID uint64) *models.KYCEntity {
	if req == nil {
		return nil
	}

	return &models.KYCEntity{
		ProfileID:    profileID,
		FullName:     req.FullName,
		IdentityCard: req.IdentityCard,
		FrontImage:   req.FrontImage,
		BackImage:    req.BackImage,
		SelfieImage:  req.SelfieImage,
		IDNumber:     req.IDNumber,
		DateOfBirth:  req.DateOfBirth,
		ExpiryDate:   req.ExpiryDate,
		Status:       10, // Pending
	}
}

func (m *KYCMapper) ToProtoList(entities []models.KYCEntity) []*sharepb.KYCV3Proto {
	if entities == nil {
		return []*sharepb.KYCV3Proto{}
	}

	result := make([]*sharepb.KYCV3Proto, len(entities))
	for i, entity := range entities {
		result[i] = m.ToProto(&entity)
	}
	return result
}
