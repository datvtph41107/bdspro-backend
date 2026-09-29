package usecase

import (
	"context"

	_errors "common/errors"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"tqd/internal"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"

	_utils "common/utils"
	"pb/clients"

	"gorm.io/datatypes"
)

type ReportUsecase interface {
	CreateAsync(ctx context.Context, userID uint64, reqDTO dto.ReportDTO) (uint64, error)
	GetStatus(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	Download(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error)
	ListUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
	Delete(ctx context.Context, userID uint64, reportID uint64) error
	AdminList(ctx context.Context, filter repo.AdminReportFilter) ([]domain.Report, int64, error)
	AdminSummary(ctx context.Context) (map[uint32]int64, error)
	AdminQueue(ctx context.Context, bucket string, severity *uint32, page, limit int, actorID uint64) ([]domain.Report, int64, error)
	AdminQueueSummary(ctx context.Context, actorID uint64) (repo.QueueSummary, error)
	AdminGet(ctx context.Context, reportID uint64) (*domain.Report, error)
	AdminCreateInternal(ctx context.Context, actorID uint64, req dto.AdminCreateReportRequest) (*domain.Report, error)
	AdminAssign(ctx context.Context, reportID, assigneeID uint64) error
	AdminUpdateSeverity(ctx context.Context, reportID uint64, severity uint32, reason string) error
	AdminUpdateQaStatus(ctx context.Context, reportID uint64, qaStatus uint32, note string) error
	AdminUpdateLinkage(ctx context.Context, reportID uint64, linkage map[string]interface{}) error
	AdminUpdateImages(ctx context.Context, reportID uint64, images []string) error
	AdminClose(ctx context.Context, reportID uint64, reject bool, note string) error
	AdminDelete(ctx context.Context, reportID uint64) error
	AdminListEvents(ctx context.Context, filter repo.AdminEventFilter) ([]domain.ReportEvent, int64, error)
	AdminListReportEvents(ctx context.Context, reportID uint64, limit int) ([]domain.ReportEvent, error)
	ResolveLinkageDisplays(ctx context.Context, reports []domain.Report) map[uint64]LinkageDisplay
	ResolveUserDisplays(ctx context.Context, reports []domain.Report) map[uint64]UserDisplay
	ResolveActorNames(ctx context.Context, actorIDs []uint64) map[uint64]string
}

type reportUsecase struct {
	repo         repo.ReportRepository
	userProvider provider.UserProvider
	authClient   *clients.AuthGrpcClient
	userClient   *clients.UserGrpcClient
}

func NewReportUsecase(
	repo repo.ReportRepository,
	userProvider provider.UserProvider,
	authClient *clients.AuthGrpcClient,
	userClient *clients.UserGrpcClient,
) ReportUsecase {
	return &reportUsecase{
		repo:         repo,
		userProvider: userProvider,
		authClient:   authClient,
		userClient:   userClient,
	}
}

func (u *reportUsecase) CreateAsync(ctx context.Context, userID uint64, reqDTO dto.ReportDTO) (uint64, error) {
	var targetSnapshot datatypes.JSON
	report := &domain.Report{
		UserID:         userID,
		ProblemReport:  reqDTO.ProblemReport,
		ReportType:     reqDTO.ReportType,
		Description:    reqDTO.Description,
		TargetID:       reqDTO.TargetId,
		TargetSnapshot: targetSnapshot,
		Status:         10,
		QaStatus:       enums.ReportQaStatusSubmitted,
		Severity:       enums.ReportSeverityMedium,
	}
	if err := u.repo.Create(ctx, report); err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (u *reportUsecase) processReport(reportID, userID uint64, reportType uint32, targetID *string, targetData interface{}, profile, format uint32, includeMapImage bool) {
	_ = u.repo.UpdateStatus(context.Background(), reportID, 20, nil)
	fileURL := "https://storage.example.com/reports/" + string(rune(reportID)) + ".pdf"
	fileSize := int64(1024)
	fileHash := "dummyhash"
	_ = u.repo.UpdateFileInfo(context.Background(), reportID, &fileURL, &fileSize, &fileHash, nil)
}

func (u *reportUsecase) GetStatus(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, _errors.ReturnError(service.ReportNotFound)
	}
	if report.UserID != userID {
		return nil, _errors.ReturnError(service.ReportPermissionDenied)
	}
	return report, nil
}

func (u *reportUsecase) Download(ctx context.Context, userID uint64, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, _errors.ReturnError(service.ReportNotFound)
	}
	if report.UserID != userID {
		return nil, _errors.ReturnError(service.ReportPermissionDenied)
	}
	if report.Status != 30 {
		return nil, _errors.ReturnError(service.ReportNotReady)
	}
	return report, nil
}

