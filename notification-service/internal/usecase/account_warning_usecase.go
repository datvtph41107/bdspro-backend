package usecase

import (
	_errors "common/errors"
	"context"
	"strings"
	"time"

	_middleware "common/middleware"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/enums"
	)

type AccountWarningUsecase struct {
	AccountWarningRepo AccountWarningStore
}

func NewAccountWarningUsecase(
	accountWarningRepo AccountWarningStore,
) *AccountWarningUsecase {
	return &AccountWarningUsecase{
		AccountWarningRepo: accountWarningRepo,
	}
}

// SendAccountWarning gửi cảnh báo tài khoản
func (u *AccountWarningUsecase) SendAccountWarning(ctx context.Context, req dto.SendAccountWarningRequest) (*dto.SendAccountWarningResponse, error) {
	// Kiểm tra quyền admin
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Kiểm tra không gửi cảnh báo cho chính mình
	profileID, err := u.getProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if profileID == req.TargetID {
		return nil, _errors.ReturnError(int32(400), "Không thể gửi cảnh báo cho chính tài khoản của mình")
	}

	// Kiểm tra nội dung tối thiểu
	if len(strings.TrimSpace(req.Content)) < 20 {
		return nil, _errors.ReturnError(int32(400), "Nội dung cảnh báo phải có ít nhất 20 ký tự")
	}

	// Kiểm tra cảnh báo trùng lặp trong 24h
	isDuplicate, err := u.AccountWarningRepo.CheckDuplicateWarning(ctx, req.TargetID, req.Content, 24)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi kiểm tra cảnh báo trùng lặp")
	}
	if isDuplicate {
		return nil, _errors.ReturnError(int32(400), "Không được gửi 2 cảnh báo trùng nội dung trong vòng 24 giờ")
	}

	// Set default severity nếu không có
	if req.Severity == "" {
		req.Severity = "medium"
	}

	// Tạo cảnh báo
	warning := &domain.AccountWarningEntity{
		TargetID:    req.TargetID,
		TargetType:  enums.TargetTypeEnum(req.TargetType),
		WarningType: enums.WarningTypeEnum(req.WarningType),
		Title:       req.Title,
		Content:     req.Content,
		Severity:    enums.SeverityEnum(req.Severity),
		Status:      enums.WarningStatusSent,
		IsRead:      false,
		SentBy:      profileID,
		EmailSent:   req.SendEmail,
		RelatedID:   req.RelatedID,
		RelatedType: req.RelatedType,
		ExpiresAt:   req.ExpiresAt,
		TemplateID:  req.TemplateID, // Lưu ID mẫu nếu có
	}

	createdWarning, err := u.AccountWarningRepo.Create(ctx, warning)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo cảnh báo")
	}

	// Tạo log
	log := &domain.AccountWarningLogEntity{
		WarningID:   createdWarning.ID,
		TargetID:    req.TargetID,
		Action:      "sent",
		Status:      "success",
		Reason:      "Cảnh báo được gửi thành công",
		PerformedBy: profileID,
	}
	u.AccountWarningRepo.CreateLog(ctx, log)

	// TODO: Gửi email nếu có chọn
	if req.SendEmail {
		// Gọi email service để gửi email
		// u.EmailService.SendWarningEmail(ctx, req.TargetID, req.Title, req.Content)
		now := time.Now()
		createdWarning.EmailSent = true
		createdWarning.EmailSentAt = &now
		u.AccountWarningRepo.Update(ctx, createdWarning)
	}

	return &dto.SendAccountWarningResponse{
		ID:          createdWarning.ID,
		TargetID:    createdWarning.TargetID,
		TargetType:  string(createdWarning.TargetType),
		WarningType: string(createdWarning.WarningType),
		Title:       createdWarning.Title,
		Content:     createdWarning.Content,
		Severity:    string(createdWarning.Severity),
		Status:      string(createdWarning.Status),
		SentBy:      createdWarning.SentBy,
		SentAt:      createdWarning.SentAt,
		EmailSent:   createdWarning.EmailSent,
		Message:     "Cảnh báo đã được gửi thành công",
	}, nil
}

