package provider

import (
	"context"
	"io"
)

// FileProvider is the legacy QH File capability shared by the planning-folder
// upload and planning-classification read paths.
//
// Keep this surface limited to operations with real production consumers.
// Generated Report owns its separate PutOwnedFile capability.
type FileProvider interface {
	UploadFile(
		ctx context.Context,
		file io.Reader,
		filename string,
		contentType string,
	) (string, error)

	GetFile(
		ctx context.Context,
		path string,
	) ([]byte, error)
}
