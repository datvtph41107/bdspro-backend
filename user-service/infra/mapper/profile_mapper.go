package mapper

import (
	_utils "common/utils"
	userpb "pb/types/user"
	"user/enums"
	"user/internal/dto"
	user_model "user/internal/models/user"
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

func (m *ProfileMapper) CertificationToPb(certification *user_model.CertificationEntity) *userpb.CertificationItem {
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

func (m *ProfileMapper) ProfessionToPb(profession *user_model.ProfessionEntity) *userpb.ProfessionItem {
	return &userpb.ProfessionItem{
		Id:        profession.ID,
		Name:      profession.Name,
		Issuer:    profession.Issuer,
		IssueDate: _utils.FormatTimeToString(profession.IssueDate),
	}
}

func (m *ProfileMapper) PublicUserInfoToPb(profile *user_model.UserProfileEntity) *userpb.ProfilePublicInfo {
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
		Birth:             _utils.FormatTimeToString(profile.Birth),
		StartDate:         _utils.FormatTimeToString(profile.StartDate),
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
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),
		TabDefault:        uint32(profile.TabDefault),
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		Following:         profile.Following,
	}

	for _, certification := range profile.Certifications {
		result.Certifications = append(result.Certifications, m.CertificationToPb(&certification))
	}

	for _, profession := range profile.Professions {
		result.Professions = append(result.Professions, m.ProfessionToPb(&profession))
	}

	// if profile.RateStats != nil {
	// 	result.RateStats = m.RateMapper.EntityToPbRateStats(profile.RateStats)
	// }

	return result
}

func (m *ProfileMapper) UserSearchItemToPbs(profiles *[]user_model.UserProfileSearch) []*userpb.UserSearchItem {
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
