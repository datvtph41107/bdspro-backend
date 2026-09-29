package enums

// ContactLabelStatus represents the status of a contact label
type ContactLabelStatus int

const (
	ContactLabelStatusActive   ContactLabelStatus = 0
	ContactLabelStatusInactive ContactLabelStatus = 1
	ContactLabelStatusArchived ContactLabelStatus = 2
	ContactLabelStatusDeleted  ContactLabelStatus = 3
)
