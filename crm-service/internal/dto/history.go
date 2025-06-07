package dto

import (
	_dto "common/dto"
	"crm/internal/enums"
	"time"
)

type HistorySearchDTO struct {
	_dto.Page
	UserID     uint64                `form:"userId"`
	FromDate   *time.Time            `form:"fromDate"`
	ToDate     *time.Time            `form:"toDate"`
	ActionType *enums.EHistory       `form:"actionType"`
	TargetType *enums.ETargetHistory `form:"targetType"`
	TargetId   *uint64               `form:"targetId"`
}
