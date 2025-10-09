package repo

import (
	"context"
	"user/internal/dto"
	models "user/internal/models"
)

type IAdminRepo interface {
	ListAllUsers(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error)
	GetUserByProfileID(profileID uint64) (*models.AdminProfile, error)
	UpdateUserStatus(profileID uint64, status uint8) error

	// Admin Management methods
	CreateAdmin(ctx context.Context, admin *models.AdminProfile) (uint64, error)
	UpdateAdmin(ctx context.Context, req *models.AdminProfile) error
	DeleteAdmin(ctx context.Context, authID uint64) error
	AssignRole(ctx context.Context, authID uint64, roleKey string) error
}
