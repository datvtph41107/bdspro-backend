package handlers

import (
	"encoding/json"
	"file/internal/ports"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService ports.FileService
}

func NewFileHandler(fileService ports.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

// Upload file
// @Summary Upload a file
// @Description Upload a file with optional metadata
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param service_id formData string false "Service ID"
// @Param reference_type formData string false "Reference type"
// @Param reference_id formData string false "Reference ID"
// @Param is_public formData bool false "Is public file"
// @Param metadata formData string false "JSON metadata"
// @Security BearerAuth
// @Success 201 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/files [post]
func (h *FileHandler) Upload(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil { // 100MB
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to parse form",
			Error:   err.Error(),
		})
		return
	}

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "No file provided",
			Error:   err.Error(),
		})
		return
	}
	defer file.Close()

	// Get user ID from context
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}
	userID, ok := val.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid user_id",
		})
		return
	}

	// Get other form values
	serviceID := c.PostForm("service_id")
	referenceType := c.PostForm("reference_type")
	referenceID := c.PostForm("reference_id")
	isPublic := c.PostForm("is_public") == "true"

	// Parse metadata
	var metadata map[string]interface{}
	if metadataStr := c.PostForm("metadata"); metadataStr != "" {
		if err := json.Unmarshal([]byte(metadataStr), &metadata); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid metadata format",
				Error:   err.Error(),
			})
			return
		}
	}

	// Create upload request
	req := ports.UploadRequest{
		File:          file,
		FileHeader:    header,
		OriginalName:  header.Filename,
		ServiceID:     serviceID,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		UserID:        userID,
		IsPublic:      isPublic,
		Metadata:      metadata,
	}

	// Call service
	resp, err := h.fileService.Upload(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "validation") {
			status = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "size limit") {
			status = http.StatusRequestEntityTooLarge
		}

		c.JSON(status, ErrorResponse{
			Code:    status,
			Message: "Upload failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Get file info
// @Summary Get file information
// @Description Get metadata for a file
// @Tags files
// @Produce json
// @Param file_id path string true "File ID"
// @Security BearerAuth
// @Success 200 {object} FileInfo
// @Failure 404 {object} ErrorResponse
// @Router /v1/files/{file_id} [get]
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	fileID := c.Param("file_id")

	info, err := h.fileService.GetFileInfo(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "File not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// Download file
// @Summary Download a file
// @Description Download file with optional byte range
// @Tags files
// @Produce octet-stream
// @Param file_id path string true "File ID"
// @Param range header string false "Byte range"
// @Param signature query string false "Access signature"
// @Security BearerAuth
// @Success 200 {file} binary
// @Success 206 {file} binary
// @Failure 404 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /v1/files/{file_id}/download [get]
func (h *FileHandler) Download(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "file_id is required",
		})
		return
	}
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}
	userID, ok := val.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid user_id",
		})
		return
	}
	signature := c.Query("signature")
	rangeHeader := c.GetHeader("Range")
	req := ports.DownloadRequest{
		FileID:    fileID,
		UserID:    userID,
		Range:     rangeHeader,
		Signature: signature,
	}

	reader, fileInfo, err := h.fileService.Download(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(err.Error(), "access denied") ||
			strings.Contains(err.Error(), "invalid signature") {
			status = http.StatusForbidden
		}

		c.JSON(status, ErrorResponse{
			Code:    status,
			Message: "Download failed",
			Error:   err.Error(),
		})
		return
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	// Set headers
	c.Header("Content-Type", fileInfo.MimeType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size, 10))
	c.Header("Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.OriginalName))

	if rangeHeader != "" {
		c.Header("Accept-Ranges", "bytes")
		c.Status(http.StatusPartialContent)
	} else {
		c.Status(http.StatusOK)
	}

	// Stream file
	io.Copy(c.Writer, reader)
}

