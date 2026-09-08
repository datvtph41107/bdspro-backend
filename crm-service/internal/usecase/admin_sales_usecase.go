package usecase

import (
	base_enum "base/enum"
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AdminSalesUsecase struct {
	leadRepo     repo.LeadRepo
	contactRepo  repo.ContactRepo
	stageRepo    repo.StageRepo
	eventRepo    repo.AdminOpportunityEventRepo
	transaction  provider.ITransaction
	ownerUsecase *OwnerUsecase
}

func NewAdminSalesUsecase(
	leadRepo repo.LeadRepo,
	contactRepo repo.ContactRepo,
	stageRepo repo.StageRepo,
	eventRepo repo.AdminOpportunityEventRepo,
	transaction provider.ITransaction,
	ownerUsecase *OwnerUsecase,
) *AdminSalesUsecase {
	return &AdminSalesUsecase{
		leadRepo:     leadRepo,
		contactRepo:  contactRepo,
		stageRepo:    stageRepo,
		eventRepo:    eventRepo,
		transaction:  transaction,
		ownerUsecase: ownerUsecase,
	}
}

func (u *AdminSalesUsecase) List(ctx context.Context, search dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	return u.leadRepo.SearchByOwnerOf(ctx, base_enum.EOwnerOfAdmin, search)
}

func (u *AdminSalesUsecase) Summary(ctx context.Context) (*dto.AdminOpportunitiesSummaryDTO, error) {
	return u.leadRepo.AdminOpportunitySummary(ctx, base_enum.EOwnerOfAdmin)
}

func (u *AdminSalesUsecase) Funnel(ctx context.Context, search dto.LeadSearchDTO) ([]dto.AdminFunnelBucketDTO, error) {
	return u.leadRepo.AdminFunnel(ctx, base_enum.EOwnerOfAdmin, search)
}

func (u *AdminSalesUsecase) Get(ctx context.Context, id uint64) (*domain.LeadEntity, error) {
	lead, err := u.leadRepo.GetByID(ctx, id)
	if err != nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy cơ hội")
	}
	if err := u.ensureAdminLead(lead); err != nil {
		return nil, err
	}
	return lead, nil
}

func (u *AdminSalesUsecase) ListEvents(ctx context.Context, opportunityID uint64) ([]dto.AdminOpportunityEventResponse, error) {
	if _, err := u.Get(ctx, opportunityID); err != nil {
		return nil, err
	}
	events, err := u.eventRepo.ListEvents(ctx, opportunityID, 50)
	if err != nil {
		return nil, err
	}
	out := make([]dto.AdminOpportunityEventResponse, len(events))
	for i, e := range events {
		out[i] = dto.AdminOpportunityEventResponse{
			ID:         e.ID,
			ActorID:    e.ActorID,
			Action:     e.Action,
			BeforeJSON: e.BeforeJSON,
			AfterJSON:  e.AfterJSON,
			Note:       e.Note,
			CreatedAt:  e.CreatedAt,
		}
	}
	return out, nil
}

func (u *AdminSalesUsecase) ListContacts(ctx context.Context, page, size int32, q string) ([]domain.ContactEntity, int64, error) {
	search := dto.ContactSearchDTO{
		Pagable: _dto.Pagable{Page: uint32(page), Size: uint32(size)},
		Text:    strings.TrimSpace(q),
		OwnerOf: base_enum.EOwnerOfAdmin,
	}
	return u.contactRepo.SearchByOwnerOf(ctx, base_enum.EOwnerOfAdmin, search)
}

