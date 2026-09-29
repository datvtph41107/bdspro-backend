package entity

import (
	_models "common/models"
	"organization/internal/enums"
)

type AttachDocument struct {
	_models.BaseEntity
	OwnerId  uint64 // DealID or GroupID do DocOwner quyết định
	DocOwner enums.DocOwner
	DocName  string
	DocPath  string
	DocType  string
	DocSize  int64
}

func (AttachDocument) TableName() string {
	return "attach_documents"
}
