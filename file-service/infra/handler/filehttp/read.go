package filehttp

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_identity "common/identity"
	_utils "common/utils"

	"file/internal/fileauthorization"
	"file/internal/filemedia"

	"github.com/gin-gonic/gin"
)

const maxConcurrentExternalPrivateReads = 10

// serviceRequestVerifier is the exact transport capability required to
// recognize an authenticated internal-service request.
type serviceRequestVerifier interface {
	Verify(encodedToken string, serviceName string, now time.Time) bool
}

// ordinaryPrivateAccess is the exact ordinary private-file access capability
// consumed by the HTTP read adapter.
type ordinaryPrivateAccess interface {
	Allows(
		decodedPath string,
		encodedToken string,
		now time.Time,
	) (bool, error)
}

// ReadHandler owns HTTP semantics for ordinary file reads.
//
// Business authorization remains in fileauthorization. Physical path resolution is consumed only through the storage capability.
type ReadHandler struct {
	serviceAuth  serviceRequestVerifier
	ordinary     ordinaryPrivateAccess
	permission   fileauthorization.Authorizer
	pathResolver storedPathResolver

	xorKey string

	externalPrivateReads chan struct{}
}

func NewReadHandler(
	serviceAuth serviceRequestVerifier,
	ordinary ordinaryPrivateAccess,
	permission fileauthorization.Authorizer,
	xorKey string,
	pathResolver storedPathResolver,
) *ReadHandler {
	if permission == nil {
		permission = fileauthorization.UnavailableAuthorizer{}
	}

	return &ReadHandler{
		serviceAuth:  serviceAuth,
		ordinary:     ordinary,
		permission:   permission,
		pathResolver: pathResolver,
		xorKey:       xorKey,

		externalPrivateReads: make(
			chan struct{},
			maxConcurrentExternalPrivateReads,
		),
	}
}

// GetFile translates the public HTTP read contract into the canonical
// authorization capabilities.
//
// @Summary Tải file
// @Description Tải file
// @Param s path string true "Parameter S"
// @Param p path string true "Parameter P"
// @Produce json
// @Router /load/{p}/{s} [get]
func (h *ReadHandler) GetFile(c *gin.Context) {
	if h == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{"error": "file read service is unavailable"},
		)
		return
	}

	p := c.Param("p")
	s := c.Param("s")

	if p == "" || s == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "p and s is required"},
		)
		return
	}

	if len(p) < 2 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Đường dẫn không chính xác"},
		)
		return
	}

	prefix := string(p[0])
	encodedPath := p[1:]

	decodedPath, err := _utils.XorDecode(
		encodedPath,
		h.xorKey,
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err},
		)
		return
	}

	// Media path interpretation is owned by filemedia.
	fileInfo, err := filemedia.ParsePath(decodedPath)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err},
		)
		return
	}

	if fileInfo.Type == "video" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Đường dẫn không dành cho video"},
		)
		return
	}

	ext := _utils.GetFileExtension(decodedPath)

	if h.pathResolver == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{"error": "file storage unavailable"},
		)
		return
	}

	physicalPath, err := h.pathResolver.ResolveStoredPath(decodedPath)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Đường dẫn không hợp lệ"},
		)
		return
	}

	if prefix == "p" {
		writeFileResponse(c, physicalPath, ext)
		return
	}

	if h.serviceAuth != nil &&
		h.serviceAuth.Verify(
			c.GetHeader("X-Service-Auth"),
			c.GetHeader("X-Service-Name"),
			time.Now(),
		) {
		writeFileResponse(c, physicalPath, ext)
		return
	}

	// Preserve the existing limiter semantics exactly: only external private
	// reads consume this admission slot. Public and authenticated internal
	// service reads bypass it.
	h.externalPrivateReads <- struct{}{}
	defer func() {
		<-h.externalPrivateReads
	}()

	permissionAuthority := h.permission

	if actor, ok := _identity.ActorFromContext(
		c.Request.Context(),
	); !ok || actor.ProfileID == 0 {
		permissionAuthority =
			fileauthorization.DeniedAuthorizer{}
	}

	authorizationErr :=
		fileauthorization.AuthorizePrivateRead(
			c.Request.Context(),
			func() (bool, error) {
				if h.ordinary == nil {
					return false, nil
				}

				return h.ordinary.Allows(
					decodedPath,
					s,
					time.Now(),
				)
			},
			permissionAuthority,
		)

	switch {
	case authorizationErr == nil:
		writeFileResponse(c, physicalPath, ext)

	case errors.Is(
		authorizationErr,
		fileauthorization.ErrUnavailable,
	):
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"error": "Permission authority unavailable",
			},
		)

	case errors.Is(
		authorizationErr,
		fileauthorization.ErrDenied,
	):
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "Bạn không có quyền truy cập file",
			},
		)

	default:
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Không thể kiểm tra quyền truy cập file",
			},
		)
	}
}

var _ fileauthorization.OrdinaryCheck = func() (bool, error) {
	return false, nil
}

// writeFileResponse owns ordinary HTTP file delivery semantics. Delivery is
// streamed through net/http so malformed or suffix Range requests use the
// standard 206/416 behavior rather than a hand-written parser/buffer.
func writeFileResponse(c *gin.Context, path string, ext string) {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	file, err := os.Open(path)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": "File không tồn tại"})
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể đọc file"})
		return
	}
	if !info.Mode().IsRegular() {
		c.JSON(http.StatusNotFound, gin.H{"error": "File không tồn tại"})
		return
	}

	contentType := ordinaryContentType(ext)
	c.Header("Content-Type", contentType)
	if ordinaryInlineExtension(ext) {
		c.Header(
			"Content-Disposition",
			fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(path)),
		)
	}

	http.ServeContent(
		c.Writer,
		c.Request,
		filepath.Base(path),
		info.ModTime(),
		file,
	)
}

func ordinaryContentType(ext string) string {
	contentTypes := map[string]string{
		"mp4":  "video/mp4",
		"mov":  "video/quicktime",
		"webm": "video/webm",
		"ogg":  "video/ogg",
		"png":  "image/png",
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"gif":  "image/gif",
		"webp": "image/webp",
		"bmp":  "image/bmp",
		"svg":  "image/svg+xml",
		"pdf":  "application/pdf",
		"txt":  "text/plain; charset=utf-8",
		"json": "application/json",
		"xml":  "application/xml",
		"html": "text/html; charset=utf-8",
		"htm":  "text/html; charset=utf-8",
		"csv":  "text/csv; charset=utf-8",
		"m3u8": "application/vnd.apple.mpegurl",
	}
	if value, ok := contentTypes[ext]; ok {
		return value
	}
	return "application/octet-stream"
}

func ordinaryInlineExtension(ext string) bool {
	switch ext {
	case "pdf", "png", "jpg", "jpeg", "gif", "webp", "bmp", "svg", "txt", "html", "htm":
		return true
	default:
		return false
	}
}
