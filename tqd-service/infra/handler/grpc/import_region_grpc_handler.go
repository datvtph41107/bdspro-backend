package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"tqd/internal"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status" // Đây là package status của gRPC

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"
)

// ImportGrpcHandler xử lý gRPC cho ImportService
type ImportGrpcHandler struct {
	tqdpb.UnimplementedImportServiceServer
	importUsecase usecase.ImportUsecase
	importMapper  *mapper.ImportMapper
}

func NewImportGrpcHandler(importUsecase usecase.ImportUsecase) *ImportGrpcHandler {
	return &ImportGrpcHandler{
		importUsecase: importUsecase,
		importMapper:  mapper.NewImportMapper(),
	}
}

// PreviewImport - Phân tích file GeoJSON trước khi import
func (h *ImportGrpcHandler) PreviewImport(ctx context.Context, req *tqdpb.PreviewImportRequest) (*tqdpb.PreviewImportResponse, error) {
	if req.GeoJson == "" {
		return nil, status.Error(codes.InvalidArgument, "geoJson is required")
	}
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	slog.InfoContext(ctx, fmt.Sprintf("[PreviewImport] LayerID: %d, GeoJSON length: %d", req.LayerId, len(req.GeoJson)))

	result, err := h.importUsecase.Preview(ctx, req.GeoJson, req.LayerId)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[PreviewImport] Error: %v", err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return h.importMapper.ToProtoPreviewResponse(result), nil
}

// ImportGeoJson - Upload file JSON/GeoJSON/NDJSON; trả processing, xử lý nền theo chunk
func (h *ImportGrpcHandler) ImportGeoJson(ctx context.Context, req *tqdpb.ImportGeoJsonRequest) (*tqdpb.ImportGeoJsonResponse, error) {
	slog.InfoContext(ctx, strings.TrimSuffix(fmt.Sprintln("ImportGeoJson", req), "\n"))
	if len(req.GetFileContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file_content is required")
	}
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(service.Unauthenticated)
	}
	sourceFileName := req.GetSourceFileName()

	importReq := &dto.ImportFileRequest{
		FileContent:      req.GetFileContent(),
		FileFormat:       req.FileFormat,
		LayerID:          req.LayerId,
		LabelField:       req.LabelField,
		LabelMappings:    req.LabelMappings,
		DefaultLabelID:   req.DefaultLabelId,
		SourceFileName:   sourceFileName,
		UserID:           userID,
		ValidateGeometry: req.ValidateGeometry,
		AutoFixGeometry:  req.AutoFixGeometry,
		SkipInvalid:      req.SkipInvalid,
	}

	result, err := h.importUsecase.EnqueueImportFromFile(ctx, importReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[ImportGeoJson] Error: %v", err))
		return nil, importGRPCError(err)
	}
	slog.InfoContext(ctx, fmt.Sprintf("[ImportGeoJson] Enqueued batchId=%s status=%s", result.BatchID, result.Status))

	return h.importMapper.ToProtoImportEnqueueResponse(result), nil
}

// GetImportBatchStatus - Kiểm tra trạng thái batch import
func (h *ImportGrpcHandler) GetImportBatchStatus(ctx context.Context, req *tqdpb.GetImportBatchStatusRequest) (*tqdpb.ImportBatchStatusResponse, error) {
	if req.BatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "batchId is required")
	}
	slog.InfoContext(ctx, fmt.Sprintf("[GetImportBatchStatus] BatchID: %s", req.BatchId))

	batchStatus, err := h.importUsecase.GetBatchStatus(ctx, req.BatchId) // ĐỔI TÊN BIẾN từ status -> batchStatus
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error()) // DÙNG status.Error của gRPC
	}

	return h.importMapper.ToProtoBatchStatusResponse(batchStatus), nil
}

// RollbackImport - Rollback batch import
func (h *ImportGrpcHandler) RollbackImport(ctx context.Context, req *tqdpb.RollbackImportRequest) (*tqdpb.RollbackImportResponse, error) {
	if req.BatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "batchId is required")
	}
	slog.InfoContext(ctx, fmt.Sprintf("[RollbackImport] BatchID: %s", req.BatchId))

	deletedCount, err := h.importUsecase.Rollback(ctx, req.BatchId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.RollbackImportResponse{
		Success:      true,
		Message:      "Rollback completed successfully",
		DeletedCount: int32(deletedCount),
	}, nil
}

// ListImportRegionErrors - Danh sách bản ghi lỗi import (qh_region_import_error_logs) theo layer (GET /v2/tqd/admin/layers/{layerId}/import/errors)
func (h *ImportGrpcHandler) ListImportRegionErrors(ctx context.Context, req *tqdpb.ListImportRegionErrorsRequest) (*tqdpb.ListImportRegionErrorsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	rows, total, err := h.importUsecase.ListImportRegionErrors(ctx, req.LayerId, pagable)
	if err != nil {
		return nil, importGRPCError(err)
	}

	return h.importMapper.ToProtoImportErrors(rows, total, int32(pagable.GetPage()), int32(pagable.GetSize())), nil
}

// RetryImportError - Retry insert từng region trong payload lỗi (POST /v2/tqd/admin/import/errors/{errorId}/retry)
func (h *ImportGrpcHandler) RetryImportError(ctx context.Context, req *tqdpb.RetryImportErrorRequest) (*tqdpb.RetryImportErrorResponse, error) {
	if req.ErrorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "errorId is required")
	}
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(service.Unauthenticated)
	}
	res, err := h.importUsecase.RetryImportError(ctx, req.ErrorId)
	if err != nil {
		return nil, importGRPCError(err)
	}
	return h.importMapper.ToProtoRetryImportErrorResponse(res), nil
}

func importGRPCError(err error) error {
	return _errors.ToGRPC(err)
}
