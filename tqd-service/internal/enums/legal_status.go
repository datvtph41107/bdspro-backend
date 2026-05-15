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
	LegalStatusReplaced:   "Bị thay thế", //"Không hiện hành",
}

func (s LegalStatus) IsValid() bool {
	return s == LegalStatusEffective || s == LegalStatusApproved || s == LegalStatusDraft || s == LegalStatusReferenced || s == LegalStatusUnknown || s == LegalStatusExpired || s == LegalStatusReplaced
}

func (s LegalStatus) AvaiableToUpdate() bool {
	return s == LegalStatusDraft || s == LegalStatusReferenced || s == LegalStatusUnknown
}

func (s LegalStatus) String() string {
	switch s {
	case LegalStatusEffective:
		return "effective"
	case LegalStatusApproved:
		return "approved"
	case LegalStatusDraft:
		return "draft"
	case LegalStatusReferenced:
		return "reference"
	case LegalStatusExpired, LegalStatusReplaced:
		return "expired"
	default:
		return "unknown"
	}
}
func (s LegalStatus) DisplayName() string {
	switch s {
	case LegalStatusEffective:
		return "Đang hiệu lực"
	case LegalStatusApproved:
		return "Đã duyệt"
	case LegalStatusDraft:
		return "Dự thảo"
	case LegalStatusReferenced:
		return "Tham chiếu"
	case LegalStatusExpired:
		return "Hết hiệu lực"
	case LegalStatusReplaced:
		return "Đã thay thế"
	default:
		return "Không xác định"
	}
}

func (s LegalStatus) GetPriority() int {
	switch s {
	case LegalStatusEffective:
		return 100
	case LegalStatusApproved:
		return 70
	case LegalStatusDraft:
		return 30
	case LegalStatusReferenced:
		return 15
	default:
		return 10
	}
}