// Generate access URL
// @Summary Generate access URL
// @Description Generate temporary access URL for file
// @Tags files
// @Produce json
// @Param file_id path string true "File ID"
// @Param expiry query int false "Expiry in seconds (default: 300)"
// @Param permission query string false "Permission (read/write/delete)"
// @Security BearerAuth
// @Success 200 {object} AccessURLResponse
// @Failure 403 {object} ErrorResponse
// @Router /v1/files/{file_id}/access-url [post]
func (h *FileHandler) GenerateAccessURL(c *gin.Context) {
	fileID := c.Param("file_id")
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}
	userID, ok := val.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid user_id",
		})
		return
	}

	expirySeconds, _ := strconv.Atoi(c.DefaultQuery("expiry", "300"))
	permission := c.DefaultQuery("permission", "read")

	req := ports.AccessURLRequest{
		FileID:     fileID,
		UserID:     userID,
		Permission: permission,
		Expiry:     time.Duration(expirySeconds) * time.Second,
	}

	url, err := h.fileService.GenerateAccessURL(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "Failed to generate access URL",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AccessURLResponse{
		URL:    url,
		Expiry: expirySeconds,
	})
}

// List files
// @Summary List files
// @Description List files with filtering and pagination
// @Tags files
// @Produce json
// @Param service_id query string false "Service ID"
// @Param reference_type query string false "Reference type"
// @Param reference_id query string false "Reference ID"
// @Param mime_type query string false "MIME type filter"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Security BearerAuth
// @Success 200 {object} ListResponse
// @Router /v1/files [get]
func (h *FileHandler) ListFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	req := ports.ListRequest{
		ServiceID:     c.Query("service_id"),
		ReferenceType: c.Query("reference_type"),
		ReferenceID:   c.Query("reference_id"),
		MimeType:      c.Query("mime_type"),
		Page:          page,
		PageSize:      pageSize,
		SortBy:        c.DefaultQuery("sort_by", "created_at"),
		SortDir:       c.DefaultQuery("sort_dir", "desc"),
	}

	resp, err := h.fileService.ListFiles(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list files",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Delete file
// @Summary Delete a file
// @Description Delete file (soft delete by default)
// @Tags files
// @Param file_id path string true "File ID"
// @Param force query bool false "Force hard delete"
// @Security BearerAuth
// @Success 204
// @Failure 403 {object} ErrorResponse
// @Router /v1/files/{file_id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("file_id")
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "unauthorized",
		})
		return
	}
	userID, ok := val.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "invalid user_id",
		})
		return
	}
	force := c.Query("force") == "true"

	req := ports.DeleteRequest{
		FileID: fileID,
		UserID: userID,
		Force:  force,
	}

	if err := h.fileService.DeleteFile(c.Request.Context(), req); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "permission denied") {
			status = http.StatusForbidden
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}

		c.JSON(status, ErrorResponse{
			Code:    status,
			Message: "Delete failed",
			Error:   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// Chunk upload
// @Summary Upload file chunk
// @Description Upload a chunk of a file (for resumable uploads)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param session_id formData string true "Upload session ID"
// @Param chunk_number formData int true "Chunk number"
// @Param total_chunks formData int true "Total chunks"
// @Param chunk formData file true "Chunk data"
// @Param file_hash formData string false "File hash"
// @Security BearerAuth
// @Success 200 {object} UploadChunkResponse
// @Router /v1/files/chunk [post]
func (h *FileHandler) UploadChunk(c *gin.Context) {
	// Parse chunk upload
	chunkNumber, _ := strconv.Atoi(c.PostForm("chunk_number"))
	totalChunks, _ := strconv.Atoi(c.PostForm("total_chunks"))
	sessionID := c.PostForm("session_id")
	fileHash := c.PostForm("file_hash")

	chunkFile, _, err := c.Request.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "No chunk data",
		})
		return
	}
	defer chunkFile.Close()

	// Read chunk data
	chunkData, err := io.ReadAll(chunkFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to read chunk",
		})
		return
	}

	req := ports.UploadChunkRequest{
		SessionID:   sessionID,
		ChunkNumber: chunkNumber,
		TotalChunks: totalChunks,
		ChunkData:   chunkData,
		FileHash:    fileHash,
	}

	resp, err := h.fileService.UploadChunk(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Chunk upload failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Response types
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type AccessURLResponse struct {
	URL    string `json:"url"`
	Expiry int    `json:"expiry_seconds"`
}