func (u *reportUsecase) ListUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error) {
	return u.repo.ListByUser(ctx, userID, reportType, status, page, limit)
}

func (u *reportUsecase) Delete(ctx context.Context, userID uint64, reportID uint64) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	if report.UserID != userID {
		return _errors.ReturnError(service.ReportPermissionDenied)
	}
	return u.repo.Delete(ctx, reportID)
}

func (u *reportUsecase) AdminList(ctx context.Context, filter repo.AdminReportFilter) ([]domain.Report, int64, error) {
	return u.repo.AdminList(ctx, filter)
}

func (u *reportUsecase) AdminSummary(ctx context.Context) (map[uint32]int64, error) {
	return u.repo.AdminSummary(ctx)
}

func (u *reportUsecase) AdminQueue(ctx context.Context, bucket string, severity *uint32, page, limit int, actorID uint64) ([]domain.Report, int64, error) {
	filter := repo.AdminReportFilter{
		Severity: severity,
		Page:     page,
		Limit:    limit,
	}
	switch bucket {
	case "unassigned":
		filter.UnassignedOnly = true
		filter.QaStatuses = []uint32{uint32(enums.ReportQaStatusSubmitted)}
	case "data_fix":
		filter.QaStatuses = []uint32{uint32(enums.ReportQaStatusWaitingDataFix)}
	default: // mine
		filter.AssigneeID = &actorID
		filter.QaStatuses = []uint32{
			uint32(enums.ReportQaStatusSubmitted),
			uint32(enums.ReportQaStatusVerifying),
			uint32(enums.ReportQaStatusNeedMoreInfo),
			uint32(enums.ReportQaStatusVerified),
		}
	}
	return u.repo.AdminList(ctx, filter)
}

func (u *reportUsecase) AdminQueueSummary(ctx context.Context, actorID uint64) (repo.QueueSummary, error) {
	return u.repo.AdminQueueSummary(ctx, actorID)
}

func (u *reportUsecase) AdminGet(ctx context.Context, reportID uint64) (*domain.Report, error) {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return nil, _errors.ReturnError(service.ReportNotFound)
	}
	return report, nil
}

func (u *reportUsecase) AdminCreateInternal(ctx context.Context, actorID uint64, req dto.AdminCreateReportRequest) (*domain.Report, error) {
	imagesJSON := datatypes.JSON([]byte("[]"))
	if len(req.Images) > 0 {
		b, err := json.Marshal(req.Images)
		if err != nil {
			return nil, err
		}
		imagesJSON = datatypes.JSON(b)
	}
	report := &domain.Report{
		UserID:          actorID,
		Title:           req.Title,
		Description:     req.Description,
		ProblemReport:   enums.ProblemReport(req.ProblemReport),
		ReportType:      enums.ReportType(req.ReportType),
		Severity:        enums.ReportSeverity(req.Severity),
		QaStatus:        enums.ReportQaStatusSubmitted,
		Status:          10,
		Profile:         10,
		Format:          10,
		SupportTicketID: req.SupportTicketID,
		TargetID:        req.TargetID,
		Images:          imagesJSON,
	}
	if req.Severity == 0 {
		report.Severity = enums.ReportSeverityMedium
	}
	if req.ProblemReport == 0 {
		report.ProblemReport = enums.ProblemReport_Other
	}
	if req.ReportType == 0 {
		report.ReportType = enums.ReportTypeParcel
	}
	if err := u.repo.Create(ctx, report); err != nil {
		return nil, err
	}
	_ = u.addEvent(ctx, report.ID, actorID, "create", "", compactJSON(map[string]interface{}{
		"title": report.Title, "qaStatus": uint32(report.QaStatus), "severity": uint32(report.Severity),
	}), "Tạo phản ánh nội bộ")
	return report, nil
}

