package handler_grpc

import (
	_dto "common/domain/dto"
	"context"
	"log"

	sharepb "pb/types/shared"
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

func (h *ParcelGrpcHandler) FindParcelsByPolygon(
	ctx context.Context,
	req *tqdpb.FindParcelsByPolygonRequest,
) (*tqdpb.FindParcelsByPolygonResponse, error) {
	// Ưu tiên dùng coordinates, nếu không có thì dùng geoJson (có thể bỏ)
	var points []dto.Point
	if len(req.Coordinates) > 0 {
		points = make([]dto.Point, len(req.Coordinates))
		for i, pt := range req.Coordinates {
			points[i] = dto.Point{Lat: pt.Lat, Lng: pt.Lng}
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "coordinates are required")
	}

	// Gọi usecase
	parcels, total, err := h.usecase.FindParcelsByPolygon(
		ctx,
		points,
		req.Intersect,
		req.MaxAreaKm2,
		int(req.Page),
		int(req.PageSize),
	)
	if err != nil {
		log.Printf("FindParcelsByPolygon error: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Chuyển domain -> proto
	data := make([]*tqdpb.Parcel, len(parcels))
	for i, p := range parcels {
		data[i] = h.mapper.DomainToProto(&p)
	}

	return &tqdpb.FindParcelsByPolygonResponse{
		Data:     data,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (h *ParcelGrpcHandler) SearchTxtClientParcels(ctx context.Context, req *sharepb.RequestV3Proto) (*tqdpb.ListParcelResponse, error) {
	searchText := req.Text

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}

	parcels, total, err := h.usecase.Search(ctx, searchText, pagable)

	if err != nil {
		return nil, err
	}

	results := make([]*tqdpb.ParcelInfo, len(parcels))

	for i := 0; i < len(parcels); i++ {
		p := parcels[i]

		// mapNumber, _ := p.MapNumber
		// landNumber, _ := strconv.ParseUint(p.LandNumber, 10, 32)

		results[i] = &tqdpb.ParcelInfo{
			Id:         p.ID,
			AddressStr: p.AddressStr,
			MapNumber:  p.MapNumber,
			LandNumber: p.LandNumber,
			Lat:        p.Latitude,
			Lng:        p.Longitude,
			Facade:     uint32(p.Facade),
			TotalArea:  p.TotalArea,

			// Direction (tuỳ bạn define enum hay item)
			// Direction: convertDirection(p.Direction),

			// Shape
			Shape: func() *sharepb.ItemV3Proto {
				if p.Shape == nil {
					return nil
				}
				return &sharepb.ItemV3Proto{
					Id:   p.ShapeId,
					Name: p.Shape.Name, // giả sử có field Name
				}
			}(),

			// LandType
			LandType: func() *sharepb.ItemV3Proto {
				if p.LandType == nil {
					return nil
				}
				return &sharepb.ItemV3Proto{
					Id:   *p.LandTypeId,
					Name: p.LandType.Name,
				}
			}(),
		}
	}

	return &tqdpb.ListParcelResponse{
		Data:  results,
		Total: uint32(total),
	}, nil
}

func (h *ParcelGrpcHandler) PreviewParcelsById(
	ctx context.Context,
	req *sharepb.RequestV3Proto,
) (*tqdpb.ParcelInfo, error) {
	result, err := h.usecase.GetParcelDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.mapper.ParcelInfoToProto(result), nil
}
