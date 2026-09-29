package usageprojection

import (
	commonmetering "common/metering"
	"context"
	"time"
	"tqd/internal/access"
	"tqd/internal/usecase/quota"
	"tqd/internal/usecase/usage"
)

/**
 * QuotaStore là phần runtime usage mà checker cần đọc/sửa.
 */
type QuotaStore interface {
	GetUsage(
		ctx context.Context,
		subject access.Subject,
		meterCode commonmetering.Code,
		periodStart time.Time,
		periodEnd time.Time,
	) (quota.Usage, error)

	CompareAndSetUsed(
		ctx context.Context,
		subject access.Subject,
		meterCode commonmetering.Code,
		periodStart time.Time,
		periodEnd time.Time,
		expected int64,
		used int64,
	) (bool, error)
}

/**
 * Result cho biết saved usage và Redis runtime có lệch nhau hay không.
 */
type Result struct {
	SavedUsed   int64
	RuntimeUsed int64
	Reserved    int64
	Fixed       bool
}

/**
 * Service so sánh usage đã lưu với Redis và sửa field used khi cần.
 *
 * Reserved không được suy lại từ DB vì nó là state đang chạy.
 */
type Service struct {
	quotaStore QuotaStore
	usageStore usage.Store
}

func NewService(quotaStore QuotaStore, usageStore usage.Store) *Service {
	return &Service{quotaStore: quotaStore, usageStore: usageStore}
}

func (s *Service) CheckAndFix(
	ctx context.Context,
	subject access.Subject,
	meterCode commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (Result, error) {
	var result Result
	for attempt := 0; attempt < 3; attempt++ {
		saved, err := s.usageStore.SumUsage(ctx, subject, meterCode, periodStart, periodEnd)
		if err != nil {
			return Result{}, err
		}

		runtime, err := s.quotaStore.GetUsage(ctx, subject, meterCode, periodStart, periodEnd)
		if err != nil {
			return Result{}, err
		}

		result = Result{
			SavedUsed:   saved,
			RuntimeUsed: runtime.Used,
			Reserved:    runtime.Reserved,
		}
		if saved == runtime.Used {
			return result, nil
		}

		updated, err := s.quotaStore.CompareAndSetUsed(
			ctx,
			subject,
			meterCode,
			periodStart,
			periodEnd,
			runtime.Used,
			saved,
		)
		if err != nil {
			return Result{}, err
		}
		if updated {
			result.Fixed = true
			return result, nil
		}
	}

	// A concurrent quota commit won every compare-and-set attempt. Leave its
	// newer value intact; the next periodic run will reconcile after the durable
	// usage event is visible.
	return result, nil
}
