package enums

type EOnlineStatus int

const (
	Online  EOnlineStatus = 10
	Offline EOnlineStatus = 20
	Busy    EOnlineStatus = 30
	// Offline2 OnlineStatus = 40
)

var OnlineStatusMap = map[EOnlineStatus]string{
	Online:  "Online",
	Offline: "Nghỉ",
	Busy:    "Bận",
}
