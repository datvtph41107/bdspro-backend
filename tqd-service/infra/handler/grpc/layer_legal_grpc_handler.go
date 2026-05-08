// infra/handler/grpc/layer_legal_grpc_handler.go
package handler_grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	tqdpb "pb/types/tqd"
)

type LayerLegalGrpcHandler struct {
	tqdpb.UnimplementedLayerLegalServiceServer
	layerLegalRepo repo.ILayerLegalRepo
}

func NewLayerLegalGrpcHandler(layerLegalRepo repo.ILayerLegalRepo) *LayerLegalGrpcHandler {
	return &LayerLegalGrpcHandler{
		layerLegalRepo: layerLegalRepo,
	}
}

func (h *LayerLegalGrpcHandler) AddLayerLegalDocs(ctx context.Context, req *tqdpb.AddLayerLegalDocsRequest) (*tqdpb.AddLayerLegalDocsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	if len(req.LegalDocs) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one legal document is required")
	}

	legals := make([]*qh_domain.QHLayerLegal, 0, len(req.LegalDocs))
	for _, input := range req.LegalDocs {
		if input.Name == "" || input.FileUrl == "" {
			continue
		}

		legal := &qh_domain.QHLayerLegal{
			LayerID:  req.LayerId,
			Name:     input.Name,
			FileURL:  input.FileUrl,
			FileType: input.FileType,
		}
		legals = append(legals, legal)
	}

	if len(legals) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no valid legal documents to add")
	}

	// Save to database
	if err := h.layerLegalRepo.CreateBatch(ctx, legals); err != nil {
		return nil, status.Error(codes.Internal, "failed to save legal documents: "+err.Error())
	}

	// Convert response
	legalDocs := make([]*tqdpb.LayerLegalDoc, 0, len(legals))
	for _, legal := range legals {
		legalDocs = append(legalDocs, &tqdpb.LayerLegalDoc{
			Id:        legal.ID,
			LayerId:   legal.LayerID,
			Name:      legal.Name,
			FileUrl:   legal.FileURL,
			FileType:  legal.FileType,
			CreatedAt: legal.CreatedAt.Format(time.RFC3339),
		})
	}

	return &tqdpb.AddLayerLegalDocsResponse{
		Success:   true,
		Message:   "Legal documents added successfully",
		LegalDocs: legalDocs,
	}, nil
}

// GetLayerLegalDocs - Lấy danh sách legal documents của layer
func (h *LayerLegalGrpcHandler) GetLayerLegalDocs(ctx context.Context, req *tqdpb.GetLayerLegalDocsRequest) (*tqdpb.GetLayerLegalDocsResponse, error) {
	if req == nil || req.LayerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "layer_id is required")
	}

	legals, err := h.layerLegalRepo.GetByLayerID(ctx, req.LayerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get legal documents: "+err.Error())
	}

	legalDocs := make([]*tqdpb.LayerLegalDoc, 0, len(legals))
	for _, legal := range legals {
		legalDocs = append(legalDocs, &tqdpb.LayerLegalDoc{
			Id:        legal.ID,
			LayerId:   legal.LayerID,
			Name:      legal.Name,
			FileUrl:   legal.FileURL,
			FileType:  legal.FileType,
			CreatedAt: legal.CreatedAt.Format(time.RFC3339),
		})
	}

	return &tqdpb.GetLayerLegalDocsResponse{
		LayerId:   req.LayerId,
		LegalDocs: legalDocs,
	}, nil
}

// DeleteLayerLegalDoc - Xóa legal document
func (h *LayerLegalGrpcHandler) DeleteLayerLegalDoc(ctx context.Context, req *tqdpb.DeleteLayerLegalDocRequest) (*tqdpb.DeleteLayerLegalDocResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := h.layerLegalRepo.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete legal document: "+err.Error())
	}

	return &tqdpb.DeleteLayerLegalDocResponse{
		Success: true,
		Message: "Legal document deleted successfully",
	}, nil
}
