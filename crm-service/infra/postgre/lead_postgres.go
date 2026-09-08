package postgre

import (
	base_enum "base/enum"
	_db "common/db"
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.LeadRepo
type LeadPostgre struct {
	*_db.BaseRepo
}

func NewPostgreCustomer(db *gorm.DB) *LeadPostgre {
	return &LeadPostgre{BaseRepo: _db.NewBaseRepo(db)}
}

func (r *LeadPostgre) Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	query := r.GetDB(c).
		Preload("Contact").
		Preload("Stage").
		Preload("Stage.Pipeline").
		Model(&domain.LeadEntity{}).
		Joins("JOIN tb_contact o ON o.id = customers.contact_id").
		Where("o.deleted_at is null and customers.deleted_at is null and o.owner_id = ? and o.owner_of = ?", ownerId, ownerType)
	return r.runLeadSearch(c, query, dto)
}

func (r *LeadPostgre) SearchByOwnerOf(c context.Context, ownerType base_enum.EOwnerOf, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	query := r.GetDB(c).
		Preload("Contact").
		Preload("Stage").
		Preload("Stage.Pipeline").
		Model(&domain.LeadEntity{}).
		Joins("JOIN tb_contact o ON o.id = customers.contact_id").
		Where("o.deleted_at is null and customers.deleted_at is null and o.owner_of = ?", ownerType)
	return r.runLeadSearch(c, query, dto)
}

