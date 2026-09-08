package enums

type NotificationType uint32

const (
	NotifTypeWatchAlert  NotificationType = 10
	NotifTypeReportReady NotificationType = 20
	NotifTypeSystem      NotificationType = 30
)

type NotificationSeverity uint32

const (
	NotifSeverityLow      NotificationSeverity = 10
	NotifSeverityMedium   NotificationSeverity = 20
	NotifSeverityHigh     NotificationSeverity = 30
	NotifSeverityCritical NotificationSeverity = 40
)
