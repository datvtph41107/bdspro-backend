package mapper

import (
	domain "mdn4/modules/crm/entity"
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
