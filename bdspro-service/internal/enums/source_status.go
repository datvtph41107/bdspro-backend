package enums

type ESourceStatus int

const (
	SourceStatusNew ESourceStatus = 10
)

var SourceStatusMap = map[ESourceStatus]string{
	SourceStatusNew: "mới",
}
