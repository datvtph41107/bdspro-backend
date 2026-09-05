package filehttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"file/internal/filecontent"
	"file/internal/serviceauth"

	"github.com/gin-gonic/gin"
)

const (
	multipartMemoryBytes = 8 << 20
	maxHTTPOverheadBytes = 1 << 20
)

type fileOperations interface {
	UploadFile(
		context.Context,
		filecontent.UploadInput,
	) (filecontent.SavedFile, error)

	PutOwnedFile(
		context.Context,
		filecontent.PutOwnedInput,
	) (filecontent.SavedFile, error)
}

type Handler struct {
	files fileOperations
	auth  serviceauth.Verifier
}

func NewHandler(
	files fileOperations,
	auth serviceauth.Verifier,
) *Handler {
	return &Handler{
		files: files,
		auth:  auth,
	}
}

type ownedFileResponse struct {
	FileID      uint64 `json:"fileId"`
	Path        string `json:"path"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func (h *Handler) PutOwnedFile(c *gin.Context) {
	if h == nil || h.files == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{"error": "owned file service is unavailable"},
		)
		return
	}

	if !h.auth.Verify(
		c.GetHeader("X-Service-Auth"),
		c.GetHeader("X-Service-Name"),
		time.Now(),
	) {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "internal service authentication failed"},
		)
		return
	}

	maxBodyBytes :=
		int64(filecontent.MaxContentBytes) +
			maxHTTPOverheadBytes

	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		maxBodyBytes,
	)

	if err :=
		c.Request.ParseMultipartForm(
			multipartMemoryBytes,
		); err != nil {
		var tooLarge *http.MaxBytesError

		if errors.As(err, &tooLarge) {
			c.JSON(
				http.StatusRequestEntityTooLarge,
				gin.H{"error": "file content exceeds limit"},
			)
			return
		}

		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid multipart request"},
		)
		return
	}

	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	ownerNamespace :=
		strings.TrimSpace(
			c.Request.FormValue("owner_namespace"),
		)

	ownerKey :=
		strings.TrimSpace(
			c.Request.FormValue("owner_key"),
		)

	description :=
		c.Request.FormValue("description")

	contentType :=
		strings.TrimSpace(
			c.Request.FormValue("content_type"),
		)

	header, err :=
		c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "file is required"},
		)
		return
	}

	reader, err :=
		header.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "cannot open uploaded file"},
		)
		return
	}
	defer reader.Close()

	content, err :=
		io.ReadAll(
			io.LimitReader(
				reader,
				int64(filecontent.MaxContentBytes)+1,
			),
		)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "cannot read uploaded file"},
		)
		return
	}

	if len(content) > filecontent.MaxContentBytes {
		c.JSON(
			http.StatusRequestEntityTooLarge,
			gin.H{"error": "file content exceeds limit"},
		)
		return
	}

	if contentType == "" {
		contentType =
			strings.TrimSpace(
				header.Header.Get("Content-Type"),
			)
	}
	if contentType == "" ||
		contentType == "application/octet-stream" {
		contentType =
			http.DetectContentType(content)
	}

	saved, err :=
		h.files.PutOwnedFile(
			c.Request.Context(),
			filecontent.PutOwnedInput{
				Content:        content,
				Filename:       filepath.Base(header.Filename),
				ContentType:    contentType,
				Description:    description,
				OwnerNamespace: ownerNamespace,
				OwnerKey:       ownerKey,
			},
		)
	if err != nil {
		switch {
		case errors.Is(
			err,
			filecontent.ErrInvalidOwnedFile,
		):
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "owned file request is invalid"},
			)

		case errors.Is(
			err,
			filecontent.ErrOwnedFileConflict,
		):
			c.JSON(
				http.StatusConflict,
				gin.H{"error": "owned file conflicts with existing effect"},
			)

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "cannot store owned file"},
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		ownedFileResponse{
			FileID:      saved.ID,
			Path:        saved.Path,
			Filename:    saved.Filename,
			ContentType: saved.ContentType,
			Size:        saved.Size,
		},
	)
}
