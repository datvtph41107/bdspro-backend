package report

import (
	commonaccess "common/access"
	"context"
	"tqd/internal/usage"
)

/**
 * SourceStore đọc parcel/region để dựng nội dung report.
 *
 * Đây chỉ là read path, không thay đổi business state.
 */
type SourceStore interface {
	FindParcel(ctx context.Context, parcelID uint64) (Parcel, bool, error)
	FindRegion(ctx context.Context, regionID uint64) (Region, bool, error)
}

/**
 * ReportStore đọc/lưu generated report theo command key.
 *
 * FindReportByCommand dùng durable state để chặn charge lại một command đã accepted.
 * AcceptReport phải lưu report và usage event trong cùng một DB transaction.
 */
type ReportStore interface {
	FindReportByCommand(
		ctx context.Context,
		userID uint64,
		commandKey string,
	) (Report, bool, error)

	AcceptReport(
		ctx context.Context,
		report Report,
		usageEvent *usage.Event,
	) (saved Report, created bool, finalizeReservation bool, err error)
}

/**
 * AccessClient lấy access hiện tại từ user-service.
 *
 * TQD không đọc plan/subscription trực tiếp và không tự quyết định subject.
 */
type AccessClient interface {
	GetReportAccess(ctx context.Context) (commonaccess.Result, error)
}
