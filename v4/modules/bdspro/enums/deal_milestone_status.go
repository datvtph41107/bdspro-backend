package enums

type MilestoneStatus uint

const (
	MilestoneStatusPending   MilestoneStatus = 10
	MilestoneStatusCompleted MilestoneStatus = 20
	MilestoneStatusOverdue   MilestoneStatus = 30
)

var MilestoneStatusMap = map[MilestoneStatus]string{
	MilestoneStatusPending:   "Đang chờ",
	MilestoneStatusCompleted: "Hoàn thành",
	MilestoneStatusOverdue:   "Quá hạn",
}
