package mapper

import (
	"google.golang.org/protobuf/types/known/structpb"

	_utils "common/utils"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

type ImportMapper struct{}

func NewImportMapper() *ImportMapper {
	return &ImportMapper{}
}

// ToProtoPreviewResponse - Convert preview result to proto response
func (m *ImportMapper) ToProtoPreviewResponse(result *dto.PreviewResult) *tqdpb.PreviewImportResponse {
	if result == nil {
		return nil
	}

	resp := &tqdpb.PreviewImportResponse{
		TotalFeatures:   int32(result.TotalFeatures),
		AvailableFields: result.AvailableFields,
	}

	// Field stats
	for _, stat := range result.FieldStats {
		resp.FieldStats = append(resp.FieldStats, &tqdpb.PreviewImportResponse_FieldStat{
			FieldName:      stat.FieldName,
			Coverage:       int32(stat.Coverage),
			DistinctValues: stat.DistinctValues,
		})
	}

	// Samples
	for _, sample := range result.Samples {
		props, _ := structpb.NewStruct(sample.Properties)
		resp.Samples = append(resp.Samples, &tqdpb.PreviewImportResponse_FeatureSample{
			Index:        int32(sample.Index),
			GeometryType: sample.GeometryType,
			Properties:   props,
		})
	}

	return resp
}

// ToProtoImportEnqueueResponse — phản hồi sau upload (processing)
func (m *ImportMapper) ToProtoImportEnqueueResponse(result *dto.ImportEnqueueResult) *tqdpb.ImportGeoJsonResponse {
	if result == nil {
		return nil
	}
	return &tqdpb.ImportGeoJsonResponse{
		Code:    result.Code,
		BatchId: result.BatchID,
		Status:  result.Status,
		Message: result.Message,
	}
}

// ToProtoBatchStatusResponse - Convert batch status to proto response
func (m *ImportMapper) ToProtoBatchStatusResponse(status *dto.BatchStatus) *tqdpb.ImportBatchStatusResponse {
	if status == nil {
		return nil
	}

	return &tqdpb.ImportBatchStatusResponse{
		BatchId:       status.BatchID,
		Status:        status.Status,
		TotalFeatures: int32(status.TotalFeatures),
		SuccessCount:  int32(status.SuccessCount),
		FailedCount:   int32(status.FailedCount),
		CreatedAt:     status.CreatedAt,
		CompletedAt:   status.CompletedAt,
		ErrorMessages: status.ErrorMessages,
	}
}

func (m *ImportMapper) ToProtoImportErrors(rows []qh_domain.QHRegionImportErrorLog, total int64, page, pageSize int32) *tqdpb.ListImportRegionErrorsResponse {
	resp := &tqdpb.ListImportRegionErrorsResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	for i := range rows {
		resp.Data = append(resp.Data, toProtoImportError(&rows[i]))
	}
	return resp
}

func toProtoImportError(row *qh_domain.QHRegionImportErrorLog) *tqdpb.ImportError {
	if row == nil {
		return nil
	}
	ca := row.CreatedAt
	ua := row.UpdatedAt
	return &tqdpb.ImportError{
		Id:             row.ID,
		ImportBatchId:  row.ImportBatchID,
		LayerId:        row.LayerID,
		ErrorMessage:   row.ErrorMessage,
		RegionsPayload: string(row.RegionsPayload),
		BatchSize:      int32(row.BatchSize),
		RetryCount:     int32(row.RetryCount),
		Status:         string(row.Status),
		CreatedAt:      _utils.FormatTimeToString(&ca),
		UpdatedAt:      _utils.FormatTimeToString(&ua),
	}
}

// ToProtoRetryImportErrorResponse maps retry result to proto
func (m *ImportMapper) ToProtoRetryImportErrorResponse(r *dto.RetryImportErrorResult) *tqdpb.RetryImportErrorResponse {
	if r == nil {
		return nil
	}
	return &tqdpb.RetryImportErrorResponse{
		Success:      r.Success,
		Message:      r.Message,
		TotalItems:   int32(r.TotalItems),
		SuccessCount: int32(r.SuccessCount),
		FailedCount:  int32(r.FailedCount),
	}
}
