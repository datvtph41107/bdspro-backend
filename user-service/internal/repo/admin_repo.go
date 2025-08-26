package repo

import (
	"context"
	"user/internal/dto"
)

// AdminRepository interface cho admin operations
type AdminRepository interface {
	ListAllUsers(ctx context.Context, req *dto.ListAllUsersRequest) ([]dto.AdminUserItem, int64, error)
}
