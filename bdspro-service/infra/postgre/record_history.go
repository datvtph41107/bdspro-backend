package postgre

import (
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind:bdspro/internal/repo.RecordHistoryRepo
type PostgreReportHistory struct {
}

func NewPostgreReportHistory(db *gorm.DB) *PostgreReportHistory {
	return &PostgreReportHistory{}
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
	return nil
}