func (u *reportUsecase) AdminAssign(ctx context.Context, reportID, assigneeID uint64) error {
	if assigneeID == 0 {
		return _errors.ReturnError(service.ReportAssigneeIDRequired)
	}
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	before := compactJSON(map[string]interface{}{
		"assigneeId": report.AssigneeID, "qaStatus": uint32(report.QaStatus),
	})
	qaStatus := report.QaStatus
	if qaStatus == enums.ReportQaStatusSubmitted {
		qaStatus = enums.ReportQaStatusVerifying
	}
	if err := u.repo.UpdateAssignee(ctx, reportID, assigneeID, uint32(qaStatus)); err != nil {
		return err
	}
	actorID := actorFromCtx(ctx)
	_ = u.addEvent(ctx, reportID, actorID, "assign", before, compactJSON(map[string]interface{}{
		"assigneeId": assigneeID, "qaStatus": uint32(qaStatus),
	}), "Gán xác minh GIS QA")
	return nil
}

func (u *reportUsecase) AdminUpdateSeverity(ctx context.Context, reportID uint64, severity uint32, reason string) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	before := compactJSON(map[string]interface{}{"severity": uint32(report.Severity)})
	report.Severity = enums.ReportSeverity(severity)
	if err := u.repo.Update(ctx, report); err != nil {
		return err
	}
	note := strings.TrimSpace(reason)
	if note == "" {
		note = "Đổi mức độ"
	}
	_ = u.addEvent(ctx, reportID, actorFromCtx(ctx), "severity", before, compactJSON(map[string]interface{}{
		"severity": severity,
	}), note)
	return nil
}

func (u *reportUsecase) AdminUpdateQaStatus(ctx context.Context, reportID uint64, qaStatus uint32, note string) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	to := enums.ReportQaStatus(qaStatus)
	if !enums.CanTransitionReportQa(report.QaStatus, to) {
		return _errors.ReturnError(service.ReportStatusTransitionInvalid)
	}
	before := compactJSON(map[string]interface{}{"qaStatus": uint32(report.QaStatus)})
	report.QaStatus = to
	if to == enums.ReportQaStatusClosed || to == enums.ReportQaStatusRejected {
		now := time.Now()
		report.ClosedAt = &now
	}
	if err := u.repo.Update(ctx, report); err != nil {
		return err
	}
	eventNote := strings.TrimSpace(note)
	if eventNote == "" {
		eventNote = "Đổi trạng thái QA"
	}
	_ = u.addEvent(ctx, reportID, actorFromCtx(ctx), "qa_status", before, compactJSON(map[string]interface{}{
		"qaStatus": uint32(to),
	}), eventNote)
	return nil
}

func (u *reportUsecase) AdminUpdateLinkage(ctx context.Context, reportID uint64, linkage map[string]interface{}) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	if linkage == nil {
		linkage = map[string]interface{}{}
	}
	before := string(report.Linkage)

	// Đồng bộ target_id nếu client gửi targetId / parcelOrRegion (parcel id dạng số).
	rawTarget := firstLinkageString(linkage, "targetId", "parcelOrRegion", "parcelId")
	var targetID *uint64
	if rawTarget != "" {
		if n, err := strconv.ParseUint(rawTarget, 10, 64); err == nil {
			targetID = &n
			report.TargetID = &n
			linkage["targetId"] = strconv.FormatUint(n, 10)
			linkage["parcelOrRegion"] = strconv.FormatUint(n, 10)
		}
	}

	// Gắn nhãn hiển thị từ DB (FE/list không cần gọi thêm API).
	if targetID != nil {
		if labels, err := u.repo.LookupParcelLabels(ctx, []uint64{*targetID}); err == nil {
			if name := labels[*targetID]; name != "" {
				linkage["parcelName"] = name
			}
		}
	}
	if lid := firstLinkageString(linkage, "layerId"); lid != "" {
		if n, err := strconv.ParseUint(lid, 10, 64); err == nil {
			if labels, err := u.repo.LookupLayerLabels(ctx, []uint64{n}); err == nil {
				if name := labels[n]; name != "" {
					linkage["layerName"] = name
				}
			}
		}
	}

	b, err := json.Marshal(linkage)
	if err != nil {
		return err
	}
	report.Linkage = datatypes.JSON(b)
	if err := u.repo.Update(ctx, report); err != nil {
		return err
	}
	_ = u.addEvent(ctx, reportID, actorFromCtx(ctx), "linkage", before, string(b), "Cập nhật liên kết thửa/lớp")
	return nil
}

