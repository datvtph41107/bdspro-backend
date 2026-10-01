package dto

// CatalogIDPath is DELETE/GET path `{id}` for catalog admin routes.
type CatalogIDPath struct {
	ID uint64 `path:"id"`
}

// CatalogSearchRequest is GET list with optional text/page/size.
type CatalogSearchRequest struct {
	Text string `form:"text"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// CatalogPageRequest is GET list with page/size only.
type CatalogPageRequest struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// MainAreaWriteRequest is POST/PUT body (+ path id on update).
type MainAreaWriteRequest struct {
	ID          uint64  `path:"id" json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// PurposeUseWriteRequest is POST/PUT body (+ path id on update).
type PurposeUseWriteRequest struct {
	ID          uint64  `path:"id" json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// ProfessionWriteRequest is POST/PUT body (+ path id on update).
type ProfessionWriteRequest struct {
	ID             uint64  `path:"id" json:"id"`
	Name           string  `json:"name"`
	Issuer         string  `json:"issuer"`
	IssueDate      string  `json:"issueDate"`
	VerifiedStatus uint32  `json:"verifiedStatus"`
	ApproveUserID  *uint64 `json:"approveUserId"`
	IsActive       *bool   `json:"isActive"`
}
