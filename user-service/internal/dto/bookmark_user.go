package dto

import (
	_dto "common/domain/dto"
	models "user/internal/models"
)

// BookmarkUserListRequest request để lấy danh sách bookmark user
type BookmarkUserListRequest struct {
	_dto.Pagable
	AdminID uint64 `json:"adminId" binding:"required"`
}

// BookmarkUserUpdateRequest request để cập nhật bookmark user
type BookmarkUserUpdateRequest struct {
	ProfileIDs []uint64 `json:"profileIds" binding:"required"`
	Bookmark   bool     `json:"bookmark" binding:"required"`
}

// BookmarkUserResponse response cho bookmark user
type BookmarkUserResponse struct {
	ID        uint64           `json:"id"`
	AdminID   uint64           `json:"adminId"`
	UserID    uint64           `json:"userId"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
	User      *models.UserItem `json:"user,omitempty"`
}

// BookmarkUserListResponse response cho danh sách bookmark user
type BookmarkUserListResponse struct {
	_dto.Pagable
	Data  []BookmarkUserResponse `json:"data"`
	Total int64                  `json:"total"`
}

// BookmarkUsersRequest request để lấy danh sách user có bookmark
type BookmarkUsersRequest struct {
	_dto.Pagable
	Search    string `json:"search" form:"search"`
	RoleType  string `json:"roleType" form:"roleType"`
	StartDate string `json:"startDate" form:"startDate"`
	EndDate   string `json:"endDate" form:"endDate"`
}
