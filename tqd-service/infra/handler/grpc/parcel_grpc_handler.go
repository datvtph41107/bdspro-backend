package handler_grpc

import (
	"context"
	"log"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParcelGrpcHandler struct {
	tqdpb.UnimplementedParcelServiceServer
	usecase usecase.ParcelUsecase
	mapper  *mapper.ParcelMapper
}

func NewParcelGrpcHandler(
	usecase usecase.ParcelUsecase,
	mapper *mapper.ParcelMapper,
) *ParcelGrpcHandler {
	return &ParcelGrpcHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

func (h *ParcelGrpcHandler) CreateBatchParcels(ctx context.Context, req *tqdpb.ParcelListDTO) (*tqdpb.ParcelListDTO, error) {

	// Convert proto to request DTO
	createReq := h.mapper.ProtoToDomains(req)

	// Call usecase
	result, err := h.usecase.CreateBatch(ctx, &dto.ParcelList{
		Data: createReq,
	})
	if err != nil {
		log.Printf("error creating parcel: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create parcel: %v", err)
	}

	return &tqdpb.ParcelListDTO{
		Data: h.mapper.ToProtoFromDomainList(result),
	}, nil
}

func (h *ParcelGrpcHandler) FindParcelByLocation(ctx context.Context, req *tqdpb.FindParcelByLocationRequest) (*tqdpb.Parcel, error) {

	// Call usecase
	result, err := h.usecase.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		log.Printf("error finding parcel by location: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to find parcel by location: %v", err)
	}

	return h.mapper.DomainToProto(result), nil
}
