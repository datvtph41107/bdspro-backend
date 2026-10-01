package dto

// AdminSalesListOpportunitiesRequest is GET /v2/crm/admin/sales/opportunities.
type AdminSalesListOpportunitiesRequest struct {
	Page                int32    `form:"page"`
	Size                int32    `form:"size"`
	Q                   string   `form:"q"`
	StageIDs            []uint64 `form:"stageIds"`
	ChargePersonIDs     []uint64 `form:"chargePersonIds"`
	PipelineID          *uint64  `form:"pipelineId"`
	Sources             []int32  `form:"sources"`
	UnassignedOnly      *bool    `form:"unassignedOnly"`
	OpportunityStatuses []int32  `form:"opportunityStatuses"`
	CustomerTypes       []int32  `form:"customerTypes"`
	AdminSources        []int32  `form:"adminSources"`
	Queue               string   `form:"queue"`
	FollowUpFrom        *string  `form:"followUpFrom"`
	FollowUpTo          *string  `form:"followUpTo"`
	MinExpectedValue    *float64 `form:"minExpectedValue"`
	MaxExpectedValue    *float64 `form:"maxExpectedValue"`
	Probabilities       []int32  `form:"probabilities"`
	ChurnRisks          []int32  `form:"churnRisks"`
	UpgradeOnly         *bool    `form:"upgradeOnly"`
	RenewalOnly         *bool    `form:"renewalOnly"`
	InterestedPlan      string   `form:"interestedPlan"`
	ViewMode            string   `form:"viewMode"`
}

// AdminSalesListContactsRequest is GET /v2/crm/admin/sales/contacts.
type AdminSalesListContactsRequest struct {
	Page int32  `form:"page"`
	Size int32  `form:"size"`
	Q    string `form:"q"`
}

// AdminSalesCreateContactRequest is POST /v2/crm/admin/sales/contacts.
type AdminSalesCreateContactRequest struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Company  string `json:"company"`
	Avatar   string `json:"avatar"`
	Address  string `json:"address"`
	Zalo     string `json:"zalo"`
	Note     string `json:"note"`
}

// AdminSalesCreateOpportunityRequest is POST /v2/crm/admin/sales/opportunities.
type AdminSalesCreateOpportunityRequest struct {
	FullName          string   `json:"fullName"`
	Phone             string   `json:"phone"`
	Email             string   `json:"email"`
	Company           string   `json:"company"`
	Avatar            string   `json:"avatar"`
	Address           string   `json:"address"`
	Zalo              string   `json:"zalo"`
	Note              string   `json:"note"`
	Source            int32    `json:"source"`
	StageID           *uint64  `json:"stageId"`
	ChargePersonID    *uint64  `json:"chargePersonId"`
	Priority          int32    `json:"priority"`
	AssignNote        string   `json:"assignNote"`
	Title             string   `json:"title"`
	CustomerType      int32    `json:"customerType"`
	NeedSummary       string   `json:"needSummary"`
	InterestedPlan    string   `json:"interestedPlan"`
	OpportunityStatus int32    `json:"opportunityStatus"`
	AdminSource       int32    `json:"adminSource"`
	ExpectedValue     *float64 `json:"expectedValue"`
	Probability       int32    `json:"probability"`
	ExpectedCloseDate *string  `json:"expectedCloseDate"`
	NextFollowUpAt    *string  `json:"nextFollowUpAt"`
	ChurnRisk         int32    `json:"churnRisk"`
	UpgradeSignal     bool     `json:"upgradeSignal"`
	RenewalSignal     bool     `json:"renewalSignal"`
	OwnerTeam         string   `json:"ownerTeam"`
	Segment           string   `json:"segment"`
	Region            string   `json:"region"`
	Tags              string   `json:"tags"`
	ContactID         *uint64  `json:"contactId"`
}

