package fileversion

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"file/dto"
	"file/models"
)

type fakeVersionMetadataStore struct {
	saveCalls   int
	deleteCalls int
	failOnSave  int
	err         error
}

func (f *fakeVersionMetadataStore) SaveWithContext(_ context.Context, file *models.FileEntity) error {
	f.saveCalls++
	if f.failOnSave == f.saveCalls && f.err != nil {
		return f.err
	}
	if file.ID == 0 {
		file.ID = 42
	}
	return nil
}

func (f *fakeVersionMetadataStore) DeleteByIDWithContext(context.Context, uint64) error {
	f.deleteCalls++
	return nil
}

type fakeVersionContentStore struct {
	calls       int
	removeCalls int
	fileID      uint64
	namespace   string
	err         error
}

func (f *fakeVersionContentStore) StoreContent(
	reader io.Reader,
	fileName string,
	contentType string,
	size int64,
	_ *uint64,
	fileID uint64,
	relativePaths ...string,
) (*dto.FileInfo, error) {
	f.calls++
	f.fileID = fileID
	if len(relativePaths) > 0 {
		f.namespace = relativePaths[0]
	}
	if f.err != nil {
		return nil, f.err
	}
	if _, err := io.ReadAll(reader); err != nil {
		return nil, err
	}
	return &dto.FileInfo{
		FileName:     fileName,
		RelativePath: "files/public/version/2026-08-26/abc/42.apk",
		AbsolutePath: "/physical/public/version/2026-08-26/abc/42.apk",
		Extension:    ".apk",
		Size:         size,
		Hash:         "hash",
		ContentType:  contentType,
	}, nil
}

func (f *fakeVersionContentStore) RemovePhysicalPath(string) error {
	f.removeCalls++
	return nil
}

type fakeVersionReferenceEncoder struct {
	value string
	err   error
}

func (f fakeVersionReferenceEncoder) Encode(string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.value, nil
}

func versionUploadInput() UploadInput {
	return UploadInput{
		OpenContent: func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("apk")), nil
		},
		FileName:    "app.apk",
		ContentType: "application/vnd.android.package-archive",
		Description: "release",
		Size:        3,
	}
}

func TestUploadOwnsVersionMetadataStorageAndReferenceEffect(t *testing.T) {
	files := &fakeVersionMetadataStore{}
	storage := &fakeVersionContentStore{}
	const reference = "v1.payload.signature"

	saved, err := NewUploader(files, storage, fakeVersionReferenceEncoder{value: reference}).Upload(
		context.Background(), versionUploadInput(),
	)
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if files.saveCalls != 2 || files.deleteCalls != 0 {
		t.Fatalf("save=%d delete=%d", files.saveCalls, files.deleteCalls)
	}
	if storage.calls != 1 || storage.fileID != 42 || storage.namespace != artifactStorageNamespace {
		t.Fatalf("storage calls=%d id=%d namespace=%q", storage.calls, storage.fileID, storage.namespace)
	}
	if saved.ID != 42 || saved.RelativePath != reference {
		t.Fatalf("saved = %#v", saved)
	}
}

func TestUploadCompensatesWhenContentOpenFails(t *testing.T) {
	files := &fakeVersionMetadataStore{}
	storage := &fakeVersionContentStore{}
	input := versionUploadInput()
	input.OpenContent = func() (io.ReadCloser, error) { return nil, errors.New("open failed") }

	_, err := NewUploader(files, storage, fakeVersionReferenceEncoder{value: "ref"}).Upload(context.Background(), input)
	if err == nil {
		t.Fatal("Upload() error = nil")
	}
	if files.deleteCalls != 1 || storage.calls != 0 {
		t.Fatalf("delete=%d storage=%d", files.deleteCalls, storage.calls)
	}
}

func TestUploadCompensatesArtifactAndMetadataWhenFinalizeFails(t *testing.T) {
	wantErr := errors.New("finalize failed")
	files := &fakeVersionMetadataStore{failOnSave: 2, err: wantErr}
	storage := &fakeVersionContentStore{}

	_, err := NewUploader(files, storage, fakeVersionReferenceEncoder{value: "v1.payload.signature"}).Upload(
		context.Background(), versionUploadInput(),
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if storage.removeCalls != 1 || files.deleteCalls != 1 {
		t.Fatalf("storage remove=%d metadata delete=%d", storage.removeCalls, files.deleteCalls)
	}
}
