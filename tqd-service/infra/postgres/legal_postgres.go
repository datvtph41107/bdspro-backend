package postgres

import (
	"context"
	"fmt"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

// LOGICAL:
// - Surface: Cần implementation để truy xuất DB
// - Root: PostgreSQL specific queries
// - Mechanism: GORM + raw SQL cho spatial queries

type LegalDocumentPostgres struct {
	DB *gorm.DB
}

func NewLegalDocumentPostgres(db *gorm.DB) repo.ILegalDocumentRepo {
	return &LegalDocumentPostgres{DB: db}
}

func (r *LegalDocumentPostgres) GetByParcelID(ctx context.Context, parcelID uint64) ([]*qh_domain.QHLegalDocument, error) {
	var documents []*qh_domain.QHLegalDocument

	// LOGICAL: Query documents áp dụng cho parcel thông qua:
	// 1. Direct parcel-legal relation (qh_parcel_legal_docs)
	// 2. Layer-legal relation (qh_layer_legal_docs) → layers chứa region cắt parcel
	// 3. Zone-legal relation (qh_zone_legal_docs) → zone cắt parcel
	query := `
		SELECT DISTINCT ld.*
		FROM qh_legal_documents ld
		WHERE ld.deleted_at IS NULL
		  AND (
			-- Direct parcel link
			EXISTS (
				SELECT 1 FROM qh_parcel_legal_docs pld
				WHERE pld.legal_document_id = ld.id
				  AND pld.parcel_id = ?
			)
			-- Layer link
			OR EXISTS (
				SELECT 1 FROM qh_layer_legal_docs lld
				JOIN qh_regions r ON r.layer_id = lld.layer_id
				WHERE lld.legal_document_id = ld.id
				  AND ST_Intersects(r.geometry, (SELECT geometry FROM parcels WHERE id = ?))
				  AND r.status = 10
				  AND r.deleted_at IS NULL
			)
			-- Zone link
			OR EXISTS (
				SELECT 1 FROM qh_zone_legal_docs zld
				WHERE zld.legal_document_id = ld.id
				  AND zld.zone_id IN (
					SELECT r.id FROM qh_regions r
					WHERE ST_Intersects(r.geometry, (SELECT geometry FROM parcels WHERE id = ?))
					  AND r.status = 10
					  AND r.deleted_at IS NULL
				  )
			)
		  )
		ORDER BY ld.issued_at DESC
	`

	err := r.DB.WithContext(ctx).
		Raw(query, parcelID, parcelID, parcelID).
		Scan(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get legal documents by parcel ID: %w", err)
	}

	return documents, nil
}

func (r *LegalDocumentPostgres) GetByIDs(ctx context.Context, ids []uint64) ([]*qh_domain.QHLegalDocument, error) {
	if len(ids) == 0 {
		return []*qh_domain.QHLegalDocument{}, nil
	}

	var documents []*qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Order("issued_at DESC").
		Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get legal documents by IDs: %w", err)
	}
	return documents, nil
}

func (r *LegalDocumentPostgres) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLegalDocument, error) {
	var document qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&document).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get legal document by ID: %w", err)
	}
	return &document, nil
}

func (r *LegalDocumentPostgres) GetByLayerIDs(ctx context.Context, layerIDs []uint64) ([]*qh_domain.QHLegalDocument, error) {
	if len(layerIDs) == 0 {
		return []*qh_domain.QHLegalDocument{}, nil
	}

	var documents []*qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where(`EXISTS (
			SELECT 1 FROM qh_layer_legal_docs lld 
			WHERE lld.legal_document_id = qh_legal_documents.id 
			AND lld.layer_id IN ?
		)`, layerIDs).
		Order("issued_at DESC").
		Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get legal documents by layer IDs: %w", err)
	}
	return documents, nil
}

func (r *LegalDocumentPostgres) GetByDocumentType(ctx context.Context, docType enums.DocumentType) ([]*qh_domain.QHLegalDocument, error) {
	var documents []*qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("document_type = ? AND deleted_at IS NULL", docType).
		Order("issued_at DESC").
		Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get legal documents by type: %w", err)
	}
	return documents, nil
}

func (r *LegalDocumentPostgres) GetActive(ctx context.Context) ([]*qh_domain.QHLegalDocument, error) {
	var documents []*qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", enums.DocStatusActive).
		Order("issued_at DESC").
		Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get active legal documents: %w", err)
	}
	return documents, nil
}

func (r *LegalDocumentPostgres) GetSupersededBy(ctx context.Context, documentID uint64) (*qh_domain.QHLegalDocument, error) {
	var document qh_domain.QHLegalDocument
	err := r.DB.WithContext(ctx).
		Where("? = ANY(supersedes) AND deleted_at IS NULL", documentID).
		First(&document).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get superseded by: %w", err)
	}
	return &document, nil
}

func (r *LegalDocumentPostgres) GetByParcelIDWithFilter(
	ctx context.Context,
	parcelID uint64,
	docType enums.DocumentType,
	status enums.DocumentStatus,
) ([]*qh_domain.QHLegalDocument, error) {
	// LOGICAL: Build query với filter dynamic
	query := r.DB.WithContext(ctx).
		Where("deleted_at IS NULL").
		Where(`EXISTS (
			SELECT 1 FROM qh_parcel_legal_docs pld
			WHERE pld.legal_document_id = qh_legal_documents.id
			AND pld.parcel_id = ?
		)`, parcelID)

	if docType != 0 {
		query = query.Where("document_type = ?", docType)
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}

	var documents []*qh_domain.QHLegalDocument
	err := query.Order("issued_at DESC").Find(&documents).Error
	if err != nil {
		return nil, fmt.Errorf("get legal documents with filter: %w", err)
	}
	return documents, nil
}