// AdminSalesUpdateOpportunityRequest is PATCH /v2/crm/admin/sales/opportunities/{id}.
type AdminSalesUpdateOpportunityRequest struct {
	ID                uint64   `path:"id" json:"id"`
	FullName          string   `json:"fullName"`
	Phone             string   `json:"phone"`
	Email             string   `json:"email"`
	Company           string   `json:"company"`
	Avatar            string   `json:"avatar"`
	Address           string   `json:"address"`
	Zalo              string   `json:"zalo"`
	Note              string   `json:"note"`
	Source            int32    `json:"source"`
	StageID           *uint64  `json:"stageId"`
	ChargePersonID    *uint64  `json:"chargePersonId"`
	Priority          int32    `json:"priority"`
	AssignNote        string   `json:"assignNote"`
	Title             string   `json:"title"`
	CustomerType      int32    `json:"customerType"`
	NeedSummary       string   `json:"needSummary"`
	InterestedPlan    string   `json:"interestedPlan"`
	OpportunityStatus int32    `json:"opportunityStatus"`
	AdminSource       int32    `json:"adminSource"`
	ExpectedValue     *float64 `json:"expectedValue"`
	Probability       int32    `json:"probability"`
	ExpectedCloseDate *string  `json:"expectedCloseDate"`
	NextFollowUpAt    *string  `json:"nextFollowUpAt"`
	ChurnRisk         int32    `json:"churnRisk"`
	UpgradeSignal     bool     `json:"upgradeSignal"`
	RenewalSignal     bool     `json:"renewalSignal"`
	OwnerTeam         string   `json:"ownerTeam"`
	Segment           string   `json:"segment"`
	Region            string   `json:"region"`
	Tags              string   `json:"tags"`
	ProposalRef       string   `json:"proposalRef"`
	PaymentRequestRef string   `json:"paymentRequestRef"`
	SubscriptionRef   string   `json:"subscriptionRef"`
	TicketRef         string   `json:"ticketRef"`
}

// AdminSalesAssignRequest is POST /v2/crm/admin/sales/opportunities/{id}/assign.
type AdminSalesAssignRequest struct {
	ID             uint64 `path:"id" json:"id"`
	ChargePersonID uint64 `json:"chargePersonId"`
	AssignNote     string `json:"assignNote"`
}

// AdminSalesSwitchStageRequest is POST /v2/crm/admin/sales/opportunities/{id}/stage.
type AdminSalesSwitchStageRequest struct {
	ID        uint64 `path:"id" json:"id"`
	StageID   uint64 `json:"stageId"`
	StageNote string `json:"stageNote"`
}

// AdminSalesUpdateNoteRequest is POST /v2/crm/admin/sales/opportunities/{id}/note.
type AdminSalesUpdateNoteRequest struct {
	ID   uint64 `path:"id" json:"id"`
	Note string `json:"note"`
}

// AdminSalesUpdateStatusRequest is PATCH /v2/crm/admin/sales/opportunities/{id}/status.
type AdminSalesUpdateStatusRequest struct {
	ID                uint64 `path:"id" json:"id"`
	OpportunityStatus int32  `json:"opportunityStatus"`
	Note              string `json:"note"`
}

// AdminSalesFollowUpRequest is POST /v2/crm/admin/sales/opportunities/{id}/follow-ups.
type AdminSalesFollowUpRequest struct {
	ID             uint64 `path:"id" json:"id"`
	NextFollowUpAt string `json:"nextFollowUpAt"`
	Content        string `json:"content"`
	NextAction     string `json:"nextAction"`
}

// AdminSalesCloseRequest is POST /v2/crm/admin/sales/opportunities/{id}/close.
type AdminSalesCloseRequest struct {
	ID          uint64 `path:"id" json:"id"`
	Result      int32  `json:"result"`
	CloseReason string `json:"closeReason"`
}

// AdminSalesExportRequest is POST /v2/crm/admin/sales/exports.
type AdminSalesExportRequest struct {
	Reason              string   `json:"reason"`
	Q                   string   `json:"q"`
	Queue               string   `json:"queue"`
	OpportunityStatuses []int32  `json:"opportunityStatuses"`
	ChargePersonIDs     []uint64 `json:"chargePersonIds"`
}
