package enums

type EOnlineStatus int

const (
	Online  EOnlineStatus = 10
	LASTEST EOnlineStatus = 20
	HIDDEN  EOnlineStatus = 30
	// Offline2 OnlineStatus = 40
)

var OnlineStatusMap = map[EOnlineStatus]string{
	Online:  "Online",
	LASTEST: "Gần đây",
	HIDDEN:  "Không hiển thị",
}
