package ports

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

// FileService interface
type FileService interface {
	Upload(ctx context.Context, req UploadRequest) (*UploadResponse, error)
	UploadChunk(ctx context.Context, req UploadChunkRequest) (*UploadChunkResponse, error)
	Download(ctx context.Context, req DownloadRequest) (io.Reader, *FileInfo, error)
	GetFileInfo(ctx context.Context, fileID string) (*FileInfo, error)
	DeleteFile(ctx context.Context, req DeleteRequest) error
	GenerateAccessURL(ctx context.Context, req AccessURLRequest) (string, error)
	ListFiles(ctx context.Context, req ListRequest) (*ListResponse, error)
	BatchUpload(ctx context.Context, files []UploadRequest) ([]*UploadResponse, error)
}

// VideoService interface
type VideoService interface {
	ProcessVideo(ctx context.Context, fileID string, options VideoProcessOptions) error
	GenerateThumbnail(ctx context.Context, fileID string, options ThumbnailOptions) (string, error)
	GetVideoStream(ctx context.Context, fileID string, quality string) (io.Reader, *VideoStreamInfo, error)
	TranscodeVideo(ctx context.Context, fileID string, profiles []VideoProfile) error
}

// StorageService interface
type StorageService interface {
	Store(ctx context.Context, key string, reader io.Reader, size int64, metadata map[string]string) error
	Retrieve(ctx context.Context, key string) (io.ReadCloser, int64, error)
	Delete(ctx context.Context, key string) error
	Copy(ctx context.Context, srcKey, dstKey string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}

// Request/Response types
type UploadRequest struct {
	File          multipart.File
	FileHeader    *multipart.FileHeader
	OriginalName  string
	ServiceID     string
	ReferenceType string
	ReferenceID   string
	UserID        uint64
	IsPublic      bool
	Metadata      map[string]interface{}
}

type UploadResponse struct {
	FileID       string    `json:"file_id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	PublicURL    string    `json:"public_url,omitempty"`
	StoragePath  string    `json:"storage_path"`
	Hash         string    `json:"hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type UploadChunkRequest struct {
	SessionID   string
	ChunkNumber int
	TotalChunks int
	ChunkData   []byte
	FileHash    string
}

type UploadChunkResponse struct {
	SessionID     string `json:"session_id"`
	ChunkReceived int    `json:"chunk_received"`
	TotalChunks   int    `json:"total_chunks"`
	IsComplete    bool   `json:"is_complete"`
	FileID        string `json:"file_id,omitempty"`
}

type DownloadRequest struct {
	FileID    string
	UserID    uint64
	Range     string
	Signature string
}

type FileInfo struct {
	ID           string                 `json:"id"`
	OriginalName string                 `json:"original_name"`
	Size         int64                  `json:"size"`
	MimeType     string                 `json:"mime_type"`
	Extension    string                 `json:"extension"`
	PublicURL    string                 `json:"public_url,omitempty"`
	ThumbnailURL string                 `json:"thumbnail_url,omitempty"`
	Width        int                    `json:"width,omitempty"`
	Height       int                    `json:"height,omitempty"`
	Duration     float64                `json:"duration,omitempty"`
	Metadata     map[string]interface{} `json:"metadata"`
	Status       string                 `json:"status"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type AccessURLRequest struct {
	FileID     string
	UserID     uint64
	Permission string
	Expiry     time.Duration
}

type DeleteRequest struct {
	FileID string
	UserID uint64
	Force  bool
}

type ListRequest struct {
	ServiceID     string
	ReferenceType string
	ReferenceID   string
	MimeType      string
	Page          int
	PageSize      int
	SortBy        string
	SortDir       string
}

type ListResponse struct {
	Files      []*FileInfo `json:"files"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type VideoProcessOptions struct {
	GenerateThumbnail bool
	CreateHLS         bool
	Qualities         []string // "360p", "720p", "1080p"
	Watermark         string
}

type ThumbnailOptions struct {
	TimeOffset string // "00:00:01"
	Width      int
	Height     int
	Quality    int
}

type VideoProfile struct {
	Name    string
	Width   int
	Height  int
	Bitrate string
}

type VideoStreamInfo struct {
	ContentType   string
	ContentLength int64
	Range         string
	IsStream      bool
}
