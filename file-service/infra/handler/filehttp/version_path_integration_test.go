package filehttp

import (
	"context"
	"io"
	"strings"
	"testing"

	"file/internal/fileversion"
	"file/models"
	"file/services"
)

type versionPathMetadataStore struct{}

func (versionPathMetadataStore) SaveWithContext(
	_ context.Context,
	file *models.FileEntity,
) error {
	if file.ID == 0 {
		file.ID = 42
	}
	return nil
}

func (versionPathMetadataStore) DeleteByIDWithContext(
	context.Context,
	uint64,
) error {
	return nil
}

func TestVersionUploadPublishedPathDownloadsFromCanonicalStorage(
	t *testing.T,
) {
	const xorKey = "version-path-contract-test-key"

	root := t.TempDir()

	storage, err := services.NewStorageService(root)
	if err != nil {
		t.Fatalf("NewStorageService() error = %v", err)
	}

	codec, err := fileversion.NewReferenceCodec(
		xorKey,
		"version-path-contract-signature-key",
	)
	if err != nil {
		t.Fatal(err)
	}

	uploader := fileversion.NewUploader(
		versionPathMetadataStore{},
		storage,
		codec,
	)

	saved, err := uploader.Upload(
		context.Background(),
		fileversion.UploadInput{
			OpenContent: func() (
				io.ReadCloser,
				error,
			) {
				return io.NopCloser(
					strings.NewReader("version-binary"),
				), nil
			},
			FileName:    "app.apk",
			ContentType: "application/vnd.android.package-archive",
			Description: "release",
			Size:        int64(len("version-binary")),
		},
	)
	if err != nil {
		t.Fatalf(
			"upload version artifact: %v",
			err,
		)
	}

	handler := NewVersionDownloadHandler(
		codec,
		storage,
	)

	rec := performVersionDownload(
		t,
		handler,
		saved.RelativePath,
		"",
	)

	if rec.Code != 200 {
		t.Fatalf(
			"download status = %d, body=%s",
			rec.Code,
			rec.Body.String(),
		)
	}

	if rec.Body.String() != "version-binary" {
		t.Fatalf(
			"download body = %q",
			rec.Body.String(),
		)
	}
}
