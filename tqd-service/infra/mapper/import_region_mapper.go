package mapper

import (
	"google.golang.org/protobuf/types/known/structpb"

	tqdpb "pb/types/tqd"
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

// ToProtoImportResponse - Convert import result to proto response
func (m *ImportMapper) ToProtoImportResponse(result *dto.ImportResult) *tqdpb.ImportGeoJsonResponse {
	if result == nil {
		return nil
	}

	resp := &tqdpb.ImportGeoJsonResponse{
		BatchId:             result.BatchID,
		TotalFeatures:       int32(result.TotalFeatures),
		SuccessCount:        int32(result.SuccessCount),
		FailedCount:         int32(result.FailedCount),
		LabelStatistics:     result.LabelStatistics,
		TotalAreaHa:         result.TotalAreaHa,
		TotalRegionsCreated: int32(result.TotalRegionsCreated),
	}

	for _, f := range result.Failures {
		resp.Failures = append(resp.Failures, &tqdpb.ImportGeoJsonResponse_FailureDetail{
			FeatureIndex: int32(f.FeatureIndex),
			Reason:       f.Reason,
		})
	}

	return resp
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
