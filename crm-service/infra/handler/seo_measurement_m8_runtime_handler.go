package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"crm/internal/dto"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	crmpb "pb/types/crm"
)

func decodeSeoM8Struct(input *structpb.Struct, output interface{}) error {
	if input == nil {
		return nil
	}
	payload, err := json.Marshal(input.AsMap())
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, output)
}

func encodeSeoM8Struct(value interface{}) (*structpb.Struct, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var object map[string]interface{}
	if err := json.Unmarshal(payload, &object); err != nil {
		return nil, err
	}
	return structpb.NewStruct(object)
}

func seoM8Invalid(err error) error {
	if err == nil {
		return nil
	}
	return status.Error(codes.InvalidArgument, err.Error())
}

func seoM8QueryPayload(input *crmpb.SeoMeasurementQueryRequest) *structpb.Struct {
	payload := map[string]interface{}{}
	if input == nil {
		result, _ := structpb.NewStruct(payload)
		return result
	}
	for key, value := range map[string]string{
		"source":       input.GetSource(),
		"status":       input.GetStatus(),
		"grade":        input.GetGrade(),
		"canonicalUrl": input.GetCanonicalUrl(),
		"range":        input.GetRange(),
	} {
		if strings.TrimSpace(value) != "" {
			payload[key] = strings.TrimSpace(value)
		}
	}
	if input.GetPage() > 0 {
		payload["page"] = float64(input.GetPage())
	}
	if input.GetSize() > 0 {
		payload["size"] = float64(input.GetSize())
	}
	if input.GetLimit() > 0 {
		payload["limit"] = float64(input.GetLimit())
	}
	result, _ := structpb.NewStruct(payload)
	return result
}

func (h *SeoDomainHandler) ImportMeasurement(ctx context.Context, input *crmpb.SeoMeasurementImportRequest) (*structpb.Struct, error) {
	payload := map[string]interface{}{}
	if input != nil && input.GetPayload() != nil {
		for key, value := range input.GetPayload().AsMap() {
			payload[key] = value
		}
	}
	if input != nil {
		if source := strings.TrimSpace(input.GetSource()); source != "" {
			payload["source"] = source
		}
	}
	requestStruct, err := structpb.NewStruct(payload)
	if err != nil {
		return nil, seoM8Invalid(err)
	}
	var req dto.SeoMeasurementImportRequest
	if err := decodeSeoM8Struct(requestStruct, &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	result, err := h.seoMeasurementUsecase.ImportMeasurement(ctx, &req)
	if err != nil {
		var validationErr *usecase.SeoMeasurementValidationError
		if errors.As(err, &validationErr) {
			return nil, seoM8Invalid(err)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(result)
}

func (h *SeoDomainHandler) GetSourceHealth(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	rows, err := h.seoMeasurementUsecase.GetSourceHealth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": len(rows)})
}

func (h *SeoDomainHandler) ListImportRuns(ctx context.Context, input *crmpb.SeoMeasurementQueryRequest) (*structpb.Struct, error) {
	var req dto.SeoMeasurementImportRunListRequest
	if err := decodeSeoM8Struct(seoM8QueryPayload(input), &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	rows, total, err := h.seoMeasurementUsecase.ListImportRuns(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": total})
}

func (h *SeoDomainHandler) RunScoring(ctx context.Context, input *crmpb.SeoMeasurementScoreRequest) (*structpb.Struct, error) {
	var req dto.SeoMeasurementScoreRequest
	var payload *structpb.Struct
	if input != nil {
		payload = input.GetPayload()
	}
	if err := decodeSeoM8Struct(payload, &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	result, err := h.seoMeasurementUsecase.RunScoring(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(result)
}

func (h *SeoDomainHandler) ListPageScores(ctx context.Context, input *crmpb.SeoMeasurementQueryRequest) (*structpb.Struct, error) {
	var req dto.SeoPageScoreListRequest
	if err := decodeSeoM8Struct(seoM8QueryPayload(input), &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	rows, total, err := h.seoMeasurementUsecase.ListPageScores(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": total})
}
