package filecontent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	_utils "common/utils"

	"file/dto"
	"file/models"
)

var ErrNoActor = errors.New("secure file requires current profile")

const (
	MaxContentBytes      = 64 << 20
	compensationTimeout  = 5 * time.Second
	ordinaryDocumentPath = "document"
)

type UploadInput struct {
	Content     io.Reader
	Size        int64
	Filename    string
	ContentType string
	Description string
	Secure      bool
	ProfileID   uint64
}

type SavedFile struct {
	ID            uint64
	Path          string
	Filename      string
	ThumbnailPath string
	Extension     string
	Size          int64
	Hash          string
	ContentType   string
}

type fileMetadataStore interface {
	SaveWithContext(ctx context.Context, file *models.FileEntity) error
	DeleteByIDWithContext(ctx context.Context, id uint64) error
	FinalizeUpload(ctx context.Context, file *models.FileEntity, access *models.AccessEntity) error
	TryCreateOwned(ctx context.Context, file *models.FileEntity) (bool, error)
	FindByOwner(ctx context.Context, namespace string, key string) (*models.FileEntity, bool, error)
}

type contentStorage interface {
	StoreContent(
		reader io.Reader,
		fileName string,
		contentType string,
		size int64,
		accessID *uint64,
		fileID uint64,
		relativePaths ...string,
	) (*dto.FileInfo, error)
	StoreOwnedContent(
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

// Service owns transport-independent ordinary and owned File effects. DB and
// filesystem implementations are consumed only through the capabilities the
// use cases require.
type Service struct {
	files   fileMetadataStore
	storage contentStorage
	xorKey  string
}

func NewService(
	files fileMetadataStore,
	storage contentStorage,
	xorKey string,
) *Service {
	return &Service{
		files:   files,
		storage: storage,
		xorKey:  xorKey,
	}
}

func (s *Service) UploadFile(ctx context.Context, input UploadInput) (SavedFile, error) {
	if ctx == nil {
		return SavedFile{}, errors.New("context is nil")
	}
	if s == nil || s.files == nil || s.storage == nil {
		return SavedFile{}, errors.New("file content dependencies are unavailable")
	}
	if input.Content == nil {
		return SavedFile{}, errors.New("file content is nil")
	}
	if input.Size < 0 {
		return SavedFile{}, errors.New("file size is invalid")
	}
	input.Filename = strings.TrimSpace(input.Filename)
	if input.Filename == "" {
		return SavedFile{}, errors.New("filename is empty")
	}
	if input.Secure && input.ProfileID == 0 {
		return SavedFile{}, ErrNoActor
	}

	fileEntity := &models.FileEntity{Description: input.Description}
	if err := s.files.SaveWithContext(ctx, fileEntity); err != nil {
		return SavedFile{}, fmt.Errorf("create file record: %w", err)
	}

	cleanupCtx, cancelCleanup := compensationContext(ctx)
	defer cancelCleanup()
	cleanupRecord := func() error {
		return s.files.DeleteByIDWithContext(cleanupCtx, fileEntity.ID)
	}

	var accessID *uint64
	if input.Secure {
		profileID := input.ProfileID
		accessID = &profileID
	}

	fileInfo, err := s.storage.StoreContent(
		input.Content,
		input.Filename,
		input.ContentType,
		input.Size,
		accessID,
		fileEntity.ID,
		ordinaryDocumentPath,
	)
	if err != nil {
		return SavedFile{}, errors.Join(
			fmt.Errorf("store file content: %w", err),
			cleanupRecord(),
		)
	}

	cleanupArtifact := func() error {
		return s.storage.RemovePhysicalPath(fileInfo.AbsolutePath)
	}

	fileEntity.SetFields(fileInfo)
	encryptedPath := _utils.XorEncode(fileInfo.RelativePath, s.xorKey)
	if input.Secure {
		fileEntity.Path = "s" + encryptedPath
	} else {
		fileEntity.Path = "p" + encryptedPath
	}

	var access *models.AccessEntity
	if input.Secure {
		access = &models.AccessEntity{
			FileID:   fileEntity.ID,
			AccessID: input.ProfileID,
		}
	}

	if err := s.files.FinalizeUpload(ctx, fileEntity, access); err != nil {
		return SavedFile{}, errors.Join(
			fmt.Errorf("finalize file metadata: %w", err),
			cleanupArtifact(),
			cleanupRecord(),
		)
	}

	return SavedFile{
		ID:            fileEntity.ID,
		Path:          fileEntity.Path,
		Filename:      fileInfo.FileName,
		ThumbnailPath: fileInfo.ThumbnailPath,
		Extension:     fileInfo.Extension,
		Size:          fileInfo.Size,
		Hash:          fileInfo.Hash,
		ContentType:   fileInfo.ContentType,
	}, nil
}

func compensationContext(ctx context.Context) (context.Context, context.CancelFunc) {
	base := context.Background()
	if ctx != nil {
		base = context.WithoutCancel(ctx)
	}
	return context.WithTimeout(base, compensationTimeout)
}
