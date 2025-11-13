package mapper

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/enums"
	"user/internal/dto"
	models "user/internal/models"
)

type ProfileMapper struct {
	// RateMapper *RateMapper
}

func NewProfileMapper() *ProfileMapper {
	return &ProfileMapper{
		// RateMapper: rateMapper,
	}
}

func (m *ProfileMapper) SetProfileToDTO(profile *userpb.SetInfoRequest) *dto.UserInfoRequest {
	result := &dto.UserInfoRequest{
		Phone:             profile.Phone,
		Phone2:            profile.Phone2,
		Email:             profile.Email,
		TaxCode:           profile.TaxCode,
		FullName:          profile.FullName,
		Address:           profile.Address,
		Avatar:            profile.Avatar,
		BackgroundImage:   profile.BackgroundImage,
		Position:          profile.Position,
		DepartmentID:      profile.DepartmentId,
		Birth:             _utils.ParseStringToTime(profile.Birth),
		StartDate:         _utils.ParseStringToTime(profile.StartDate),
		FrontIdentify:     profile.FrontIdentify,
		BackIdentify:      profile.BackIdentify,
		Slogan:            profile.Slogan,
		Website:           profile.Website,
		Facebook:          profile.Facebook,
		Instagram:         profile.Instagram,
		Twitter:           profile.Twitter,
		Linkedin:          profile.Linkedin,
		Youtube:           profile.Youtube,
		Introduction:      profile.Introduction,
		Visibility:        uint8(profile.Visibility),
		ProfileVisibility: uint8(profile.ProfileVisibility),
		StatusOnline:      enums.EOnlineStatus(profile.StatusOnline),
		TabDefault:        enums.ESettingTab(profile.TabDefault),
		RoleRealEstate:    enums.ERoleRealEstate(profile.RoleRealEstate),

		CertificationRemoveIds: profile.CertificationRemoveIds,
		ProfessionRemoveIds:    profile.ProfessionRemoveIds,
		PurposeUseIDs:          profile.PurposeUseIds,
		MainAreaIds:            profile.MainAreaIds,
	}

	for _, certification := range profile.Certifications {
		result.Certifications = append(result.Certifications, *m.CertificationToDTO(certification))
	}

	for _, profession := range profile.Professions {
		result.Professions = append(result.Professions, *m.ProfessionToDTO(profession))
	}

	return result
}

func (m *ProfileMapper) CertificationToDTO(certification *userpb.CertificationItem) *dto.CertificationBatchSave {
	return &dto.CertificationBatchSave{
		ID:        certification.Id,
		Name:      certification.Name,
		Issuer:    certification.Issuer,
		IssueDate: _utils.ParseStringToTime(certification.IssueDate),
		FileName:  certification.FileName,
		FileURL:   certification.FileUrl,
		FileType:  certification.FileType,
	}
}

func (m *ProfileMapper) ProfessionToDTO(profession *userpb.ProfessionItem) *dto.ProfessionBatchSave {
	return &dto.ProfessionBatchSave{
		ID:        profession.Id,
		Name:      profession.Name,
		Issuer:    profession.Issuer,
		IssueDate: _utils.ParseStringToTime(profession.IssueDate),
		// FileName:  profession.FileName,
		// FileURL:   profession.FileUrl,
		// FileType:  profession.FileType,
	}
}

func (m *ProfileMapper) CertificationToPb(certification *models.CertificationEntity) *userpb.CertificationItem {
	return &userpb.CertificationItem{
		Id:        certification.ID,
		Name:      certification.Name,
		Issuer:    certification.Issuer,
		IssueDate: _utils.FormatTimeToString(certification.IssueDate),
		FileName:  certification.FileName,
		FileUrl:   certification.FileURL,
		FileType:  certification.FileType,
	}
}

func (m *ProfileMapper) ProfessionToPb(profession *models.ProfessionEntity) *userpb.ProfessionItem {
	return &userpb.ProfessionItem{
		Id:        profession.ID,
		Name:      profession.Name,
		Issuer:    profession.Issuer,
		IssueDate: _utils.FormatTimeToString(profession.IssueDate),
	}
}

func (m *ProfileMapper) PurposeUseToPb(purpose *models.PurposeUseEntity) *userpb.PurposeUseItem {
	result := &userpb.PurposeUseItem{
		Id:   purpose.ID,
		Name: purpose.Name,
	}
	if purpose.Description != "" {
		result.Description = &purpose.Description
	}
	return result
}

