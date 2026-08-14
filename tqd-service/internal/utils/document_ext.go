package utils

import (
	"path/filepath"
	"strings"
)

// ResolveDocumentExt lấy đuôi file từ relativePath hoặc filepath (không có dấu chấm).
func ResolveDocumentExt(relativePath, filepathValue string) string {
	path := strings.TrimSpace(relativePath)
	if path == "" {
		path = strings.TrimSpace(filepathValue)
	}
	if path == "" {
		return ""
	}
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
}
