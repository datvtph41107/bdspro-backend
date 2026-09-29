package quota

import (
	commonmetering "common/metering"
	commonoperation "common/operation"
	"context"
	"errors"
	"time"
	"tqd/internal/access"
)

const GeneratedReportOperation commonoperation.Code = "workspace.report.generate"

type ProfileAccessReader interface {
	GetAccess(context.Context, commonoperation.Code) (access.Result, error)
}

type DurableUsageReader interface {
	SumUsage(context.Context, access.Subject, commonmetering.Code, time.Time, time.Time) (int64, error)
}

type ProfileUsageProjection struct {
	Access           access.Result
	DurableUsed      int64
	RuntimeUsed      int64
	RuntimeReserved  int64
	Remaining        int64
	RuntimeAvailable bool
	Reconciliation   string
}

type ProfileUsageService struct {
	access  ProfileAccessReader
	durable DurableUsageReader
	runtime RuntimeUsageReader
}

func NewProfileUsageService(accessReader ProfileAccessReader, durable DurableUsageReader, runtime RuntimeUsageReader) *ProfileUsageService {
	return &ProfileUsageService{access: accessReader, durable: durable, runtime: runtime}
}

// GetGeneratedReportUsage projects entitlement, durable accounting and
// runtime admission without collapsing their distinct authorities.
func (s *ProfileUsageService) GetGeneratedReportUsage(ctx context.Context) (ProfileUsageProjection, error) {
	if s == nil || s.access == nil || s.durable == nil {
		return ProfileUsageProjection{}, errors.New("profile usage projection is not configured")
	}
	decision, err := s.access.GetAccess(ctx, GeneratedReportOperation)
	if err != nil {
		return ProfileUsageProjection{}, err
	}
	projection := ProfileUsageProjection{Access: decision, Reconciliation: "not_required"}
	if !decision.Allowed || !decision.UsesQuota() {
		return projection, nil
	}
	durableUsed, err := s.durable.SumUsage(ctx, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil {
		return ProfileUsageProjection{}, err
	}
	projection.DurableUsed = durableUsed
	projection.Remaining = decision.Limit - durableUsed
	if projection.Remaining < 0 {
		projection.Remaining = 0
	}
	projection.Reconciliation = "runtime_unavailable"
	if s.runtime == nil {
		return projection, nil
	}
	runtimeUsage, err := s.runtime.GetUsage(ctx, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd)
	if err != nil {
		return projection, nil
	}
	projection.RuntimeAvailable = true
	projection.RuntimeUsed = runtimeUsage.Used
	projection.RuntimeReserved = runtimeUsage.Reserved
	projection.Remaining = decision.Limit - runtimeUsage.Used - runtimeUsage.Reserved
	if projection.Remaining < 0 {
		projection.Remaining = 0
	}
	projection.Reconciliation = "in_sync"
	if durableUsed != runtimeUsage.Used {
		projection.Reconciliation = "drift"
	}
	return projection, nil
}
