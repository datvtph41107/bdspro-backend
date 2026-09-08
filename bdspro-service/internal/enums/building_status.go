package enums

type EBuildStatus uint32

const (
	EBuildStatusEmpty     EBuildStatus = 10
	EBuildStatusBuilding  EBuildStatus = 20
	EBuildStatusCompleted EBuildStatus = 30
)

var BuildingStatusNames = map[EBuildStatus]string{
	EBuildStatusEmpty:     "Đất trống",
	EBuildStatusBuilding:  "Đang xây",
	EBuildStatusCompleted: "Hoàn thiện",
}
