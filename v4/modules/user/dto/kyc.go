package dto

// ListKYCsRequest is GET /v2/user/admin/kycs.
type ListKYCsRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Search    string `form:"search"`
	Status    uint32 `form:"status"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

// ApproveKYCRequest is POST /v2/user/admin/kyc/approve/{id}.
type ApproveKYCRequest struct {
	ID uint64 `path:"id" json:"id"`
}

// RejectKYCRequest is POST /v2/user/admin/kyc/reject/{id}.
type RejectKYCRequest struct {
	ID     uint64 `path:"id" json:"id"`
	Reason string `json:"reason"`
}
