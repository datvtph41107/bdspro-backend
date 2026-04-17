package handler_grpc

import (
	"context"
	"tqd/infra/mapper"
	"tqd/internal/usecase"

	pb "pb/types/tqd"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *OneHouseGrpcHandler) GetOneHouseDetail(ctx context.Context, req *pb.GetOneHouseDetailRequest) (*pb.OneHouseDetailResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ohDomain, err := h.oneHouseUsecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	protoOneHouse, err := mapper.ToProtoOneHouse(ohDomain)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &pb.OneHouseDetailResponse{
		Id:  ohDomain.ID,
		Lat: ohDomain.Lat,
		Lng: ohDomain.Lng,
		Res: &pb.OneHouseRes{
			Data: protoOneHouse,
			Meta: &pb.Meta{Code: "200"},
		},
	}
	return res, nil
}
