// Package errorresponse chuyển kết quả usecase TQD sang transport ở một nơi.
// Chỉ lỗi nghiệp vụ đã khai báo mới được trả chi tiết; lỗi hạ tầng được log cùng request.
package errorresponse

import (
	commonerrors "common/errors"
	"common/requestlog"
	"common/validation"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"runtime"
	"tqd/internal/usecase"
	discovery "tqd/internal/usecase/discovery/application"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC không phân loại bằng message: đổi câu chữ hoặc thêm %w không đổi status.
func GRPC(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	// Caller trỏ tới handler; validation.Source chỉ quy tắc gốc tại usecase.
	_, file, line, _ := runtime.Caller(1)
	method, _ := grpc.Method(ctx)
	log := requestlog.FromContext(ctx, nil).With("handler_source", fmt.Sprintf("%s:%d", file, line), "grpc_method", method)
	if v, ok := validation.As(err); ok {
		log.WarnContext(ctx, "TQD request rejected", "category", "application", "event_name", "tqd.validation.failed", "field", v.Field, "error_reason", v.Reason, "error_source", v.Source)
		result, _ := commonerrors.ValidationStatus(err, "tqd")
		return result
	}
	rules := []struct {
		cause  error
		code   codes.Code
		reason string
	}{
		{usecase.ErrLandUseNotFound, codes.NotFound, "LAND_USE_NOT_FOUND"},
		{usecase.ErrLandUseConflict, codes.AlreadyExists, "LAND_USE_CONFLICT"},
		{usecase.ErrLegendNotFound, codes.NotFound, "LEGEND_NOT_FOUND"},
		{usecase.ErrLegendConflict, codes.AlreadyExists, "LEGEND_CONFLICT"},
		{usecase.ErrWorkspaceItemNotFound, codes.NotFound, "WORKSPACE_ITEM_NOT_FOUND"},
		{usecase.ErrReportNotReadyToShare, codes.FailedPrecondition, "REPORT_NOT_READY_TO_SHARE"},
		{discovery.ErrProvidersUnavailable, codes.Unavailable, "DISCOVERY_PROVIDERS_UNAVAILABLE"},
		{usecase.ErrLayerNotFound, codes.NotFound, "LAYER_NOT_FOUND"},
		{usecase.ErrLayerFamilyNotFound, codes.NotFound, "LAYER_FAMILY_NOT_FOUND"},
		{usecase.ErrAuthorityNotFound, codes.NotFound, "AUTHORITY_NOT_FOUND"},
		{usecase.ErrAuthorityCodeConflict, codes.AlreadyExists, "AUTHORITY_CODE_CONFLICT"},
		{usecase.ErrLabelNotFound, codes.NotFound, "LABEL_NOT_FOUND"},
		{usecase.ErrLabelNameConflict, codes.AlreadyExists, "LABEL_NAME_CONFLICT"},
		{usecase.ErrRegionNotFound, codes.NotFound, "REGION_NOT_FOUND"},
		{usecase.ErrRegionSyncReferenceNotFound, codes.NotFound, "REGION_SYNC_REFERENCE_NOT_FOUND"},
		{usecase.ErrImportBatchNotFound, codes.NotFound, "IMPORT_BATCH_NOT_FOUND"},
		{usecase.ErrImportErrorNotFound, codes.NotFound, "IMPORT_ERROR_NOT_FOUND"},
		{usecase.ErrImportInProgress, codes.FailedPrecondition, "IMPORT_IN_PROGRESS"},
		{usecase.ErrImportNotRetryable, codes.FailedPrecondition, "IMPORT_NOT_RETRYABLE"},
		{usecase.ErrImportRetryConflict, codes.Aborted, "IMPORT_RETRY_CONFLICT"},
		{context.Canceled, codes.Canceled, "REQUEST_CANCELED"},
		{context.DeadlineExceeded, codes.DeadlineExceeded, "REQUEST_DEADLINE_EXCEEDED"},
	}
	for _, rule := range rules {
		if errors.Is(err, rule.cause) {
			log.WarnContext(ctx, "TQD operation rejected", "category", "application", "event_name", "tqd.operation.rejected", "error_reason", rule.reason)
			s := status.New(rule.code, rule.cause.Error())
			detailed, e := s.WithDetails(&errdetails.ErrorInfo{Domain: "tqd", Reason: rule.reason})
			if e != nil {
				return s.Err()
			}
			return detailed.Err()
		}
	}
	// Log cause trước khi che chi tiết; không biến lỗi repository thành lỗi đầu vào.
	log.ErrorContext(ctx, "TQD operation failed", "category", "application", "event_name", "tqd.operation.failed", "error", err.Error())
	return status.Error(codes.Internal, "internal server error")
}
