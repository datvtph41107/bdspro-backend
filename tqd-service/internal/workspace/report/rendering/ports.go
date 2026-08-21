package rendering

import (
	"context"
	"io"
)

// Source is owned by the rendering consumer. Implementations may read durable
// Report state, but they must not mutate acceptance/quota/job state.
type Source interface {
	LoadForRender(ctx context.Context, reportID, userID uint64, jobID string) (Document, error)
}

// PDFEngine converts a self-contained HTML document into PDF bytes.
type PDFEngine interface {
	RenderPDF(ctx context.Context, html []byte) ([]byte, error)
}

// Uploader is the narrow File-service capability required after rendering.
type Uploader interface {
	UploadFile(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
}
