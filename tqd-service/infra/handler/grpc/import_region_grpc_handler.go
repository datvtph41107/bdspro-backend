package handler_grpc

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"log"

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

// ImportGeoJson - Import file GeoJSON vào hệ thống
func (h *ImportGrpcHandler) ImportGeoJson(ctx context.Context, req *tqdpb.ImportGeoJsonRequest) (*tqdpb.ImportGeoJsonResponse, error) {
	if req.GeoJson == "" {
		return nil, status.Error(codes.InvalidArgument, "geoJson is required")
	}
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	// Convert *string sang string
	sourceFileName := ""
	if req.SourceFileName != nil {
		sourceFileName = *req.SourceFileName
	}

	importReq := &dto.ImportRequest{
		GeoJSON:          req.GeoJson,
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

	result, err := h.importUsecase.Import(ctx, importReq)
	if err != nil {
		log.Printf("[ImportGeoJson] Error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("[ImportGeoJson] Completed: success=%d, failed=%d, batchId=%s",
		result.SuccessCount, result.FailedCount, result.BatchID)

	return h.importMapper.ToProtoImportResponse(result), nil
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
