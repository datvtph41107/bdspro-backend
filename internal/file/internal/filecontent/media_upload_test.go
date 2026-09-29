package filecontent

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	_utils "common/utils"

	"file/dto"
	"file/internal/filemedia"
	"file/models"
)

type fakeMediaFileStore struct {
	saveCalls   int
	deleteCalls int
	failOnSave  int
	err         error
}

func (f *fakeMediaFileStore) SaveWithContext(_ context.Context, file *models.FileEntity) error {
	f.saveCalls++
	if f.failOnSave == f.saveCalls && f.err != nil {
		return f.err
	}
	if file.ID == 0 {
		file.ID = 42
	}
	return nil
}

func (f *fakeMediaFileStore) DeleteByIDWithContext(context.Context, uint64) error {
	f.deleteCalls++
	return nil
}

type fakeMediaProcessor struct {
	err           error
	sourcePath    string
	sourceExisted bool
	cleanupCalls  int
}

func (f *fakeMediaProcessor) Process(input filemedia.VideoInput) (*dto.FileInfo, error) {
	f.sourcePath = input.SourcePath
	if _, err := os.Stat(input.SourcePath); err == nil {
		f.sourceExisted = true
	}
	if f.err != nil {
		return nil, f.err
	}
	return &dto.FileInfo{
		FileName:      input.FileName,
		RelativePath:  "files/public/video/2026-08-25/abc/42",
		AbsolutePath:  "/physical/public/video/2026-08-25/abc/42",
		ThumbnailPath: "files/public/document/2026-08-25/abc/42.jpg",
		Extension:     ".mp4",
		Size:          input.Size,
		Hash:          "hash",
		ContentType:   input.ContentType,
	}, nil
}

func (f *fakeMediaProcessor) Cleanup(*dto.FileInfo) error {
	f.cleanupCalls++
	return nil
}

func mediaUploadInput() MediaUploadInput {
	return MediaUploadInput{
		Content:     strings.NewReader("video"),
		Filename:    "video.mp4",
		ContentType: "video/mp4",
		Size:        5,
		Description: "demo",
	}
}

func TestMediaUploadOwnsTemporarySourceAndStableReferences(t *testing.T) {
	files := &fakeMediaFileStore{}
	processor := &fakeMediaProcessor{}
	const xorKey = "media-xor-key"

	saved, err := NewMediaUploadService(files, processor, xorKey).UploadVideo(
		context.Background(), mediaUploadInput(),
	)
	if err != nil {
		t.Fatalf("UploadVideo() error = %v", err)
	}
	if !processor.sourceExisted {
		t.Fatal("processor did not receive materialized source file")
	}
	if _, err := os.Stat(processor.sourcePath); !os.IsNotExist(err) {
		t.Fatalf("temporary source still exists: %v", err)
	}
	if files.saveCalls != 2 || files.deleteCalls != 0 {
		t.Fatalf("save=%d delete=%d", files.saveCalls, files.deleteCalls)
	}

	wantPath := "p" + _utils.XorEncode("files/public/video/2026-08-25/abc/42", xorKey)
	wantThumbnail := "p" + _utils.XorEncode("files/public/document/2026-08-25/abc/42.jpg", xorKey)
	if saved.ID != 42 || saved.RelativePath != wantPath || saved.ThumbnailPath != wantThumbnail {
		t.Fatalf("saved = %#v", saved)
	}
}

func TestMediaUploadCompensatesMetadataWhenProcessingFails(t *testing.T) {
	files := &fakeMediaFileStore{}
	processor := &fakeMediaProcessor{err: errors.New("ffmpeg failed")}

	_, err := NewMediaUploadService(files, processor, "media-xor-key").UploadVideo(
		context.Background(), mediaUploadInput(),
	)
	if err == nil {
		t.Fatal("UploadVideo() error = nil")
	}
	if files.deleteCalls != 1 {
		t.Fatalf("delete calls = %d, want 1", files.deleteCalls)
	}
	// Processor owns cleanup of its own partial FFmpeg effect on Process error.
	if processor.cleanupCalls != 0 {
		t.Fatalf("cleanup calls = %d, want 0", processor.cleanupCalls)
	}
}

func TestMediaUploadCompensatesFilesystemAndMetadataWhenFinalizeFails(t *testing.T) {
	wantErr := errors.New("finalize failed")
	files := &fakeMediaFileStore{failOnSave: 2, err: wantErr}
	processor := &fakeMediaProcessor{}

	_, err := NewMediaUploadService(files, processor, "media-xor-key").UploadVideo(
		context.Background(), mediaUploadInput(),
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if processor.cleanupCalls != 1 || files.deleteCalls != 1 {
		t.Fatalf("processor cleanup=%d metadata delete=%d", processor.cleanupCalls, files.deleteCalls)
	}
}