func (u *AdminSalesUsecase) CreateContact(ctx context.Context, in *dto.AdminContactSaveDTO) (*domain.ContactEntity, error) {
	if strings.TrimSpace(in.FullName) == "" {
		return nil, _errors.ReturnError(400, "Tên liên hệ là bắt buộc")
	}
	ownerID, err := u.ownerUsecase.GetOwnerInfo(ctx, base_enum.EOwnerOfAdmin, 0)
	if err != nil {
		return nil, err
	}
	contact := &domain.ContactEntity{
		FullName: strings.TrimSpace(in.FullName),
		Phone:    strings.TrimSpace(in.Phone),
		Email:    strings.TrimSpace(in.Email),
		Company:  strings.TrimSpace(in.Company),
		Avatar:   in.Avatar,
		Address:  in.Address,
		Zalo:     in.Zalo,
		Note:     in.Note,
		OwnerID:  ownerID,
		OwnerOf:  base_enum.EOwnerOfAdmin,
	}
	created, err := u.contactRepo.Create(ctx, contact)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (u *AdminSalesUsecase) Create(ctx context.Context, in *dto.AdminOpportunitySaveDTO) (*domain.LeadEntity, error) {
	if in.ContactID == nil || *in.ContactID == 0 {
		if strings.TrimSpace(in.FullName) == "" {
			return nil, _errors.ReturnError(400, "Vui lòng chọn liên hệ hoặc nhập tên khách hàng")
		}
	}

	ownerID, err := u.ownerUsecase.GetOwnerInfo(ctx, base_enum.EOwnerOfAdmin, 0)
	if err != nil {
		return nil, err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)

	var stage *domain.StageEntity
	if in.StageID != nil && *in.StageID > 0 {
		stage, err = u.stageRepo.GetByID(ctx, *in.StageID)
		if err != nil {
			return nil, err
		}
	}

	priority := enums.PriorityLow
	if in.Priority.IsValid() {
		priority = in.Priority
	}
	status := enums.OpportunityStatusNew
	if in.OpportunityStatus.IsValid() {
		status = in.OpportunityStatus
	} else if in.ChargePersonID == nil || *in.ChargePersonID == 0 {
		status = enums.OpportunityStatusUnassigned
	} else {
		status = enums.OpportunityStatusAssigned
	}
	customerType := enums.CustomerTypeIndividual
	if in.CustomerType.IsValid() {
		customerType = in.CustomerType
	}
	adminSource := enums.AdminSourceManual
	if in.AdminSource.IsValid() {
		adminSource = in.AdminSource
	}

	var lead *domain.LeadEntity
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		var contact *domain.ContactEntity
		if in.ContactID != nil && *in.ContactID > 0 {
			existing, err := u.contactRepo.GetByID(txCtx, *in.ContactID)
			if err != nil || existing == nil {
				return _errors.ReturnError(404, "Không tìm thấy liên hệ")
			}
			if existing.OwnerOf != base_enum.EOwnerOfAdmin {
				return _errors.ReturnError(403, "Liên hệ không thuộc phạm vi Admin")
			}
			contact = existing
		} else {
			createdContact, err := u.contactRepo.Create(txCtx, &domain.ContactEntity{
				FullName: strings.TrimSpace(in.FullName),
				Phone:    strings.TrimSpace(in.Phone),
				Email:    strings.TrimSpace(in.Email),
				Company:  strings.TrimSpace(in.Company),
				Avatar:   in.Avatar,
				Address:  in.Address,
				Zalo:     in.Zalo,
				Note:     in.Note,
				OwnerID:  ownerID,
				OwnerOf:  base_enum.EOwnerOfAdmin,
			})
			if err != nil {
				return err
			}
			contact = createdContact
		}

		entity := &domain.LeadEntity{
			ContactID:          contact.ID,
			Source:             in.Source,
			AssignNote:         in.AssignNote,
			Priority:           priority,
			ChargePersonID:     in.ChargePersonID,
			ChargePersonType:   enums.EOwnerOfAdmin,
			Note:               in.Note,
			Title:              adminSalesFirstNonEmpty(in.Title, contact.FullName),
			CustomerType:       customerType,
			NeedSummary:        in.NeedSummary,
			InterestedPlan:     in.InterestedPlan,
			OpportunityStatus:  status,
			AdminSource:        adminSource,
			ExpectedValue:      in.ExpectedValue,
			Probability:        in.Probability,
			ExpectedCloseDate:  in.ExpectedCloseDate,
			NextFollowUpAt:     in.NextFollowUpAt,
			ChurnRisk:          in.ChurnRisk,
			UpgradeSignal:      in.UpgradeSignal,
			RenewalSignal:      in.RenewalSignal,
			OwnerTeam:          in.OwnerTeam,
			Segment:            in.Segment,
			Region:             in.Region,
			Tags:               in.Tags,
			ProposalRef:        in.ProposalRef,
			PaymentRequestRef:  in.PaymentRequestRef,
			SubscriptionRef:    in.SubscriptionRef,
			TicketRef:          in.TicketRef,
		}
		if stage != nil {
			entity.StageID = &stage.ID
			entity.PipelineID = &stage.PipelineID
		}
		lead, err = u.leadRepo.Create(txCtx, entity)
		if err != nil {
			return err
		}
		lead.Code = fmt.Sprintf("KD-%d", lead.ID)
		if _, err := u.leadRepo.Update(txCtx, lead.ID, lead); err != nil {
			return err
		}
		// link contact.lead_id for convenience
		if _, err := u.contactRepo.UpdateContactID(txCtx, contact.ID, lead.ID); err != nil {
			return err
		}
		lead.Contact = contact
		lead.Stage = stage
		return u.addEvent(txCtx, lead.ID, actorID, "create", "", adminSalesMustJSON(map[string]any{
			"code":              lead.Code,
			"contactId":         contact.ID,
			"fullName":          contact.FullName,
			"opportunityStatus": status,
			"adminSource":       adminSource,
			"expectedValue":     in.ExpectedValue,
			"priority":          priority,
		}), "Tạo cơ hội kinh doanh")
	})
	if err != nil {
		return nil, err
	}
	return u.Get(ctx, lead.ID)
}

