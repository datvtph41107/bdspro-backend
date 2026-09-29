package enums

// PushConfigKey định nghĩa các khóa cấu hình thông báo cá nhân hóa
type PushConfigKey string

const (
	PushConfigNotifyLogin     PushConfigKey = "notifyLogin"
	PushConfigNotifySecurity  PushConfigKey = "notifySecurity"
	PushConfigNotifyMessage   PushConfigKey = "notifyMessage"
	PushConfigNotifyReview    PushConfigKey = "notifyReview"
	PushConfigNotifySystem    PushConfigKey = "notifySystem"
	PushConfigNotifyMarketing PushConfigKey = "notifyMarketing"
	PushConfigMuteZalo        PushConfigKey = "muteZaloMarketing"
)
