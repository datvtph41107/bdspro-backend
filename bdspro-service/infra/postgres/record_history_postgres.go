package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind:bdspro/internal/repo.RecordHistoryRepo
type PostgreReportHistory struct {
	DB *gorm.DB
}

func NewPostgreReportHistory(db *gorm.DB) *PostgreReportHistory {
	return &PostgreReportHistory{
		DB: db,
	}
}

func (r *PostgreReportHistory) CreateHistory(
	c context.Context,
	action enums.EHistoryAction,
	recordID *uint64,
	parentID *uint64,
	childIDs []uint64,
	date string,
	description string,
	performedBy uint64,
) error {
	entity := &domain.RecordHistory{
		ActionType: action,
		RecordID:   recordID,
		// ParentID:    parentID,
		// ChildIDs:    ConvertToPQArray(childIDs),
		// Date:        date,
		// Description: description,
		// PerformedBy: performedBy,
	}
	return r.DB.WithContext(c).Create(entity).Error
}
