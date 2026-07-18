package service

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	_jwt "common/jwt"
	_middleware "common/middleware"

	tqdpb "pb/types/tqd"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

// tqdPlanningFolderMaxMemory — khớp TqdGRPCMaxMsgBytes, cho phép upload nhiều file trong 1 folder.
const tqdPlanningFolderMaxMemory = TqdGRPCMaxMsgBytes

// RegisterTqdPlanningFromFolderMultipartOnV2 đăng ký POST multipart trước mux (giống import GeoJSON,
// dùng chung group "/form" để tránh đụng route wildcard "/v2/*any") để grpc-gateway không cố decode
// multipart form thành JSON.
//
// Route thực tế: POST /form/tqd/admin/planning-projects/from-folder
// (proto khai báo /v2/tqd/admin/planning-projects/from-folder chỉ mang tính mô tả REST, giống ImportGeoJson).
//
// Form fields:
//   - folder_name (bắt buộc), code, name, planning_type, planning_level, validity_status, authority, summary, jurisdiction_id
//   - files: nhiều file (multipart field "files"); mỗi file dùng chính Filename làm relative_path trong folder
//     (FE nên set filename = đường dẫn tương đối trong folder, ví dụ "thuyet-minh/TM.pdf", khi append vào FormData)
func RegisterTqdPlanningFromFolderMultipartOnV2(v2 gin.IRoutes, env string, opts []grpc.DialOption) {
	conn, err := DialTQDGRPC(env, opts)
	if err != nil {
		log.Printf("[TQD Planning from-folder multipart] bỏ qua đăng ký route: %v", err)
		return
	}
	v2.POST("/tqd/admin/planning-projects/from-folder", createPlanningProjectFromFolderMultipartHandler(conn))
}

func createPlanningProjectFromFolderMultipartHandler(conn *grpc.ClientConn) gin.HandlerFunc {
	client := tqdpb.NewQHPlanningServiceClient(conn)
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(tqdPlanningFolderMaxMemory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "multipart: " + err.Error()})
			return
		}

		folderName := firstNonEmptyGin(c, "folder_name", "folderName")
		if folderName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "folder_name is required"})
			return
		}

		var fileHeaders []*multipart.FileHeader
		if c.Request.MultipartForm != nil {
			fileHeaders = c.Request.MultipartForm.File["files"]
		}
		if len(fileHeaders) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "files is required"})
			return
		}

		files := make([]*tqdpb.FolderFileEntry, 0, len(fileHeaders))
		for _, fh := range fileHeaders {
			src, err := fh.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "open file " + fh.Filename + ": " + err.Error()})
				return
			}
			content, err := io.ReadAll(src)
			src.Close()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "read file " + fh.Filename + ": " + err.Error()})
				return
			}
			contentType := fh.Header.Get("Content-Type")
			files = append(files, &tqdpb.FolderFileEntry{
				RelativePath: fh.Filename,
				Content:      content,
				ContentType:  contentType,
			})
		}

		req := &tqdpb.CreatePlanningProjectFromFolderRequest{
			FolderName:     folderName,
			Code:           firstNonEmptyGin(c, "code", "code"),
			Name:           firstNonEmptyGin(c, "name", "name"),
			ValidityStatus: firstNonEmptyGin(c, "validity_status", "validityStatus"),
			Authority:      firstNonEmptyGin(c, "authority", "authority"),
			Summary:        firstNonEmptyGin(c, "summary", "summary"),
			Files:          files,
		}

		if s := firstNonEmptyGin(c, "planning_type", "planningType"); s != "" {
			if v, err := strconv.ParseUint(s, 10, 32); err == nil {
				req.PlanningType = uint32(v)
			}
		}
		if s := firstNonEmptyGin(c, "planning_level", "planningLevel"); s != "" {
			if v, err := strconv.ParseUint(s, 10, 32); err == nil {
				req.PlanningLevel = uint32(v)
			}
		}
		if v, err := parseUint64FormGin(c, "jurisdiction_id", "jurisdictionId"); err == nil {
			req.JurisdictionId = v
		}

		md := _middleware.InjectGrpcMetadataContextMiddleware(c.Request.Context(), c.Request)
		outCtx := metadata.NewOutgoingContext(c.Request.Context(), md)

		resp, err := client.CreatePlanningProjectFromFolder(outCtx, req)
		if err != nil {
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
