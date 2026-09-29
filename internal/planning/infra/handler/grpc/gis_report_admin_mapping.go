package handler_grpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tqdpb "pb/types/tqd"
	"tqd/internal/domain"
	"tqd/internal/usecase"

	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/datatypes"
)

func adminGisReportListResponse(
	reports []domain.Report,
	total int64,
	page, size int32,
	bucket string,
	displays map[uint64]usecase.LinkageDisplay,
	users map[uint64]usecase.UserDisplay,
) *tqdpb.AdminGisReportListResponse {
	items := make([]*tqdpb.AdminGisReportListItem, 0, len(reports))
	for _, r := range reports {
		items = append(items, toAdminGisReportListItem(r, displays[r.ID], users[r.ID]))
	}
	resp := &tqdpb.AdminGisReportListResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  size,
	}
	if bucket != "" {
		resp.Bucket = &bucket
	}
	return resp
}

func toAdminGisReportListItem(r domain.Report, display usecase.LinkageDisplay, user usecase.UserDisplay) *tqdpb.AdminGisReportListItem {
	item := &tqdpb.AdminGisReportListItem{
		Id:            r.ID,
		UserId:        r.UserID,
		Title:         r.Title,
		ProblemReport: uint32(r.ProblemReport),
		ReportType:    uint32(r.ReportType),
		Status:        uint32(r.Status),
		QaStatus:      uint32(r.QaStatus),
		Severity:      uint32(r.Severity),
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
	}
	if r.AssigneeID != nil {
		item.AssigneeId = r.AssigneeID
	}
	if r.SupportTicketID != nil {
		item.SupportTicketId = r.SupportTicketID
	}
	if r.TargetID != nil {
		item.TargetId = r.TargetID
	}
	if r.ClosedAt != nil {
		closed := r.ClosedAt.Format(time.RFC3339)
		item.ClosedAt = &closed
	}
	applyLinkageSummaryListItem(item, r.Linkage, r.TargetID, display)
	applyUserDisplayListItem(item, user)
	return item
}

func toAdminGisReportDetail(r domain.Report, display usecase.LinkageDisplay, user usecase.UserDisplay) *tqdpb.AdminGisReportDetail {
	detail := &tqdpb.AdminGisReportDetail{
		Id:            r.ID,
		UserId:        r.UserID,
		Title:         r.Title,
		Description:   r.Description,
		ProblemReport: uint32(r.ProblemReport),
		ReportType:    uint32(r.ReportType),
		Status:        uint32(r.Status),
		QaStatus:      uint32(r.QaStatus),
		Severity:      uint32(r.Severity),
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
		Images:        reportImagesProto(r.Images),
	}
	if st, err := linkageToStruct(r.Linkage); err == nil && st != nil {
		detail.Linkage = st
	}
	if r.AssigneeID != nil {
		detail.AssigneeId = r.AssigneeID
	}
	if r.TargetID != nil {
		detail.TargetId = r.TargetID
	}
	if r.SupportTicketID != nil {
		detail.SupportTicketId = r.SupportTicketID
	}
	if r.ClosedAt != nil {
		closed := r.ClosedAt.Format(time.RFC3339)
		detail.ClosedAt = &closed
	}
	applyLinkageSummaryDetail(detail, r.Linkage, r.TargetID, display)
	applyUserDisplayDetail(detail, user)
	return detail
}

func toAdminGisReportEvents(
	events []domain.ReportEvent,
	total int64,
	names map[uint64]string,
) *tqdpb.AdminGisReportEventListResponse {
	out := make([]*tqdpb.AdminGisReportEvent, 0, len(events))
	for _, e := range events {
		ev := &tqdpb.AdminGisReportEvent{
			Id:        e.ID,
			ReportId:  e.ReportID,
			ActorId:   e.ActorID,
			Action:    e.Action,
			Note:      e.Note,
			CreatedAt: e.CreatedAt.Format(time.RFC3339),
		}
		if e.BeforeJSON != "" {
			ev.BeforeJson = &e.BeforeJSON
		}
		if e.AfterJSON != "" {
			ev.AfterJson = &e.AfterJSON
		}
		if n := names[e.ActorID]; n != "" {
			ev.ActorName = &n
		}
		out = append(out, ev)
	}
	return &tqdpb.AdminGisReportEventListResponse{
		Data:  out,
		Total: total,
	}
}

func linkageToStruct(raw datatypes.JSON) (*structpb.Struct, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return structpb.NewStruct(m)
}

func structToMap(st *structpb.Struct) map[string]interface{} {
	if st == nil {
		return nil
	}
	return st.AsMap()
}