func (u *AdminSalesUsecase) Update(ctx context.Context, id uint64, in *dto.AdminOpportunitySaveDTO) (*domain.LeadEntity, error) {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if lead.Contact == nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy liên hệ")
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(snapshotLead(lead))

	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		contact := lead.Contact
		if strings.TrimSpace(in.FullName) != "" {
			contact.FullName = strings.TrimSpace(in.FullName)
		}
		contact.Phone = strings.TrimSpace(in.Phone)
		contact.Email = strings.TrimSpace(in.Email)
		contact.Company = strings.TrimSpace(in.Company)
		contact.Avatar = in.Avatar
		contact.Address = in.Address
		contact.Zalo = in.Zalo
		contact.Note = in.Note
		if _, err := u.contactRepo.Update(txCtx, contact.ID, contact); err != nil {
			return err
		}

		if in.Source != 0 {
			lead.Source = in.Source
		}
		lead.AssignNote = in.AssignNote
		lead.Note = in.Note
		if in.Priority.IsValid() {
			lead.Priority = in.Priority
		}
		if in.ChargePersonID != nil {
			lead.ChargePersonID = in.ChargePersonID
			lead.ChargePersonType = enums.EOwnerOfAdmin
		}
		if in.StageID != nil && *in.StageID > 0 {
			stage, err := u.stageRepo.GetByID(txCtx, *in.StageID)
			if err != nil {
				return err
			}
			lead.StageID = &stage.ID
			lead.PipelineID = &stage.PipelineID
		}
		applyOpportunityFields(lead, in)
		if _, err := u.leadRepo.Update(txCtx, lead.ID, lead); err != nil {
			return err
		}
		return u.addEvent(txCtx, lead.ID, actorID, "update", before, adminSalesMustJSON(snapshotLead(lead)), "Cập nhật cơ hội")
	})
	if err != nil {
		return nil, err
	}
	return u.Get(ctx, id)
}

func applyOpportunityFields(lead *domain.LeadEntity, in *dto.AdminOpportunitySaveDTO) {
	if strings.TrimSpace(in.Title) != "" {
		lead.Title = strings.TrimSpace(in.Title)
	}
	if in.CustomerType.IsValid() {
		lead.CustomerType = in.CustomerType
	}
	lead.NeedSummary = in.NeedSummary
	lead.InterestedPlan = in.InterestedPlan
	if in.OpportunityStatus.IsValid() {
		lead.OpportunityStatus = in.OpportunityStatus
	}
	if in.AdminSource.IsValid() {
		lead.AdminSource = in.AdminSource
	}
	if in.ExpectedValue != nil {
		lead.ExpectedValue = in.ExpectedValue
	}
	lead.Probability = in.Probability
	if in.ExpectedCloseDate != nil {
		lead.ExpectedCloseDate = in.ExpectedCloseDate
	}
	if in.NextFollowUpAt != nil {
		lead.NextFollowUpAt = in.NextFollowUpAt
	}
	lead.ChurnRisk = in.ChurnRisk
	lead.UpgradeSignal = in.UpgradeSignal
	lead.RenewalSignal = in.RenewalSignal
	lead.OwnerTeam = in.OwnerTeam
	lead.Segment = in.Segment
	lead.Region = in.Region
	lead.Tags = in.Tags
	if in.ProposalRef != "" {
		lead.ProposalRef = in.ProposalRef
	}
	if in.PaymentRequestRef != "" {
		lead.PaymentRequestRef = in.PaymentRequestRef
	}
	if in.SubscriptionRef != "" {
		lead.SubscriptionRef = in.SubscriptionRef
	}
	if in.TicketRef != "" {
		lead.TicketRef = in.TicketRef
	}
}

