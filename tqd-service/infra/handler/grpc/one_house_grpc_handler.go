package handler_grpc

import (
	"context"
	"fmt"

	_errors "common/errors"

	pb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal"
	"tqd/internal/usecase"
)

type OneHouseGrpcHandler struct {
	pb.UnimplementedOneHouseServiceServer
	oneHouseUsecase usecase.OneHouseUsecase
}

func NewOneHouseGrpcHandler(oneHouseUsecase usecase.OneHouseUsecase) *OneHouseGrpcHandler {
	return &OneHouseGrpcHandler{
		oneHouseUsecase: oneHouseUsecase,
	}
}

func (h *OneHouseGrpcHandler) GetOneHouseDetail(ctx context.Context, req *pb.GetOneHouseDetailRequest) (*pb.OneHouse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(service.OneHouseIDRequired)
	}

	ohDomain, err := h.oneHouseUsecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if ohDomain == nil {
		return nil, _errors.ReturnError(service.OneHouseNotFound)
	}

	protoOneHouse, err := mapper.ToProtoOneHouse(ohDomain)
	if err != nil {
		return nil, fmt.Errorf("map one house response: %w", err)
	}

	return protoOneHouse, nil
}
