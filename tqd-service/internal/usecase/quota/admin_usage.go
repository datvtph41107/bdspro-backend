package quota

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	commonmetering "common/metering"
	"tqd/internal/access"
	quotadomain "tqd/internal/domain/quota"
)

const PermissionUsageView = "COMMERCIAL_USAGE_VIEW"

type AdminUsageQuery struct {
	Page      uint32
	PageSize  uint32
	ProfileID uint64
	Operation string
	MeterCode string
}

type AdminUsagePage struct {
	Usage    []quotadomain.AdminUsageProjection
	Total    uint64
	Page     uint32
	PageSize uint32
}

type AdminUsageRepository interface {
	List(context.Context, AdminUsageQuery) (AdminUsagePage, error)
}

type RuntimeUsageReader interface {
	GetUsage(context.Context, access.Subject, commonmetering.Code, time.Time, time.Time) (Usage, error)
}

type AdminUsageService struct {
	repository AdminUsageRepository
	runtime    RuntimeUsageReader
}

func NewAdminUsageService(repository AdminUsageRepository, runtime RuntimeUsageReader) *AdminUsageService {
	return &AdminUsageService{repository: repository, runtime: runtime}
}

func (s *AdminUsageService) List(ctx context.Context, query AdminUsageQuery) (AdminUsagePage, error) {
	if s == nil || s.repository == nil {
		return AdminUsagePage{}, errors.New("admin usage projection is not configured")
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		return AdminUsagePage{}, fmt.Errorf("page size must be between 1 and 100")
	}
	query.Operation = strings.TrimSpace(query.Operation)
	query.MeterCode = strings.TrimSpace(query.MeterCode)
	page, err := s.repository.List(ctx, query)
	if err != nil {
		return AdminUsagePage{}, err
	}
	for i := range page.Usage {
		s.attachRuntime(ctx, &page.Usage[i])
	}
	return page, nil
}

func (s *AdminUsageService) attachRuntime(ctx context.Context, item *quotadomain.AdminUsageProjection) {
	item.Reconciliation = "runtime_unavailable"
	if s.runtime == nil {
		return
	}
	subject := access.Subject{Type: access.SubjectType(item.SubjectType), ID: item.SubjectID}
	meter := commonmetering.Code(item.MeterCode)
	if !subject.IsValid() || !meter.IsValid() || item.PeriodStart.IsZero() || !item.PeriodEnd.After(item.PeriodStart) {
		return
	}
	runtimeUsage, err := s.runtime.GetUsage(ctx, subject, meter, item.PeriodStart, item.PeriodEnd)
	if err != nil {
		return
	}
	item.RuntimeAvailable = true
	item.RuntimeUsed = runtimeUsage.Used
	item.RuntimeReserved = runtimeUsage.Reserved
	item.Reconciliation = "in_sync"
	if item.RuntimeUsed != item.DurableUsed {
		item.Reconciliation = "drift"
	}
	if item.LimitKnown {
		item.Remaining = item.LimitSnapshot - item.RuntimeUsed - item.RuntimeReserved
		if item.Remaining < 0 {
			item.Remaining = 0
		}
	}
}
