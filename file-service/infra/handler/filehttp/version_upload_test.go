package filehttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"file/dto"
	"file/internal/fileversion"
	"file/internal/versionauth"

	"github.com/gin-gonic/gin"
)

type fakeVersionVerifier struct {
	calls int
	key   string
	err   error
}

func (f *fakeVersionVerifier) Verify(
	_ context.Context,
	apiKey string,
) error {
	f.calls++
	f.key = apiKey
	return f.err
}

type fakeVersionUploader struct {
	calls   int
	input   fileversion.UploadInput
	content string
	saved   *dto.FileInfo
	err     error
}

func (f *fakeVersionUploader) Upload(
	_ context.Context,
	input fileversion.UploadInput,
) (*dto.FileInfo, error) {
	f.calls++
	f.input = input

	if f.err != nil {
		return nil, f.err
	}

	content, err := input.OpenContent()
	if err != nil {
		return nil, err
	}
	defer content.Close()

	data, err := io.ReadAll(content)
	if err != nil {
		return nil, err
	}
	f.content = string(data)

	if f.saved == nil {
		f.saved = &dto.FileInfo{
			ID:           42,
			RelativePath: "p-version",
		}
	}

	return f.saved, nil
}

func performVersionUploadWithHeaders(
	t *testing.T,
	handler *VersionUploadHandler,
	headers map[string]string,
	includeFile bool,
) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if includeFile {
		part, err := writer.CreateFormFile(
			"file",
			"app.apk",
		)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := part.Write([]byte("apk")); err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.WriteField(
		"description",
		"release",
	); err != nil {
		t.Fatal(err)
	}

	contentType := writer.FormDataContentType()

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/upload",
		&body,
	)
	req.Header.Set(
		"Content-Type",
		contentType,
	)

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	rec := httptest.NewRecorder()

	router := gin.New()
	router.POST(
		"/upload",
		handler.UploadVersion,
	)
	router.ServeHTTP(rec, req)

	return rec
}

func performVersionUpload(
	t *testing.T,
	handler *VersionUploadHandler,
	apiKey string,
	includeFile bool,
) *httptest.ResponseRecorder {
	t.Helper()

	headers := map[string]string{}
	if apiKey != "" {
		headers["API-KEY"] = apiKey
	}

	return performVersionUploadWithHeaders(
		t,
		handler,
		headers,
		includeFile,
	)
}

func TestVersionUploadRejectsMissingAPIKeyBeforeEffects(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"",
		true,
	)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusUnauthorized,
		)
	}
	if verifier.calls != 0 {
		t.Fatalf(
			"verifier calls = %d, want 0",
			verifier.calls,
		)
	}
	if uploader.calls != 0 {
		t.Fatalf(
			"uploader calls = %d, want 0",
			uploader.calls,
		)
	}
}

func TestVersionUploadRejectsInvalidAPIKey(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{
		err: versionauth.ErrInvalid,
	}
	uploader := &fakeVersionUploader{}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"bad-version-key-0001",
		true,
	)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusUnauthorized,
		)
	}
	if uploader.calls != 0 {
		t.Fatalf(
			"uploader calls = %d, want 0",
			uploader.calls,
		)
	}
}

func TestVersionUploadFailsClosedWhenVerifierUnavailable(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{
		err: errors.Join(
			versionauth.ErrUnavailable,
			errors.New("hub down"),
		),
	}
	uploader := &fakeVersionUploader{}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"version-api-key-0001",
		true,
	)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusServiceUnavailable,
		)
	}
	if uploader.calls != 0 {
		t.Fatalf(
			"uploader calls = %d, want 0",
			uploader.calls,
		)
	}
}

func TestVersionUploadRequiresMultipartFile(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"version-api-key-0001",
		false,
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
	if uploader.calls != 0 {
		t.Fatalf(
			"uploader calls = %d, want 0",
			uploader.calls,
		)
	}
}

func TestVersionUploadMapsUploaderFailureToBadRequest(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{
		err: errors.New("upload failed"),
	}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"version-api-key-0001",
		true,
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusBadRequest,
		)
	}
}

func TestVersionUploadAcceptsCanonicalAPIKeyHeader(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{}

	rec := performVersionUploadWithHeaders(
		t,
		NewVersionUploadHandler(verifier, uploader),
		map[string]string{
			"X-API-Key": "canonical-version-key-0001",
		},
		true,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d body=%s",
			rec.Code,
			http.StatusOK,
			rec.Body.String(),
		)
	}
	if verifier.key != "canonical-version-key-0001" {
		t.Fatalf("verified key = %q", verifier.key)
	}
}

func TestVersionUploadRejectsAmbiguousAPIKeyHeadersBeforeEffects(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{}

	rec := performVersionUploadWithHeaders(
		t,
		NewVersionUploadHandler(verifier, uploader),
		map[string]string{
			"X-API-Key": "canonical-version-key-0001",
			"API-KEY":   "legacy-version-key-0002",
		},
		true,
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d body=%s",
			rec.Code,
			http.StatusBadRequest,
			rec.Body.String(),
		)
	}
	if verifier.calls != 0 {
		t.Fatalf("verifier calls = %d, want 0", verifier.calls)
	}
	if uploader.calls != 0 {
		t.Fatalf("uploader calls = %d, want 0", uploader.calls)
	}
}

func TestVersionUploadPreservesHTTPContract(
	t *testing.T,
) {
	verifier := &fakeVersionVerifier{}
	uploader := &fakeVersionUploader{}

	rec := performVersionUpload(
		t,
		NewVersionUploadHandler(
			verifier,
			uploader,
		),
		"version-api-key-0001",
		true,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusOK,
		)
	}

	if verifier.calls != 1 {
		t.Fatalf(
			"verifier calls = %d, want 1",
			verifier.calls,
		)
	}
	if verifier.key != "version-api-key-0001" {
		t.Fatalf(
			"verified key = %q",
			verifier.key,
		)
	}

	if uploader.calls != 1 {
		t.Fatalf(
			"uploader calls = %d, want 1",
			uploader.calls,
		)
	}

	if uploader.input.FileName != "app.apk" {
		t.Fatalf(
			"file name = %q",
			uploader.input.FileName,
		)
	}
	if uploader.input.Description != "release" {
		t.Fatalf(
			"description = %q",
			uploader.input.Description,
		)
	}
	if uploader.input.Size != 3 {
		t.Fatalf(
			"size = %d, want 3",
			uploader.input.Size,
		)
	}
	if uploader.content != "apk" {
		t.Fatalf(
			"content = %q",
			uploader.content,
		)
	}

	var response struct {
		Code int `json:"code"`
	}

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatal(err)
	}

	if response.Code != 0 {
		t.Fatalf(
			"code = %d, want 0",
			response.Code,
		)
	}
}
