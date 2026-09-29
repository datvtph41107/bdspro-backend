package quota

import (
	"context"
	"errors"
	"testing"
	"time"

	commonmetering "common/metering"
	commonoperation "common/operation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"tqd/internal/access"
)

type mockProfileAccessReader struct{ mock.Mock }

func (m *mockProfileAccessReader) GetAccess(ctx context.Context, operation commonoperation.Code) (access.Result, error) {
	args := m.Called(ctx, operation)
	result, _ := args.Get(0).(access.Result)
	return result, args.Error(1)
}

type mockDurableUsageReader struct{ mock.Mock }

func (m *mockDurableUsageReader) SumUsage(
	ctx context.Context,
	subject access.Subject,
	meter commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (int64, error) {
	args := m.Called(ctx, subject, meter, periodStart, periodEnd)
	return args.Get(0).(int64), args.Error(1)
}

type mockRuntimeUsageReader struct{ mock.Mock }

func (m *mockRuntimeUsageReader) GetUsage(
	ctx context.Context,
	subject access.Subject,
	meter commonmetering.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (Usage, error) {
	args := m.Called(ctx, subject, meter, periodStart, periodEnd)
	usage, _ := args.Get(0).(Usage)
	return usage, args.Error(1)
}

func meteredProfileAccess() access.Result {
	return access.Result{
		Subject: access.Subject{Type: access.SubjectProfile, ID: "42"}, Operation: GeneratedReportOperation,
		Allowed: true, Limit: 100, Period: access.PeriodSubscriptionCycle,
		PeriodStart: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), PeriodEnd: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		Metering: access.Metering{FeatureCode: "generated_report", MeterCode: "workspace.report_generation.accepted", UnitsPerAction: 1, PolicyVersion: "1"},
	}
}

func TestProfileUsageProjection(t *testing.T) {
	t.Parallel()

	runtimeUnavailable := errors.New("redis unavailable")
	accessUnavailable := errors.New("user service unavailable")
	testCases := []struct {
		name  string
		setup func(*mockProfileAccessReader, *mockDurableUsageReader, *mockRuntimeUsageReader)
		check func(*testing.T, ProfileUsageProjection, error)
	}{
		{
			name: "keeps durable and runtime evidence distinct",
			setup: func(accessReader *mockProfileAccessReader, durable *mockDurableUsageReader, runtime *mockRuntimeUsageReader) {
				decision := meteredProfileAccess()
				accessReader.On("GetAccess", mock.Anything, GeneratedReportOperation).Return(decision, nil).Once()
				durable.On("SumUsage", mock.Anything, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd).Return(int64(42), nil).Once()
				runtime.On("GetUsage", mock.Anything, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd).Return(Usage{Used: 41, Reserved: 1}, nil).Once()
			},
			check: func(t *testing.T, projection ProfileUsageProjection, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(42), projection.DurableUsed)
				require.Equal(t, int64(41), projection.RuntimeUsed)
				require.Equal(t, int64(1), projection.RuntimeReserved)
				require.Equal(t, int64(58), projection.Remaining)
				require.Equal(t, "drift", projection.Reconciliation)
			},
		},
		{
			name: "preserves durable truth when runtime is unavailable",
			setup: func(accessReader *mockProfileAccessReader, durable *mockDurableUsageReader, runtime *mockRuntimeUsageReader) {
				decision := meteredProfileAccess()
				accessReader.On("GetAccess", mock.Anything, GeneratedReportOperation).Return(decision, nil).Once()
				durable.On("SumUsage", mock.Anything, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd).Return(int64(42), nil).Once()
				runtime.On("GetUsage", mock.Anything, decision.Subject, decision.Metering.MeterCode, decision.PeriodStart, decision.PeriodEnd).Return(Usage{}, runtimeUnavailable).Once()
			},
			check: func(t *testing.T, projection ProfileUsageProjection, err error) {
				require.NoError(t, err)
				require.Equal(t, int64(42), projection.DurableUsed)
				require.False(t, projection.RuntimeAvailable)
				require.Equal(t, int64(58), projection.Remaining)
				require.Equal(t, "runtime_unavailable", projection.Reconciliation)
			},
		},
		{
			name: "does not read usage when access is unmetered",
			setup: func(accessReader *mockProfileAccessReader, _ *mockDurableUsageReader, _ *mockRuntimeUsageReader) {
				decision := meteredProfileAccess()
				decision.Allowed = false
				accessReader.On("GetAccess", mock.Anything, GeneratedReportOperation).Return(decision, nil).Once()
			},
			check: func(t *testing.T, projection ProfileUsageProjection, err error) {
				require.NoError(t, err)
				require.False(t, projection.Access.Allowed)
				require.Equal(t, "not_required", projection.Reconciliation)
			},
		},
		{
			name: "propagates entitlement failure",
			setup: func(accessReader *mockProfileAccessReader, _ *mockDurableUsageReader, _ *mockRuntimeUsageReader) {
				accessReader.On("GetAccess", mock.Anything, GeneratedReportOperation).Return(access.Result{}, accessUnavailable).Once()
			},
			check: func(t *testing.T, projection ProfileUsageProjection, err error) {
				require.ErrorIs(t, err, accessUnavailable)
				require.Zero(t, projection)
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			accessReader := new(mockProfileAccessReader)
			durable := new(mockDurableUsageReader)
			runtime := new(mockRuntimeUsageReader)
			accessReader.Test(t)
			durable.Test(t)
			runtime.Test(t)
			t.Cleanup(func() {
				accessReader.AssertExpectations(t)
				durable.AssertExpectations(t)
				runtime.AssertExpectations(t)
			})

			testCase.setup(accessReader, durable, runtime)
			service := NewProfileUsageService(accessReader, durable, runtime)
			projection, err := service.GetGeneratedReportUsage(context.Background())
			testCase.check(t, projection, err)
		})
	}
}
