package handler_grpc

import (
	_errors "common/errors"
	"context"

	pb "pb/types/tqd"
	"tqd/infra/mapper"
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
		return nil, _errors.ReturnError(400, "id is required")
	}

	ohDomain, err := h.oneHouseUsecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}

	if ohDomain == nil {
		return nil, _errors.ReturnError(404, "one_house not found")
	}

	protoOneHouse, err := mapper.ToProtoOneHouse(ohDomain)
	if err != nil {
		return nil, _errors.ReturnError(500, "internal error: "+err.Error())
	}

	return protoOneHouse, nil
}
