package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"file/dto"
)

// BuildOwnedFilePath deliberately excludes date/time/random components.
// One durable File row always maps to one physical owned-effect path.
func (s *StorageService) BuildOwnedFilePath(
	fileName string,
	accessID *uint64,
	fileID uint64,
	relativePaths ...string,
) (string, string, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		return "", "", errors.New("file name is empty")
	}
	if fileID == 0 {
		return "", "", errors.New("file id is required")
	}

	ext := filepath.Ext(filepath.Base(fileName))
	secure := accessID != nil
	location := s.resolve(secure, relativePaths...)

	filePath := filepath.Join(
		location,
		fmt.Sprintf("%d%s", fileID, ext),
	)

	if err := validatePath(filePath, s.rootLocation); err != nil {
		return "", ext, err
	}

	return filePath, ext, nil
}

// StoreOwnedContent writes through a temporary file in the same directory and
// atomically renames it to the stable owned-effect path.
func (s *StorageService) StoreOwnedContent(
	reader io.Reader,
	fileName string,
	contentType string,
	size int64,
	accessID *uint64,
	fileID uint64,
	relativePaths ...string,
) (*dto.FileInfo, error) {
	if reader == nil {
		return nil, errors.New("file content is nil")
	}

	filePath, ext, err :=
		s.BuildOwnedFilePath(
			fileName,
			accessID,
			fileID,
			relativePaths...,
		)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf(
			"create owned file directory: %w",
			err,
		)
	}

	tmp, err := os.CreateTemp(
		dir,
		".owned-*",
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create owned file temp: %w",
			err,
		)
	}

	tmpPath := tmp.Name()
	removeTemp := true

	defer func() {
		_ = tmp.Close()
		if removeTemp {
			_ = os.Remove(tmpPath)
		}
	}()

	written, err := io.Copy(tmp, reader)
	if err != nil {
		return nil, fmt.Errorf(
			"write owned file temp: %w",
			err,
		)
	}

	if written != size {
		return nil, fmt.Errorf(
			"owned file size mismatch: expected=%d actual=%d",
			size,
			written,
		)
	}

	if err := tmp.Sync(); err != nil {
		return nil, fmt.Errorf(
			"sync owned file temp: %w",
			err,
		)
	}

	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf(
			"close owned file temp: %w",
			err,
		)
	}

	if err := os.Rename(tmpPath, filePath); err != nil {
		return nil, fmt.Errorf(
			"publish owned file: %w",
			err,
		)
	}

	removeTemp = false

	hash, err := computeFileHash(filePath)
	if err != nil {
		return nil, err
	}

	storedPath, err := s.StoredPath(filePath)
	if err != nil {
		return nil, err
	}

	return &dto.FileInfo{
		FileName:     filepath.Base(fileName),
		RelativePath: storedPath,
		AbsolutePath: filePath,
		Extension:    ext,
		Size:         written,
		Hash:         hash,
		ContentType:  contentType,
	}, nil
}
