package domain

import "strings"

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

// planningNewsVisibleOnByFeed hỗ trợ cả canonical code của domain
// và alias mà các frontend hiện đang sử dụng.
//
// Alias không tạo thêm giá trị visible_on mới.
var planningNewsVisibleOnByFeed = map[string]PlanningNewsVisibleOn{
	// visible_on = 10
	"top":      PlanningNewsVisibleOnTop,
	"pinned":   PlanningNewsVisibleOnTop,
	"featured": PlanningNewsVisibleOnTop,

	// visible_on = 20
	"notable":   PlanningNewsVisibleOnNotable,
	"attention": PlanningNewsVisibleOnNotable,

	// visible_on = 30
	"latest": PlanningNewsVisibleOnLatest,

	// visible_on = 40
	"most_read": PlanningNewsVisibleOnMostRead,

	// Compatibility alias của frontend hiện tại.
	// Không tạo thêm enum ngoài contract 10/20/30/40/50.
	"trending": PlanningNewsVisibleOnMostRead,

	// visible_on = 50
	"quick_news": PlanningNewsVisibleOnQuickNews,
	"breaking":   PlanningNewsVisibleOnQuickNews,
}

func normalizePlanningNewsCode(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
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

func (value PlanningNewsVisibleOn) Uint() uint {
	return uint(value)
}

func (value PlanningNewsVisibleOn) Uint32() uint32 {
	return uint32(value)
}

func (value PlanningNewsVisibleOn) String() string {
	if code, ok := PlanningNewsVisibleOnCodeByValue[value]; ok {
		return code
	}

	return ""
}

func (value PlanningNewsVisibleOn) Label() string {
	if label, ok := PlanningNewsVisibleOnLabelByValue[value]; ok {
		return label
	}

	return ""
}

// ParsePlanningNewsVisibleOn phân tích canonical domain code:
// top, notable, latest, most_read, quick_news.
func ParsePlanningNewsVisibleOn(
	code string,
) (PlanningNewsVisibleOn, bool) {
	value, ok := PlanningNewsVisibleOnValueByCode[normalizePlanningNewsCode(code)]

	return value, ok
}

// ResolvePlanningNewsFeed phân tích tên feed public.
//
// Hàm này là nguồn sự thật duy nhất cho mapping:
//
//	latest -> visible_on 30
//	attention -> visible_on 20
//	most_read -> visible_on 40
//
// Repository và application không tự hard-code lại mapping.
func ResolvePlanningNewsFeed(
	feed string,
) (PlanningNewsVisibleOn, bool) {
	value, ok := planningNewsVisibleOnByFeed[normalizePlanningNewsCode(feed)]

	return value, ok
}

// CanonicalPlanningNewsFeed trả canonical code của enum.
//
// Ví dụ:
//
//	featured -> top
//	attention -> notable
//	breaking -> quick_news
func CanonicalPlanningNewsFeed(feed string) (string, bool) {
	value, ok := ResolvePlanningNewsFeed(feed)
	if !ok {
		return "", false
	}

	return value.String(), true
}