func (r *LeadPostgre) runLeadSearch(c context.Context, query *gorm.DB, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error) {
	total := int64(0)
	customers := []domain.LeadEntity{}

	if dto.FullName != "" {
		q := "%" + strings.ToLower(dto.FullName) + "%"
		query = query.Where("(lower(o.full_name) like ? OR o.phone like ? OR lower(o.email) like ?)", q, q, q)
	}
	if dto.Phone != "" {
		query = query.Where("o.phone like ?", "%"+dto.Phone+"%")
	}

	if len(dto.Steps) > 0 {
		query = query.Where("step_id in (?)", dto.Steps)
	}

	if len(dto.StageIDs) > 0 {
		query = query.Where("stage_id in (?)", dto.StageIDs)
	}

	if len(dto.Sources) > 0 {
		query = query.Where("source in (?)", dto.Sources)
	}

	if len(dto.ChargePersonIds) > 0 {
		query = query.Where("charge_person_id in (?)", dto.ChargePersonIds)
	}

	if dto.FromDate != nil {
		query = query.Where("customers.created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("customers.created_at <= ?", dto.ToDate)
	}

	if dto.LastUpdatedAt != nil {
		query = query.Where("customers.updated_at >= ?", dto.LastUpdatedAt)
	}

	if dto.PipelineID != nil {
		query = query.Where("pipeline_id = ?", dto.PipelineID)
	}

	if dto.UnassignedOnly {
		query = query.Where("(charge_person_id is null OR charge_person_id = 0)")
	}

	if len(dto.OpportunityStatuses) > 0 {
		query = query.Where("opportunity_status in (?)", dto.OpportunityStatuses)
	}
	if len(dto.CustomerTypes) > 0 {
		query = query.Where("customer_type in (?)", dto.CustomerTypes)
	}
	if len(dto.AdminSources) > 0 {
		query = query.Where("admin_source in (?)", dto.AdminSources)
	}
	if len(dto.Probabilities) > 0 {
		query = query.Where("probability in (?)", dto.Probabilities)
	}
	if len(dto.ChurnRisks) > 0 {
		query = query.Where("churn_risk in (?)", dto.ChurnRisks)
	}
	if dto.UpgradeOnly {
		query = query.Where("upgrade_signal = true")
	}
	if dto.RenewalOnly {
		query = query.Where("renewal_signal = true")
	}
	if dto.InterestedPlan != "" {
		query = query.Where("interested_plan = ?", dto.InterestedPlan)
	}
	if dto.MinExpectedValue != nil {
		query = query.Where("expected_value >= ?", *dto.MinExpectedValue)
	}
	if dto.MaxExpectedValue != nil {
		query = query.Where("expected_value <= ?", *dto.MaxExpectedValue)
	}
	if dto.FollowUpFrom != nil {
		query = query.Where("next_follow_up_at >= ?", *dto.FollowUpFrom)
	}
	if dto.FollowUpTo != nil {
		query = query.Where("next_follow_up_at <= ?", *dto.FollowUpTo)
	}
	switch dto.Queue {
	case "overdue_followup":
		query = query.Where("next_follow_up_at is not null and next_follow_up_at < now() and opportunity_status not in (?)",
			[]enums.EOpportunityStatus{enums.OpportunityStatusWon, enums.OpportunityStatusLost, enums.OpportunityStatusCancelled})
	case "unassigned":
		query = query.Where("(charge_person_id is null OR charge_person_id = 0)")
	case "waiting_quote":
		query = query.Where("opportunity_status = ?", enums.OpportunityStatusWaitingQuote)
	case "waiting_approval":
		query = query.Where("opportunity_status = ?", enums.OpportunityStatusWaitingApproval)
	case "waiting_payment":
		query = query.Where("opportunity_status = ?", enums.OpportunityStatusWaitingPayment)
	case "renewal":
		query = query.Where("renewal_signal = true")
	case "upgrade":
		query = query.Where("upgrade_signal = true")
	case "churn":
		query = query.Where("churn_risk >= ?", enums.ChurnRiskMedium)
	case "new":
		query = query.Where("opportunity_status = ?", enums.OpportunityStatusNew)
	case "consulting":
		query = query.Where("opportunity_status = ?", enums.OpportunityStatusConsulting)
	case "business":
		query = query.Where("customer_type = ?", enums.CustomerTypeBusiness)
	case "api":
		query = query.Where("customer_type = ?", enums.CustomerTypeAPIPartner)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order(`
			case
				when next_follow_up_at is not null and next_follow_up_at < now() then 0
				when charge_person_id is null or charge_person_id = 0 then 1
				when opportunity_status in (60,70) then 2
				when opportunity_status = 100 then 3
				when renewal_signal then 4
				when upgrade_signal then 5
				when churn_risk >= 20 then 6
				when opportunity_status = 10 then 7
				else 8
			end,
			customers.priority desc,
			customers.updated_at desc
		`).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&customers).
		Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *LeadPostgre) adminBaseQuery(c context.Context, ownerType base_enum.EOwnerOf) *gorm.DB {
	return r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Joins("JOIN tb_contact o ON o.id = customers.contact_id").
		Where("o.deleted_at is null and customers.deleted_at is null and o.owner_of = ?", ownerType)
}

func (r *LeadPostgre) AdminOpportunitySummary(c context.Context, ownerType base_enum.EOwnerOf) (*dto.AdminOpportunitiesSummaryDTO, error) {
	base := r.adminBaseQuery(c, ownerType)
	out := &dto.AdminOpportunitiesSummaryDTO{}
	closed := []enums.EOpportunityStatus{
		enums.OpportunityStatusWon, enums.OpportunityStatusLost, enums.OpportunityStatusCancelled,
	}

	countWhere := func(where string, args ...any) (int64, error) {
		var n int64
		q := base.Session(&gorm.Session{})
		if where != "" {
			q = q.Where(where, args...)
		}
		err := q.Count(&n).Error
		return n, err
	}
	var err error
	if out.Total, err = countWhere(""); err != nil {
		return nil, err
	}
	if out.OverdueFollowUp, err = countWhere(
		"next_follow_up_at is not null and next_follow_up_at < now() and opportunity_status not in (?)", closed,
	); err != nil {
		return nil, err
	}
	if out.Unassigned, err = countWhere("(charge_person_id is null OR charge_person_id = 0)"); err != nil {
		return nil, err
	}
	if out.WaitingQuote, err = countWhere("opportunity_status = ?", enums.OpportunityStatusWaitingQuote); err != nil {
		return nil, err
	}
	if out.WaitingApproval, err = countWhere("opportunity_status = ?", enums.OpportunityStatusWaitingApproval); err != nil {
		return nil, err
	}
	if out.WaitingPayment, err = countWhere("opportunity_status = ?", enums.OpportunityStatusWaitingPayment); err != nil {
		return nil, err
	}
	if out.Renewal, err = countWhere("renewal_signal = true"); err != nil {
		return nil, err
	}
	if out.Upgrade, err = countWhere("upgrade_signal = true"); err != nil {
		return nil, err
	}
	if out.ChurnRisk, err = countWhere("churn_risk >= ?", enums.ChurnRiskMedium); err != nil {
		return nil, err
	}
	if out.StatusNew, err = countWhere("opportunity_status = ?", enums.OpportunityStatusNew); err != nil {
		return nil, err
	}
	if out.Consulting, err = countWhere("opportunity_status = ?", enums.OpportunityStatusConsulting); err != nil {
		return nil, err
	}
	if out.Won, err = countWhere("opportunity_status = ?", enums.OpportunityStatusWon); err != nil {
		return nil, err
	}
	if out.Lost, err = countWhere("opportunity_status = ?", enums.OpportunityStatusLost); err != nil {
		return nil, err
	}

	sumVal := func(where string, args ...any) (float64, error) {
		var sum *float64
		q := base.Session(&gorm.Session{}).Select("COALESCE(SUM(expected_value),0)")
		if where != "" {
			q = q.Where(where, args...)
		}
		err := q.Scan(&sum).Error
		if err != nil {
			return 0, err
		}
		if sum == nil {
			return 0, nil
		}
		return *sum, nil
	}
	if out.ExpectedRevenue, err = sumVal("opportunity_status not in (?)", closed); err != nil {
		return nil, err
	}
	if out.WaitingPaymentValue, err = sumVal("opportunity_status = ?", enums.OpportunityStatusWaitingPayment); err != nil {
		return nil, err
	}
	if out.WonValue, err = sumVal("opportunity_status = ?", enums.OpportunityStatusWon); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *LeadPostgre) AdminFunnel(c context.Context, ownerType base_enum.EOwnerOf, search dto.LeadSearchDTO) ([]dto.AdminFunnelBucketDTO, error) {
	search.Page = 0
	search.Size = 100000
	leads, _, err := r.SearchByOwnerOf(c, ownerType, search)
	if err != nil {
		return nil, err
	}
	counts := map[enums.EOpportunityStatus]dto.AdminFunnelBucketDTO{}
	for _, lead := range leads {
		b := counts[lead.OpportunityStatus]
		b.Status = lead.OpportunityStatus
		b.Label = enums.OpportunityStatusMap[lead.OpportunityStatus]
		b.Count++
		if lead.ExpectedValue != nil {
			b.ExpectedValueSum += *lead.ExpectedValue
		}
		counts[lead.OpportunityStatus] = b
	}
	order := []enums.EOpportunityStatus{
		enums.OpportunityStatusNew,
		enums.OpportunityStatusUnassigned,
		enums.OpportunityStatusAssigned,
		enums.OpportunityStatusContacted,
		enums.OpportunityStatusConsulting,
		enums.OpportunityStatusWaitingQuote,
		enums.OpportunityStatusWaitingApproval,
		enums.OpportunityStatusQuoteSent,
		enums.OpportunityStatusNegotiating,
		enums.OpportunityStatusWaitingPayment,
		enums.OpportunityStatusWon,
		enums.OpportunityStatusLost,
		enums.OpportunityStatusCancelled,
	}
	out := make([]dto.AdminFunnelBucketDTO, 0, len(order))
	for _, st := range order {
		if b, ok := counts[st]; ok {
			out = append(out, b)
		} else {
			out = append(out, dto.AdminFunnelBucketDTO{
				Status: st,
				Label:  enums.OpportunityStatusMap[st],
			})
		}
	}
	return out, nil
}

func (r *LeadPostgre) Create(ctx context.Context, entity *domain.LeadEntity) (*domain.LeadEntity, error) {
	err := r.GetDB(ctx).Create(entity).Error
	return entity, err
}

func (r *LeadPostgre) Update(ctx context.Context, id uint64, entity *domain.LeadEntity) (*domain.LeadEntity, error) {
	err := r.GetDB(ctx).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Updates(map[string]any{
			"source":              entity.Source,
			"assign_note":         entity.AssignNote,
			"stage_id":            entity.StageID,
			"priority":            entity.Priority,
			"charge_person_id":    entity.ChargePersonID,
			"charge_person_type":  entity.ChargePersonType,
			"pipeline_id":         entity.PipelineID,
			"stage_note":          entity.StageNote,
			"note":                entity.Note,
			"code":                entity.Code,
			"title":               entity.Title,
			"customer_type":       entity.CustomerType,
			"need_summary":        entity.NeedSummary,
			"interested_plan":     entity.InterestedPlan,
			"opportunity_status":   entity.OpportunityStatus,
			"admin_source":        entity.AdminSource,
			"expected_value":      entity.ExpectedValue,
			"probability":         entity.Probability,
			"expected_close_date": entity.ExpectedCloseDate,
			"next_follow_up_at":   entity.NextFollowUpAt,
			"churn_risk":          entity.ChurnRisk,
			"upgrade_signal":      entity.UpgradeSignal,
			"renewal_signal":      entity.RenewalSignal,
			"owner_team":          entity.OwnerTeam,
			"related_user_id":     entity.RelatedUserID,
			"related_business_id": entity.RelatedBusinessID,
			"proposal_ref":        entity.ProposalRef,
			"payment_request_ref": entity.PaymentRequestRef,
			"subscription_ref":    entity.SubscriptionRef,
			"ticket_ref":          entity.TicketRef,
			"segment":             entity.Segment,
			"region":              entity.Region,
			"tags":                entity.Tags,
			"closed_at":           entity.ClosedAt,
			"close_reason":        entity.CloseReason,
		}).Error
	return entity, err
}

func (r *LeadPostgre) Delete(ctx context.Context, id uint64) error {
	return r.GetDB(ctx).
		WithContext(ctx).
		Model(&domain.LeadEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).
		Error
}

func (r *LeadPostgre) GetByID(ctx context.Context, id uint64) (*domain.LeadEntity, error) {
	var customer domain.LeadEntity
	err := r.GetDB(ctx).
		Preload("Stage").
		Preload("Stage.Pipeline").
		Preload("Documents").
		Preload("ProductCares").
		Preload("Contact").
		Where("id = ? and deleted_at is null", id).
		First(&customer).Error

	if customer.Stage != nil {
		customer.Pipeline = customer.Stage.Pipeline
	}

	return &customer, err
}

func (r *LeadPostgre) Existed(c context.Context, entity *domain.LeadEntity) (bool, error) {
	return false, nil
}

func (r *LeadPostgre) UpdateNote(c context.Context, id uint64, note string) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("note", note).
		Error
	return &customer, err
}