func firstLinkageString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		v, ok := m[k]
		if !ok || v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			if s := strings.TrimSpace(t); s != "" {
				return s
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
	return ""
}

// LinkageDisplay nhãn hiển thị thửa/lớp cho list + detail (không N+1 API phía FE).
type LinkageDisplay struct {
	ParcelName string
	LayerName  string
}

// ResolveLinkageDisplays batch-resolve parcel/layer names cho một trang reports.
func (u *reportUsecase) ResolveLinkageDisplays(ctx context.Context, reports []domain.Report) map[uint64]LinkageDisplay {
	out := make(map[uint64]LinkageDisplay, len(reports))
	parcelIDs := make([]uint64, 0)
	layerIDs := make([]uint64, 0)
	seenP := map[uint64]struct{}{}
	seenL := map[uint64]struct{}{}

	type ids struct {
		parcelID     uint64
		layerID      uint64
		hasP         bool
		hasL         bool
		cachedParcel string
		cachedLayer  string
	}
	perReport := make(map[uint64]ids, len(reports))

	for _, r := range reports {
		info := ids{}
		if r.TargetID != nil && *r.TargetID > 0 {
			info.parcelID = *r.TargetID
			info.hasP = true
		}
		var m map[string]interface{}
		if len(r.Linkage) > 0 {
			_ = json.Unmarshal(r.Linkage, &m)
		}
		if m != nil {
			info.cachedParcel = firstLinkageString(m, "parcelName")
			info.cachedLayer = firstLinkageString(m, "layerName", "layer")
			if !info.hasP {
				if raw := firstLinkageString(m, "targetId", "parcelOrRegion", "parcelId"); raw != "" {
					if n, err := strconv.ParseUint(raw, 10, 64); err == nil {
						info.parcelID = n
						info.hasP = true
					}
				}
			}
			if raw := firstLinkageString(m, "layerId"); raw != "" {
				if n, err := strconv.ParseUint(raw, 10, 64); err == nil {
					info.layerID = n
					info.hasL = true
				}
			}
		}
		perReport[r.ID] = info
		if info.hasP && info.cachedParcel == "" {
			if _, ok := seenP[info.parcelID]; !ok {
				seenP[info.parcelID] = struct{}{}
				parcelIDs = append(parcelIDs, info.parcelID)
			}
		}
		if info.hasL && info.cachedLayer == "" {
			if _, ok := seenL[info.layerID]; !ok {
				seenL[info.layerID] = struct{}{}
				layerIDs = append(layerIDs, info.layerID)
			}
		}
	}

	parcelLabels, _ := u.repo.LookupParcelLabels(ctx, parcelIDs)
	layerLabels, _ := u.repo.LookupLayerLabels(ctx, layerIDs)

	for _, r := range reports {
		info := perReport[r.ID]
		d := LinkageDisplay{
			ParcelName: info.cachedParcel,
			LayerName:  info.cachedLayer,
		}
		if d.ParcelName == "" && info.hasP {
			d.ParcelName = parcelLabels[info.parcelID]
		}
		if d.LayerName == "" && info.hasL {
			d.LayerName = layerLabels[info.layerID]
		}
		out[r.ID] = d
	}
	return out
}

// UserDisplay nhãn người tạo / người được giao (batch proto, FE không gọi thêm API).
type UserDisplay struct {
	UserName         string
	UserFullName     string
	UserUsername     string
	AssigneeName     string
	AssigneeFullName string
	AssigneeUsername string
}

