package _enum

type EChannelNotification uint32

const (
	EChannelNotiEmail   EChannelNotification = 10
	EChannelNotiZNS_SMS EChannelNotification = 20
	EChannelNotiApp     EChannelNotification = 30
)

func (e EChannelNotification) IsValid() bool {
	switch e {
	case EChannelNotiEmail, EChannelNotiZNS_SMS, EChannelNotiApp:
		return true
	default:
		return false
	}
}
