package dto

// AppointmentListRequest mirrors GET /v2.1/appointment query params for GET /v4/appointment/me.
type AppointmentListRequest struct {
	Title              string  `form:"title"`
	Text               string  `form:"text"`
	Status             uint32  `form:"status"`
	DealIDs            string  `form:"dealIds"` // "1,2,3"
	TransactionID      *uint64 `form:"transactionId"`
	Page               uint32  `form:"page"`
	Size               uint32  `form:"size"`
	Sort               string  `form:"sort"` // "startTime,desc" | "startTime,asc"
	ParticipantIDs     string  `form:"participantIds"`
	FromDate           string  `form:"fromDate"` // YYYY-MM-DD
	ToDate             string  `form:"toDate"`
	ContainProduct     *bool   `form:"containProduct"`
	ContainDeal        *bool   `form:"containDeal"`
	ContainTransaction *bool   `form:"containTransaction"`
	StatusFilters      string  `form:"statusFilters"` // "10,20"
	Timestamp          int64   `form:"timestamp"`
}