// ResolveUserDisplays batch-resolve tên người tạo + người được giao qua user/auth proto.
func (u *reportUsecase) ResolveUserDisplays(ctx context.Context, reports []domain.Report) map[uint64]UserDisplay {
	out := make(map[uint64]UserDisplay, len(reports))
	if len(reports) == 0 {
		return out
	}

	idSet := map[uint64]struct{}{}
	for _, r := range reports {
		if r.UserID > 0 {
			idSet[r.UserID] = struct{}{}
		}
		if r.AssigneeID != nil && *r.AssigneeID > 0 {
			idSet[*r.AssigneeID] = struct{}{}
		}
	}
	ids := make([]uint64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	type person struct {
		fullName string
		username string
	}
	people := make(map[uint64]person, len(ids))

	if u.userProvider != nil && len(ids) > 0 {
		if profiles, err := u.userProvider.GetProfilesByIDs(ctx, ids); err == nil {
			for id, p := range profiles {
				if p == nil {
					continue
				}
				cur := people[id]
				if fn := strings.TrimSpace(p.FullName); fn != "" {
					cur.fullName = fn
				}
				people[id] = cur
			}
		}
	}

	if u.userClient != nil && u.userClient.InternalClient != nil && len(ids) > 0 {
		if resp, err := u.userClient.GetAuthAdminByIds(ctx, ids); err == nil && resp != nil {
			for _, a := range resp.GetData() {
				if a == nil {
					continue
				}
				// GetAuthAdminByIds query theo auth_method.user_id (= profile/admin id).
				key := a.GetUserId()
				if key == 0 {
					key = a.GetId()
				}
				if key == 0 {
					continue
				}
				cur := people[key]
				if fn := strings.TrimSpace(a.GetFullName()); fn != "" {
					cur.fullName = fn
				}
				if un := strings.TrimSpace(a.GetUsername()); un != "" {
					cur.username = un
				}
				people[key] = cur
				// Cũng index theo auth_method.id nếu khác user_id (match createdBy kiểu authId).
				if aid := a.GetId(); aid > 0 && aid != key {
					cur2 := people[aid]
					if cur2.fullName == "" {
						cur2.fullName = cur.fullName
					}
					if cur2.username == "" {
						cur2.username = cur.username
					}
					people[aid] = cur2
				}
			}
		}
	}

	displayName := func(p person, id uint64) string {
		if p.fullName != "" {
			return p.fullName
		}
		if p.username != "" {
			return p.username
		}
		if id > 0 {
			return fmt.Sprintf("User #%d", id)
		}
		return ""
	}

	for _, r := range reports {
		d := UserDisplay{}
		if r.UserID > 0 {
			p := people[r.UserID]
			d.UserFullName = p.fullName
			d.UserUsername = p.username
			d.UserName = displayName(p, r.UserID)
		}
		if r.AssigneeID != nil && *r.AssigneeID > 0 {
			aid := *r.AssigneeID
			p := people[aid]
			d.AssigneeFullName = p.fullName
			d.AssigneeUsername = p.username
			d.AssigneeName = displayName(p, aid)
		}
		out[r.ID] = d
	}
	return out
}

func (u *reportUsecase) AdminUpdateImages(ctx context.Context, reportID uint64, images []string) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	before := compactJSON(map[string]interface{}{"count": len(unmarshalReportImages(report.Images))})
	next := sanitizeImageList(images)
	b, err := json.Marshal(next)
	if err != nil {
		return err
	}
	report.Images = datatypes.JSON(b)
	if err := u.repo.Update(ctx, report); err != nil {
		return err
	}
	_ = u.addEvent(ctx, reportID, actorFromCtx(ctx), "images", before, compactJSON(map[string]interface{}{
		"count": len(next),
	}), "Đã tải ảnh lên")
	return nil
}

