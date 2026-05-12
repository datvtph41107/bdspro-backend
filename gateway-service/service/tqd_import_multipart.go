package service

import (
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	_jwt "common/jwt"
	_middleware "common/middleware"

	tqdpb "pb/types/tqd"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Khớp TqdGRPCMaxMsgBytes: form phải đọc được file trước khi đóng gói một message gRPC.
const tqdImportMultipartMaxMemory = TqdGRPCMaxMsgBytes

// RegisterTqdImportGeoJSONMultipartOnV2 đăng ký POST multipart trước mux để không bị grpc-gateway
// cố decode body thành JSON (lỗi "invalid character '-' in numeric literal").
// Form: file hoặc file_content (file), file_format, layer_id, label_field, label_mappings (JSON string), ...
func RegisterTqdImportGeoJSONMultipartOnV2(v2 gin.IRoutes, env string, opts []grpc.DialOption) {
	conn, err := DialTQDGRPC(env, opts)
	if err != nil {
		log.Printf("[TQD Import multipart] bỏ qua đăng ký route: %v", err)
		return
	}
	v2.POST("/tqd/admin/import/geojson", importGeoJSONMultipartHandler(conn))
}

func importGeoJSONMultipartHandler(conn *grpc.ClientConn) gin.HandlerFunc {
	client := tqdpb.NewImportServiceClient(conn)
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(tqdImportMultipartMaxMemory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "multipart: " + err.Error()})
			return
		}

		var fh *multipart.FileHeader
		if f, err := c.FormFile("file"); err == nil {
			fh = f
		} else if f, err2 := c.FormFile("file_content"); err2 == nil {
			fh = f
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "form field \"file\" or \"file_content\" is required"})
			return
		}

		src, err := fh.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		fileBytes, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "read file: " + err.Error()})
			return
		}
		if len(fileBytes) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
			return
		}

		layerID, err := parseUint64FormGin(c, "layer_id", "layerId")
		if err != nil || layerID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "layer_id is required"})
			return
		}

		fileFormat := firstNonEmptyGin(c, "file_format", "fileFormat")
		labelField := firstNonEmptyGin(c, "label_field", "labelField")
		labelMappings, err := parseLabelMappingsFormGin(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "label_mappings: " + err.Error()})
			return
		}

		req := &tqdpb.ImportGeoJsonRequest{
			FileContent:      fileBytes,
			FileFormat:       fileFormat,
			LayerId:          layerID,
			LabelField:       labelField,
			LabelMappings:    labelMappings,
			ValidateGeometry: parseBoolFormGin(c, "validate_geometry", "validateGeometry"),
			AutoFixGeometry:  parseBoolFormGin(c, "auto_fix_geometry", "autoFixGeometry"),
			SkipInvalid:      parseBoolFormGin(c, "skip_invalid", "skipInvalid"),
		}

		if s := firstNonEmptyGin(c, "default_label_id", "defaultLabelId"); s != "" {
			v, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "default_label_id invalid"})
				return
			}
			req.DefaultLabelId = proto.Uint64(v)
		}

		sourceName := firstNonEmptyGin(c, "source_file_name", "sourceFileName")
		if sourceName == "" {
			sourceName = fh.Filename
		}
		if sourceName != "" {
			req.SourceFileName = proto.String(sourceName)
		}

		md := _middleware.InjectGrpcMetadataContextMiddleware(c.Request.Context(), c.Request)
		outCtx := metadata.NewOutgoingContext(c.Request.Context(), md)

		resp, err := client.ImportGeoJson(outCtx, req)
		if err != nil {
			st, ok := status.FromError(err)
			if ok && st.Code() == codes.FailedPrecondition {
				c.JSON(http.StatusConflict, gin.H{"error": st.Message()})
				return
			}
			if ok && st.Code() == codes.NotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": st.Message()})
				return
			}
			_jwt.ErrorToGatewayResponse(c.Request.Context(), c.Writer, err)
			return
		}

		out, err := protojson.MarshalOptions{UseProtoNames: false, EmitUnpopulated: false}.Marshal(resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", out)
	}
}

func firstNonEmptyGin(c *gin.Context, snake, camel string) string {
	a := strings.TrimSpace(c.PostForm(snake))
	if a != "" {
		return a
	}
	return strings.TrimSpace(c.PostForm(camel))
}

func parseUint64FormGin(c *gin.Context, snake, camel string) (uint64, error) {
	s := firstNonEmptyGin(c, snake, camel)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseUint(s, 10, 64)
}

func parseBoolFormGin(c *gin.Context, snake, camel string) bool {
	s := strings.ToLower(firstNonEmptyGin(c, snake, camel))
	return s == "1" || s == "true" || s == "yes"
}

func parseLabelMappingsFormGin(c *gin.Context) (map[string]uint64, error) {
	raw := firstNonEmptyGin(c, "label_mappings", "labelMappings")
	if raw == "" {
		return nil, nil
	}
	var m map[string]uint64
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	return m, nil
}
