package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"log"
	"strings"

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

	log.Printf("[PreviewImport] LayerID: %d, GeoJSON length: %d", req.LayerId, len(req.GeoJson))

	result, err := h.importUsecase.Preview(ctx, req.GeoJson, req.LayerId)
	if err != nil {
		log.Printf("[PreviewImport] Error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return h.importMapper.ToProtoPreviewResponse(result), nil
}

// ImportGeoJson - Upload file JSON/GeoJSON/NDJSON; trả processing, xử lý nền theo chunk
func (h *ImportGrpcHandler) ImportGeoJson(ctx context.Context, req *tqdpb.ImportGeoJsonRequest) (*tqdpb.ImportGeoJsonResponse, error) {
	log.Println("ImportGeoJson", req)
	if len(req.GetFileContent()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "file_content is required")
	}
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
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
		log.Printf("[ImportGeoJson] Error: %v", err)
		if strings.Contains(err.Error(), "already has an import") {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[ImportGeoJson] Enqueued batchId=%s status=%s", result.BatchID, result.Status)

	return h.importMapper.ToProtoImportEnqueueResponse(result), nil
}

// GetImportBatchStatus - Kiểm tra trạng thái batch import
func (h *ImportGrpcHandler) GetImportBatchStatus(ctx context.Context, req *tqdpb.GetImportBatchStatusRequest) (*tqdpb.ImportBatchStatusResponse, error) {
	if req.BatchId == "" {
		return nil, status.Error(codes.InvalidArgument, "batchId is required")
	}

	log.Printf("[GetImportBatchStatus] BatchID: %s", req.BatchId)

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

	log.Printf("[RollbackImport] BatchID: %s", req.BatchId)

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
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	res, err := h.importUsecase.RetryImportError(ctx, req.ErrorId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if strings.Contains(err.Error(), "already in progress") {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		if strings.Contains(err.Error(), "cannot be retried") {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		if strings.Contains(err.Error(), "could not be locked") {
			return nil, status.Error(codes.Aborted, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return h.importMapper.ToProtoRetryImportErrorResponse(res), nil
}
