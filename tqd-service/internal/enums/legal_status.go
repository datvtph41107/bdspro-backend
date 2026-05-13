package enums

// LegalStatus represents the legal status of a layer
type LegalStatus uint32

const (
	LegalStatusEffective  LegalStatus = 1000 // Effective legal status
	LegalStatusApproved   LegalStatus = 900  // Valid legal status
	LegalStatusDraft      LegalStatus = 800  // Draft legal status
	LegalStatusReferenced LegalStatus = 700  // Referenced legal status
	LegalStatusUnknown    LegalStatus = 600  // Default legal status
	LegalStatusExpired    LegalStatus = 400  // Expired legal status
	LegalStatusReplaced   LegalStatus = 300  // Replaced legal status
)

var LegalStatusMap = map[LegalStatus]string{
	LegalStatusEffective:  "Hiệu lực",
	LegalStatusApproved:   "Đã duyệt",
	LegalStatusDraft:      "Dự thảo",
	LegalStatusReferenced: "Tham chiếu",
	LegalStatusUnknown:    "Không xác định",
	LegalStatusExpired:    "Hết hiệu lực",
	LegalStatusReplaced:   "Không hiện hành",
}
