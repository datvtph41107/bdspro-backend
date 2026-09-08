package filehttp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	_utils "common/utils"

	"file/internal/fileauthorization"
	"file/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHTTPReadXORKey = "http-read-xor-key"

type fakeServiceRequestVerifier struct {
	allowed bool
	calls   int
}

func (f *fakeServiceRequestVerifier) Verify(
	_ string,
	_ string,
	_ time.Time,
) bool {
	f.calls++
	return f.allowed
}

type fakeOrdinaryPrivateAccess struct {
	allowed bool
	err     error
	calls   int
}

func (f *fakeOrdinaryPrivateAccess) Allows(
	_ string,
	_ string,
	_ time.Time,
) (bool, error) {
	f.calls++
	return f.allowed, f.err
}

type fakeGlobalReadAuthorizer struct {
	err   error
	calls int
}

type fakeStoredPathResolver struct {
	storedPath   string
	physicalPath string
	err          error
	calls        int
}

func (f *fakeStoredPathResolver) ResolveStoredPath(
	storedPath string,
) (string, error) {
	f.calls++
	if f.err != nil {
		return "", f.err
	}
	if storedPath != f.storedPath {
		return "", errors.New("unexpected stored path")
	}
	return f.physicalPath, nil
}

func (f *fakeGlobalReadAuthorizer) Require(
	context.Context,
	fileauthorization.Capability,
) error {
	f.calls++
	return f.err
}

func testReadableFile(
	t *testing.T,
	content string,
) (string, string) {
	t.Helper()

	storedPath := filepath.ToSlash(filepath.Join(
		"files",
		"public",
		"document",
		"2026-08-25",
		"uuid",
		"42.pdf",
	))
	physicalPath := filepath.Join(
		t.TempDir(),
		filepath.FromSlash(storedPath),
	)

	require.NoError(
		t,
		os.MkdirAll(filepath.Dir(physicalPath), 0o755),
	)
	require.NoError(
		t,
		os.WriteFile(physicalPath, []byte(content), 0o600),
	)

	return storedPath, physicalPath
}

func runReadRequest(
	t *testing.T,
	h *ReadHandler,
	p string,
	s string,
) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)

	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	c.Params = gin.Params{
		{Key: "p", Value: p},
		{Key: "s", Value: s},
	}

	h.GetFile(c)

	return recorder
}

func TestReadHandlerPublicReadBypassesAuthorization(t *testing.T) {
	storedPath, physicalPath := testReadableFile(t, "public-content")

	serviceAuth := &fakeServiceRequestVerifier{}
	ordinary := &fakeOrdinaryPrivateAccess{}
	global := &fakeGlobalReadAuthorizer{}

	h := NewReadHandler(
		serviceAuth,
		ordinary,
		global,
		testHTTPReadXORKey,
		&fakeStoredPathResolver{
			storedPath:   storedPath,
			physicalPath: physicalPath,
		},
	)

	p := "p" + _utils.XorEncode(
		storedPath,
		testHTTPReadXORKey,
	)

	response := runReadRequest(
		t,
		h,
		p,
		"unused",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "public-content", response.Body.String())
	assert.Zero(t, serviceAuth.calls)
	assert.Zero(t, ordinary.calls)
	assert.Zero(t, global.calls)
}

func TestReadHandlerResolvesStableWirePathAgainstConfiguredRoot(
	t *testing.T,
) {
	root := filepath.Join(t.TempDir(), "mounted-file-root")
	storage, err := services.NewStorageService(root)
	require.NoError(t, err)

	storedPath := "files/public/document/2026-08-25/uuid/42.pdf"
	physicalPath, err := storage.ResolveStoredPath(storedPath)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(physicalPath), 0o755))
	require.NoError(
		t,
		os.WriteFile(physicalPath, []byte("configured-root"), 0o600),
	)

	h := NewReadHandler(
		&fakeServiceRequestVerifier{},
		&fakeOrdinaryPrivateAccess{},
		&fakeGlobalReadAuthorizer{},
		testHTTPReadXORKey,
		storage,
	)

	p := "p" + _utils.XorEncode(storedPath, testHTTPReadXORKey)
	response := runReadRequest(t, h, p, "unused")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "configured-root", response.Body.String())
}