// GetAccountWarningList lấy danh sách cảnh báo
func (u *AccountWarningUsecase) GetAccountWarningList(ctx context.Context, req dto.GetAccountWarningListRequest) (*dto.GetAccountWarningListResponse, error) {
	// Log context values for debugging
	// _utils.LogContextValues(ctx, "GetAccountWarningList")

	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Set default pagination - page bắt đầu từ 0
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	// Build filters
	filters := make(map[string]interface{})
	if req.TargetID != nil {
		filters["target_id"] = *req.TargetID
	}
	if req.TargetType != "" {
		filters["target_type"] = req.TargetType
	}
	if req.WarningType != "" {
		filters["warning_type"] = req.WarningType
	}
	if req.Severity != "" {
		filters["severity"] = req.Severity
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.SentBy != nil {
		filters["sent_by"] = *req.SentBy
	}
	if req.IsRead != nil {
		filters["is_read"] = *req.IsRead
	}

	// Lấy danh sách - giữ nguyên page = 0
	warnings, total, err := u.AccountWarningRepo.FindAll(ctx, req.Page, req.Size, filters)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách cảnh báo")
	}

	// Convert to response
	items := make([]*dto.AccountWarningItem, len(warnings))
	for i, warning := range warnings {
		items[i] = &dto.AccountWarningItem{
			ID:          warning.ID,
			TargetID:    warning.TargetID,
			TargetType:  string(warning.TargetType),
			WarningType: string(warning.WarningType),
			Title:       warning.Title,
			Content:     warning.Content,
			Severity:    string(warning.Severity),
			Status:      string(warning.Status),
			IsRead:      warning.IsRead,
			ReadAt:      warning.ReadAt,
			SentBy:      warning.SentBy,
			SentAt:      warning.SentAt,
			EmailSent:   warning.EmailSent,
			EmailSentAt: warning.EmailSentAt,
			RelatedID:   warning.RelatedID,
			RelatedType: warning.RelatedType,
			ExpiresAt:   warning.ExpiresAt,
			CreatedAt:   *warning.BaseEntity.CreatedAt,
			UpdatedAt:   *warning.BaseEntity.UpdatedAt,
		}
	}

	return &dto.GetAccountWarningListResponse{
		Data:  items,
		Total: total,
		Page:  req.Page, // Trả về page gốc (0-based)
		Size:  req.Size,
	}, nil
}

// MarkWarningAsRead đánh dấu đã đọc
func (u *AccountWarningUsecase) MarkWarningAsRead(ctx context.Context, req dto.MarkWarningAsReadRequest) (*dto.MarkWarningAsReadResponse, error) {
	// Kiểm tra quyền
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}
	profileID, err := u.getProfileID(ctx)
	if err != nil {
		return nil, err
	}
	// Kiểm tra cảnh báo tồn tại
	warning, err := u.AccountWarningRepo.FindByID(ctx, req.WarningID)
	if err != nil || warning == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy cảnh báo")
	}

	// Kiểm tra quyền đọc (chỉ người nhận mới được đánh dấu đã đọc)
	if warning.TargetID != req.TargetID {
		return nil, _errors.ReturnError(int32(403), "Không có quyền đánh dấu đã đọc cảnh báo này")
	}

	err = u.AccountWarningRepo.MarkAsRead(ctx, req.WarningID, req.TargetID)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi đánh dấu đã đọc")
	}

	// Tạo log
	log := &domain.AccountWarningLogEntity{
		WarningID:   req.WarningID,
		TargetID:    req.TargetID,
		Action:      "read",
		Status:      "success",
		Reason:      "Cảnh báo được đánh dấu đã đọc",
		PerformedBy: profileID,
	}
	u.AccountWarningRepo.CreateLog(ctx, log)

	return &dto.MarkWarningAsReadResponse{
		Success: true,
		Message: "Đã đánh dấu đã đọc",
	}, nil
}

