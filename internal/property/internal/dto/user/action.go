package user_dto

import "time"

type ActionRequest struct {
	TargetID uint64 `json:"targetId" binding:"required"`
	Action   uint   `json:"action" binding:"required"`
}

type PostBarRequest struct {
	PostID   string     `form:"content"`
	FromDate *time.Time `form:"fromDate"`
	ToDate   *time.Time `form:"toDate"`
	Group    string     `form:"group"` // day, month
}

type DashboardPostResponse struct {
	NumView  int64 `form:"view"`
	NumClick int64 `form:"click"`
}
