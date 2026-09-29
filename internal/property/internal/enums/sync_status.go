package enums

type ESyncStatus uint32

const (
	ESyncStatusPending    ESyncStatus = 10
	ESyncStatusProcessing ESyncStatus = 20
	ESyncStatusDone       ESyncStatus = 30
	ESyncStatusFailed     ESyncStatus = 40
)

var ESyncStatusNames = map[ESyncStatus]string{
	ESyncStatusPending:    "Chờ xử lý",
	ESyncStatusProcessing: "Đang xử lý",
	ESyncStatusDone:       "Hoàn thành",
	ESyncStatusFailed:     "Thất bại",
}

func (e ESyncStatus) String() string {
	return ESyncStatusNames[e]
}

func (e ESyncStatus) IsValid() bool {
	return e == ESyncStatusPending || e == ESyncStatusProcessing || e == ESyncStatusDone || e == ESyncStatusFailed
}
