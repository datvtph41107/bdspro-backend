package tqdmultipart

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"gateway/internal/httperror"
	"gateway/internal/tqdtransport"
	tqdpb "pb/types/tqd"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) importGeoJSON(c *gin.Context) {
	if !h.acquire(c) {
		return
	}
	defer h.release()
	cleanup, ok := h.parse(c)
	if !ok {
		return
	}
	defer cleanup()

	fh, err := c.FormFile("file")
	if err != nil {
		fh, err = c.FormFile("file_content")
	}
	if err != nil || fh == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field \"file\" or \"file_content\" is required"})
		return
	}
	src, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "open file failed"})
		return
	}
	defer src.Close()
	fileBytes, err := io.ReadAll(io.LimitReader(src, h.limits.MaxBodyBytes+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "read file failed"})
		return
	}
	if len(fileBytes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}
	if int64(len(fileBytes)) > h.limits.MaxBodyBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
		return
	}

	layerID, present, err := parseUint64(c, "layer_id", "layerId")
	if err != nil || !present || layerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "layer_id is required and must be valid"})
		return
	}
	labelMappings, err := parseLabelMappings(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label_mappings invalid"})
		return
	}
	req := &tqdpb.ImportGeoJsonRequest{
		FileContent:      fileBytes,
		FileFormat:       firstNonEmpty(c, "file_format", "fileFormat"),
		LayerId:          layerID,
		LabelField:       firstNonEmpty(c, "label_field", "labelField"),
		LabelMappings:    labelMappings,
		ValidateGeometry: parseBool(c, "validate_geometry", "validateGeometry"),
		AutoFixGeometry:  parseBool(c, "auto_fix_geometry", "autoFixGeometry"),
		SkipInvalid:      parseBool(c, "skip_invalid", "skipInvalid"),
	}
	if value, present, err := parseUint64(c, "default_label_id", "defaultLabelId"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if present {
		req.DefaultLabelId = proto.Uint64(value)
	}
	sourceName := firstNonEmpty(c, "source_file_name", "sourceFileName")
	if sourceName == "" {
		sourceName = strings.TrimSpace(fh.Filename)
	}
	if sourceName != "" {
		req.SourceFileName = proto.String(sourceName)
	}
	if proto.Size(req) > tqdtransport.MaxGRPCMessageBytes {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "encoded TQD request too large"})
		return
	}

	resp, err := tqdpb.NewImportServiceClient(h.conn).ImportGeoJson(c.Request.Context(), req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.FailedPrecondition:
				c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
				return
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
				return
			}
		}
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

func parseLabelMappings(c *gin.Context) (map[string]uint64, error) {
	raw := firstNonEmpty(c, "label_mappings", "labelMappings")
	if raw == "" {
		return nil, nil
	}
	var result map[string]uint64
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, err
	}
	return result, nil
}