func (u *AdminSalesUsecase) Delete(ctx context.Context, id uint64) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(snapshotLead(lead))
	return u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.addEvent(txCtx, lead.ID, actorID, "delete", before, "", "Xóa cơ hội"); err != nil {
			return err
		}
		if err := u.leadRepo.Delete(txCtx, lead.ID); err != nil {
			return err
		}
		if lead.Contact != nil {
			return u.contactRepo.Delete(txCtx, lead.Contact.ID)
		}
		return nil
	})
}

func (u *AdminSalesUsecase) Assign(ctx context.Context, id uint64, in *dto.AdminOpportunityAssignDTO) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	if in.ChargePersonID == 0 {
		return _errors.ReturnError(400, "Người phụ trách không hợp lệ")
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(map[string]any{
		"chargePersonId":    lead.ChargePersonID,
		"assignNote":        lead.AssignNote,
		"opportunityStatus": lead.OpportunityStatus,
	})
	if _, err := u.leadRepo.Assign(ctx, id, &dto.LeadDTO{
		ChargePersonId: &in.ChargePersonID,
		AssignNote:     in.AssignNote,
	}); err != nil {
		return err
	}
	lead.ChargePersonID = &in.ChargePersonID
	lead.AssignNote = in.AssignNote
	if lead.OpportunityStatus == enums.OpportunityStatusNew ||
		lead.OpportunityStatus == enums.OpportunityStatusUnassigned ||
		lead.OpportunityStatus == 0 {
		lead.OpportunityStatus = enums.OpportunityStatusAssigned
	}
	if _, err := u.leadRepo.Update(ctx, id, lead); err != nil {
		return err
	}
	return u.addEvent(ctx, id, actorID, "assign", before, adminSalesMustJSON(map[string]any{
		"chargePersonId":    in.ChargePersonID,
		"assignNote":        in.AssignNote,
		"opportunityStatus": lead.OpportunityStatus,
	}), "Gán người phụ trách")
}

func (u *AdminSalesUsecase) SwitchStage(ctx context.Context, id uint64, stageID uint64, note string) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	if stageID == 0 {
		return _errors.ReturnError(400, "Stage không hợp lệ")
	}
	stage, err := u.stageRepo.GetByID(ctx, stageID)
	if err != nil {
		return err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(map[string]any{"stageId": lead.StageID, "pipelineId": lead.PipelineID})
	if _, err := u.leadRepo.SwitchStage(ctx, id, &stage.ID, note); err != nil {
		return err
	}
	lead.StageID = &stage.ID
	lead.PipelineID = &stage.PipelineID
	if _, err = u.leadRepo.Update(ctx, id, lead); err != nil {
		return err
	}
	return u.addEvent(ctx, id, actorID, "stage", before, adminSalesMustJSON(map[string]any{
		"stageId":    stage.ID,
		"pipelineId": stage.PipelineID,
		"stageName":  stage.StageName,
	}), adminSalesFirstNonEmpty(note, "Đổi giai đoạn"))
}

func (u *AdminSalesUsecase) UpdateNote(ctx context.Context, id uint64, note string) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(map[string]any{"note": lead.Note})
	if _, err := u.leadRepo.UpdateNote(ctx, id, note); err != nil {
		return err
	}
	return u.addEvent(ctx, id, actorID, "note", before, adminSalesMustJSON(map[string]any{"note": note}), "Cập nhật ghi chú")
}

func (u *AdminSalesUsecase) UpdateStatus(ctx context.Context, id uint64, in *dto.AdminOpportunityStatusDTO) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	if !in.OpportunityStatus.IsValid() {
		return _errors.ReturnError(400, "Trạng thái cơ hội không hợp lệ")
	}
	if lead.OpportunityStatus.IsClosed() {
		return _errors.ReturnError(400, "Trạng thái cơ hội không cho phép thao tác")
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(map[string]any{"opportunityStatus": lead.OpportunityStatus})
	lead.OpportunityStatus = in.OpportunityStatus
	if in.OpportunityStatus.IsClosed() {
		now := time.Now()
		lead.ClosedAt = &now
		lead.CloseReason = in.Note
	}
	if _, err := u.leadRepo.Update(ctx, id, lead); err != nil {
		return err
	}
	return u.addEvent(ctx, id, actorID, "status", before, adminSalesMustJSON(map[string]any{
		"opportunityStatus": in.OpportunityStatus,
		"note":              in.Note,
	}), adminSalesFirstNonEmpty(in.Note, "Cập nhật trạng thái cơ hội"))
}

