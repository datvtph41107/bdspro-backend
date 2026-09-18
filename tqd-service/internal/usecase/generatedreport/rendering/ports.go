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
	FindOwnedFile(ctx context.Context, ownerNamespace, ownerKey string) (path string, size int64, found bool, err error)
	PutOwnedFile(
		ctx context.Context,
		ownerNamespace string,
		ownerKey string,
		file io.Reader,
		filename string,
		contentType string,
	) (string, error)
}

// DeliveryReference translates a durable File path into the public reference
// that Report is allowed to persist and expose.
type DeliveryReference interface {
	PublicURL(path string) (string, error)
}