func (r *LeadPostgre) Assign(c context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("charge_person_id", dto.ChargePersonId).
		Update("assign_note", dto.AssignNote).
		Error
	return &customer, err
}

func (r *LeadPostgre) SwitchStage(c context.Context, id uint64, stageID *uint64, note string) (*domain.LeadEntity, error) {
	customer := domain.LeadEntity{}
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("stage_id", stageID).
		Update("stage_note", note).
		Error
	return &customer, err
}

func (r *LeadPostgre) ExistedWithStageID(c context.Context, stageID uint64) (bool, error) {
	var count int64
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("stage_id = ? and deleted_at is null", stageID).
		Count(&count).Error
	return count > 0, err
}

func (r *LeadPostgre) ExistedWithPipelineID(c context.Context, pipelineID uint64) (bool, error) {
	var count int64
	err := r.GetDB(c).
		Model(&domain.LeadEntity{}).
		Where("pipeline_id = ? and deleted_at is null", pipelineID).
		Count(&count).Error
	return count > 0, err
}

func (r *LeadPostgre) GetOwnerWithRule(c context.Context) ([]dto.CustomerWithRuleDTO, error) {
	var customers []dto.CustomerWithRuleDTO

	query := `
	SELECT
		l.id AS lead_id,
		l.contact_id,
		l.pipeline_id,
		l.stage_id,
		l.charge_person_id,
		l.owner_id,
		l.owner_type,
		l.charge_person_type,
		l.priority,
		l.assign_note,
		l.stage_note,
		l.updated_at,

		-- Contact fields
		c.full_name,
		c.phone,
		c.email,
		c.address,
		c.birthday,
		c.note,

		-- Stage fields
		s.stage_name,

		-- Rule fields
		r.id AS rule_id,
		r.rule_name,
		r.condition,
		r.condition_value,
		r.trigger,
		r.trigger_value

	FROM customers l
	JOIN tb_contact c ON l.contact_id = c.id
	JOIN stage s ON l.stage_id = s.id
	JOIN rules r ON s.rule_id = r.id

	WHERE l.stage_id IS NOT NULL
	AND s.rule_id IS NOT NULL
	AND r.active = true;
	`
	err := r.GetDB(c).
		Raw(query).
		Scan(&customers).Error
	return customers, err
}

