package enums

type EPropertyStatus uint32

const (
	EPropertyStatusActive   EPropertyStatus = 10 // Hoạt động
	EPropertyStatusInactive EPropertyStatus = 20 // Không hoạt động
	EPropertyStatusArchived EPropertyStatus = 30 // Đã lưu trữ
)

var EPropertyStatusNames = map[EPropertyStatus]string{
	EPropertyStatusActive:   "Hoạt động",
	EPropertyStatusInactive: "Không hoạt động",
	EPropertyStatusArchived: "Đã lưu trữ",
}

func (e EPropertyStatus) String() string {
	return EPropertyStatusNames[e]
}

func (e EPropertyStatus) IsValid() bool {
	return e == EPropertyStatusActive || e == EPropertyStatusInactive || e == EPropertyStatusArchived
}
