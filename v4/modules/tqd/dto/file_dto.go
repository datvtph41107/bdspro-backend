package dto

// FileInfoDTO represents file information
type FileInfoDTO struct {
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	URL         string `json:"url"`
}

// FileUploadDTO represents file upload request
type FileUploadDTO struct {
	File        []byte `json:"file"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	MaxSize     int64  `json:"max_size"`
}

// FileValidationDTO represents file validation request
type FileValidationDTO struct {
	MaxSize      int64    `json:"max_size"`
	AllowedTypes []string `json:"allowed_types"`
	Filename     string   `json:"filename"`
}
