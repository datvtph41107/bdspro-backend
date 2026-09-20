package qh_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

// classifyPromptFilenameTemplate — phân loại chỉ dựa tên file / path / đuôi (fallback).
const classifyPromptFilenameTemplate = `Bạn là hệ thống phân loại tài liệu trong hồ sơ đồ án quy hoạch xây dựng/sử dụng đất tại Việt Nam.

Thông tin tài liệu:
- Tên file: %s
- Đường dẫn trong hồ sơ: %s
- Phần mở rộng: %s

Hãy phân loại tài liệu này vào MỘT trong các loại sau (chỉ chọn đúng 1 giá trị, giữ nguyên tiếng Việt):
- "Pháp lý": quyết định, văn bản pháp lý, giấy phép
- "Thuyết minh": báo cáo thuyết minh, thuyết minh tổng hợp
- "Bản đồ": bản đồ quy hoạch, sơ đồ quy hoạch (ảnh/pdf bản đồ)
- "CAD": file thiết kế CAD (.dwg, .dxf)
- "GIS gốc": dữ liệu GIS gốc (.shp, .geojson, .kml, .gpkg)
- "Phụ lục": phụ lục, bảng biểu kèm theo
- "Tệp khác": không xác định được loại nào phù hợp

Trả về CHỈ MỘT JSON thuần (không markdown, không code fence, không giải thích), đúng format:
{"documentType":"<một trong các loại trên>","title":"<tên tài liệu ngắn gọn, có nghĩa>","code":"","description":"<mô tả ngắn 1 câu>","confidence":<số 0.0-1.0>}`

// classifyPromptContentTemplate — phân loại dựa nội dung tài liệu đính kèm (PDF/ảnh/text).
const classifyPromptContentTemplate = `Bạn là hệ thống phân loại tài liệu trong hồ sơ đồ án quy hoạch xây dựng/sử dụng đất tại Việt Nam.

Thông tin tài liệu:
- Tên file: %s
- Đường dẫn trong hồ sơ: %s
- Phần mở rộng: %s

Hãy ĐỌC NỘI DUNG tài liệu đính kèm và phân loại vào MỘT trong các loại sau (chỉ chọn đúng 1 giá trị, giữ nguyên tiếng Việt):
- "Pháp lý": quyết định, văn bản pháp lý, giấy phép
- "Thuyết minh": báo cáo thuyết minh, thuyết minh tổng hợp
- "Bản đồ": bản đồ quy hoạch, sơ đồ quy hoạch (ảnh/pdf bản đồ)
- "CAD": file thiết kế CAD (.dwg, .dxf)
- "GIS gốc": dữ liệu GIS gốc (.shp, .geojson, .kml, .gpkg)
- "Phụ lục": phụ lục, bảng biểu kèm theo
- "Tệp khác": không xác định được loại nào phù hợp

Trả về CHỈ MỘT JSON thuần (không markdown, không code fence, không giải thích), đúng format:
{"documentType":"<một trong các loại trên>","title":"<tên tài liệu ngắn gọn, có nghĩa>","code":"","description":"<mô tả ngắn 1 câu>","confidence":<số 0.0-1.0>}`

const (
	classifyModeContent   = "content"
	classifyModeFilename  = "filename"
	classifyModeExtension = "extension"

	defaultMaxInlineBytes = 15 << 20 // 15MB
)

// ClassifyPolicy is process-owned configuration injected into the business usecase.
type ClassifyPolicy struct {
	ReadContent    bool
	MaxInlineBytes int
}

// classifyResult — kết quả JSON Gemini trả về, parse từ content thô.
type classifyResult struct {
	DocumentType string  `json:"documentType"`
	Title        string  `json:"title"`
	Code         string  `json:"code"`
	Description  string  `json:"description"`
	Confidence   float64 `json:"confidence"`
}

// IQHPlanningClassifyUsecase quét các QHPlanningDocument đang Pending, gọi Gemini (qua assistant-service)
// để phân loại DocumentType, cập nhật trạng thái tài liệu/đồ án. Được gọi lặp lại bởi 1 worker định kỳ.
type IQHPlanningClassifyUsecase interface {
	RunOnce(ctx context.Context, batchSize int) (int, error)
	Close() error
}

type qhPlanningClassifyUsecase struct {
	documentRepo repo.IQHPlanningDocumentRepo
	projectRepo  repo.IQHPlanningProjectRepo
	aiJobRepo    repo.IProAIJobRepo
	assistant    provider.AssistantProvider
	fileProvider provider.FileProvider
	policy       ClassifyPolicy
}

