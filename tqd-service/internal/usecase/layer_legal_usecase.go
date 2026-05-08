// internal/usecase/layer_legal_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

type LayerLegalUsecase struct {
	layerLegalRepo repo.ILayerLegalRepo
}

func NewLayerLegalUsecase(layerLegalRepo repo.ILayerLegalRepo) *LayerLegalUsecase {
	return &LayerLegalUsecase{
		layerLegalRepo: layerLegalRepo,
	}
}

// AddLegalDocs - Thêm nhiều legal documents cho layer
func (u *LayerLegalUsecase) AddLegalDocs(ctx context.Context, layerID uint64, inputs []*dto.LayerLegalInput) (*dto.AddLayerLegalDocsResponse, error) {
	// Validate
	if layerID == 0 {
		return nil, fmt.Errorf("layer_id is required")
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("at least one legal document is required")
	}

	// Convert to domain models
	legals := make([]*qh_domain.QHLayerLegal, 0, len(inputs))
	for _, input := range inputs {
		if input.Name == "" || input.FileUrl == "" {
			continue // Skip invalid entries
		}

		legal := &qh_domain.QHLayerLegal{
			LayerID:  layerID,
			Name:     input.Name,
			FileURL:  input.FileUrl,
			FileType: input.FileType,
		}
		legals = append(legals, legal)
	}

	if len(legals) == 0 {
		return nil, fmt.Errorf("no valid legal documents to add")
	}

	// Save to database
	if err := u.layerLegalRepo.CreateBatch(ctx, legals); err != nil {
		return nil, fmt.Errorf("failed to save legal documents: %w", err)
	}

	// Convert to DTOs for response
	legalDTOs := make([]*dto.LayerLegalDTO, 0, len(legals))
	for _, legal := range legals {
		legalDTOs = append(legalDTOs, &dto.LayerLegalDTO{
			ID:        legal.ID,
			LayerID:   legal.LayerID,
			Name:      legal.Name,
			FileUrl:   legal.FileURL,
			FileType:  legal.FileType,
			CreatedAt: legal.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.AddLayerLegalDocsResponse{
		Success:   true,
		Message:   fmt.Sprintf("Successfully added %d legal document(s)", len(legalDTOs)),
		LegalDocs: legalDTOs,
	}, nil
}

// GetLegalDocsByLayer - Lấy danh sách legal documents của layer
func (u *LayerLegalUsecase) GetLegalDocsByLayer(ctx context.Context, layerID uint64) (*dto.GetLayerLegalDocsResponse, error) {
	if layerID == 0 {
		return nil, fmt.Errorf("layer_id is required")
	}

	legals, err := u.layerLegalRepo.GetByLayerID(ctx, layerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get legal documents: %w", err)
	}

	legalDTOs := make([]*dto.LayerLegalDTO, 0, len(legals))
	for _, legal := range legals {
		legalDTOs = append(legalDTOs, &dto.LayerLegalDTO{
			ID:        legal.ID,
			LayerID:   legal.LayerID,
			Name:      legal.Name,
			FileUrl:   legal.FileURL,
			FileType:  legal.FileType,
			CreatedAt: legal.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.GetLayerLegalDocsResponse{
		LayerID:   layerID,
		Total:     len(legalDTOs),
		LegalDocs: legalDTOs,
	}, nil
}

// GetLegalDocsByLayerIDs - Lấy legal documents cho nhiều layers (dùng trong parcel response)
func (u *LayerLegalUsecase) GetLegalDocsByLayerIDs(ctx context.Context, layerIDs []uint64) (map[uint64][]*dto.LayerLegalDTO, error) {
	if len(layerIDs) == 0 {
		return make(map[uint64][]*dto.LayerLegalDTO), nil
	}

	legalsMap, err := u.layerLegalRepo.GetByLayerIDs(ctx, layerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get legal documents: %w", err)
	}

	result := make(map[uint64][]*dto.LayerLegalDTO)
	for layerID, legals := range legalsMap {
		dtos := make([]*dto.LayerLegalDTO, 0, len(legals))
		for _, legal := range legals {
			dtos = append(dtos, &dto.LayerLegalDTO{
				ID:        legal.ID,
				LayerID:   legal.LayerID,
				Name:      legal.Name,
				FileUrl:   legal.FileURL,
				FileType:  legal.FileType,
				CreatedAt: legal.CreatedAt.Format(time.RFC3339),
			})
		}
		result[layerID] = dtos
	}

	return result, nil
}

// DeleteLegalDoc - Xóa legal document
func (u *LayerLegalUsecase) DeleteLegalDoc(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}

	// Check if exists
	existing, err := u.layerLegalRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check legal document: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("legal document not found")
	}

	// Soft delete
	if err := u.layerLegalRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete legal document: %w", err)
	}

	return nil
}

// DeleteLegalDocsByLayer - Xóa tất cả legal documents của layer
func (u *LayerLegalUsecase) DeleteLegalDocsByLayer(ctx context.Context, layerID uint64) error {
	if layerID == 0 {
		return fmt.Errorf("layer_id is required")
	}

	if err := u.layerLegalRepo.DeleteByLayerID(ctx, layerID); err != nil {
		return fmt.Errorf("failed to delete legal documents: %w", err)
	}

	return nil
}
