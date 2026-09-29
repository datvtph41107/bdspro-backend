package rendering

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"
)

type templateData struct {
	Title       string
	Subtitle    string
	ReportID    uint64
	CreatedAt   string
	Address     string
	Province    string
	ParcelID    string
	RegionID    string
	MapNumber   string
	LandNumber  string
	Center      string
	Bounds      string
	Comparison  string
	HasParcel   bool
	HasRegion   bool
	HasCenter   bool
	HasBounds   bool
	HasMetadata bool
	HasCompare  bool
}

var reportTemplate = template.Must(template.New("generated-report").Parse(`<!doctype html>
<html lang="vi">
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
<style>
@page { size: A4; margin: 18mm 16mm 18mm 16mm; }
body { font-family: "Arial", "Noto Sans", sans-serif; color: #172033; font-size: 12px; line-height: 1.55; }
h1 { font-size: 23px; margin: 0 0 4px; color: #0b3b66; }
h2 { font-size: 14px; margin: 18px 0 8px; color: #0b3b66; border-bottom: 1px solid #d9e2ec; padding-bottom: 4px; }
.sub { color: #52606d; margin-bottom: 14px; }
.meta { width: 100%; border-collapse: collapse; }
.meta td { border: 1px solid #d9e2ec; padding: 7px 9px; vertical-align: top; }
.meta td:first-child { width: 34%; font-weight: 600; background: #f6f8fa; }
.note { margin-top: 20px; color: #7b8794; font-size: 10px; }
pre { white-space: pre-wrap; word-break: break-word; font: inherit; margin: 0; }
</style>
</head>
<body>
<h1>{{.Title}}</h1>
{{if .Subtitle}}<div class="sub">{{.Subtitle}}</div>{{end}}
<table class="meta">
<tr><td>Mã báo cáo</td><td>{{.ReportID}}</td></tr>
<tr><td>Thời điểm ghi nhận</td><td>{{.CreatedAt}}</td></tr>
{{if .Address}}<tr><td>Địa chỉ</td><td>{{.Address}}</td></tr>{{end}}
{{if .Province}}<tr><td>Tỉnh/Thành phố</td><td>{{.Province}}</td></tr>{{end}}
{{if .HasParcel}}<tr><td>Thửa đất</td><td>{{.ParcelID}}</td></tr>{{end}}
{{if .HasRegion}}<tr><td>Vùng quy hoạch</td><td>{{.RegionID}}</td></tr>{{end}}
{{if .MapNumber}}<tr><td>Tờ bản đồ</td><td>{{.MapNumber}}</td></tr>{{end}}
{{if .LandNumber}}<tr><td>Số thửa</td><td>{{.LandNumber}}</td></tr>{{end}}
{{if .HasCenter}}<tr><td>Tâm tọa độ</td><td>{{.Center}}</td></tr>{{end}}
{{if .HasBounds}}<tr><td>Phạm vi tọa độ</td><td>{{.Bounds}}</td></tr>{{end}}
</table>
{{if .HasCompare}}<h2>So sánh quy hoạch</h2><pre>{{.Comparison}}</pre>{{end}}
<div class="note">Tài liệu được tạo từ snapshot đã được hệ thống QHPro chấp nhận. Trạng thái quota và retry không phải nội dung của tài liệu.</div>
</body>
</html>`))

// RenderHTML is deterministic for the same accepted Document.
func RenderHTML(document Document) ([]byte, error) {
	if document.ReportID == 0 || document.UserID == 0 || strings.TrimSpace(document.JobID) == "" {
		return nil, fmt.Errorf("render document identity is invalid")
	}
	if strings.TrimSpace(document.Title) == "" {
		return nil, fmt.Errorf("render document title is required")
	}

	data := templateData{
		Title:     strings.TrimSpace(document.Title),
		Subtitle:  strings.TrimSpace(document.Subtitle),
		ReportID:  document.ReportID,
		CreatedAt: stableTime(document.CreatedAt),
		Address:   strings.TrimSpace(document.Address),
		Province:  strings.TrimSpace(document.Province),
	}
	if document.ParcelID != nil && *document.ParcelID > 0 {
		data.HasParcel = true
		data.ParcelID = fmt.Sprintf("%d", *document.ParcelID)
	}
	if document.RegionID != nil && *document.RegionID > 0 {
		data.HasRegion = true
		data.RegionID = fmt.Sprintf("%d", *document.RegionID)
	}
	if document.CenterLat != 0 || document.CenterLon != 0 {
		data.HasCenter = true
		data.Center = fmt.Sprintf("%.6f, %.6f", document.CenterLat, document.CenterLon)
	}
	if document.MinLon != 0 || document.MinLat != 0 || document.MaxLon != 0 || document.MaxLat != 0 {
		data.HasBounds = true
		data.Bounds = fmt.Sprintf("%.6f, %.6f → %.6f, %.6f", document.MinLon, document.MinLat, document.MaxLon, document.MaxLat)
	}

	metadata := decodeObject(document.MetadataJSON)
	data.MapNumber = scalar(metadata["mapNumber"])
	data.LandNumber = scalar(metadata["landNumber"])
	data.HasMetadata = data.MapNumber != "" || data.LandNumber != ""

	comparison := decodeObject(document.ComparisonJSON)
	if len(comparison) > 0 {
		fromPlan := scalar(comparison["fromPlanName"])
		toPlan := scalar(comparison["toPlanName"])
		fromYear := scalar(comparison["fromYear"])
		toYear := scalar(comparison["toYear"])
		parts := make([]string, 0, 2)
		if fromPlan != "" || fromYear != "" {
			parts = append(parts, strings.TrimSpace(fromPlan+" "+fromYear))
		}
		if toPlan != "" || toYear != "" {
			parts = append(parts, strings.TrimSpace(toPlan+" "+toYear))
		}
		if len(parts) > 0 {
			data.HasCompare = true
			data.Comparison = strings.Join(parts, " → ")
		}
	}

	var output bytes.Buffer
	if err := reportTemplate.Execute(&output, data); err != nil {
		return nil, fmt.Errorf("render report HTML: %w", err)
	}
	return output.Bytes(), nil
}

func stableTime(value time.Time) string {
	if value.IsZero() {
		return "—"
	}
	return value.UTC().Format("02/01/2006 15:04:05 UTC")
}

func decodeObject(raw []byte) map[string]any {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil
	}
	return object
}

func scalar(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		if typed == float64(int64(typed)) {
			return fmt.Sprintf("%d", int64(typed))
		}
		return fmt.Sprintf("%g", typed)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}
