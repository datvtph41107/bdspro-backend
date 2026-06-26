package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
)

// ============================================================
// LEGAL REPOSITORY INTERFACE
// ============================================================
// LOGICAL:
// - Surface: Cần truy xuất legal documents từ DB
// - Root: Tách biệt persistence khỏi business logic
// - Mechanism: Interface cho phép mock và swap implementation

type ILegalDocumentRepo interface {
	// GetByParcelID lấy tất cả legal documents liên quan đến parcel
	// LOGICAL: Query documents qua layer/zone relationships
	GetByParcelID(ctx context.Context, parcelID uint64) ([]*qh_domain.QHLegalDocument, error)

	// GetByIDs lấy documents theo danh sách ID
	// LOGICAL: Bulk fetch cho registry references
	GetByIDs(ctx context.Context, ids []uint64) ([]*qh_domain.QHLegalDocument, error)

	// GetByID lấy 1 document theo ID
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLegalDocument, error)

	// GetByLayerIDs lấy documents liên quan đến layers
	GetByLayerIDs(ctx context.Context, layerIDs []uint64) ([]*qh_domain.QHLegalDocument, error)

	// GetByDocumentType lấy documents theo loại
	GetByDocumentType(ctx context.Context, docType enums.DocumentType) ([]*qh_domain.QHLegalDocument, error)

	// GetActive lấy documents đang có hiệu lực
	GetActive(ctx context.Context) ([]*qh_domain.QHLegalDocument, error)

	// GetSupersededBy lấy document thay thế
	GetSupersededBy(ctx context.Context, documentID uint64) (*qh_domain.QHLegalDocument, error)

	// GetByParcelIDWithFilter lấy documents với filter
	GetByParcelIDWithFilter(ctx context.Context, parcelID uint64, docType enums.DocumentType, status enums.DocumentStatus) ([]*qh_domain.QHLegalDocument, error)
}