func unmarshalReportImages(raw datatypes.JSON) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func sanitizeImageList(images []string) []string {
	seen := make(map[string]struct{}, len(images))
	out := make([]string, 0, len(images))
	for _, s := range images {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func mergeUniqueStrings(existing, incoming []string) []string {
	return sanitizeImageList(append(append([]string{}, existing...), incoming...))
}

func (u *reportUsecase) AdminClose(ctx context.Context, reportID uint64, reject bool, note string) error {
	report, err := u.repo.GetByID(ctx, reportID)
	if err != nil || report == nil {
		return _errors.ReturnError(service.ReportNotFound)
	}
	to := enums.ReportQaStatusClosed
	action := "close"
	eventNote := "Đóng phản ánh"
	if reject {
		to = enums.ReportQaStatusRejected
		action = "reject"
		eventNote = "Từ chối phản ánh"
	}
	if !enums.CanTransitionReportQa(report.QaStatus, to) {
		return _errors.ReturnError(service.ReportStatusTransitionInvalid)
	}
	before := compactJSON(map[string]interface{}{"qaStatus": uint32(report.QaStatus)})
	report.QaStatus = to
	now := time.Now()
	report.ClosedAt = &now
	if err := u.repo.Update(ctx, report); err != nil {
		return err
	}
	if strings.TrimSpace(note) != "" {
		eventNote = strings.TrimSpace(note)
	}
	_ = u.addEvent(ctx, reportID, actorFromCtx(ctx), action, before, compactJSON(map[string]interface{}{
		"qaStatus": uint32(to),
	}), eventNote)
	return nil
}

func (u *reportUsecase) AdminDelete(ctx context.Context, reportID uint64) error {
	return u.repo.Delete(ctx, reportID)
}

func (u *reportUsecase) AdminListEvents(ctx context.Context, filter repo.AdminEventFilter) ([]domain.ReportEvent, int64, error) {
	return u.repo.ListEvents(ctx, filter)
}

func (u *reportUsecase) AdminListReportEvents(ctx context.Context, reportID uint64, limit int) ([]domain.ReportEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	return u.repo.ListEventsByReport(ctx, reportID, limit)
}

func actorFromCtx(ctx context.Context) uint64 {
	return _utils.GetProfileIdWithContext(ctx)
}

func compactJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func (u *reportUsecase) addEvent(ctx context.Context, reportID, actorID uint64, action, before, after, note string) error {
	if actorID == 0 {
		actorID = actorFromCtx(ctx)
	}
	return u.repo.AddEvent(ctx, &domain.ReportEvent{
		ReportID:   reportID,
		ActorID:    actorID,
		Action:     action,
		BeforeJSON: before,
		AfterJSON:  after,
		Note:       note,
		CreatedAt:  time.Now(),
	})
}

// ResolveActorNames batch tên người thực hiện cho nhật ký.
func (u *reportUsecase) ResolveActorNames(ctx context.Context, actorIDs []uint64) map[uint64]string {
	out := make(map[uint64]string, len(actorIDs))
	if len(actorIDs) == 0 {
		return out
	}
	seen := map[uint64]struct{}{}
	ids := make([]uint64, 0, len(actorIDs))
	for _, id := range actorIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	people := make(map[uint64]struct{ full, user string })
	if u.userProvider != nil {
		if profiles, err := u.userProvider.GetProfilesByIDs(ctx, ids); err == nil {
			for id, p := range profiles {
				if p == nil {
					continue
				}
				people[id] = struct{ full, user string }{full: strings.TrimSpace(p.FullName)}
			}
		}
	}
	if u.userClient != nil && u.userClient.InternalClient != nil {
		if resp, err := u.userClient.GetAuthAdminByIds(ctx, ids); err == nil && resp != nil {
			for _, a := range resp.GetData() {
				if a == nil {
					continue
				}
				key := a.GetUserId()
				if key == 0 {
					key = a.GetId()
				}
				cur := people[key]
				if fn := strings.TrimSpace(a.GetFullName()); fn != "" {
					cur.full = fn
				}
				if un := strings.TrimSpace(a.GetUsername()); un != "" {
					cur.user = un
				}
				people[key] = cur
				if aid := a.GetId(); aid > 0 && aid != key {
					people[aid] = cur
				}
			}
		}
	}
	for _, id := range ids {
		p := people[id]
		if p.full != "" {
			out[id] = p.full
		} else if p.user != "" {
			out[id] = p.user
		} else {
			out[id] = fmt.Sprintf("User #%d", id)
		}
	}
	return out
}
