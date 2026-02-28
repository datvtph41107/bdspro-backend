package domain

import (
	"strings"
	"time"
)

type File struct {
	ID            uint64    `json:"id" gorm:"primaryKey"`
	UUID          string    `json:"uuid" gorm:"uniqueIndex;size:36"`
	OriginalName  string    `json:"original_name"`
	StoragePath   string    `json:"storage_path"`
	PublicPath    string    `json:"public_path"`
	ThumbnailPath string    `json:"thumbnail_path,omitempty"`
	Size          int64     `json:"size"`
	MimeType      string    `json:"mime_type"`
	Extension     string    `json:"extension"`
	Hash          string    `json:"hash" gorm:"uniqueIndex"`
	Width         int       `json:"width,omitempty"`
	Height        int       `json:"height,omitempty"`
	Duration      float64   `json:"duration,omitempty"`
	Metadata      JSONMap   `json:"metadata" gorm:"type:jsonb"`
	Status        string    `json:"status" gorm:"default:'active'"`
	UploadedBy    uint64    `json:"uploaded_by,omitempty"`
	ServiceID     string    `json:"service_id,omitempty"`
	ReferenceID   string    `json:"reference_id,omitempty"`
	ReferenceType string    `json:"reference_type,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type FileAccess struct {
	ID         uint64    `json:"id" gorm:"primaryKey"`
	FileID     uint64    `json:"file_id" gorm:"index"`
	UserID     uint64    `json:"user_id" gorm:"index"`
	Permission string    `json:"permission"` // "read", "write", "delete"
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type UploadSession struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	FileID       uint64    `json:"file_id,omitempty"`
	TotalChunks  int       `json:"total_chunks"`
	ChunkSize    int64     `json:"chunk_size"`
	TotalSize    int64     `json:"total_size"`
	UploadedSize int64     `json:"uploaded_size"`
	Status       string    `json:"status"` // "pending", "uploading", "completed", "failed"
	Metadata     JSONMap   `json:"metadata" gorm:"type:jsonb"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type JSONMap map[string]interface{}

func (j JSONMap) GormDataType() string {
	return "jsonb"
}

// File status constants
const (
	FileStatusPending   = "pending"
	FileStatusUploading = "uploading"
	FileStatusActive    = "active"
	FileStatusDeleted   = "deleted"
	FileStatusFailed    = "failed"
)

// File types
const (
	FileTypeImage    = "image"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeDocument = "document"
	FileTypeArchive  = "archive"
	FileTypeOther    = "other"
)

// Permissions
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
	PermissionAdmin  = "admin"
)

// GetFileType returns file type based on mime type
func GetFileType(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return FileTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return FileTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return FileTypeAudio
	case strings.HasPrefix(mimeType, "application/pdf"):
		return FileTypeDocument
	case strings.HasPrefix(mimeType, "application/zip") ||
		strings.HasPrefix(mimeType, "application/x-rar-compressed"):
		return FileTypeArchive
	default:
		return FileTypeOther
	}
}