func TestReadHandlerInternalServiceBypassesPrivateAuthorization(t *testing.T) {
	storedPath, physicalPath := testReadableFile(t, "internal-content")

	serviceAuth := &fakeServiceRequestVerifier{
		allowed: true,
	}
	ordinary := &fakeOrdinaryPrivateAccess{}
	global := &fakeGlobalReadAuthorizer{}

	h := NewReadHandler(
		serviceAuth,
		ordinary,
		global,
		testHTTPReadXORKey,
		&fakeStoredPathResolver{
			storedPath:   storedPath,
			physicalPath: physicalPath,
		},
	)

	p := "s" + _utils.XorEncode(
		storedPath,
		testHTTPReadXORKey,
	)

	response := runReadRequest(
		t,
		h,
		p,
		"unused",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "internal-content", response.Body.String())
	assert.Equal(t, 1, serviceAuth.calls)
	assert.Zero(t, ordinary.calls)
	assert.Zero(t, global.calls)
}

func TestReadHandlerOrdinarySignedAccessAllowsPrivateRead(t *testing.T) {
	storedPath, physicalPath := testReadableFile(t, "private-content")

	serviceAuth := &fakeServiceRequestVerifier{}
	ordinary := &fakeOrdinaryPrivateAccess{
		allowed: true,
	}
	global := &fakeGlobalReadAuthorizer{}

	h := NewReadHandler(
		serviceAuth,
		ordinary,
		global,
		testHTTPReadXORKey,
		&fakeStoredPathResolver{
			storedPath:   storedPath,
			physicalPath: physicalPath,
		},
	)

	p := "s" + _utils.XorEncode(
		storedPath,
		testHTTPReadXORKey,
	)

	response := runReadRequest(
		t,
		h,
		p,
		"signed-token",
	)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "private-content", response.Body.String())
	assert.Equal(t, 1, serviceAuth.calls)
	assert.Equal(t, 1, ordinary.calls)

	// Anonymous callers must not fan out to the global IAM authority.
	assert.Zero(t, global.calls)
}

func TestReadHandlerAnonymousPrivateReadDeniesWithoutOrdinaryEvidence(t *testing.T) {
	storedPath, physicalPath := testReadableFile(t, "private-content")

	h := NewReadHandler(
		&fakeServiceRequestVerifier{},
		&fakeOrdinaryPrivateAccess{},
		&fakeGlobalReadAuthorizer{},
		testHTTPReadXORKey,
		&fakeStoredPathResolver{
			storedPath:   storedPath,
			physicalPath: physicalPath,
		},
	)

	p := "s" + _utils.XorEncode(
		storedPath,
		testHTTPReadXORKey,
	)

	response := runReadRequest(
		t,
		h,
		p,
		"invalid-or-unowned-token",
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		response.Code,
	)
}

func TestReadHandlerPreservesOrdinaryAccessFailure(t *testing.T) {
	storedPath, physicalPath := testReadableFile(t, "private-content")

	lookupErr := errors.New("postgres unavailable")

	h := NewReadHandler(
		&fakeServiceRequestVerifier{},
		&fakeOrdinaryPrivateAccess{
			err: lookupErr,
		},
		&fakeGlobalReadAuthorizer{},
		testHTTPReadXORKey,
		&fakeStoredPathResolver{
			storedPath:   storedPath,
			physicalPath: physicalPath,
		},
	)

	p := "s" + _utils.XorEncode(
		storedPath,
		testHTTPReadXORKey,
	)

	response := runReadRequest(
		t,
		h,
		p,
		"signed-token",
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		response.Code,
	)
}
