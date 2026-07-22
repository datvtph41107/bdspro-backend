package handler

import (
	"context"
	"encoding/json"
	"errors"

	"crm/internal/dto"
	"crm/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
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

func (h *SeoDomainHandler) ImportMeasurement(ctx context.Context, input *structpb.Struct) (*structpb.Struct, error) {
	var req dto.SeoMeasurementImportRequest
	if err := decodeSeoM8Struct(input, &req); err != nil {
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

func (h *SeoDomainHandler) GetSourceHealth(ctx context.Context, _ *structpb.Struct) (*structpb.Struct, error) {
	rows, err := h.seoMeasurementUsecase.GetSourceHealth(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": len(rows)})
}

func (h *SeoDomainHandler) ListImportRuns(ctx context.Context, input *structpb.Struct) (*structpb.Struct, error) {
	var req dto.SeoMeasurementImportRunListRequest
	if err := decodeSeoM8Struct(input, &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	rows, total, err := h.seoMeasurementUsecase.ListImportRuns(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": total})
}

func (h *SeoDomainHandler) RunScoring(ctx context.Context, input *structpb.Struct) (*structpb.Struct, error) {
	var req dto.SeoMeasurementScoreRequest
	if err := decodeSeoM8Struct(input, &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	result, err := h.seoMeasurementUsecase.RunScoring(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(result)
}

func (h *SeoDomainHandler) ListPageScores(ctx context.Context, input *structpb.Struct) (*structpb.Struct, error) {
	var req dto.SeoPageScoreListRequest
	if err := decodeSeoM8Struct(input, &req); err != nil {
		return nil, seoM8Invalid(err)
	}
	rows, total, err := h.seoMeasurementUsecase.ListPageScores(ctx, &req)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return encodeSeoM8Struct(map[string]interface{}{"data": rows, "total": total})
}
