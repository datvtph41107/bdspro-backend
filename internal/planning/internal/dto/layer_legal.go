package dto

type LayerLegalInput struct {
	Name     string `json:"name" binding:"required"`
	FileUrl  string `json:"fileUrl" binding:"required"`
	FileType string `json:"fileType"`
}

type LayerLegalDTO struct {
	ID        uint64 `json:"id"`
	LayerID   uint64 `json:"layerId"`
	Name      string `json:"name"`
	FileUrl   string `json:"fileUrl"`
	FileType  string `json:"fileType"`
	CreatedAt string `json:"createdAt"`
}

type AddLayerLegalDocsRequest struct {
	LayerID   uint64             `json:"layerId"`
	LegalDocs []*LayerLegalInput `json:"legalDocs" binding:"required,min=1"`
}

type AddLayerLegalDocsResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message"`
	LegalDocs []*LayerLegalDTO `json:"legalDocs"`
}

type GetLayerLegalDocsResponse struct {
	LayerID   uint64           `json:"layerId"`
	Total     int              `json:"total"`
	LegalDocs []*LayerLegalDTO `json:"legalDocs"`
}
