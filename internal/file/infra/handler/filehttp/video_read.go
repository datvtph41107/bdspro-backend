package filehttp

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	_utils "common/utils"

	"file/internal/filemedia"

	"github.com/gin-gonic/gin"
)

var hlsSegmentPattern = regexp.MustCompile(`^[0-9]{8}\.ts$`)

type VideoReadHandler struct {
	xorKey       string
	pathResolver storedPathResolver
}

func NewVideoReadHandler(
	xorKey string,
	pathResolver storedPathResolver,
) *VideoReadHandler {
	return &VideoReadHandler{xorKey: xorKey, pathResolver: pathResolver}
}

// GetVideo serves one HLS playlist/segment from an opaque public File video
// reference. Segment selection is intentionally constrained to the exact names
// emitted by filemedia.Processor.
func (h *VideoReadHandler) GetVideo(c *gin.Context) {
	p := c.Param("p")
	segmentID := c.Param("segmentId")
	if p == "" || segmentID == "" || len(p) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video reference"})
		return
	}
	if segmentID != "index" && !hlsSegmentPattern.MatchString(segmentID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video segment"})
		return
	}

	decodedPath, err := _utils.XorDecode(p[1:], h.xorKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video reference"})
		return
	}

	fileInfo, err := filemedia.ParsePath(decodedPath)
	if err != nil || fileInfo.Type != "video" || fileInfo.Scope != "public" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video path"})
		return
	}
	if h == nil || h.pathResolver == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "video storage unavailable"})
		return
	}

	physicalPath, err := h.pathResolver.ResolveStoredPath(decodedPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video path"})
		return
	}
	writeVideoResponse(c, physicalPath, segmentID)
}

func writeVideoResponse(c *gin.Context, videoDirectory string, segmentID string) {
	name := segmentID
	contentType := "video/mp2t"
	if segmentID == "index" {
		name = "playlist.m3u8"
		contentType = "application/vnd.apple.mpegurl"
	}
	if name != "playlist.m3u8" && !hlsSegmentPattern.MatchString(name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video segment"})
		return
	}

	path := filepath.Join(videoDirectory, name)
	file, err := os.Open(path)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": "video segment not found"})
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		c.JSON(http.StatusNotFound, gin.H{"error": "video segment not found"})
		return
	}

	c.Header("Content-Type", contentType)
	http.ServeContent(c.Writer, c.Request, name, info.ModTime(), file)
}
