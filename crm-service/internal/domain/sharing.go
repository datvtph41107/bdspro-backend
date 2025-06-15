package domain

import "crm/internal/enums"

type SharingEntity struct {
	// BaseEntity
	ContactID    uint64           `json:"contactId"`
	ReceiverID   uint64           `json:"receiverId"`
	ReceiverType enums.EOwnerType `json:"receiverType"`
	Permissions  uint32           `json:"permissions"`
	Note         string           `json:"note"`
}
