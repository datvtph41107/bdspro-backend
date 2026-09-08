package quota

import (
	"context"
	"testing"
	"time"

	commonmetering "common/metering"
	"github.com/stretchr/testify/require"
	"tqd/internal/access"
	quotadomain "tqd/internal/domain/quota"
)

type adminUsageRepositoryStub struct{ page AdminUsagePage }

func (r adminUsageRepositoryStub) List(context.Context, AdminUsageQuery) (AdminUsagePage, error) {
	return r.page, nil
}

type adminRuntimeUsageStub struct {
	usage Usage
	err   error
}

func (r adminRuntimeUsageStub) GetUsage(context.Context, access.Subject, commonmetering.Code, time.Time, time.Time) (Usage, error) {
	return r.usage, r.err
}

func TestListSeparatesDurableAndRuntimeUsageAndFlagsDrift(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	service := NewAdminUsageService(adminUsageRepositoryStub{page: AdminUsagePage{Usage: []quotadomain.AdminUsageProjection{{
		SubjectType: "profile", SubjectID: "42", MeterCode: "workspace.report_generation.accepted",
		DurableUsed: 3, LimitKnown: true, LimitSnapshot: 10, PeriodStart: start, PeriodEnd: start.AddDate(0, 1, 0),
	}}}}, adminRuntimeUsageStub{usage: Usage{Used: 2, Reserved: 1}})
	page, err := service.List(context.Background(), AdminUsageQuery{})
	require.NoError(t, err)
	require.Len(t, page.Usage, 1)
	item := page.Usage[0]
	require.True(t, item.RuntimeAvailable)
	require.Equal(t, int64(2), item.RuntimeUsed)
	require.Equal(t, int64(1), item.RuntimeReserved)
	require.Equal(t, int64(7), item.Remaining)
	require.Equal(t, "drift", item.Reconciliation)
}

func TestListExposesRuntimeUnavailable(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	service := NewAdminUsageService(adminUsageRepositoryStub{page: AdminUsagePage{Usage: []quotadomain.AdminUsageProjection{{
		SubjectType: "profile", SubjectID: "42", MeterCode: "meter", PeriodStart: start, PeriodEnd: start.Add(time.Hour),
	}}}}, nil)
	page, err := service.List(context.Background(), AdminUsageQuery{})
	require.NoError(t, err)
	require.Len(t, page.Usage, 1)
	require.False(t, page.Usage[0].RuntimeAvailable)
	require.Equal(t, "runtime_unavailable", page.Usage[0].Reconciliation)
}