func reportImagesProto(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func applyLinkageSummaryListItem(item *tqdpb.AdminGisReportListItem, raw datatypes.JSON, targetID *uint64, display usecase.LinkageDisplay) {
	summary := linkageSummaryFields(raw, targetID, display)
	if v, ok := summary["layerName"].(string); ok && v != "" {
		item.LayerName = &v
	}
	if v, ok := summary["layerId"].(string); ok && v != "" {
		item.LayerId = &v
	}
	if v, ok := summary["version"].(string); ok && v != "" {
		item.Version = &v
	}
	if v, ok := summary["legalStatus"].(string); ok && v != "" {
		item.LegalStatus = &v
	}
	if v, ok := summary["parcelName"].(string); ok && v != "" {
		item.ParcelName = &v
	}
	if v, ok := summary["parcelOrRegion"].(string); ok && v != "" {
		item.ParcelOrRegion = &v
	}
}

func applyLinkageSummaryDetail(detail *tqdpb.AdminGisReportDetail, raw datatypes.JSON, targetID *uint64, display usecase.LinkageDisplay) {
	summary := linkageSummaryFields(raw, targetID, display)
	if v, ok := summary["layerName"].(string); ok && v != "" {
		detail.LayerName = &v
	}
	if v, ok := summary["layerId"].(string); ok && v != "" {
		detail.LayerId = &v
	}
	if v, ok := summary["version"].(string); ok && v != "" {
		detail.Version = &v
	}
	if v, ok := summary["legalStatus"].(string); ok && v != "" {
		detail.LegalStatus = &v
	}
	if v, ok := summary["parcelName"].(string); ok && v != "" {
		detail.ParcelName = &v
	}
	if v, ok := summary["parcelOrRegion"].(string); ok && v != "" {
		detail.ParcelOrRegion = &v
	}
}

func linkageSummaryFields(raw datatypes.JSON, targetID *uint64, display usecase.LinkageDisplay) map[string]interface{} {
	out := map[string]interface{}{}
	parcelOrRegion := ""
	if targetID != nil {
		parcelOrRegion = strconv.FormatUint(*targetID, 10)
	}
	var m map[string]interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &m)
	}
	if m != nil {
		if parcelOrRegion == "" {
			parcelOrRegion = firstLinkageString(m, "parcelOrRegion", "parcelId", "regionId", "targetId")
		}
		if v := firstLinkageString(m, "layerName", "layer"); v != "" {
			out["layerName"] = v
		}
		if v := firstLinkageString(m, "layerId"); v != "" {
			out["layerId"] = v
		}
		if v := firstLinkageString(m, "version", "layerVersion", "versionNo"); v != "" {
			out["version"] = v
		}
		if v := firstLinkageString(m, "legalStatusName", "legalStatus"); v != "" {
			out["legalStatus"] = v
		}
		if v := firstLinkageString(m, "parcelName"); v != "" {
			out["parcelName"] = v
		}
	}
	if parcelOrRegion != "" {
		out["parcelOrRegion"] = parcelOrRegion
	}
	if display.ParcelName != "" {
		out["parcelName"] = display.ParcelName
	}
	if display.LayerName != "" {
		out["layerName"] = display.LayerName
	}
	return out
}

func firstLinkageString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch t := v.(type) {
			case string:
				if strings.TrimSpace(t) != "" {
					return strings.TrimSpace(t)
				}
			case float64:
				return strconv.FormatInt(int64(t), 10)
			case json.Number:
				return t.String()
			default:
				s := strings.TrimSpace(fmt.Sprint(t))
				if s != "" && s != "<nil>" {
					return s
				}
			}
		}
	}
	return ""
}

func applyUserDisplayListItem(item *tqdpb.AdminGisReportListItem, user usecase.UserDisplay) {
	if user.UserName != "" {
		item.UserName = &user.UserName
	}
	if user.UserFullName != "" {
		item.UserFullName = &user.UserFullName
	}
	if user.UserUsername != "" {
		item.UserUsername = &user.UserUsername
	}
	if user.AssigneeName != "" {
		item.AssigneeName = &user.AssigneeName
	}
	if user.AssigneeFullName != "" {
		item.AssigneeFullName = &user.AssigneeFullName
	}
	if user.AssigneeUsername != "" {
		item.AssigneeUsername = &user.AssigneeUsername
	}
}

func applyUserDisplayDetail(detail *tqdpb.AdminGisReportDetail, user usecase.UserDisplay) {
	if user.UserName != "" {
		detail.UserName = &user.UserName
	}
	if user.UserFullName != "" {
		detail.UserFullName = &user.UserFullName
	}
	if user.UserUsername != "" {
		detail.UserUsername = &user.UserUsername
	}
	if user.AssigneeName != "" {
		detail.AssigneeName = &user.AssigneeName
	}
	if user.AssigneeFullName != "" {
		detail.AssigneeFullName = &user.AssigneeFullName
	}
	if user.AssigneeUsername != "" {
		detail.AssigneeUsername = &user.AssigneeUsername
	}
}