// MarkWarningAsAcknowledged đánh dấu đã xác nhận
func (u *AccountWarningUsecase) MarkWarningAsAcknowledged(ctx context.Context, req dto.MarkWarningAsAcknowledgedRequest) (*dto.MarkWarningAsAcknowledgedResponse, error) {
	// Kiểm tra quyền
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}
	profileID, err := u.getProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Kiểm tra cảnh báo tồn tại
	warning, err := u.AccountWarningRepo.FindByID(ctx, req.WarningID)
	if err != nil || warning == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy cảnh báo")
	}

	// Kiểm tra quyền xác nhận (chỉ người nhận mới được xác nhận)
	if warning.TargetID != req.TargetID {
		return nil, _errors.ReturnError(int32(403), "Không có quyền xác nhận cảnh báo này")
	}

	err = u.AccountWarningRepo.MarkAsAcknowledged(ctx, req.WarningID, req.TargetID)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi đánh dấu đã xác nhận")
	}

	// Tạo log
	log := &domain.AccountWarningLogEntity{
		WarningID:   req.WarningID,
		TargetID:    req.TargetID,
		Action:      "acknowledged",
		Status:      "success",
		Reason:      "Cảnh báo được đánh dấu đã xác nhận",
		PerformedBy: profileID,
	}
	u.AccountWarningRepo.CreateLog(ctx, log)

	return &dto.MarkWarningAsAcknowledgedResponse{
		Success: true,
		Message: "Đã đánh dấu đã xác nhận",
	}, nil
}

// CreateWarningTemplate tạo mẫu cảnh báo
func (u *AccountWarningUsecase) CreateWarningTemplate(ctx context.Context, req dto.CreateWarningTemplateRequest) (*dto.CreateWarningTemplateResponse, error) {
	// Kiểm tra quyền admin
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Set default severity nếu không có
	if req.Severity == "" {
		req.Severity = "medium"
	}

	// Tạo mẫu
	template := &domain.AccountWarningTemplateEntity{
		Name:        req.Name,
		WarningType: enums.WarningTypeEnum(req.WarningType),
		Title:       req.Title,
		Content:     req.Content,
		Severity:    enums.SeverityEnum(req.Severity),
		IsActive:    req.IsActive,
	}

	createdTemplate, err := u.AccountWarningRepo.CreateTemplate(ctx, template)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo mẫu cảnh báo")
	}

	return &dto.CreateWarningTemplateResponse{
		ID:          createdTemplate.ID,
		Name:        createdTemplate.Name,
		WarningType: string(createdTemplate.WarningType),
		Title:       createdTemplate.Title,
		Content:     createdTemplate.Content,
		Severity:    string(createdTemplate.Severity),
		IsActive:    createdTemplate.IsActive,
		CreatedBy:   createdTemplate.CreatedBy,
		UpdatedBy:   createdTemplate.UpdatedBy,
		CreatedAt:   *createdTemplate.BaseEntity.CreatedAt,
		UpdatedAt:   *createdTemplate.BaseEntity.UpdatedAt,
	}, nil
}

// GetWarningTemplateList lấy danh sách mẫu cảnh báo
func (u *AccountWarningUsecase) GetWarningTemplateList(ctx context.Context, req dto.GetWarningTemplateListRequest) (*dto.GetWarningTemplateListResponse, error) {
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Set default pagination - page bắt đầu từ 0
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	var templates []*domain.AccountWarningTemplateEntity
	var total int64

	// Lấy mẫu theo loại hoặc tất cả
	if req.WarningType != "" {
		templates, err = u.AccountWarningRepo.FindTemplatesByType(ctx, req.WarningType)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách mẫu")
		}
		total = int64(len(templates))

		// Xử lý pagination thủ công
		start := req.Page * req.Size
		end := start + req.Size
		if start >= len(templates) {
			templates = []*domain.AccountWarningTemplateEntity{}
		} else if end > len(templates) {
			templates = templates[start:]
		} else {
			templates = templates[start:end]
		}
	} else {
		templates, total, err = u.AccountWarningRepo.FindAllTemplates(ctx, req.Page, req.Size)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách mẫu")
		}
	}

	// Convert to response
	items := make([]*dto.WarningTemplateItem, len(templates))
	for i, template := range templates {
		items[i] = &dto.WarningTemplateItem{
			ID:          template.ID,
			Name:        template.Name,
			WarningType: string(template.WarningType),
			Title:       template.Title,
			Content:     template.Content,
			Severity:    string(template.Severity),
			IsActive:    template.IsActive,
			CreatedBy:   template.CreatedBy,
			UpdatedBy:   template.UpdatedBy,
			CreatedAt:   *template.BaseEntity.CreatedAt,
			UpdatedAt:   *template.BaseEntity.UpdatedAt,
		}
	}

	return &dto.GetWarningTemplateListResponse{
		Data:  items,
		Total: total,
		Page:  req.Page, // Trả về page gốc (0-based)
		Size:  req.Size,
	}, nil
}

