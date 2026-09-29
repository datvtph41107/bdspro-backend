package filecontent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	_utils "common/utils"

	"file/dto"
	"file/internal/filemedia"
	"file/models"
)

type mediaFileStore interface {
	SaveWithContext(ctx context.Context, file *models.FileEntity) error
	DeleteByIDWithContext(ctx context.Context, id uint64) error
}

type mediaProcessor interface {
	Process(input filemedia.VideoInput) (*dto.FileInfo, error)
	Cleanup(info *dto.FileInfo) error
}

type MediaUploadInput struct {
	Content     io.Reader
	Filename    string
	ContentType string
	Description string
	Size        int64
}

type MediaUploadService struct {
	files  mediaFileStore
	media  mediaProcessor
	xorKey string
}

func NewMediaUploadService(
	files mediaFileStore,
	media mediaProcessor,
	xorKey string,
) *MediaUploadService {
	return &MediaUploadService{files: files, media: media, xorKey: xorKey}
}

// UploadVideo owns one media effect from metadata allocation through FFmpeg
// processing and durable publication. Partial metadata/filesystem effects are
// compensated before the error is returned.
func (s *MediaUploadService) UploadVideo(
	ctx context.Context,
	input MediaUploadInput,
) (*dto.FileInfo, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if s == nil || s.files == nil || s.media == nil {
		return nil, errors.New("media upload dependencies are unavailable")
	}
	if input.Content == nil {
		return nil, errors.New("video content is nil")
	}
	if input.Size < 0 {
		return nil, errors.New("video size is invalid")
	}
	input.Filename = strings.TrimSpace(input.Filename)
	if input.Filename == "" {
		return nil, errors.New("filename is empty")
	}

	fileEntity := &models.FileEntity{Description: input.Description}
	if err := s.files.SaveWithContext(ctx, fileEntity); err != nil {
		return nil, fmt.Errorf("create file record: %w", err)
	}

	cleanupCtx, cancelCleanup := compensationContext(ctx)
	defer cancelCleanup()
	cleanupRecord := func() error {
		return s.files.DeleteByIDWithContext(cleanupCtx, fileEntity.ID)
	}

	// FFmpeg requires a real pathname; materialization belongs to the media
	// capability rather than the HTTP adapter.
	tempFile, err := os.CreateTemp("", fmt.Sprintf("file-media-%d-*.mp4", fileEntity.ID))
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("create temporary media input: %w", err),
			cleanupRecord(),
		)
	}

	tempPath := tempFile.Name()
	defer os.Remove(tempPath)

	if _, err := io.Copy(tempFile, input.Content); err != nil {
		_ = tempFile.Close()
		return nil, errors.Join(
			fmt.Errorf("write temporary media input: %w", err),
			cleanupRecord(),
		)
	}
	if err := tempFile.Close(); err != nil {
		return nil, errors.Join(
			fmt.Errorf("close temporary media input: %w", err),
			cleanupRecord(),
		)
	}

	fileInfo, err := s.media.Process(filemedia.VideoInput{
		SourcePath:  tempPath,
		FileName:    input.Filename,
		ContentType: input.ContentType,
		Size:        input.Size,
		FileID:      fileEntity.ID,
	})
	if err != nil {
		return nil, errors.Join(err, cleanupRecord())
	}
	if fileInfo == nil {
		return nil, errors.Join(
			errors.New("media processor returned no file info"),
			cleanupRecord(),
		)
	}

	// Both references are stable File namespace paths before they are encoded;
	// the configured physical storage root never enters the wire token.
	fileEntity.SetFields(fileInfo)
	fileEntity.Path = "p" + _utils.XorEncode(fileInfo.RelativePath, s.xorKey)
	fileEntity.ThumbnailPath = "p" + _utils.XorEncode(fileInfo.ThumbnailPath, s.xorKey)

	if err := s.files.SaveWithContext(ctx, fileEntity); err != nil {
		return nil, errors.Join(
			fmt.Errorf("save media file info: %w", err),
			s.media.Cleanup(fileInfo),
			cleanupRecord(),
		)
	}

	fileInfo.ID = fileEntity.ID
	fileInfo.RelativePath = fileEntity.Path
	fileInfo.ThumbnailPath = fileEntity.ThumbnailPath
	return fileInfo, nil
}
