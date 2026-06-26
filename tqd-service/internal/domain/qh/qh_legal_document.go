package qh_domain

import (
	_models "common/domain/entity"
	"time"
	"tqd/internal/enums"
)

// QHLegalDocument — Registry văn bản pháp lý
type QHLegalDocument struct {
	_models.BaseEntity

	// Identity — lưu dạng uint32 trong DB
	DocumentType   enums.DocumentType `gorm:"column:document_type;type:int;not null;index" json:"documentType"`
	DocumentCode   string             `gorm:"column:document_code;type:varchar(100);index" json:"documentCode"`
	DocumentNumber string             `gorm:"column:document_number;type:varchar(100)" json:"documentNumber"`

	// Metadata
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	FullName    string `gorm:"column:full_name;type:text" json:"fullName"`
	Description string `gorm:"column:description;type:text" json:"description"`

	// Issuance
	IssuedBy    string     `gorm:"column:issued_by;type:varchar(255);not null" json:"issuedBy"`
	IssuedAt    *time.Time `gorm:"column:issued_at;type:date" json:"issuedAt"`
	EffectiveAt *time.Time `gorm:"column:effective_at;type:date;index" json:"effectiveAt"`
	ExpiresAt   *time.Time `gorm:"column:expires_at;type:date" json:"expiresAt"`

	// Lifecycle — lưu dạng uint32 trong DB
	Status       enums.DocumentStatus `gorm:"column:status;type:int;not null;index" json:"status"`
	SupersededBy *uint64              `gorm:"column:superseded_by;index" json:"supersededBy"`
	Supersedes   []uint64             `gorm:"column:supersedes;type:jsonb;default:'[]'" json:"supersedes"`

	// File
	FileType    string `gorm:"column:file_type;type:varchar(20)" json:"fileType"`
	FileSizeKb  int    `gorm:"column:file_size_kb;default:0" json:"fileSizeKb"`
	PageCount   int    `gorm:"column:page_count;default:0" json:"pageCount"`
	FileURL     string `gorm:"column:file_url;type:text" json:"fileUrl"`
	PreviewURL  string `gorm:"column:preview_url;type:text" json:"previewUrl"`
	DownloadURL string `gorm:"column:download_url;type:text" json:"downloadUrl"`

	// Access Control
	IsPublic     bool `gorm:"column:is_public;default:true" json:"isPublic"`
	RequiresAuth bool `gorm:"column:requires_auth;default:false" json:"requiresAuth"`

	// References
	RelevantArticles  []string `gorm:"column:relevant_articles;type:jsonb;default:'[]'" json:"relevantArticles"`
	Tags              []string `gorm:"column:tags;type:jsonb;default:'[]'" json:"tags"`
	AffectedLayerIDs  []uint64 `gorm:"column:affected_layer_ids;type:jsonb;default:'[]'" json:"affectedLayerIds"`
	AffectedZoneIDs   []uint64 `gorm:"column:affected_zone_ids;type:jsonb;default:'[]'" json:"affectedZoneIds"`
	AffectedParcelIDs []uint64 `gorm:"column:affected_parcel_ids;type:jsonb;default:'[]'" json:"affectedParcelIds"`
}

func (QHLegalDocument) TableName() string {
	return "qh_legal_documents"
}

func (d *QHLegalDocument) IsActive() bool {
	if d.Status != enums.DocStatusActive {
		return false
	}
	if d.ExpiresAt != nil && d.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

func (d *QHLegalDocument) IsSuperseded() bool {
	return d.Status == enums.DocStatusSuperseded || d.SupersededBy != nil
}

func (d *QHLegalDocument) GetStatusLabel() string {
	return enums.DocumentStatusLabel[d.Status]
}

func (d *QHLegalDocument) GetStatusColor() string {
	return enums.DocumentStatusColor[d.Status]
}

func (d *QHLegalDocument) GetTypeName() string {
	return enums.DocumentTypeName[d.DocumentType]
}

func (d *QHLegalDocument) GetTypeIcon() string {
	switch d.DocumentType {
	case enums.DocTypePlanningDecision:
		return "file-certificate"
	case enums.DocTypeLandCertificate:
		return "certificate"
	case enums.DocTypeLaw:
		return "gavel"
	case enums.DocTypeDecree:
		return "scroll"
	case enums.DocTypeRegulation:
		return "book"
	case enums.DocTypeCircular:
		return "mail"
	default:
		return "file"
	}
}

func (d *QHLegalDocument) GetTypeString() string {
	return enums.DocumentTypeString[d.DocumentType]
}

func (d *QHLegalDocument) GetStatusString() string {
	return enums.DocumentStatusString[d.Status]
}
