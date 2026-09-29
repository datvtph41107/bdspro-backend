package fileversion

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"file/dto"
	"file/models"
)

const (
	artifactStorageNamespace   = "version"
	versionCompensationTimeout = 5 * time.Second
)

type metadataStore interface {
	SaveWithContext(ctx context.Context, file *models.FileEntity) error
	DeleteByIDWithContext(ctx context.Context, id uint64) error
}

type artifactStorage interface {
	StoreContent(
		reader io.Reader,
		fileName string,
		contentType string,
		size int64,
		accessID *uint64,
		fileID uint64,
		relativePaths ...string,
	) (*dto.FileInfo, error)
	RemovePhysicalPath(path string) error
}

type referenceEncoder interface {
	Encode(storedPath string) (string, error)
}

type UploadInput struct {
	OpenContent func() (io.ReadCloser, error)
	FileName    string
	ContentType string
	Description string
	Size        int64
}

// Uploader owns one version-artifact effect from metadata allocation through
// physical storage and durable reference publication. Failed partial effects
// are compensated; HTTP/API-key policy is deliberately outside this type.
type Uploader struct {
	metadata  metadataStore
	storage   artifactStorage
	reference referenceEncoder
}

func NewUploader(
	metadata metadataStore,
	storage artifactStorage,
	reference referenceEncoder,
) *Uploader {
	return &Uploader{metadata: metadata, storage: storage, reference: reference}
}

func (u *Uploader) Upload(ctx context.Context, input UploadInput) (*dto.FileInfo, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if u == nil || u.metadata == nil || u.storage == nil || u.reference == nil {
		return nil, errors.New("version upload dependencies are unavailable")
	}
	if input.OpenContent == nil {
		return nil, errors.New("version content source is unavailable")
	}
	input.FileName = strings.TrimSpace(input.FileName)
	if input.FileName == "" {
		return nil, errors.New("version filename is empty")
	}
	if input.Size < 0 {
		return nil, errors.New("version size is invalid")
	}

	fileEntity := &models.FileEntity{Description: input.Description}
	if err := u.metadata.SaveWithContext(ctx, fileEntity); err != nil {
		return nil, fmt.Errorf("allocate version metadata: %w", err)
	}

	cleanupCtx, cancelCleanup := versionCompensationContext(ctx)
	defer cancelCleanup()
	cleanupRecord := func() error {
		return u.metadata.DeleteByIDWithContext(cleanupCtx, fileEntity.ID)
	}

	content, err := input.OpenContent()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("open version content: %w", err), cleanupRecord())
	}
	defer content.Close()

	fileInfo, err := u.storage.StoreContent(
		content,
		input.FileName,
		input.ContentType,
		input.Size,
		nil,
		fileEntity.ID,
		artifactStorageNamespace,
	)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("store version content: %w", err), cleanupRecord())
	}
	if fileInfo == nil {
		return nil, errors.Join(errors.New("version storage returned no file info"), cleanupRecord())
	}

	cleanupArtifact := func() error {
		return u.storage.RemovePhysicalPath(fileInfo.AbsolutePath)
	}

	reference, err := u.reference.Encode(fileInfo.RelativePath)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("encode version reference: %w", err),
			cleanupArtifact(),
			cleanupRecord(),
		)
	}

	fileEntity.SetFields(fileInfo)
	fileEntity.Path = reference
	if err := u.metadata.SaveWithContext(ctx, fileEntity); err != nil {
		return nil, errors.Join(
			fmt.Errorf("finalize version metadata: %w", err),
			cleanupArtifact(),
			cleanupRecord(),
		)
	}

	fileInfo.ID = fileEntity.ID
	fileInfo.RelativePath = fileEntity.Path
	return fileInfo, nil
}

func versionCompensationContext(ctx context.Context) (context.Context, context.CancelFunc) {
	base := context.Background()
	if ctx != nil {
		base = context.WithoutCancel(ctx)
	}
	return context.WithTimeout(base, versionCompensationTimeout)
}
