package enums

type EPropertyStatus uint32

const (
	EPropertyStatusActive   EPropertyStatus = 10 // Hoạt động
	EPropertyStatusHide     EPropertyStatus = 20 // Ẩn
	EPropertyStatusArchived EPropertyStatus = 30 // Đã lưu trữ
)

var EPropertyStatusNames = map[EPropertyStatus]string{
	EPropertyStatusActive:   "Hoạt động",
	EPropertyStatusHide:     "Tạm ẩn",
	EPropertyStatusArchived: "Lưu trữ",
}

func (e EPropertyStatus) String() string {
	return EPropertyStatusNames[e]
}

func (e EPropertyStatus) IsValid() bool {
	return e == EPropertyStatusActive || e == EPropertyStatusHide || e == EPropertyStatusArchived
}
