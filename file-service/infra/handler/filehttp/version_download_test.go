package filehttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_utils "common/utils"
	"file/internal/fileversion"

	"github.com/gin-gonic/gin"
)

type identityVersionPathResolver struct{}

func (identityVersionPathResolver) ResolveStoredPath(relativePath string) (string, error) {
	return relativePath, nil
}

func testVersionCodec(t *testing.T) *fileversion.ReferenceCodec {
	t.Helper()
	codec, err := fileversion.NewReferenceCodec("version-test-xor-key", "version-test-signature-key")
	if err != nil {
		t.Fatal(err)
	}
	return codec
}

func performVersionDownload(
	t *testing.T,
	handler *VersionDownloadHandler,
	param string,
	rangeHeader string,
) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	c.Request = req
	c.Params = gin.Params{{Key: "p", Value: param}}
	handler.DownloadVersion(c)
	return rec
}

func encodedVersionPath(t *testing.T, path string) string {
	t.Helper()
	reference, err := testVersionCodec(t).Encode(path)
	if err != nil {
		t.Fatal(err)
	}
	return reference
}

func TestVersionDownloadRejectsInvalidOrTamperedReference(t *testing.T) {
	h := NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{})
	for _, reference := range []string{"", "p", "v1.invalid.invalid"} {
		rec := performVersionDownload(t, h, reference, "")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("reference %q status = %d, want %d", reference, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestVersionDownloadReadsLegacyReference(t *testing.T) {
	const xorKey = "version-test-xor-key"
	path := filepath.Join(t.TempDir(), "legacy.apk")
	if err := os.WriteFile(path, []byte("legacy"), 0o600); err != nil {
		t.Fatal(err)
	}
	legacy := "p" + _utils.XorEncode(path, xorKey)
	rec := performVersionDownload(
		t,
		NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{}),
		legacy,
		"",
	)
	if rec.Code != http.StatusOK || rec.Body.String() != "legacy" {
		t.Fatalf("status=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestVersionDownloadReturnsNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.apk")
	rec := performVersionDownload(
		t,
		NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{}),
		encodedVersionPath(t, path),
		"",
	)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestVersionDownloadFullFileContract(t *testing.T) {
	path := filepath.Join(t.TempDir(), "app.apk")
	if err := os.WriteFile(path, []byte("abcdef"), 0o600); err != nil {
		t.Fatal(err)
	}
	rec := performVersionDownload(
		t,
		NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{}),
		encodedVersionPath(t, path),
		"",
	)
	if rec.Code != http.StatusOK || rec.Body.String() != "abcdef" {
		t.Fatalf("status=%d body=%q", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/vnd.android.package-archive" {
		t.Fatalf("content type = %q", got)
	}
	if got := rec.Header().Get("Accept-Ranges"); got != "bytes" {
		t.Fatalf("Accept-Ranges = %q", got)
	}
}

func TestVersionDownloadRangeAndSuffixContracts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "release.zip")
	if err := os.WriteFile(path, []byte("abcdef"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{})
	reference := encodedVersionPath(t, path)

	for _, tc := range []struct {
		rangeHeader  string
		body         string
		contentRange string
	}{
		{"bytes=1-3", "bcd", "bytes 1-3/6"},
		{"bytes=-2", "ef", "bytes 4-5/6"},
		{"bytes=3-", "def", "bytes 3-5/6"},
	} {
		rec := performVersionDownload(t, h, reference, tc.rangeHeader)
		if rec.Code != http.StatusPartialContent {
			t.Fatalf("%s status=%d", tc.rangeHeader, rec.Code)
		}
		if rec.Body.String() != tc.body {
			t.Fatalf("%s body=%q want=%q", tc.rangeHeader, rec.Body.String(), tc.body)
		}
		if got := rec.Header().Get("Content-Range"); got != tc.contentRange {
			t.Fatalf("%s Content-Range=%q", tc.rangeHeader, got)
		}
	}
}

func TestVersionDownloadInvalidRangeReturns416(t *testing.T) {
	path := filepath.Join(t.TempDir(), "release.zip")
	if err := os.WriteFile(path, []byte("abcdef"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{})
	reference := encodedVersionPath(t, path)

	for _, rangeHeader := range []string{"invalid", "bytes=99-100", "bytes=2-1", "bytes=1-2,4-5"} {
		rec := performVersionDownload(t, h, reference, rangeHeader)
		if rec.Code != http.StatusRequestedRangeNotSatisfiable {
			t.Fatalf("%s status = %d, want 416", rangeHeader, rec.Code)
		}
		if got := rec.Header().Get("Content-Range"); got != "bytes */6" {
			t.Fatalf("%s Content-Range = %q", rangeHeader, got)
		}
	}
}

func TestVersionDownloadEmptyArtifactIsValid(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.zip")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	rec := performVersionDownload(
		t,
		NewVersionDownloadHandler(testVersionCodec(t), identityVersionPathResolver{}),
		encodedVersionPath(t, path),
		"",
	)
	if rec.Code != http.StatusOK || rec.Body.Len() != 0 {
		t.Fatalf("status=%d body=%q", rec.Code, rec.Body.String())
	}
}
