package filehttp

import (
	"errors"
	"net/http"

	_jwt "common/jwt"
	"file/dto"
	"file/internal/filecontent"

	"github.com/gin-gonic/gin"
)

// UploadFile owns the public HTTP upload contract while delegating the
// file lifecycle itself to the canonical transport-independent capability.
func (h *Handler) UploadFile(c *gin.Context) {
	if h == nil || h.files == nil {
		c.JSON(
			http.StatusServiceUnavailable,
			gin.H{"error": "file service is unavailable"},
		)
		return
	}

	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "File is required"},
		)
		return
	}

	reader, err := header.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "cannot open uploaded file"},
		)
		return
	}
	defer reader.Close()

	// The public upload contract supports public content by default and binds the current profile when secure=true.
	secure := c.PostForm("secure") == "true"
	description := c.DefaultQuery("description", "")

	var profileID uint64
	if secure {
		profileID = _jwt.GetProfileId(c)
	}

	saved, err := h.files.UploadFile(
		c.Request.Context(),
		filecontent.UploadInput{
			Content:     reader,
			Size:        header.Size,
			Filename:    header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			Description: description,
			Secure:      secure,
			ProfileID:   profileID,
		},
	)
	if err != nil {
		if errors.Is(err, filecontent.ErrNoActor) {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": "secure file requires current profile"},
			)
			return
		}

		// The established public upload contract classifies persistence/storage failures as 400.
		// Preserve that transport contract;
		// error taxonomy can be corrected independently later.
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "cannot upload file"},
		)
		return
	}

	// Preserve the established {"file": dto.FileInfo} response shape. Some fields
	// remain zero-valued exactly as they did when StorageService.Store returned
	// dto.FileInfo directly.
	fileInfo := dto.FileInfo{
		ID:            saved.ID,
		FileName:      saved.Filename,
		RelativePath:  saved.Path,
		AbsolutePath:  "",
		ThumbnailPath: saved.ThumbnailPath,
		Extension:     saved.Extension,
		Size:          saved.Size,
		Hash:          saved.Hash,
		ContentType:   saved.ContentType,
	}

	c.JSON(http.StatusOK, gin.H{"file": fileInfo})
}
