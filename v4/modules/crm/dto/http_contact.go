package dto

// ContactListMeRequest is GET /v2/crm/contact/list/me.
type ContactListMeRequest struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Text      string `form:"text"`
	Timestamp int64  `form:"timestamp"`
}

// ContactSearchHTTPRequest is GET /v2/crm/contact/list/{ownerOf}/{ownerId}.
type ContactSearchHTTPRequest struct {
	OwnerOf      string  `path:"ownerOf"`
	OwnerID      uint64  `path:"ownerId"`
	Page         uint32  `form:"page"`
	Size         uint32  `form:"size"`
	Text         string  `form:"text"`
	FriendStatus string  `form:"friendStatus"`
	Following    *bool   `form:"following"`
	Roles        []int32 `form:"roles"`
	Leaded       *bool   `form:"leaded"`
	Connected    *bool   `form:"connected"`
	AppInstalled *bool   `form:"appInstalled"`
	Timestamp    int64   `form:"timestamp"`
}

// ContactDetailHTTPRequest is GET /v2/crm/contact/{id}/detail.
type ContactDetailHTTPRequest struct {
	ID        uint64 `path:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

// ContactCreateHTTPRequest is POST /v2/crm/contact (no path id — matches ContactSaveDTO).
type ContactCreateHTTPRequest struct {
	Phone       string   `json:"phone"`
	FullName    string   `json:"fullName"`
	OwnerOf     int32    `json:"ownerOf"`
	OwnerID     uint64   `json:"ownerId"`
	Email       string   `json:"email"`
	Zalo        string   `json:"zalo"`
	Company     string   `json:"company"`
	Note        string   `json:"note"`
	Avatar      string   `json:"avatar"`
	Visibility  int32    `json:"visibility"`
	AssignedCrm bool     `json:"assignedCrm"`
	ProductIds  []uint64 `json:"productIds"`
}

// ContactSaveHTTPRequest is PUT /v2/crm/contact/{id}/edit and sync payload rows.
// @delivery
type ContactSaveHTTPRequest struct {
	ID          uint64   `path:"id" json:"id"`
	Phone       string   `json:"phone"`
	FullName    string   `json:"fullName"`
	OwnerOf     int32    `json:"ownerOf"`
	OwnerID     uint64   `json:"ownerId"`
	Email       string   `json:"email"`
	Zalo        string   `json:"zalo"`
	Company     string   `json:"company"`
	Note        string   `json:"note"`
	Avatar      string   `json:"avatar"`
	Visibility  int32    `json:"visibility"`
	AssignedCrm bool     `json:"assignedCrm"`
	ProductIds  []uint64 `json:"productIds"`
}

// ContactSyncHTTPRequest is POST /v2/crm/contact/sync.
type ContactSyncHTTPRequest struct {
	Contacts  []ContactSaveHTTPRequest `json:"contacts"`
	OwnerID   uint64                   `json:"ownerId"`
	OwnerType int32                    `json:"ownerType"`
}

// ContactNoteHTTPRequest is PUT /v2/crm/contact/{id}/note.
type ContactNoteHTTPRequest struct {
	ID   uint64 `path:"id" json:"id"`
	Note string `json:"note"`
}

// ContactCheckPhoneHTTPRequest is GET /v2/crm/contact/phone/{ownerOf}/{ownerId}/{phone}.
type ContactCheckPhoneHTTPRequest struct {
	OwnerOf string `path:"ownerOf"`
	OwnerID uint64 `path:"ownerId"`
	Phone   string `path:"phone"`
}

// ContactByProductHTTPRequest is GET /v2/crm/contact/{id}/by-product ({id} is the product id).
type ContactByProductHTTPRequest struct {
	ID        uint64 `path:"id"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	Timestamp int64  `form:"timestamp"`
}

// ContactUpdateProductsHTTPRequest is PUT /v2/crm/contact/{id}/products.
type ContactUpdateProductsHTTPRequest struct {
	ContactID  uint64   `path:"id" json:"contactId"`
	ProductIDs []uint64 `json:"productIds"`
}

// ContactProductsHTTPRequest is GET /v2/crm/contact/{id}/products.
type ContactProductsHTTPRequest struct {
	ID   uint64 `path:"id"`
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
}

// ContactWithMemberHTTPRequest is GET /v2/crm/contact/with-member and /friends.
type ContactWithMemberHTTPRequest struct {
	Page           uint32  `form:"page"`
	Size           uint32  `form:"size"`
	Text           string  `form:"text"`
	FriendStatus   string  `form:"friendStatus"`
	DealID         *uint64 `form:"dealId"`
	OrganizationID *uint64 `form:"organizationId"`
	GroupID        *uint64 `form:"groupId"`
	ShareProductID *uint64 `form:"shareProductId"`
}

// ContactInterestedHTTPRequest is GET /v2/crm/contact/interested/product/{productId}.
type ContactInterestedHTTPRequest struct {
	ProductID uint64 `path:"productId"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
}

// ContactPinHTTPRequest is POST /v2/crm/contact/pin.
type ContactPinHTTPRequest struct {
	Target    uint32            `json:"target"`
	TargetIDs map[string]uint64 `json:"targetIds"`
	Namespace string            `json:"namespace"`
	Action    string            `json:"action"`
}
