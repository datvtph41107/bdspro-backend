package filehttp

import (
	"context"
	"net/http"

	"file/dto"
	"file/internal/filecontent"

	"github.com/gin-gonic/gin"
)

type videoUploader interface {
	UploadVideo(
		ctx context.Context,
		input filecontent.MediaUploadInput,
	) (*dto.FileInfo, error)
}

type VideoUploadHandler struct {
	files videoUploader
}

func NewVideoUploadHandler(
	files videoUploader,
) *VideoUploadHandler {
	return &VideoUploadHandler{
		files: files,
	}
}

// UploadVideo translates the existing multipart HTTP contract into the
// transport-independent media upload capability.
//
// @Summary Upload video HLS
// @Description Upload video HLS
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /upload/video [post]
func (h *VideoUploadHandler) UploadVideo(
	c *gin.Context,
) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "No file uploaded"},
		)
		return
	}

	content, err := file.Open()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "failed to open uploaded file"},
		)
		return
	}
	defer content.Close()

	saved, err := h.files.UploadVideo(
		c,
		filecontent.MediaUploadInput{
			Content:     content,
			Filename:    file.Filename,
			ContentType: file.Header.Get("Content-Type"),
			Description: c.DefaultQuery("description", ""),
			Size:        file.Size,
		},
	)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"code":    0,
			"message": "success",
			"data":    externalFileInfo(saved),
		},
	)
}
