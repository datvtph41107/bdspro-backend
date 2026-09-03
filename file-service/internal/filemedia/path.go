package filemedia

import (
	"fmt"
	"path/filepath"
	"strings"

	"file/dto"
)

// ParsePath interprets the stable File media namespace:
// files/<scope>/<type>/<date>/<uuid>/<fileID>[.<ext>].
func ParsePath(path string) (*dto.FileInfo, error) {
	path = strings.TrimSpace(path)
	if path == "" || filepath.IsAbs(path) {
		return nil, fmt.Errorf("path không hợp lệ: %s", path)
	}

	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return nil, fmt.Errorf("path không hợp lệ: %s", path)
	}
	parts := strings.Split(clean, "/")
	if len(parts) != 6 || parts[0] != "files" {
		return nil, fmt.Errorf("path không hợp lệ: %s", path)
	}
	if parts[1] != "public" && parts[1] != "secure" {
		return nil, fmt.Errorf("scope không hợp lệ: %s", parts[1])
	}
	for _, part := range parts[2:] {
		if strings.TrimSpace(part) == "" || part == "." || part == ".." {
			return nil, fmt.Errorf("path không hợp lệ: %s", path)
		}
	}

	fileName := parts[5]
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	fileID := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	if fileID == "" {
		return nil, fmt.Errorf("file id không hợp lệ: %s", path)
	}

	return &dto.FileInfo{
		AbsolutePath: path,
		Scope:        parts[1],
		Type:         parts[2],
		Date:         parts[3],
		UUID:         parts[4],
		FileID:       fileID,
		Ext:          ext,
	}, nil
}
