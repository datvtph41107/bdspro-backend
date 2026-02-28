package storage

import (
	"file/internal/config"
	"file/internal/ports"
	"fmt"
)

type StorageType string

const (
	StorageLocal StorageType = "local"
	StorageS3    StorageType = "s3"
	StorageGCS   StorageType = "gcs"
	StorageMinIO StorageType = "minio"
)

// StorageFactory creates storage instances
type StorageFactory struct {
	config *config.StorageConfig
}

func NewStorageFactory(cfg *config.StorageConfig) *StorageFactory {
	return &StorageFactory{config: cfg}
}

func (f *StorageFactory) CreateStorage(storageType StorageType) (ports.StorageService, error) {
	switch storageType {
	case StorageLocal:
		return NewLocalStorage(f.config.LocalPath), nil
	// case StorageS3:
	// return NewS3Storage(
	// 	f.config.S3Bucket,
	// 	f.config.S3Region,
	// 	f.config.S3Endpoint,
	// ), nil
	case StorageGCS:
		return nil, fmt.Errorf("GCS storage not implemented")
	case StorageMinIO:
		return nil, fmt.Errorf("MinIO storage not implemented")
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func (f *StorageFactory) GetDefaultStorage() (ports.StorageService, error) {
	return f.CreateStorage(StorageType(f.config.Type))
}

// Base storage with common functionality
type baseStorage struct {
	baseURL string
}

func (s *baseStorage) GenerateKey(fileName, fileHash string) string {
	// Generate storage key: hash/filename
	return fmt.Sprintf("%s/%s", fileHash[:2], fileName)
}

func (s *baseStorage) GetPublicURL(key string) string {
	if s.baseURL == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.baseURL, key)
}
