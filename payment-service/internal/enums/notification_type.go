package enums

type NotificationType int32

const (
	NotificationTypeDeposit  NotificationType = 5001
	NotificationTypeWithdraw NotificationType = 5002
	NotificationTypePayment  NotificationType = 5003
)

var NotificationTypeMap = map[NotificationType]string{
	NotificationTypeDeposit:  "Nạp tiền",
	NotificationTypeWithdraw: "Rút tiền",
	NotificationTypePayment:  "Thanh toán",
}