func (m *ProfileMapper) PublicUserInfoToPb(profile *models.UserProfileEntity) *userpb.ProfilePublicInfo {
	result := &userpb.ProfilePublicInfo{
		ProfileId:         profile.ProfileID,
		Phone:             profile.Phone,
		Phone2:            profile.Phone2,
		Email:             profile.Email,
		TaxCode:           profile.TaxCode,
		FullName:          profile.FullName,
		Address:           profile.Address,
		Avatar:            profile.Avatar,
		BackgroundImage:   profile.BackgroundImage,
		Position:          profile.Position,
		DepartmentId:      profile.DepartmentID,
		Gender:            profile.Gender,
		Birth:             _utils.FormatTimeToString(profile.Birth),
		StartDate:         _utils.FormatTimeToString(profile.StartDate),
		Slogan:            profile.Slogan,
		Website:           profile.Website,
		ProvinceId:        profile.ProvinceID,
		WardId:            profile.WardID,
		Facebook:          profile.Facebook,
		Instagram:         profile.Instagram,
		Twitter:           profile.Twitter,
		Linkedin:          profile.Linkedin,
		Youtube:           profile.Youtube,
		Introduction:      profile.Introduction,
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),
		TabDefault:        uint32(profile.TabDefault),
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		Following:         profile.Following,
		ReferralCode:      profile.ReferralCode,
		// FrontIdentify:     profile.FrontIdentify,
		// BackIdentify:      profile.BackIdentify,
	}

	for _, certification := range profile.Certifications {
		result.Certifications = append(result.Certifications, m.CertificationToPb(&certification))
	}

	for _, profession := range profile.Professions {
		result.Professions = append(result.Professions, m.ProfessionToPb(&profession))
	}

	for _, purpose := range profile.PurposeUses {
		result.PurposeUses = append(result.PurposeUses, m.PurposeUseToPb(&purpose))
	}

	for _, area := range profile.MainAreas {
		item := &userpb.MainAreaItem{
			Id:   area.ID,
			Name: area.Name,
		}
		if area.Description != "" {
			desc := area.Description
			item.Description = &desc
		}
		result.MainAreas = append(result.MainAreas, item)
	}

	// if profile.RateStats != nil {
	// 	result.RateStats = m.RateMapper.EntityToPbRateStats(profile.RateStats)
	// }

	return result
}

func (m *ProfileMapper) UserSearchItemToPbs(profiles *[]models.UserProfileSearch) []*userpb.UserSearchItem {
	infos := []*userpb.UserSearchItem{}
	for _, profile := range *profiles {
		infos = append(infos, &userpb.UserSearchItem{
			ProfileId: profile.ProfileID,
			FullName:  profile.FullName,
			Avatar:    profile.Avatar,
		})
	}
	return infos
}

func (m *ProfileMapper) SetPersonToDTO(req *userpb.SetPersonRequest) *dto.UserInfoRequest {
	result := &dto.UserInfoRequest{}
	if req.FullName != nil {
		result.FullName = *req.FullName
	}
	if req.Gender != nil {
		result.Gender = *req.Gender
	}
	if req.Birth != nil {
		result.Birth = _utils.ParseStringToTime(*req.Birth)
	}
	if req.Phone != nil {
		result.Phone = *req.Phone
	}
	if req.Email != nil {
		result.Email = *req.Email
	}
	if req.ProvinceId != nil {
		result.ProvinceID = req.ProvinceId
	}
	if req.WardId != nil {
		result.WardID = req.WardId
	}
	return result
}

func (m *ProfileMapper) UserV3ToDTO(req *sharepb.UserV3Proto) *_dto.UserV3DTO {
	result := &_dto.UserV3DTO{
		ProfileID:      req.ProfileId,
		FullName:       req.FullName,
		Avatar:         req.Avatar,
		RoleRealEstate: req.RoleRealEstate,
		MainAreaIds:    req.MainAreaIds,
		PurposeUseIDs:  req.PurposeUseIds,
		// MainAreaNews:   req.MainAreaNews,
	}

	if len(req.MainAreaNews) > 0 {
		for _, news := range req.MainAreaNews {
			result.MainAreaNews = append(result.MainAreaNews, _dto.ItemDTO{
				ID:   news.Id,
				Name: news.Name,
			})
		}
	}

	return result
}

func (m *ProfileMapper) OwnerProfileToUserV3Proto(profile *models.UserProfileEntity) *sharepb.UserV3Proto {
	if profile == nil {
		return nil
	}
	result := &sharepb.UserV3Proto{
		ProfileId:       profile.ProfileID,
		FullName:        profile.FullName,
		Avatar:          profile.Avatar,
		BackgroundImage: profile.BackgroundImage,
		Gender:          profile.Gender,
		Birth:           _utils.FormatTimeToString(profile.Birth),
		Phone:           profile.Phone,
		Email:           profile.Email,
		Address: &sharepb.AddressV3Proto{
			Detail:     profile.Address,
			ProvinceId: profile.ProvinceID,
			WardId:     profile.WardID,
		},
		RoleRealEstate: uint32(profile.RoleRealEstate),
		MainAreaIds:    nil,
		PurposeUseIds:  nil,
		CreatedAt:      _utils.FormatTimeToString(profile.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(profile.UpdatedAt),
	}

	if len(profile.MainAreas) > 0 {
		for _, area := range profile.MainAreas {
			result.MainAreaIds = append(result.MainAreaIds, area.ID)
		}
	}
	if len(profile.PurposeUses) > 0 {
		for _, purpose := range profile.PurposeUses {
			result.PurposeUseIds = append(result.PurposeUseIds, purpose.ID)
		}
	}
	return result
}
