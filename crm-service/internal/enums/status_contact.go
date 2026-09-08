package enums

type EStatusContact uint32

const (
	StatusContactUnknown      EStatusContact = 10
	StatusContactFriend       EStatusContact = 20
	StatusContactAppUsed      EStatusContact = 30
	StatusContactAppInstalled EStatusContact = 40
	StatusContactNotSaved     EStatusContact = 50
)

var StatusContactMap = map[EStatusContact]string{
	StatusContactUnknown:      "Chưa dùng app",
	StatusContactFriend:       "Bạn bè",
	StatusContactAppUsed:      "Đã dùng app",
	StatusContactAppInstalled: "Đã cài app",
	StatusContactNotSaved:     "Chưa lưu",
}