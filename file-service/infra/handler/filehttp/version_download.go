package filehttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxVersionDownloadChunkSize = int64(10 * 1024 * 1024)

var errVersionRangeNotSatisfiable = errors.New("version range not satisfiable")

type versionReferenceDecoder interface {
	Decode(reference string) (string, error)
}

// VersionDownloadHandler owns version-reference translation and HTTP delivery.
// The physical storage location remains behind storedPathResolver.
type VersionDownloadHandler struct {
	reference    versionReferenceDecoder
	pathResolver storedPathResolver
}

func NewVersionDownloadHandler(
	reference versionReferenceDecoder,
	pathResolver storedPathResolver,
) *VersionDownloadHandler {
	return &VersionDownloadHandler{reference: reference, pathResolver: pathResolver}
}

func (h *VersionDownloadHandler) DownloadVersion(c *gin.Context) {
	if h == nil || h.reference == nil || h.pathResolver == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "version delivery unavailable"})
		return
	}

	decodedPath, err := h.reference.Decode(c.Param("p"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không hợp lệ"})
		return
	}

	physicalPath, err := h.pathResolver.ResolveStoredPath(decodedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Đường dẫn không hợp lệ"})
		return
	}

	if err := writeVersionDownload(c, physicalPath); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
	}
}

func writeVersionDownload(c *gin.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("version artifact is not a regular file")
	}

	fileSize := info.Size()
	rangeHeader := strings.TrimSpace(c.GetHeader("Range"))
	start, end, ranged, err := parseVersionDownloadRange(rangeHeader, fileSize)
	if err != nil {
		if errors.Is(err, errVersionRangeNotSatisfiable) {
			c.Header("Accept-Ranges", "bytes")
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			c.Writer.WriteHeaderNow()
			return nil
		}
		return err
	}

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", versionDownloadContentType(path))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(path)))

	if fileSize == 0 {
		c.Header("Content-Length", "0")
		c.Status(http.StatusOK)
		return nil
	}

	if end-start+1 > maxVersionDownloadChunkSize {
		end = start + maxVersionDownloadChunkSize - 1
		ranged = true
	}
	if end >= fileSize {
		end = fileSize - 1
	}
	length := end - start + 1

	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return err
	}

	status := http.StatusOK
	if ranged {
		status = http.StatusPartialContent
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	}
	c.Header("Content-Length", strconv.FormatInt(length, 10))
	c.Status(status)

	_, err = io.CopyN(c.Writer, file, length)
	return err
}

func parseVersionDownloadRange(
	rangeHeader string,
	fileSize int64,
) (start int64, end int64, ranged bool, err error) {
	if rangeHeader == "" {
		if fileSize == 0 {
			return 0, -1, false, nil
		}
		return 0, fileSize - 1, false, nil
	}
	if fileSize <= 0 || !strings.HasPrefix(rangeHeader, "bytes=") {
		return 0, 0, false, errVersionRangeNotSatisfiable
	}

	spec := strings.TrimSpace(strings.TrimPrefix(rangeHeader, "bytes="))
	if spec == "" || strings.Contains(spec, ",") {
		return 0, 0, false, errVersionRangeNotSatisfiable
	}
	parts := strings.Split(spec, "-")
	if len(parts) != 2 {
		return 0, 0, false, errVersionRangeNotSatisfiable
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])
	if left == "" {
		// RFC 7233 suffix-byte-range-spec: bytes=-N
		suffix, parseErr := strconv.ParseInt(right, 10, 64)
		if parseErr != nil || suffix <= 0 {
			return 0, 0, false, errVersionRangeNotSatisfiable
		}
		if suffix > fileSize {
			suffix = fileSize
		}
		return fileSize - suffix, fileSize - 1, true, nil
	}

	start, parseErr := strconv.ParseInt(left, 10, 64)
	if parseErr != nil || start < 0 || start >= fileSize {
		return 0, 0, false, errVersionRangeNotSatisfiable
	}

	if right == "" {
		return start, fileSize - 1, true, nil
	}
	end, parseErr = strconv.ParseInt(right, 10, 64)
	if parseErr != nil || end < start {
		return 0, 0, false, errVersionRangeNotSatisfiable
	}
	if end >= fileSize {
		end = fileSize - 1
	}
	return start, end, true, nil
}

func versionDownloadContentType(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".apk":
		return "application/vnd.android.package-archive"
	case ".ipa":
		return "application/octet-stream"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}