// GetWarningTemplateDetail lấy chi tiết mẫu cảnh báo theo ID
func (u *AccountWarningUsecase) GetWarningTemplateDetail(ctx context.Context, id uint64) (*dto.WarningTemplateItem, error) {
	// Kiểm tra quyền admin
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Lấy template theo ID
	template, err := u.AccountWarningRepo.FindTemplateByID(ctx, id)
	if err != nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy mẫu cảnh báo")
	}

	return &dto.WarningTemplateItem{
		ID:          template.ID,
		Name:        template.Name,
		WarningType: string(template.WarningType),
		Title:       template.Title,
		Content:     template.Content,
		Severity:    string(template.Severity),
		IsActive:    template.IsActive,
		CreatedBy:   template.CreatedBy,
		UpdatedBy:   template.UpdatedBy,
		CreatedAt:   *template.BaseEntity.CreatedAt,
		UpdatedAt:   *template.BaseEntity.UpdatedAt,
	}, nil
}

// GetWarningLogList lấy danh sách log cảnh báo
func (u *AccountWarningUsecase) GetWarningLogList(ctx context.Context, req dto.GetWarningLogListRequest) (*dto.GetWarningLogListResponse, error) {
	// Kiểm tra quyền admin
	hasAdminRole, err := u.hasAdminRole(ctx)
	if err != nil {
		return nil, err
	}
	if !hasAdminRole {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}
	profileID, err := u.getProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Set default pagination - page bắt đầu từ 0
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	var logs []*domain.AccountWarningLogEntity
	var total int64

	// Lấy log theo WarningID hoặc TargetID - giữ nguyên page = 0
	if req.WarningID != nil {
		logs, total, err = u.AccountWarningRepo.FindLogsByWarningID(ctx, *req.WarningID, req.Page, req.Size)
	} else if req.TargetID != nil {
		logs, total, err = u.AccountWarningRepo.FindLogsByTargetID(ctx, *req.TargetID, req.Page, req.Size)
	} else {
		return nil, _errors.ReturnError(int32(400), "Cần cung cấp WarningID hoặc TargetID")
	}

	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách log")
	}

	// Convert to response
	items := make([]*dto.WarningLogItem, len(logs))
	for i, log := range logs {
		items[i] = &dto.WarningLogItem{
			ID:          log.ID,
			WarningID:   log.WarningID,
			TargetID:    log.TargetID,
			Action:      log.Action,
			Status:      log.Status,
			Reason:      log.Reason,
			PerformedBy: log.PerformedBy,
			PerformedAt: log.PerformedAt,
			IPAddress:   log.IPAddress,
			UserAgent:   log.UserAgent,
		}
	}

	return &dto.GetWarningLogListResponse{
		Data:  items,
		Total: total,
		Page:  req.Page, // Trả về page gốc (0-based)
		Size:  req.Size,
	}, nil
}

func (u *AccountWarningUsecase) hasAdminRole(ctx context.Context) (bool, error) {
	principal, err := _middleware.PrincipalFromContext(ctx)
	if err != nil {
		return false, err
	}

	return principal.Role == "ROLE_ADMIN", nil
}

func (u *AccountWarningUsecase) getProfileID(ctx context.Context) (uint64, error) {
	principal, err := _middleware.PrincipalFromContext(ctx)
	if err != nil {
		return 0, err
	}

	return principal.ProfileId, nil
}
