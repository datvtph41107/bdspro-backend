package mapper

import (
	"crm/internal/domain"
	crmpb "pb/types/crm"
)

func DocumentDomainToPb(doc *domain.DocumentEntity) *crmpb.DocumentDTO {
	return &crmpb.DocumentDTO{
		Id:       doc.ID,
		FileName: doc.FileName,
		FileUrl:  doc.FileUrl,
		FileType: doc.FileType,
	}
}