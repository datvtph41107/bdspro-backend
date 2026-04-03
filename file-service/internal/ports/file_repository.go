package ports

import (
	"context"
	"file/internal/core/domain"
	"time"
)

// FileRepository interface
type FileRepository interface {
	Create(ctx context.Context, file *domain.File) error
	FindByID(ctx context.Context, id uint64) (*domain.File, error)
	FindByUUID(ctx context.Context, uuid string) (*domain.File, error)
	FindByHash(ctx context.Context, hash string) (*domain.File, error)
	FindByReference(ctx context.Context, serviceID, referenceType, referenceID string) ([]*domain.File, error)
	Update(ctx context.Context, file *domain.File) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, filter FileFilter, pagination Pagination) ([]*domain.File, int64, error)
	UpdateThumbnail(ctx context.Context, fileID uint64, thumbnailPath string) error
	UpdateStatus(ctx context.Context, fileID uint64, status string) error
}

// FileAccessRepository interface
type FileAccessRepository interface {
	GrantAccess(ctx context.Context, access *domain.FileAccess) error
	RevokeAccess(ctx context.Context, fileID, userID uint64) error
	HasAccess(ctx context.Context, fileID, userID uint64, permission string) (bool, error)
	GetAccessibleFiles(ctx context.Context, userID uint64, permission string) ([]uint64, error)
}

// UploadSessionRepository interface
type UploadSessionRepository interface {
	CreateSession(ctx context.Context, session *domain.UploadSession) error
	GetSession(ctx context.Context, sessionID string) (*domain.UploadSession, error)
	UpdateSession(ctx context.Context, session *domain.UploadSession) error
	DeleteSession(ctx context.Context, sessionID string) error
	CleanupExpiredSessions(ctx context.Context) error
}

// CacheRepository interface
type CacheRepository interface {
	// Basic operations
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)

	// Extended operations
	GetJSON(ctx context.Context, key string, v interface{}) error
	SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error
	HSet(ctx context.Context, key string, values map[string]interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)

	// Key generators
	FileKey(fileID string) string
	FileAccessKey(fileID string, userID uint64) string
	SignatureKey(signature string) string
	UploadSessionKey(sessionID string) string
	RateLimitKey(ip string, endpoint string) string
}

// QueueRepository interface
type QueueRepository interface {
	Publish(ctx context.Context, queue string, message interface{}) error
	Subscribe(ctx context.Context, queue string, handler func(message []byte) error) error
}

// Types
type FileFilter struct {
	ServiceID     string
	ReferenceType string
	ReferenceID   string
	MimeType      string
	Status        string
	UploadedBy    uint64
	StartDate     time.Time
	EndDate       time.Time
}

type Pagination struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
}