func (u *AdminSalesUsecase) CreateFollowUp(ctx context.Context, id uint64, in *dto.AdminOpportunityFollowUpDTO) error {
	lead, err := u.Get(ctx, id)
	if err != nil {
		return err
	}
	if in.NextFollowUpAt.IsZero() {
		return _errors.ReturnError(400, "Thời điểm follow-up không hợp lệ")
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	before := adminSalesMustJSON(map[string]any{"nextFollowUpAt": lead.NextFollowUpAt})
	t := in.NextFollowUpAt
	lead.NextFollowUpAt = &t
	if strings.TrimSpace(in.Content) != "" {
		lead.Note = strings.TrimSpace(in.Content)
	}
	if _, err := u.leadRepo.Update(ctx, id, lead); err != nil {
		return err
	}
	return u.addEvent(ctx, id, actorID, "follow_up", before, adminSalesMustJSON(map[string]any{
		"nextFollowUpAt": in.NextFollowUpAt,
		"content":        in.Content,
		"nextAction":     in.NextAction,
	}), adminSalesFirstNonEmpty(in.NextAction, in.Content, "Tạo nhắc việc follow-up"))
}

func (u *AdminSalesUsecase) Close(ctx context.Context, id uint64, in *dto.AdminOpportunityCloseDTO) error {
	if in.Result != enums.OpportunityStatusWon &&
		in.Result != enums.OpportunityStatusLost &&
		in.Result != enums.OpportunityStatusCancelled {
		return _errors.ReturnError(400, "Kết quả đóng không hợp lệ")
	}
	if strings.TrimSpace(in.CloseReason) == "" {
		return _errors.ReturnError(400, "Lý do đóng là bắt buộc")
	}
	return u.UpdateStatus(ctx, id, &dto.AdminOpportunityStatusDTO{
		OpportunityStatus: in.Result,
		Note:              in.CloseReason,
	})
}

func (u *AdminSalesUsecase) Export(ctx context.Context, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return _errors.ReturnError(400, "Vui lòng nhập lý do xuất dữ liệu vì có thể chứa thông tin nhạy cảm")
	}
	actorID := _utils.GetProfileIdWithContext(ctx)
	// Audit stub — file export wiring deferred; record intent
	return u.eventRepo.AddEvent(ctx, &domain.AdminOpportunityEvent{
		OpportunityID: 0,
		ActorID:       actorID,
		Action:        "export",
		BeforeJSON:    "",
		AfterJSON:     adminSalesMustJSON(map[string]any{"reason": reason}),
		Note:          "Yêu cầu xuất danh sách cơ hội: " + strings.TrimSpace(reason),
	})
}

func (u *AdminSalesUsecase) addEvent(ctx context.Context, opportunityID, actorID uint64, action, before, after, note string) error {
	return u.eventRepo.AddEvent(ctx, &domain.AdminOpportunityEvent{
		OpportunityID: opportunityID,
		ActorID:       actorID,
		Action:        action,
		BeforeJSON:    before,
		AfterJSON:     after,
		Note:          note,
	})
}

func (u *AdminSalesUsecase) ensureAdminLead(lead *domain.LeadEntity) error {
	if lead == nil || lead.ID == 0 {
		return _errors.ReturnError(404, "Không tìm thấy cơ hội")
	}
	if lead.Contact == nil {
		return _errors.ReturnError(404, "Không tìm thấy liên hệ")
	}
	if lead.Contact.OwnerOf != base_enum.EOwnerOfAdmin {
		return _errors.ReturnError(403, "Cơ hội không thuộc phạm vi Admin")
	}
	return nil
}

