package mapper

import (
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type DocumentMapper struct{}

func NewDocumentMapper() *DocumentMapper {
	return &DocumentMapper{}
}

func (m *DocumentMapper) EntityToPb(document *domain.AttachDocument) *bdspropb.DocumentItem {
	return &bdspropb.DocumentItem{
		Name: document.DocName,
		Path: document.DocPath,
		Type: document.DocType,
		Size: document.DocSize,
	}
}

func (m *DocumentMapper) PbToEntity(document *bdspropb.DocumentItem) *domain.AttachDocument {
	return &domain.AttachDocument{
		DocName: document.Name,
		DocPath: document.Path,
		DocType: document.Type,
		DocSize: document.Size,
	}
}
