package mapper

import (
	_models "common/models"
	_utils "common/utils"
	authpb "pb/types/auth"
	userpb "pb/types/user"

	_dto "common/domain/dto"
	"mdn4/modules/user/dto"
	entity "mdn4/modules/user/entity"
)

type AdminMapper struct{}

func NewAdminMapper() *AdminMapper {
	return &AdminMapper{}
}

// ListAdminRequestToDTO chuyển đổi từ proto request sang DTO request
func (m *AdminMapper) ListAdminRequestToDTO(req *userpb.AdminListRequest) *dto.AdminListRequest {
	return &dto.AdminListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		RoleID: req.RoleId,
		Name:   req.Name,
	}
}

// ListAllUsersResponseToProto chuyển đổi từ DTO response sang proto response
func (m *AdminMapper) ListAdminResponseToProto(result []entity.AdminProfile) []*userpb.AdminProfile {
	response := []*userpb.AdminProfile{}

	// Convert data items
	for _, item := range result {
		protoItem := m.ModelToPb(item)
		response = append(response, protoItem)
	}

	return response
}

// AdminUserItemToProto chuyển đổi từ DTO item sang proto item
func (m *AdminMapper) ModelToPb(item entity.AdminProfile) *userpb.AdminProfile {
	pbAdmin := &userpb.AdminProfile{
		Id:               item.ID,
		FullName:         item.FullName,
		Email:            item.Email,
		Phone:            item.Phone,
		Avatar:           item.Avatar,
		Address:          item.Address,
		Gender:           item.Gender,
		Birth:            _utils.FormatTimeToString(item.Birth),
		Status:           uint32(item.Status),
		TotalLogin:       item.TotalLogin,
		VerifiedAt:       _utils.FormatTimeToString(item.VerifiedAt),
		LockedAt:         _utils.FormatTimeToString(item.LockedAt),
		LastLoginAt:      _utils.FormatTimeToString(item.LastLoginAt),
		CreatedAt:        _utils.FormatTimeToString(item.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(item.UpdatedAt),
		JobTitle:         item.JobTitle,
		WorkAt:           _utils.FormatTimeToString(item.WorkAt),
		RoleKey:          item.RoleKey,
		RoleType:         item.RoleType,
		Attachments:      item.Attachments,
		InternalNotes:    item.InternalNotes,
		SendNotification: item.SendNotification,
		Username:         item.Username,
		SuperAdmin:       uint32(item.SuperAdmin),
	}

	// Map roleId nếu có
	if item.RoleID != nil {
		pbAdmin.RoleId = item.RoleID
	}

	// Map role info nếu có
	if item.Role != nil {
		pbAdmin.Role = &authpb.Role{
			Id:              item.Role.ID,
			RoleName:        item.Role.RoleName,
			RoleKey:         item.Role.RoleKey,
			RoleDescription: item.Role.RoleDescription,
		}
	}

	return pbAdmin
}

// CreateAdminRequestToModel chuyển đổi từ proto request sang model
func (m *AdminMapper) PbToDomain(req *userpb.AdminProfile) *entity.AdminProfile {
	admin := &entity.AdminProfile{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		Avatar:           req.Avatar,
		Address:          req.Address,
		Gender:           req.Gender,
		Birth:            _utils.ParseStringToTime(req.Birth),
		WorkAt:           _utils.ParseStringToTime(req.WorkAt),
		JobTitle:         req.JobTitle,
		RoleKey:          req.RoleKey,
		RoleType:         req.RoleType,
		Attachments:      req.Attachments,
		InternalNotes:    req.InternalNotes,
		SendNotification: req.SendNotification,
		Username:         req.Username,
		Password:         req.Password,
		SuperAdmin:       normalizeSuperAdmin(req.SuperAdmin),
	}

	// Map roleId nếu có
	if req.RoleId != nil {
		admin.RoleID = req.RoleId
	}

	return admin
}

// AdminProfileResponseToProto chuyển đổi từ model response sang proto response
func (m *AdminMapper) AdminProfileResponseToProto(result *entity.AdminProfile) *userpb.AdminProfile {
	pbAdmin := &userpb.AdminProfile{
		Id:               result.ID,
		FullName:         result.FullName,
		Email:            result.Email,
		Phone:            result.Phone,
		Avatar:           result.Avatar,
		Address:          result.Address,
		Gender:           result.Gender,
		Birth:            _utils.FormatTimeToString(result.Birth),
		CreatedAt:        _utils.FormatTimeToString(result.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(result.UpdatedAt),
		Status:           uint32(result.Status),
		TotalLogin:       result.TotalLogin,
		VerifiedAt:       _utils.FormatTimeToString(result.VerifiedAt),
		LockedAt:         _utils.FormatTimeToString(result.LockedAt),
		LastLoginAt:      _utils.FormatTimeToString(result.LastLoginAt),
		JobTitle:         result.JobTitle,
		WorkAt:           _utils.FormatTimeToString(result.WorkAt),
		RoleKey:          result.RoleKey,
		RoleType:         result.RoleType,
		Attachments:      result.Attachments,
		InternalNotes:    result.InternalNotes,
		SendNotification: result.SendNotification,
		Username:         result.Username,
		SuperAdmin:       uint32(result.SuperAdmin),
	}

	// Map roleId nếu có
	if result.RoleID != nil {
		pbAdmin.RoleId = result.RoleID
	}

	// Map role info nếu có
	if result.Role != nil {
		pbAdmin.Role = &authpb.Role{
			Id:              result.Role.ID,
			RoleName:        result.Role.RoleName,
			RoleKey:         result.Role.RoleKey,
			RoleDescription: result.Role.RoleDescription,
		}
	}

	return pbAdmin
}

func normalizeSuperAdmin(value uint32) uint8 {
	if value != 0 {
		return 1
	}
	return 0
}

// DeleteAdminRequestToDTO chuyển đổi từ proto request sang DTO
func (m *AdminMapper) DeleteAdminRequestToDTO(req *authpb.DeleteAdminRequest) *dto.DeleteAdminRequest {
	return &dto.DeleteAdminRequest{
		AuthID: req.AuthId,
	}
}
