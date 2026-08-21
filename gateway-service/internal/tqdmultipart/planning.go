package tqdmultipart

import (
	"io"
	"mime/multipart"
	"net/http"

	"gateway/internal/httperror"
	"gateway/internal/tqdtransport"
	tqdpb "pb/types/tqd"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) createPlanningProjectFromFolder(c *gin.Context) {
	if !h.acquire(c) {
		return
	}
	defer h.release()
	cleanup, ok := h.parse(c)
	if !ok {
		return
	}
	defer cleanup()

	folderName := firstNonEmpty(c, "folder_name", "folderName")
	if folderName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder_name is required"})
		return
	}
	var headers []*multipart.FileHeader
	if c.Request.MultipartForm != nil {
		headers = c.Request.MultipartForm.File["files"]
	}
	if len(headers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "files is required"})
		return
	}

	remaining := h.limits.MaxBodyBytes
	files := make([]*tqdpb.FolderFileEntry, 0, len(headers))
	for _, fh := range headers {
		if remaining <= 0 {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "folder payload too large"})
			return
		}
		src, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "open file failed"})
			return
		}
		content, readErr := io.ReadAll(io.LimitReader(src, remaining+1))
		_ = src.Close()
		if readErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "read file failed"})
			return
		}
		if int64(len(content)) > remaining {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "folder payload too large"})
			return
		}
		remaining -= int64(len(content))
		files = append(files, &tqdpb.FolderFileEntry{RelativePath: fh.Filename, Content: content, ContentType: fh.Header.Get("Content-Type")})
	}
	req := &tqdpb.CreatePlanningProjectFromFolderRequest{
		FolderName:     folderName,
		Code:           firstNonEmpty(c, "code", "code"),
		Name:           firstNonEmpty(c, "name", "name"),
		ValidityStatus: firstNonEmpty(c, "validity_status", "validityStatus"),
		Authority:      firstNonEmpty(c, "authority", "authority"),
		Summary:        firstNonEmpty(c, "summary", "summary"),
		Files:          files,
	}
	if value, present, err := parseUint32(c, "planning_type", "planningType"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if present {
		req.PlanningType = value
	}
	if value, present, err := parseUint32(c, "planning_level", "planningLevel"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if present {
		req.PlanningLevel = value
	}
	if value, present, err := parseUint64(c, "jurisdiction_id", "jurisdictionId"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if present {
		req.JurisdictionId = value
	}
	if value, present, err := parseUint32(c, "job_type", "jobType"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if present {
		req.JobType = value
	} else {
		req.JobType = 10
	}
	if proto.Size(req) > tqdtransport.MaxGRPCMessageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "encoded TQD request too large"})
		return
	}
	resp, err := tqdpb.NewQHPlanningServiceClient(h.conn).CreatePlanningProjectFromFolder(c.Request.Context(), req)
	if err != nil {
		httperror.WriteGRPC(c.Writer, err)
		return
	}
	out, err := protojson.MarshalOptions{UseProtoNames: false, EmitUnpopulated: false}.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal response failed"})
		return
	}
	c.Data(http.StatusOK, "application/json", out)
}
