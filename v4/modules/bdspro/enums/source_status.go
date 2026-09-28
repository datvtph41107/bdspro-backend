package enums

type ESourceStatus int

const (
	SourceStatusNew      ESourceStatus = 10
	SourceStatusTracking ESourceStatus = 20
	SourceStatusOwn      ESourceStatus = 30
)

var SourceStatusNames = map[ESourceStatus]string{
	SourceStatusNew:      "Chờ xác minh",
	SourceStatusTracking: "Thiết giấy tờ",
	SourceStatusOwn:      "Đã xác minh số",
}
