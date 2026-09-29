package mapper

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_models "common/models"
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
		Workplace:         profile.Workplace,
		RoleTitle:         profile.RoleTitle,
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
		// SignatureVisible:  profile.SignatureVisible,

		CertificationRemoveIds: profile.CertificationRemoveIds,
		ProfessionRemoveIds:    profile.ProfessionRemoveIds,
		PurposeUseIDs:          profile.PurposeUseIds,
		MainAreaIds:            profile.MainAreaIds,
		UpdateTags:             profile.UpdateTags,
	}

	for _, certification := range profile.Certifications {
		result.Certifications = append(result.Certifications, *m.CertificationToDTO(certification))
	}

	for _, profession := range profile.Professions {
		result.Professions = append(result.Professions, *m.ProfessionToDTO(profession))
	}

	for _, tag := range profile.Tags {
		if model := m.TagToModel(tag); model != nil {
			result.Tags = append(result.Tags, *model)
		}
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
		IsActive:  profession.GetIsActive(),
		// FileName:  profession.FileName,
		// FileURL:   profession.FileUrl,
		// FileType:  profession.FileType,
	}
}

func (m *ProfileMapper) TagToModel(tag *userpb.TagItem) *models.TagEntity {
	if tag == nil {
		return nil
	}
	return &models.TagEntity{
		BaseEntity: _models.BaseEntity{
			ID: tag.Id,
		},
		Name:      tag.Name,
		TagType:   enums.TagType(tag.TagType),
		IsDefault: tag.IsDefault,
		IsActive:  tag.IsActive,
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
	result := &userpb.ProfessionItem{
		Id:             profession.ID,
		Name:           profession.Name,
		Issuer:         profession.Issuer,
		IssueDate:      _utils.FormatTimeToString(profession.IssueDate),
		VerifiedStatus: uint32(profession.VerifiedStatus),
	}

	isActive := profession.IsActive
	result.IsActive = &isActive

	return result
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

func (m *ProfileMapper) TagToPb(tag *models.TagEntity) *userpb.TagItem {
	isActive := tag.IsActive
	return &userpb.TagItem{
		Id:        tag.ID,
		Name:      tag.Name,
		TagType:   uint32(tag.TagType),
		IsDefault: tag.IsDefault,
		IsActive:  isActive,
	}
}

func (m *ProfileMapper) PublicUserInfoToPb(profile *models.UserProfileEntity, kyc *sharepb.KYCV3Proto) *userpb.ProfilePublicInfo {
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
		Workplace:         profile.Workplace,
		RoleTitle:         profile.RoleTitle,
		DepartmentId:      profile.DepartmentID,
		Gender:            profile.Gender,
		Birth:             _utils.FormatTimeToString(profile.Birth),
		StartDate:         _utils.FormatTimeToString(profile.StartDate),
		Slogan:            profile.Slogan,
		Website:           profile.WebsiteURL,
		ProvinceId:        profile.ProvinceID,
		WardId:            profile.WardID,
		Facebook:          profile.FacebookURL,
		Instagram:         profile.InstagramURL,
		Twitter:           profile.TwitterURL,
		Linkedin:          profile.LinkedinURL,
		Youtube:           profile.YoutubeURL,
		Introduction:      profile.Introduction,
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),
		TabDefault:        uint32(profile.TabDefault),
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		Following:         profile.Following,
		ReferralCode:      profile.ReferralCode,
		// VisibilityIntroduce:  uint32(profile.VisibilityIntroduce),
		// VisibilityProfession: uint32(profile.VisibilityProfession),
		// VisibilityMainArea:   uint32(profile.VisibilityMainArea),
		// VisibilityFriends:    uint32(profile.VisibilityFriends),
		// VisibilitySignature:  uint32(profile.VisibilitySignature),
		// FrontIdentify:     profile.FrontIdentify,
		// BackIdentify:      profile.BackIdentify,
	}

	// Convert ViewRoles from []enums.EViewRole to []uint32
	// for _, role := range profile.ViewRoles {
	// 	result.ViewRoles = append(result.ViewRoles, uint32(role))
	// }

	for _, certification := range profile.Certifications {
		result.Certifications = append(result.Certifications, m.CertificationToPb(&certification))
	}

	for _, profession := range profile.Professions {
		result.Professions = append(result.Professions, m.ProfessionToPb(&profession))
	}

	for _, purpose := range profile.PurposeUses {
		result.PurposeUses = append(result.PurposeUses, m.PurposeUseToPb(&purpose))
	}

	// for _, area := range profile.TagMainAreas {
	// 	item := &userpb.MainAreaItem{
	// 		Id:   area.ID,
	// 		Name: area.Name,
	// 	}
	// 	result.MainAreas = append(result.MainAreas, item)
	// }

	for _, tag := range profile.Tags {
		result.Tags = append(result.Tags, m.TagToPb(&tag))
	}

	// if profile.RateStats != nil {
	// 	result.RateStats = m.RateMapper.EntityToPbRateStats(profile.RateStats)
	// }

	// Thêm KYC vào response
	if kyc != nil {
		result.Kyc = kyc
	}

	return result
}

