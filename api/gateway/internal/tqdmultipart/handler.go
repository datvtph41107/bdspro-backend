// Package tqdmultipart owns bounded HTTP multipart compatibility adapters for TQD.
package tqdmultipart

import (
	"errors"
	"fmt"
	"net/http"

	"gateway/internal/tqdtransport"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Limits struct {
	MaxBodyBytes       int64
	MaxMemoryBytes     int64
	MaxConcurrentCalls int
}

func DefaultLimits() Limits {
	return Limits{
		MaxBodyBytes:       tqdtransport.MaxGRPCMessageBytes,
		MaxMemoryBytes:     32 << 20,
		MaxConcurrentCalls: 1,
	}
}

type Handler struct {
	conn   *grpc.ClientConn
	limits Limits
	slots  chan struct{}
}

func New(conn *grpc.ClientConn, limits Limits) (*Handler, error) {
	if conn == nil {
		return nil, errors.New("TQD multipart connection is required")
	}
	if limits.MaxBodyBytes <= 0 || limits.MaxMemoryBytes <= 0 || limits.MaxConcurrentCalls <= 0 {
		return nil, errors.New("TQD multipart limits must be positive")
	}
	if limits.MaxBodyBytes > tqdtransport.MaxGRPCMessageBytes {
		return nil, fmt.Errorf("TQD multipart body limit exceeds gRPC message limit")
	}
	return &Handler{conn: conn, limits: limits, slots: make(chan struct{}, limits.MaxConcurrentCalls)}, nil
}

func (h *Handler) Register(routes gin.IRoutes) {
	routes.POST("/tqd/admin/import/geojson", h.importGeoJSON)
	routes.POST("/tqd/admin/planning-projects/from-folder", h.createPlanningProjectFromFolder)
}

func (h *Handler) acquire(c *gin.Context) bool {
	select {
	case h.slots <- struct{}{}:
		return true
	default:
		c.Header("Retry-After", "1")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "multipart capacity unavailable"})
		return false
	}
}

func (h *Handler) release() { <-h.slots }

func (h *Handler) parse(c *gin.Context) (func(), bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.limits.MaxBodyBytes)
	if err := c.Request.ParseMultipartForm(h.limits.MaxMemoryBytes); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "multipart body too large"})
			return nil, false
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart body"})
		return nil, false
	}
	cleanup := func() {}
	if c.Request.MultipartForm != nil {
		form := c.Request.MultipartForm
		cleanup = func() { _ = form.RemoveAll() }
	}
	return cleanup, true
}
