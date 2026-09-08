package dto

import (
	"crm/internal/enums"
	"time"
)

type AdminOpportunitySaveDTO struct {
	FullName       string
	Phone          string
	Email          string
	Company        string
	Avatar         string
	Address        string
	Zalo           string
	Note           string
	Source         enums.ESourceLead
	StageID        *uint64
	ChargePersonID *uint64
	Priority       enums.EPriority
	AssignNote     string
	ContactID      *uint64 // if set: attach opportunity to existing admin contact

	Title             string
	CustomerType      enums.ECustomerType
	NeedSummary       string
	InterestedPlan    string
	OpportunityStatus enums.EOpportunityStatus
	AdminSource       enums.EAdminSalesSource
	ExpectedValue     *float64
	Probability       enums.EWinProbability
	ExpectedCloseDate *time.Time
	NextFollowUpAt    *time.Time
	ChurnRisk         enums.EChurnRisk
	UpgradeSignal     bool
	RenewalSignal     bool
	OwnerTeam         string
	Segment           string
	Region            string
	Tags              string
	ProposalRef       string
	PaymentRequestRef string
	SubscriptionRef   string
	TicketRef         string
}

type AdminContactSaveDTO struct {
	FullName string
	Phone    string
	Email    string
	Company  string
	Avatar   string
	Address  string
	Zalo     string
	Note     string
}

type AdminOpportunityAssignDTO struct {
	ChargePersonID uint64
	AssignNote     string
}

type AdminOpportunityStatusDTO struct {
	OpportunityStatus enums.EOpportunityStatus
	Note              string
}

type AdminOpportunityFollowUpDTO struct {
	NextFollowUpAt time.Time
	Content        string
	NextAction     string
}

type AdminOpportunityCloseDTO struct {
	Result      enums.EOpportunityStatus
	CloseReason string
}

type AdminOpportunityEventResponse struct {
	ID         uint64     `json:"id"`
	ActorID    uint64     `json:"actorId"`
	ActorName  string     `json:"actorName,omitempty"`
	Action     string     `json:"action"`
	BeforeJSON string     `json:"beforeJson,omitempty"`
	AfterJSON  string     `json:"afterJson,omitempty"`
	Note       string     `json:"note,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
}

type AdminOpportunitiesSummaryDTO struct {
	Total               int64   `json:"total"`
	OverdueFollowUp     int64   `json:"overdueFollowUp"`
	Unassigned          int64   `json:"unassigned"`
	WaitingQuote        int64   `json:"waitingQuote"`
	WaitingApproval     int64   `json:"waitingApproval"`
	WaitingPayment      int64   `json:"waitingPayment"`
	Renewal             int64   `json:"renewal"`
	Upgrade             int64   `json:"upgrade"`
	ChurnRisk           int64   `json:"churnRisk"`
	StatusNew           int64   `json:"statusNew"`
	Consulting          int64   `json:"consulting"`
	Won                 int64   `json:"won"`
	Lost                int64   `json:"lost"`
	ExpectedRevenue     float64 `json:"expectedRevenue"`
	WaitingPaymentValue float64 `json:"waitingPaymentValue"`
	WonValue            float64 `json:"wonValue"`
}

type AdminFunnelBucketDTO struct {
	Status           enums.EOpportunityStatus `json:"status"`
	Label            string                   `json:"label"`
	Count            int64                    `json:"count"`
	ExpectedValueSum float64                  `json:"expectedValueSum"`
}