// NewQHPlanningClassifyUsecase @bind: internal/usecase/qh.IQHPlanningClassifyUsecase
func NewQHPlanningClassifyUsecase(
	documentRepo repo.IQHPlanningDocumentRepo,
	projectRepo repo.IQHPlanningProjectRepo,
	aiJobRepo repo.IProAIJobRepo,
	assistant provider.AssistantProvider,
	fileProvider provider.FileProvider,
	policy ClassifyPolicy,
) IQHPlanningClassifyUsecase {
	return &qhPlanningClassifyUsecase{
		documentRepo: documentRepo,
		projectRepo:  projectRepo,
		aiJobRepo:    aiJobRepo,
		assistant:    assistant,
		fileProvider: fileProvider,
		policy:       normalizeClassifyPolicy(policy),
	}
}

func normalizeClassifyPolicy(policy ClassifyPolicy) ClassifyPolicy {
	if policy.MaxInlineBytes <= 0 {
		policy.MaxInlineBytes = defaultMaxInlineBytes
	}
	return policy
}

func (u *qhPlanningClassifyUsecase) Close() error {
	if u == nil || u.assistant == nil {
		return nil
	}
	return u.assistant.Close()
}

func (u *qhPlanningClassifyUsecase) RunOnce(ctx context.Context, batchSize int) (int, error) {
	if batchSize <= 0 {
		batchSize = 5
	}

	docs, err := u.documentRepo.ListByProcessStatus(ctx, enums.PlanningProcessStatusPending, batchSize)
	if err != nil {
		return 0, err
	}
	// Reclaim tài liệu kẹt Processing (crash/restart worker) để không treo mãi.
	if len(docs) < batchSize {
		stale, errStale := u.documentRepo.ListByProcessStatus(ctx, enums.PlanningProcessStatusProcessing, batchSize-len(docs))
		if errStale == nil && len(stale) > 0 {
			docs = append(docs, stale...)
		}
	}

	processed := 0
	touchedProjects := make(map[uint64]struct{})

	for i := range docs {
		doc := docs[i]
		doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusProcessing)
		_ = u.documentRepo.Update(ctx, doc.ID, &doc)

		u.classifyOne(ctx, &doc)

		if err := u.documentRepo.Update(ctx, doc.ID, &doc); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] update document %d lỗi: %v", doc.ID, err))
			continue
		}
		processed++
		touchedProjects[doc.PlanningProjectID] = struct{}{}
	}

	for projectID := range touchedProjects {
		u.refreshProjectStatus(ctx, projectID)
	}

	return processed, nil
}

func (u *qhPlanningClassifyUsecase) classifyOne(ctx context.Context, doc *qh_domain.QHPlanningDocument) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(doc.RelativePath), "."))
	baseName := filepath.Base(doc.RelativePath)

	// CAD / GIS / binary: map thẳng theo đuôi, không gọi AI.
	if typeVal, ok := classifyByExtension(ext); ok {
		doc.DocumentType = enums.PlanningDocumentType(typeVal)
		if strings.TrimSpace(doc.Title) == "" {
			doc.Title = baseName
		}
		metaBytes, _ := json.Marshal(map[string]interface{}{
			"suggestedType": enums.GetPlanningDocumentTypeLabel(typeVal),
			"confidence":    1.0,
			"classifyMode":  classifyModeExtension,
			"extension":     ext,
		})
		doc.Metadata = string(metaBytes)
		doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified)
		doc.ClassifyError = ""
		return
	}

	readContent := u.policy.ReadContent

	mimeType, readable := mimeTypeForExt(ext)
	var (
		content       string
		model         string
		err           error
		classifyMode  = classifyModeFilename
		fileBytes     []byte
		usedContentAI bool
	)

	if readContent && readable && doc.Filepath != "" {
		fileBytes, err = u.fileProvider.GetFile(ctx, doc.Filepath)
		maxInline := u.policy.MaxInlineBytes
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] GetFile document %d lỗi, fallback filename: %v", doc.ID, err))
		} else if len(fileBytes) > maxInline {
			slog.WarnContext(ctx, fmt.Sprintf("[QHPlanningClassify] document %d vượt max_inline_bytes (%d > %d), fallback filename", doc.ID, len(fileBytes), maxInline))
			fileBytes = nil
		} else if len(fileBytes) > 0 {
			prompt := fmt.Sprintf(classifyPromptContentTemplate, baseName, doc.RelativePath, ext)
			content, model, err = u.assistant.ClassifyDocument(ctx, prompt, fileBytes, mimeType)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] ClassifyDocument lỗi document %d, fallback filename: %v", doc.ID, err))
				fileBytes = nil
			} else {
				usedContentAI = true
				classifyMode = classifyModeContent
			}
		}
	}

	if !usedContentAI {
		prompt := fmt.Sprintf(classifyPromptFilenameTemplate, baseName, doc.RelativePath, ext)
		content, model, err = u.assistant.GenerateContent(ctx, prompt, 512, 0.2)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] gọi Gemini lỗi cho document %d: %v", doc.ID, err))
			doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed)
			doc.ClassifyError = "gemini: " + err.Error()
			return
		}
		classifyMode = classifyModeFilename
	}

	result, err := parseClassifyResult(content)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] parse JSON Gemini lỗi cho document %d: %v (raw=%s)", doc.ID, err, content))
		doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed)
		doc.ClassifyError = "parse kết quả gemini: " + err.Error()
		return
	}

	typeVal := enums.GetPlanningDocumentTypeValue(result.DocumentType)
	if typeVal == 0 {
		typeVal = enums.PlanningDocumentTypeOther
	}
	doc.DocumentType = enums.PlanningDocumentType(typeVal)
	if strings.TrimSpace(result.Title) != "" {
		doc.Title = result.Title
	}
	if strings.TrimSpace(result.Code) != "" {
		doc.Code = result.Code
	}
	if strings.TrimSpace(result.Description) != "" {
		doc.Description = result.Description
	}

	metaBytes, _ := json.Marshal(map[string]interface{}{
		"suggestedType": result.DocumentType,
		"confidence":    result.Confidence,
		"model":         model,
		"classifyMode":  classifyMode,
		"raw":           content,
	})
	doc.Metadata = string(metaBytes)
	doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusClassified)
	doc.ClassifyError = ""
}

