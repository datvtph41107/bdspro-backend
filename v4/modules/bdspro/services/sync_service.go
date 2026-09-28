package services

import (
	"context"

	"userv4/modules/bdspro/dto"
)

// SyncService owns the app-wide sync timestamp aggregator (GET /v2/bdspro/ts).
type SyncService interface {
	GetAppTimestamps(ctx context.Context) (*dto.AppTimestampsResponse, error)
}
