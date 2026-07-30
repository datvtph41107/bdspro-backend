package domain

import "strings"

// PlanningNewsVisibleOn là taxonomy nghiệp vụ của vị trí hiển thị Tin quy hoạch.
// Proto và database dùng số; domain sở hữu tập giá trị hợp lệ và mã chuỗi.
type PlanningNewsVisibleOn uint

const (
	PlanningNewsVisibleOnTop       PlanningNewsVisibleOn = 10
	PlanningNewsVisibleOnNotable   PlanningNewsVisibleOn = 20
	PlanningNewsVisibleOnLatest    PlanningNewsVisibleOn = 30
	PlanningNewsVisibleOnMostRead  PlanningNewsVisibleOn = 40
	PlanningNewsVisibleOnQuickNews PlanningNewsVisibleOn = 50
)

var PlanningNewsVisibleOnCodeByValue = map[PlanningNewsVisibleOn]string{
	PlanningNewsVisibleOnTop:       "top",
	PlanningNewsVisibleOnNotable:   "notable",
	PlanningNewsVisibleOnLatest:    "latest",
	PlanningNewsVisibleOnMostRead:  "most_read",
	PlanningNewsVisibleOnQuickNews: "quick_news",
}

var PlanningNewsVisibleOnValueByCode = map[string]PlanningNewsVisibleOn{
	"top":        PlanningNewsVisibleOnTop,
	"notable":    PlanningNewsVisibleOnNotable,
	"latest":     PlanningNewsVisibleOnLatest,
	"most_read":  PlanningNewsVisibleOnMostRead,
	"quick_news": PlanningNewsVisibleOnQuickNews,
}

var PlanningNewsVisibleOnLabelByValue = map[PlanningNewsVisibleOn]string{
	PlanningNewsVisibleOnTop:       "Hiển thị lên đầu",
	PlanningNewsVisibleOnNotable:   "Tin đáng chú ý",
	PlanningNewsVisibleOnLatest:    "Tin mới nhất",
	PlanningNewsVisibleOnMostRead:  "Tin đọc nhiều",
	PlanningNewsVisibleOnQuickNews: "Tin nhanh",
}

func (value PlanningNewsVisibleOn) IsValid() bool {
	switch value {
	case PlanningNewsVisibleOnTop,
		PlanningNewsVisibleOnNotable,
		PlanningNewsVisibleOnLatest,
		PlanningNewsVisibleOnMostRead,
		PlanningNewsVisibleOnQuickNews:
		return true
	default:
		return false
	}
}

func (value PlanningNewsVisibleOn) Uint() uint     { return uint(value) }
func (value PlanningNewsVisibleOn) Uint32() uint32 { return uint32(value) }
func (value PlanningNewsVisibleOn) String() string { return PlanningNewsVisibleOnCodeByValue[value] }
func (value PlanningNewsVisibleOn) Label() string  { return PlanningNewsVisibleOnLabelByValue[value] }

func ParsePlanningNewsVisibleOn(code string) (PlanningNewsVisibleOn, bool) {
	value, ok := PlanningNewsVisibleOnValueByCode[strings.ToLower(strings.TrimSpace(code))]
	return value, ok
}
