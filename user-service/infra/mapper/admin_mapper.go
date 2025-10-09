package mapper

import (
	_utils "common/utils"
	authpb "pb/types/auth"
	userpb "pb/types/user"

	_dto "common/domain/dto"
	"user/internal/dto"
	models "user/internal/models"
)

type AdminMapper struct{}

func NewAdminMapper() *AdminMapper {
	return &AdminMapper{}
}

// ListAdminRequestToDTO chuyển đổi từ proto request sang DTO request
func (m *AdminMapper) ListAdminRequestToDTO(req *userpb.ListAdminRequest) *dto.AdminListRequest {
	return &dto.AdminListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}
}

// ListAllUsersResponseToProto chuyển đổi từ DTO response sang proto response
func (m *AdminMapper) ListAdminResponseToProto(result []models.AdminProfile) []*userpb.AdminProfile {
	response := []*userpb.AdminProfile{}

	// Convert data items
	for _, item := range result {
		protoItem := m.AdminUserItemToProto(item)
		response = append(response, protoItem)
	}

	return response
}

// AdminUserItemToProto chuyển đổi từ DTO item sang proto item
func (m *AdminMapper) AdminUserItemToProto(item models.AdminProfile) *userpb.AdminProfile {
	return &userpb.AdminProfile{
		Id:       item.ProfileID,
		FullName: item.FullName,
		Email:    item.Email,
		Phone:    item.Phone,
		Avatar:   item.Avatar,
		// Address:        item.Address,
		// Gender:         uint32(item.Gender),
		// Birth:          _utils.FormatTimeToString(item.Birth),
		Status: uint32(item.Status),
		// RoleType:       item.RoleType,
		// RoleRealEstate: item.RoleRealEstate,
		// Position:       item.Position,
		// DepartmentId:   item.DepartmentID,
		// Bookmark:       item.Bookmark,
		// TickVerified:   item.TickVerified,
		TotalLogin:  item.TotalLogin,
		VerifiedAt:  _utils.FormatTimeToString(item.VerifiedAt),
		LockedAt:    _utils.FormatTimeToString(item.LockedAt),
		LastLoginAt: _utils.FormatTimeToString(item.LastLoginAt),
		CreatedAt:   _utils.FormatTimeToString(item.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(item.UpdatedAt),
	}
}

// CreateAdminRequestToModel chuyển đổi từ proto request sang model
func (m *AdminMapper) CreateAdminRequestToModel(req *authpb.CreateAdminRequest) *models.AdminProfile {
	return &models.AdminProfile{
		Username:         req.Username,
		Password:         req.Password,
		FullName:         req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		Avatar:           req.Avatar,
		JobTitle:         req.JobTitle,
		RoleKey:          req.RoleKey,
		Attachments:      req.Attachments,
		InternalNotes:    req.InternalNotes,
		SendNotification: req.SendNotification,
	}
}

// AdminProfileResponseToProto chuyển đổi từ model response sang proto response
func (m *AdminMapper) AdminProfileResponseToProto(result *models.AdminProfile) *authpb.AuthSuccessResponse {
	return &authpb.AuthSuccessResponse{
		AuthId:    result.AuthID,
		ProfileId: result.ProfileID,
		Role:      result.RoleKey,
	}
}

// UpdateAdminRequestToDTO chuyển đổi từ proto request sang DTO
func (m *AdminMapper) UpdateAdminRequestToDTO(req *authpb.UpdateAdminRequest) *dto.UpdateAdminRequest {
	return &dto.UpdateAdminRequest{
		AuthID:   req.AuthId,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	}
}

// DeleteAdminRequestToDTO chuyển đổi từ proto request sang DTO
func (m *AdminMapper) DeleteAdminRequestToDTO(req *authpb.DeleteAdminRequest) *dto.DeleteAdminRequest {
	return &dto.DeleteAdminRequest{
		AuthID: req.Id,
	}
}

// AssignRoleRequestToDTO chuyển đổi từ proto request sang DTO
func (m *AdminMapper) AssignRoleRequestToDTO(req *authpb.AssignRoleRequest) *dto.AssignRoleRequest {
	return &dto.AssignRoleRequest{
		AuthID:  req.AuthId,
		RoleKey: req.RoleKey,
	}
}
