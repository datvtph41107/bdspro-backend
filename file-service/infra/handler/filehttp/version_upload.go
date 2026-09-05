package filehttp

import (
	"context"
	"errors"
	"io"
	"net/http"

	_jwt "common/jwt"
	"file/dto"
	"file/internal/fileversion"
	"file/internal/versionauth"

	"github.com/gin-gonic/gin"
)

type versionUploader interface {
	Upload(
		ctx context.Context,
		input fileversion.UploadInput,
	) (*dto.FileInfo, error)
}

// VersionUploadHandler owns the HTTP contract for publishing version
// artifacts: API-key evidence, multipart translation, and HTTP response
// mapping. Durable metadata/storage effects belong to fileversion.Uploader.
type VersionUploadHandler struct {
	verifier versionauth.Verifier
	uploader versionUploader
}

func NewVersionUploadHandler(
	verifier versionauth.Verifier,
	uploader versionUploader,
) *VersionUploadHandler {
	return &VersionUploadHandler{
		verifier: verifier,
		uploader: uploader,
	}
}

func (h *VersionUploadHandler) UploadVersion(
	c *gin.Context,
) {
	apiKey, hasAPIKey, headerErr := _jwt.APIKeyFromRequest(c.Request)
	if headerErr != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid API key header"},
		)
		return
	}
	if !hasAPIKey {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Missing API key"},
		)
		return
	}

	var verifyErr error
	if h == nil || h.verifier == nil {
		verifyErr = versionauth.ErrUnavailable
	} else {
		verifyErr = h.verifier.Verify(
			c.Request.Context(),
			apiKey,
		)
	}

	if verifyErr != nil {
		switch {
		case errors.Is(
			verifyErr,
			versionauth.ErrInvalid,
		):
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid API key"},
			)

		case errors.Is(
			verifyErr,
			versionauth.ErrUnavailable,
		):
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"error": "API key verifier unavailable",
				},
			)

		default:
			c.JSON(
				http.StatusServiceUnavailable,
				gin.H{
					"error": "API key verifier unavailable",
				},
			)
		}
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "File is required"},
		)
		return
	}

	description := c.DefaultPostForm(
		"description",
		"",
	)

	saved, err := h.uploader.Upload(
		c.Request.Context(),
		fileversion.UploadInput{
			OpenContent: func() (
				io.ReadCloser,
				error,
			) {
				content, err := file.Open()
				if err != nil {
					return nil, err
				}
				return content, nil
			},
			FileName:    file.Filename,
			ContentType: file.Header.Get("Content-Type"),
			Description: description,
			Size:        file.Size,
		},
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"data": externalFileInfo(saved),
		},
	)
}
