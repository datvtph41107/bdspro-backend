// file/internal/infrastructure/storage/local_storage.go
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"file/internal/ports"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath string) *LocalStorage {
	// Create base directory if it doesn't exist
	os.MkdirAll(basePath, 0755)

	return &LocalStorage{
		basePath: basePath,
	}
}

func (s *LocalStorage) WithBaseURL(baseURL string) *LocalStorage {
	s.baseURL = baseURL
	return s
}

func (s *LocalStorage) Store(ctx context.Context, key string, reader io.Reader, size int64, metadata map[string]string) error {
	// Create full path
	fullPath := filepath.Join(s.basePath, key)
	dir := filepath.Dir(fullPath)

	// Create directory if not exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data
	if _, err := io.Copy(file, reader); err != nil {
		// Clean up on error
		os.Remove(fullPath)
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Retrieve(ctx context.Context, key string) (io.ReadCloser, int64, error) {
	fullPath := filepath.Join(s.basePath, key)

	file, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, fmt.Errorf("file not found: %s", key)
		}
		return nil, 0, fmt.Errorf("failed to open file: %w", err)
	}

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return file, stat.Size(), nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)
	return os.Remove(fullPath)
}

func (s *LocalStorage) Copy(ctx context.Context, srcKey, dstKey string) error {
	srcPath := filepath.Join(s.basePath, srcKey)
	dstPath := filepath.Join(s.basePath, dstKey)

	// Create destination directory
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy data
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		// Clean up on error
		os.Remove(dstPath)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.basePath, key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *LocalStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if s.baseURL == "" {
		return "", fmt.Errorf("base URL not configured for local storage")
	}
	return fmt.Sprintf("%s/%s", s.baseURL, key), nil
}

func (s *LocalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	fullPath := filepath.Join(s.basePath, prefix)

	var files []string
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(s.basePath, path)
			files = append(files, relPath)
		}
		return nil
	})

	return files, err
}

func (s *LocalStorage) GetSize(ctx context.Context, key string) (int64, error) {
	fullPath := filepath.Join(s.basePath, key)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

// Implement ports.StorageService interface
var _ ports.StorageService = (*LocalStorage)(nil)
