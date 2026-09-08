package enums

type ESourceStatus int

const (
	SourceStatusNew      ESourceStatus = 10
	SourceStatusTracking ESourceStatus = 20
	SourceStatusOwn      ESourceStatus = 30
)

var SourceStatusNames = map[ESourceStatus]string{
	SourceStatusNew:      "Mới",
	SourceStatusTracking: "Đang chăm",
	SourceStatusOwn:      "Đang sở hữu",
}
