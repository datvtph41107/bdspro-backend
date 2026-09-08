package domain

import (
	_models "common/models"
	"crm/internal/enums"
	"time"
)

type LeadEntity struct {
	_models.BaseEntity
	ContactID        uint64            `json:"contactId"`
	Source           enums.ESourceLead `json:"customerSource"`
	AssignNote       string            `json:"assignNote"`
	PipelineID       *uint64           `json:"pipelineId"`
	StageID          *uint64           `json:"stageId"`
	Priority         enums.EPriority   `json:"priority"`
	ChargePersonID   *uint64           `json:"chargePersonId"`
	ChargePersonType enums.EOwnerOf    `json:"chargePersonType"`
	StageNote        string            `json:"stageNote"`
	Note             string            `json:"note"`

	// Admin Sales / IV.10.9 opportunity fields
	Code               string                   `gorm:"column:code;type:varchar(64);index" json:"code"`
	Title              string                   `gorm:"column:title;type:varchar(255)" json:"title"`
	CustomerType       enums.ECustomerType      `gorm:"column:customer_type;default:10" json:"customerType"`
	NeedSummary        string                   `gorm:"column:need_summary;type:text" json:"needSummary"`
	InterestedPlan     string                   `gorm:"column:interested_plan;type:varchar(128)" json:"interestedPlan"`
	OpportunityStatus  enums.EOpportunityStatus `gorm:"column:opportunity_status;default:10;index" json:"opportunityStatus"`
	AdminSource        enums.EAdminSalesSource  `gorm:"column:admin_source;default:70" json:"adminSource"`
	ExpectedValue      *float64                 `gorm:"column:expected_value;type:numeric(18,2)" json:"expectedValue"`
	Probability        enums.EWinProbability    `gorm:"column:probability;default:0" json:"probability"`
	ExpectedCloseDate  *time.Time               `gorm:"column:expected_close_date;type:date" json:"expectedCloseDate"`
	NextFollowUpAt     *time.Time               `gorm:"column:next_follow_up_at;index" json:"nextFollowUpAt"`
	ChurnRisk          enums.EChurnRisk         `gorm:"column:churn_risk;default:0" json:"churnRisk"`
	UpgradeSignal      bool                     `gorm:"column:upgrade_signal;default:false" json:"upgradeSignal"`
	RenewalSignal      bool                     `gorm:"column:renewal_signal;default:false" json:"renewalSignal"`
	OwnerTeam          string                   `gorm:"column:owner_team;type:varchar(128)" json:"ownerTeam"`
	RelatedUserID      *uint64                  `gorm:"column:related_user_id" json:"relatedUserId"`
	RelatedBusinessID  *uint64                  `gorm:"column:related_business_id" json:"relatedBusinessId"`
	ProposalRef        string                   `gorm:"column:proposal_ref;type:varchar(128)" json:"proposalRef"`
	PaymentRequestRef  string                   `gorm:"column:payment_request_ref;type:varchar(128)" json:"paymentRequestRef"`
	SubscriptionRef    string                   `gorm:"column:subscription_ref;type:varchar(128)" json:"subscriptionRef"`
	TicketRef          string                   `gorm:"column:ticket_ref;type:varchar(128)" json:"ticketRef"`
	Segment            string                   `gorm:"column:segment;type:varchar(128)" json:"segment"`
	Region             string                   `gorm:"column:region;type:varchar(255)" json:"region"`
	Tags               string                   `gorm:"column:tags;type:text" json:"tags"`
	ClosedAt           *time.Time               `gorm:"column:closed_at" json:"closedAt"`
	CloseReason        string                   `gorm:"column:close_reason;type:text" json:"closeReason"`

	Stage        *StageEntity        `gorm:"references:ID" json:"stage"`
	Pipeline     *PipelineEntity     `gorm:"-" json:"pipeline"`
	Documents    []DocumentEntity    `gorm:"foreignKey:LeadID" json:"documents"`
	ProductCares []ProductCareEntity `gorm:"foreignKey:LeadID" json:"productCares"`
	Contact      *ContactEntity      `gorm:"foreignKey:ContactID;references:ID" json:"contact"`
}

func (e *LeadEntity) TableName() string {
	return "customers"
}

type LeadContactEntity struct {
	ID         uint64            `json:"id"`
	ContactID  *uint64           `json:"contactId"`
	Source     enums.ESourceLead `json:"customerSource"`
	FullName   string            `json:"fullName" binding:"required"`
	Phone      string            `json:"phone"`
	Email      string            `json:"email"`
	Address    string            `json:"address"`
	Birthday   *time.Time        `json:"birthday"`
	Note       string            `json:"note"`
	AssignNote string            `json:"assignNote"`
	PipelineID *uint64           `json:"pipelineId"`
	StageID    *uint64           `json:"stageId"`
	ChargeId   *uint64           `json:"chargeId"`
	ChargeType enums.EOwnerOf    `json:"chargeType"`
	OwnerID    *uint64           `json:"ownerId" binding:"required"`
	OwnerType  enums.EOwnerOf    `json:"ownerType" binding:"required"`

	// Step       enums.EStep           `gorm:"default:10" json:"step"`
	// todo: sắp trễ/ chưa phân công
}