func (m *ProfileMapper) UserSearchItemToPbs(profiles *[]models.UserProfileSearch) []*userpb.UserSearchItem {
	infos := []*userpb.UserSearchItem{}
	if profiles != nil {
		for _, profile := range *profiles {
			infos = append(infos, &userpb.UserSearchItem{
				ProfileId: profile.ProfileID,
				FullName:  profile.FullName,
				Avatar:    profile.Avatar,
			})
		}
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

func (m *ProfileMapper) PrivacyRequestToDTO(req *userpb.UpdatePrivacySettingRequest) *dto.ProfilePrivacyRequest {
	if req == nil {
		return nil
	}

	result := &dto.ProfilePrivacyRequest{
		Visibility:        uint8(req.Visibility),
		ProfileVisibility: uint8(req.ProfileVisibility),
		StatusOnline:      enums.EOnlineStatus(req.StatusOnline),
	}

	if len(req.PersonConfigs) > 0 {
		result.PersonConfigs = make([]*dto.PersonConfigDTO, len(req.PersonConfigs))
		for i, cfg := range req.PersonConfigs {
			result.PersonConfigs[i] = &dto.PersonConfigDTO{
				ID:        cfg.Id,
				Key:       cfg.Key,
				Checked:   cfg.Checked,
				IsDefault: cfg.IsDefault,
				Channel:   _enum.EChannelNotification(cfg.Channel),
			}
		}
	}

	return result
}

func (m *ProfileMapper) UserV3ToDTO(req *sharepb.UserV3Proto) *_dto.UserV3DTO {
	result := &_dto.UserV3DTO{
		ProfileID:        req.ProfileId,
		FullName:         req.FullName,
		Avatar:           req.Avatar,
		BackgroundImage:  req.BackgroundImage,
		RoleRealEstate:   req.RoleRealEstate,
		MainAreaIds:      req.MainAreaIds,
		PurposeUseIDs:    req.PurposeUseIds,
		FacebookURL:      req.FacebookUrl,
		InstagramURL:     req.InstagramUrl,
		TwitterURL:       req.TwitterUrl,
		LinkedinURL:      req.LinkedinUrl,
		YoutubeURL:       req.YoutubeUrl,
		WebsiteURL:       req.WebsiteUrl,
		ZaloURL:          req.ZaloUrl,
		SignatureVisible: req.SignatureVisible,
		Gender:           req.Gender,
		Birth:            _utils.ParseStringToTime(req.Birth),
		UpdateTags:       req.UpdateTags,
		Introduction:     req.Introduction,
		Email:            req.Email,
		Phone2:           req.Phone2,
		Workplace:        req.Workplace,
		RoleTitle:        req.RoleTitle,

		VisibilityIntroduce:  req.VisibilityIntroduce,
		VisibilityProfession: req.VisibilityProfession,
		VisibilityMainArea:   req.VisibilityMainArea,
		VisibilityFriends:    req.VisibilityFriends,
		VisibilitySignature:  req.VisibilitySignature,
		ViewRoles:            req.ViewRoles,
		FieldSets:            req.FieldSets,
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

	if req.Address != nil {
		result.Address = &_dto.AddressV3DTO{
			Detail:       req.Address.Detail,
			ProvinceID:   req.Address.ProvinceId,
			ProvinceName: req.Address.ProvinceName,
			DistrictID:   req.Address.DistrictId,
			DistrictName: req.Address.DistrictName,
			WardID:       req.Address.WardId,
			WardName:     req.Address.WardName,
		}
	}

	if req.Tags != nil && result.UpdateTags {
		for _, tag := range req.Tags {
			result.Tags = append(result.Tags, _dto.ItemDTO{
				ID:       tag.Id,
				Name:     tag.Name,
				ItemType: tag.ItemType,
			})
		}
	}

	// Map certifications từ ItemV3Proto sang ItemDTO
	if len(req.Certifications) > 0 {
		for _, cert := range req.Certifications {
			result.Certifications = append(result.Certifications, _dto.ItemDTO{
				ID:       cert.Id,
				Name:     cert.Name,
				FileName: cert.FileName,
				FileURL:  cert.FileUrl,
				FileType: cert.FileType,
				Status:   cert.Status,
				Issuer:   cert.Issuer,
				IssuedAt: cert.IssuedAt,
			})
		}
	}

	// Map medias từ ItemV3Proto sang ItemDTO
	if len(req.Medias) > 0 {
		for _, media := range req.Medias {
			result.Medias = append(result.Medias, _dto.ItemDTO{
				ID:       media.Id,
				Name:     media.Name,
				FileName: media.FileName,
				FileURL:  media.FileUrl,
				FileType: media.FileType,
			})
		}
	}

	return result
}

func (m *ProfileMapper) ModelToUserV3Proto(profile *models.UserProfileEntity, kyc *sharepb.KYCV3Proto) *sharepb.UserV3Proto {
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
		Workplace:       profile.Workplace,
		RoleTitle:       profile.RoleTitle,
		Address: &sharepb.AddressV3Proto{
			Detail:     profile.Address,
			ProvinceId: profile.ProvinceID,
			WardId:     profile.WardID,
		},
		Phone2:            profile.Phone2,
		ZaloUrl:           profile.ZaloURL,
		FacebookUrl:       profile.FacebookURL,
		InstagramUrl:      profile.InstagramURL,
		TwitterUrl:        profile.TwitterURL,
		LinkedinUrl:       profile.LinkedinURL,
		YoutubeUrl:        profile.YoutubeURL,
		WebsiteUrl:        profile.WebsiteURL,
		SignatureVisible:  profile.SignatureVisible,
		RoleRealEstate:    uint32(profile.RoleRealEstate),
		MainAreaIds:       nil,
		PurposeUseIds:     nil,
		Visibility:        uint32(profile.Visibility),
		ProfileVisibility: uint32(profile.ProfileVisibility),
		StatusOnline:      uint32(profile.StatusOnline),

		VisibilityIntroduce:  uint32(profile.VisibilityIntroduce),
		VisibilityProfession: uint32(profile.VisibilityProfession),
		VisibilityMainArea:   uint32(profile.VisibilityMainArea),
		VisibilityFriends:    uint32(profile.VisibilityFriends),
		VisibilitySignature:  uint32(profile.VisibilitySignature),
		CreatedAt:            _utils.FormatTimeToString(profile.CreatedAt),
		UpdatedAt:            _utils.FormatTimeToString(profile.UpdatedAt),

		// TabDefault:        enums.ESettingTab(profile.TabDefault),
	}

	// Convert ViewRoles from pq.Int32Array to []uint32
	if len(profile.ViewRoles) > 0 {
		for _, role := range profile.ViewRoles {
			result.ViewRoles = append(result.ViewRoles, uint32(role))
		}
	}

	if len(profile.TagJobs) > 0 {
		for _, profession := range profile.TagJobs {
			result.TagJobs = append(result.TagJobs, &sharepb.ItemV3Proto{
				Id:   profession.ID,
				Name: profession.Name,
			})
		}
	}
	if len(profile.TagMainAreas) > 0 {
		for _, area := range profile.TagMainAreas {
			result.TagMainAreas = append(result.TagMainAreas, &sharepb.ItemV3Proto{
				Id:   area.ID,
				Name: area.Name,
			})
		}
	}
	if len(profile.TagSpecialities) > 0 {
		for _, specialty := range profile.TagSpecialities {
			result.TagSpecialities = append(result.TagSpecialities, &sharepb.ItemV3Proto{
				Id:   specialty.ID,
				Name: specialty.Name,
			})
		}
	}
	if len(profile.TagProjects) > 0 {
		for _, project := range profile.TagProjects {
			result.TagProjects = append(result.TagProjects, &sharepb.ItemV3Proto{
				Id:   project.ID,
				Name: project.Name,
			})
		}
	}
	if len(profile.PurposeUses) > 0 {
		for _, purpose := range profile.PurposeUses {
			result.PurposeUseIds = append(result.PurposeUseIds, purpose.ID)
		}
	}

	if len(profile.Tags) > 0 {
		for _, tag := range profile.Tags {
			result.MainAreaNews = append(result.MainAreaNews, &sharepb.ItemV3Proto{
				Id:   tag.ID,
				Name: tag.Name,
			})
		}
	}

	// Thêm KYC vào response
	if kyc != nil {
		result.Kyc = kyc
	}

	return result
}

func (m *ProfileMapper) DTOToUserV3Proto(profile *_dto.UserV3DTO, kyc *sharepb.KYCV3Proto) *sharepb.UserV3Proto {
	if profile == nil {
		return nil
	}
	result := &sharepb.UserV3Proto{
		ProfileId:        profile.ProfileID,
		FullName:         profile.FullName,
		Avatar:           profile.Avatar,
		BackgroundImage:  profile.BackgroundImage,
		Gender:           profile.Gender,
		Birth:            _utils.FormatTimeToString(profile.Birth),
		Phone:            profile.Phone,
		Phone2:           profile.Phone2,
		Email:            profile.Email,
		Workplace:        profile.Workplace,
		RoleTitle:        profile.RoleTitle,
		ZaloUrl:          profile.ZaloURL,
		FacebookUrl:      profile.FacebookURL,
		InstagramUrl:     profile.InstagramURL,
		TwitterUrl:       profile.TwitterURL,
		LinkedinUrl:      profile.LinkedinURL,
		YoutubeUrl:       profile.YoutubeURL,
		WebsiteUrl:       profile.WebsiteURL,
		Introduction:     profile.Introduction,
		SignatureVisible: profile.SignatureVisible,
		RoleRealEstate:   uint32(profile.RoleRealEstate),
		MainAreaIds:      nil,
		PurposeUseIds:    nil,

		Visibility:           profile.Visibility,
		ProfileVisibility:    profile.ProfileVisibility,
		StatusOnline:         profile.StatusOnline,
		VisibilityIntroduce:  profile.VisibilityIntroduce,
		VisibilityProfession: profile.VisibilityProfession,
		VisibilityMainArea:   profile.VisibilityMainArea,
		VisibilityFriends:    profile.VisibilityFriends,
		VisibilitySignature:  profile.VisibilitySignature,
		CreatedAt:            _utils.FormatTimeToString(profile.CreatedAt),
		UpdatedAt:            _utils.FormatTimeToString(profile.UpdatedAt),
	}

	// Convert ViewRoles from pq.Int32Array to []uint32
	if len(profile.ViewRoles) > 0 {
		for _, role := range profile.ViewRoles {
			result.ViewRoles = append(result.ViewRoles, uint32(role))
		}
	}

	if profile.Address != nil {
		result.Address = &sharepb.AddressV3Proto{
			Detail:       profile.Address.Detail,
			ProvinceId:   profile.Address.ProvinceID,
			ProvinceName: profile.Address.ProvinceName,
			WardId:       profile.Address.WardID,
			WardName:     profile.Address.WardName,
		}
	}

	if len(profile.TagJobs) > 0 {
		for _, profession := range profile.TagJobs {
			result.TagJobs = append(result.TagJobs, &sharepb.ItemV3Proto{
				Id:   profession.ID,
				Name: profession.Name,
			})
		}
	}
	if len(profile.TagMainAreas) > 0 {
		for _, area := range profile.TagMainAreas {
			result.TagMainAreas = append(result.TagMainAreas, &sharepb.ItemV3Proto{
				Id:   area.ID,
				Name: area.Name,
			})
		}
	}
	if len(profile.TagSpecialities) > 0 {
		for _, specialty := range profile.TagSpecialities {
			result.TagSpecialities = append(result.TagSpecialities, &sharepb.ItemV3Proto{
				Id:   specialty.ID,
				Name: specialty.Name,
			})
		}
	}
	if len(profile.TagProjects) > 0 {
		for _, project := range profile.TagProjects {
			result.TagProjects = append(result.TagProjects, &sharepb.ItemV3Proto{
				Id:   project.ID,
				Name: project.Name,
			})
		}
	}
	// if len(profile.PurposeUses) > 0 {
	// 	for _, purpose := range profile.PurposeUses {
	// 		result.PurposeUseIds = append(result.PurposeUseIds, purpose.ID)
	// 	}
	// }

	if len(profile.Tags) > 0 {
		for _, tag := range profile.Tags {
			result.MainAreaNews = append(result.MainAreaNews, &sharepb.ItemV3Proto{
				Id:   tag.ID,
				Name: tag.Name,
			})
		}
	}

	// Map certifications từ ItemDTO sang ItemV3Proto
	if len(profile.Certifications) > 0 {
		for _, cert := range profile.Certifications {
			result.Certifications = append(result.Certifications, &sharepb.ItemV3Proto{
				Id:       cert.ID,
				Name:     cert.Name,
				ItemType: cert.ItemType,
				FileName: cert.FileName,
				FileUrl:  cert.FileURL,
				FileType: cert.FileType,
				Status:   cert.Status,
				Issuer:   cert.Issuer,
				IssuedAt: cert.IssuedAt,
			})
		}
	}

	// Map medias từ ItemDTO sang ItemV3Proto
	if len(profile.Medias) > 0 {
		for _, media := range profile.Medias {
			itemProto := &sharepb.ItemV3Proto{
				Id:   media.ID,
				Name: media.Name,
			}
			if media.FileName != nil {
				itemProto.FileName = media.FileName
			}
			if media.FileURL != nil {
				itemProto.FileUrl = media.FileURL
			}
			if media.FileType != nil {
				itemProto.FileType = media.FileType
			}
			result.Medias = append(result.Medias, itemProto)
		}
	}

	// Thêm KYC vào response
	if kyc != nil {
		result.Kyc = kyc
	}

	return result
}

func (m *ProfileMapper) ModelToUserV3ProtoWithMedias(profile *models.UserProfileEntity, kyc *sharepb.KYCV3Proto) *sharepb.UserV3Proto {
	result := m.ModelToUserV3Proto(profile, kyc)
	if result == nil {
		return nil
	}

	// Map medias từ ProfileMediaEntity sang ItemV3Proto
	if len(profile.ProfileMedias) > 0 {
		for _, media := range profile.ProfileMedias {
			fileName := media.FileName
			fileURL := media.FileURL
			fileType := media.FileType
			result.Medias = append(result.Medias, &sharepb.ItemV3Proto{
				Id:       media.ID,
				Name:     media.Name,
				FileName: &fileName,
				FileUrl:  &fileURL,
				FileType: &fileType,
			})
		}
	}

	return result
}

func (m *ProfileMapper) DTOToFriendV3Proto(friend *_dto.FriendV3DTO) *sharepb.FriendV3Proto {
	if friend == nil {
		return nil
	}
	result := &sharepb.FriendV3Proto{
		Id:         friend.ID,
		ReceiverId: friend.ReceiverID,
		CreatedBy:  friend.CreatedBy,
		Status:     friend.Status,
	}
	return result
}

func (m *ProfileMapper) DTOToRateStatsV3Proto(rateStats *_dto.RateStateV3DTO) *sharepb.RateStatsV3Proto {
	if rateStats == nil {
		return nil
	}
	result := &sharepb.RateStatsV3Proto{
		TotalReviews: rateStats.TotalReviews,
		AverageScore: rateStats.AverageScore,
		Star1Count:   rateStats.Star1Count,
		Star2Count:   rateStats.Star2Count,
		Star3Count:   rateStats.Star3Count,
		Star4Count:   rateStats.Star4Count,
		Star5Count:   rateStats.Star5Count,
	}
	return result
}

func (m *ProfileMapper) DTOToKYCProto(kyc *_dto.KYCV3DTO) *sharepb.KYCV3Proto {
	return m.DTOToKYCProtoWithOwnerCheck(kyc, false)
}

// DTOToKYCProtoWithOwnerCheck trả về KYC proto từ DTO, chỉ trả về 4 trường nhạy cảm nếu isOwner = true
func (m *ProfileMapper) DTOToKYCProtoWithOwnerCheck(kyc *_dto.KYCV3DTO, isOwner bool) *sharepb.KYCV3Proto {
	if kyc == nil {
		return nil
	}
	result := &sharepb.KYCV3Proto{
		Id:           kyc.ID,
		ProfileId:    kyc.ProfileID,
		FullName:     kyc.FullName,
		IdentityCard: kyc.IdentityCard,
		FrontImage:   kyc.FrontImage,
		BackImage:    kyc.BackImage,
		SelfieImage:  kyc.SelfieImage,
		Status:       uint32(kyc.Status),
		RejectReason: kyc.RejectReason,
		ReviewedAt:   _utils.FormatTimeToString(kyc.ReviewedAt),
		CreatedAt:    _utils.FormatTimeToString(kyc.CreatedAt),
		UpdatedAt:    _utils.FormatTimeToString(kyc.UpdatedAt),
	}

	// Chỉ trả về 4 trường nhạy cảm nếu là owner request
	if isOwner {
		if kyc.IDNumber != "" {
			result.IdNumber = &kyc.IDNumber
		}
		if kyc.DateOfBirth != "" {
			result.DateOfBirth = &kyc.DateOfBirth
		}
		if kyc.ExpiryDate != "" {
			result.ExpiryDate = &kyc.ExpiryDate
		}
	}

	return result
}
