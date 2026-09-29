package filehttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type fakeVideoPathResolver struct {
	physicalPath string
}

func (f fakeVideoPathResolver) ResolveStoredPath(
	string,
) (string, error) {
	return f.physicalPath, nil
}

func TestVideoReadHandlerRejectsMalformedPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := NewVideoReadHandler("test-xor-key", nil)

	r := gin.New()
	r.GET("/video/:p/:segmentId", h.GetVideo)

	req := httptest.NewRequest(
		http.MethodGet,
		"/video/x/index",
		nil,
	)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			resp.Code,
			http.StatusBadRequest,
		)
	}
}

func TestVideoReadHandlerAcceptsCurrentEncodedVideoPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	const xorKey = "test-xor-key"

	h := NewVideoReadHandler(
		xorKey,
		fakeVideoPathResolver{
			physicalPath: filepath.Join(t.TempDir(), "missing-video"),
		},
	)

	r := gin.New()
	r.GET("/video/:p/:segmentId", h.GetVideo)

	encoded := _utils.XorEncode(
		"files/public/video/2026-08-25/123/42",
		xorKey,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/video/p"+encoded+"/index",
		nil,
	)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	// The path contract was accepted and delivery was reached.
	// The physical test artifact intentionally does not exist.
	if resp.Code != http.StatusNotFound {
		t.Fatalf(
			"status = %d, want %d",
			resp.Code,
			http.StatusNotFound,
		)
	}
}

func TestWriteVideoResponseServesPlaylist(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dir := t.TempDir()
	want := []byte("#EXTM3U\n")

	if err := os.WriteFile(
		filepath.Join(dir, "playlist.m3u8"),
		want,
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = httptest.NewRequest(http.MethodGet, "/video/index", nil)

	writeVideoResponse(c, dir, "index")

	if resp.Code != http.StatusOK {
		t.Fatalf(
			"status = %d, want %d",
			resp.Code,
			http.StatusOK,
		)
	}

	if got := resp.Header().Get("Content-Type"); got != "application/vnd.apple.mpegurl" {
		t.Fatalf("Content-Type = %q", got)
	}

	if got := resp.Body.String(); got != string(want) {
		t.Fatalf("body = %q, want %q", got, want)
	}
}
