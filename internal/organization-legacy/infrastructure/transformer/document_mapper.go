package transformer

import (
	"organization/internal/domain/entity"
	organizationpb "pb/types/organization"
)

type DocumentMapper struct{}

func NewDocumentMapper() *DocumentMapper {
	return &DocumentMapper{}
}

func (m *DocumentMapper) EntityToPb(document *entity.AttachDocument) *organizationpb.DocumentItem {
	return &organizationpb.DocumentItem{
		Name: document.DocName,
		Path: document.DocPath,
		Type: document.DocType,
		Size: document.DocSize,
	}
}

func (m *DocumentMapper) PbToEntity(document *organizationpb.DocumentItem) *entity.AttachDocument {
	return &entity.AttachDocument{
		DocName: document.Name,
		DocPath: document.Path,
		DocType: document.Type,
		DocSize: document.Size,
	}
}