func snapshotLead(lead *domain.LeadEntity) map[string]any {
	m := map[string]any{
		"code":              lead.Code,
		"title":             lead.Title,
		"source":            lead.Source,
		"adminSource":       lead.AdminSource,
		"priority":          lead.Priority,
		"opportunityStatus": lead.OpportunityStatus,
		"customerType":      lead.CustomerType,
		"needSummary":       lead.NeedSummary,
		"interestedPlan":    lead.InterestedPlan,
		"expectedValue":     lead.ExpectedValue,
		"probability":       lead.Probability,
		"nextFollowUpAt":    lead.NextFollowUpAt,
		"churnRisk":         lead.ChurnRisk,
		"upgradeSignal":     lead.UpgradeSignal,
		"renewalSignal":     lead.RenewalSignal,
		"stageId":           lead.StageID,
		"pipelineId":        lead.PipelineID,
		"chargePersonId":    lead.ChargePersonID,
		"note":              lead.Note,
		"assignNote":        lead.AssignNote,
	}
	if lead.Contact != nil {
		m["fullName"] = lead.Contact.FullName
		m["phone"] = lead.Contact.Phone
		m["email"] = lead.Contact.Email
		m["company"] = lead.Contact.Company
	}
	return m
}

func adminSalesMustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

func adminSalesFirstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// BuildAdminLeadSearch maps list query params.
func BuildAdminLeadSearch(
	page, size int32,
	q string,
	stageIDs, chargeIDs []uint64,
	pipelineID *uint64,
	sources []int32,
	unassignedOnly bool,
	opportunityStatuses, customerTypes, adminSources, probabilities, churnRisks []int32,
	queue string,
	followUpFrom, followUpTo *time.Time,
	minEV, maxEV *float64,
	upgradeOnly, renewalOnly bool,
	interestedPlan string,
) dto.LeadSearchDTO {
	search := dto.LeadSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(page),
			Size: uint32(size),
		},
		FullName:        strings.TrimSpace(q),
		StageIDs:        stageIDs,
		ChargePersonIds: chargeIDs,
		PipelineID:      pipelineID,
		UnassignedOnly:  unassignedOnly,
		Queue:           queue,
		FollowUpFrom:    followUpFrom,
		FollowUpTo:      followUpTo,
		MinExpectedValue: minEV,
		MaxExpectedValue: maxEV,
		UpgradeOnly:     upgradeOnly,
		RenewalOnly:     renewalOnly,
		InterestedPlan:  interestedPlan,
		OwnerOf:         base_enum.EOwnerOfAdmin,
	}
	for _, s := range sources {
		search.Sources = append(search.Sources, enums.ESourceLead(s))
	}
	for _, s := range opportunityStatuses {
		search.OpportunityStatuses = append(search.OpportunityStatuses, enums.EOpportunityStatus(s))
	}
	for _, s := range customerTypes {
		search.CustomerTypes = append(search.CustomerTypes, enums.ECustomerType(s))
	}
	for _, s := range adminSources {
		search.AdminSources = append(search.AdminSources, enums.EAdminSalesSource(s))
	}
	for _, s := range probabilities {
		search.Probabilities = append(search.Probabilities, enums.EWinProbability(s))
	}
	for _, s := range churnRisks {
		search.ChurnRisks = append(search.ChurnRisks, enums.EChurnRisk(s))
	}
	return search
}

// BuildOpportunityAlerts derives badge alerts for list/quick-view.
func BuildOpportunityAlerts(lead *domain.LeadEntity) []string {
	if lead == nil {
		return nil
	}
	var alerts []string
	now := time.Now()
	if lead.NextFollowUpAt != nil && lead.NextFollowUpAt.Before(now) && !lead.OpportunityStatus.IsClosed() {
		alerts = append(alerts, "followup_overdue")
	} else if lead.NextFollowUpAt != nil && !lead.OpportunityStatus.IsClosed() {
		if lead.NextFollowUpAt.Sub(now) <= 24*time.Hour {
			alerts = append(alerts, "followup_due")
		}
	}
	if lead.ChargePersonID == nil || *lead.ChargePersonID == 0 {
		alerts = append(alerts, "unassigned")
	}
	switch lead.OpportunityStatus {
	case enums.OpportunityStatusNew:
		alerts = append(alerts, "new")
	case enums.OpportunityStatusWaitingApproval:
		alerts = append(alerts, "proposal_pending")
	case enums.OpportunityStatusWaitingPayment:
		alerts = append(alerts, "payment_pending")
	}
	if lead.RenewalSignal {
		alerts = append(alerts, "renewal")
	}
	if lead.UpgradeSignal {
		alerts = append(alerts, "upgrade")
	}
	if lead.ChurnRisk >= enums.ChurnRiskMedium {
		alerts = append(alerts, "churn")
	}
	return alerts
}