func (r *LeadPostgre) GetByContactID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error) {
	var customer domain.LeadEntity
	err := r.GetDB(ctx).
		Where("contact_id = ? and deleted_at is null", contactID).
		First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &_routes.Except{
				Code:    404,
				Message: "Liên hệ không tồn tại",
			}
		}
		return nil, err
	}
	return &customer, err
}

// func (r *LeadPostgre) DetailByID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error) {
// 	var customer domain.LeadEntity
// 	err := r.GetDB(ctx).
// 		Preload("Stage").
// 		Where("contact_id = ? and deleted_at is null", contactID).
// 		First(&customer).Error
// 	return &customer, err
// }

func (r *LeadPostgre) CareExpiredTime(c context.Context, organizationId uint64) (int64, error) {
	var count int64

	from := time.Now().AddDate(0, 0, -7) // 7 ngày trước
	to := from.AddDate(0, 0, 3)          // từ ngày đó + 3 ngày = khoảng 3 ngày

	err := r.GetDB(c).
		Debug().
		Model(&domain.LeadEntity{}).
		Where(`
			owner_id = ? 
			AND owner_type = ? 
			AND deleted_at IS NULL 
			AND updated_at BETWEEN ? AND ?`,
			organizationId, base_enum.EOwnerOfOrgnization, from, to).
		Count(&count).Error

	return count, err
}