// classifyByExtension map CAD/GIS theo đuôi; ok=false nếu không thuộc nhóm này.
func classifyByExtension(ext string) (uint32, bool) {
	switch ext {
	case "dwg", "dxf":
		return enums.PlanningDocumentTypeCAD, true
	case "shp", "geojson", "kml", "kmz", "gpkg":
		return enums.PlanningDocumentTypeGIS, true
	default:
		return 0, false
	}
}

// mimeTypeForExt trả mimeType Gemini đọc được + readable flag.
// .doc/.docx/.xls/.xlsx: Gemini không nhận trực tiếp → không readable.
func mimeTypeForExt(ext string) (string, bool) {
	switch ext {
	case "pdf":
		return "application/pdf", true
	case "png":
		return "image/png", true
	case "jpg", "jpeg":
		return "image/jpeg", true
	case "webp":
		return "image/webp", true
	case "gif":
		return "image/gif", true
	case "txt", "md", "csv":
		return "text/plain", true
	default:
		return "", false
	}
}

// refreshProjectStatus cập nhật trạng thái đồ án khi không còn tài liệu Pending/Processing.
func (u *qhPlanningClassifyUsecase) refreshProjectStatus(ctx context.Context, projectID uint64) {
	remaining, err := u.documentRepo.CountByProjectAndStatuses(ctx, projectID, []uint32{
		enums.PlanningProcessStatusPending, enums.PlanningProcessStatusProcessing,
	})
	if err != nil || remaining > 0 {
		return
	}

	classifiedCount, err := u.documentRepo.CountByProjectAndStatuses(ctx, projectID, []uint32{
		enums.PlanningProcessStatusClassified, enums.PlanningProcessStatusApproved,
	})
	if err != nil {
		return
	}

	newStatus := enums.PlanningProcessStatusFailed
	if classifiedCount > 0 {
		newStatus = enums.PlanningProcessStatusClassified
	}
	if err := u.projectRepo.UpdateProcessStatus(ctx, projectID, newStatus); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] cập nhật trạng thái đồ án %d lỗi: %v", projectID, err))
		return
	}
	if err := u.aiJobRepo.UpdateProcessStatusByProjectID(ctx, projectID, newStatus); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningClassify] sync job AI đồ án %d lỗi: %v", projectID, err))
	}
}

// parseClassifyResult trích JSON từ content trả về (đề phòng Gemini bọc thêm code fence ```json ... ```).
func parseClassifyResult(content string) (*classifyResult, error) {
	trimmed := strings.TrimSpace(content)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || end < start {
		return nil, fmt.Errorf("không tìm thấy JSON hợp lệ trong nội dung trả về")
	}
	trimmed = trimmed[start : end+1]

	var result classifyResult
	if err := json.Unmarshal([]byte(trimmed), &result); err != nil {
		return nil, err
	}
	return &result, nil
}
