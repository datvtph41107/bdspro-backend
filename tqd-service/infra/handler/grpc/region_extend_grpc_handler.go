package handler_grpc

import (
	"context"
	"encoding/json"
	"log"

	_dto "common/domain/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

type RegionExtendGrpcHandler struct {
	tqdpb.UnimplementedRegionExtendServiceServer
	extendUsecase usecase.RegionExtendUsecase
}

func NewRegionExtendGrpcHandler(extendUsecase usecase.RegionExtendUsecase) *RegionExtendGrpcHandler {
	return &RegionExtendGrpcHandler{extendUsecase: extendUsecase}
}

func structToRawJSON(s *structpb.Struct) ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.Marshal(s.AsMap())
}

func (h *RegionExtendGrpcHandler) CreateRegionExtend(ctx context.Context, req *tqdpb.CreateRegionExtendRequest) (*tqdpb.RegionExtendResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layerId is required")
	}
	if len(req.MergeSourceIds) == 0 && len(req.Geometry) == 0 {
		return nil, status.Error(codes.InvalidArgument, "geometry is required when mergeSourceIds is empty")
	}
	props, err := structToRawJSON(req.OriginalProperties)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "originalProperties is invalid")
	}

	rec, err := h.extendUsecase.Create(ctx, usecase.RegionExtendCreateInput{
		LayerID:            req.LayerId,
		Name:               req.Name,
		DisplayName:        req.DisplayName,
		Description:        req.Description,
		LabelID:            req.LabelId,
		Geometry:           req.Geometry,
		LegalDoc:           req.LegalDoc,
		PlanningName:       req.PlanningName,
		OriginalProperties: props,
		SourceFile:         req.SourceFile,
		ImportBatchID:      req.ImportBatchId,
		Status:             req.Status,
		MergeSourceIDs:     req.MergeSourceIds,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toRegionExtendProto(rec), nil
}

func (h *RegionExtendGrpcHandler) UpdateRegionExtend(ctx context.Context, req *tqdpb.UpdateRegionExtendRequest) (*tqdpb.RegionExtendResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	props, err := structToRawJSON(req.OriginalProperties)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "originalProperties is invalid")
	}

	rec, err := h.extendUsecase.Update(ctx, usecase.RegionExtendUpdatePatch{
		ID:                 req.Id,
		LayerID:            req.LayerId,
		Name:               req.Name,
		DisplayName:        req.DisplayName,
		Description:        req.Description,
		LabelID:            req.LabelId,
		Geometry:           req.Geometry,
		LegalDoc:           req.LegalDoc,
		PlanningName:       req.PlanningName,
		OriginalProperties: props,
		SourceFile:         req.SourceFile,
		ImportBatchID:      req.ImportBatchId,
		Status:             req.Status,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return toRegionExtendProto(rec), nil
}

func (h *RegionExtendGrpcHandler) DeleteRegionExtend(ctx context.Context, req *tqdpb.DeleteRegionExtendRequest) (*emptypb.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := h.extendUsecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *RegionExtendGrpcHandler) ListRegionExtends(ctx context.Context, req *tqdpb.ListRegionExtendsRequest) (*tqdpb.ListRegionExtendsResponse, error) {
	if req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	pagable := &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	}

	log.Printf("[ListRegionExtends] LayerID: %d geomType: %v", req.LayerId, req.Size)

	records, total, err := h.extendUsecase.List(ctx, req.LayerId, req.GeomType, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.RegionExtendResponse, 0, len(records))
	for i := range records {
		data = append(data, toRegionExtendProto(&records[i]))
	}

	return &tqdpb.ListRegionExtendsResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

func (h *RegionExtendGrpcHandler) GetRegionExtend(ctx context.Context, req *tqdpb.GetRegionExtendRequest) (*tqdpb.RegionExtendResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log.Printf("[GetRegionExtend] ID: %d", req.Id)

	rec, err := h.extendUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return toRegionExtendProto(rec), nil
}

func toRegionExtendProto(r *qh_domain.QHRegionExtend) *tqdpb.RegionExtendResponse {
	resp := &tqdpb.RegionExtendResponse{
		Id:            r.ID,
		LayerId:       r.LayerID,
		Name:          r.Name,
		DisplayName:   r.DisplayName,
		Description:   r.Description,
		GeomType:      r.GeomType,
		LegalDoc:      r.LegalDoc,
		PlanningName:  r.PlanningName,
		SourceFile:    r.SourceFile,
		ImportBatchId: r.ImportBatchID,
		Status:        uint32(r.Status),
		Version:       uint32(r.Version),
		IsLatest:      r.IsLatest,
		CreatedAt:     r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if r.LabelID != nil {
		resp.LabelId = *r.LabelID
	}

	if len(r.Geometry.Raw) > 0 {
		// resp.GeometryGeoJson = string(r.Geometry.Raw)
		resp.Geometry = r.Geometry.Raw
	}

	if len(r.OriginalProperties) > 0 {
		var m map[string]interface{}
		if err := json.Unmarshal(r.OriginalProperties, &m); err == nil {
			if s, err := structpb.NewStruct(m); err == nil {
				resp.OriginalProperties = s
			}
		}
	}

	return resp
}
