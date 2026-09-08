package application

import (
	commonoperation "common/operation"
	"context"
	"tqd/internal/access"

	"tqd/internal/usecase/quota"
)

// SourceStore đọc dữ liệu cần để dựng Report.
//
// Đây là read dependency.
// Nó không sở hữu Report acceptance.
type SourceStore interface {
	FindParcel(
		ctx context.Context,
		parcelID uint64,
	) (Parcel, bool, error)

	FindRegion(
		ctx context.Context,
		regionID uint64,
	) (Region, bool, error)
}

// AcceptanceStore là durable boundary của create-report.
//
// FindReportByCommand phục vụ durable replay.
// Accept phải commit Report + optional Usage + Job trong cùng transaction.
// Interface mô tả đúng business checkpoint; nó không expose PostgreSQL mechanics.
type AcceptanceStore interface {
	FindReportByCommand(
		ctx context.Context,
		userID uint64,
		commandKey string,
	) (Report, bool, error)

	AcceptReport(
		ctx context.Context,
		input Acceptance,
	) (AcceptanceResult, error)
}

// AccessClient là dependency mà Report cần từ Access.
//
// Report không đọc Plan/Subscription trực tiếp.
type AccessClient interface {
	GetAccess(
		ctx context.Context,
		operation commonoperation.Code,
	) (access.Result, error)
}

// Quota là phần behavior Quota mà Report create-flow cần.
//
// Ownership:
//
//	consumer = Report
//
// Implementation:
//
//	quota.Service hiện tại thỏa interface này.
//
// Interface này cố ý không expose:
//
//	GetUsage
//	FindExpired
//	CompareAndSetUsed
//
// vì Report không có quyền repair quota projection.
type Quota interface {
	ReserveQuota(
		ctx context.Context,
		input quota.ReserveInput,
	) (quota.Reservation, error)

	CommitQuota(
		ctx context.Context,
		reservation quota.Reservation,
	) (quota.Reservation, error)

	CancelQuota(
		ctx context.Context,
		reservation quota.Reservation,
	) (quota.Reservation, error)
}
