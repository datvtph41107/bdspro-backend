package _utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAccents loại bỏ dấu tiếng Việt và các ký tự đặc biệt
func RemoveAccents(s string) string {
	// Chuẩn hóa chuỗi sang dạng NFD (Decomposition)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	// Thay thế các ký tự đặc biệt của tiếng Việt không được xử lý bởi NFD
	result = strings.ReplaceAll(result, "đ", "d")
	result = strings.ReplaceAll(result, "Đ", "D")

	return result
}

// CleanSearchText loại bỏ dấu và chỉ giữ lại chữ cái, số và khoảng trắng
func CleanSearchText(s string) string {
	// Loại bỏ dấu
	s = RemoveAccents(s)

	// Chỉ giữ lại chữ cái, số và khoảng trắng
	reg := regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	s = reg.ReplaceAllString(s, "")

	// Loại bỏ khoảng trắng thừa
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// ParcelSearchInfo chứa thông tin search thửa đất đã được bóc tách
type ParcelSearchInfo struct {
	CleanedText string
	MapNumber   string
	LandNumber  string
}

// InferParcelSearch bóc tách số tờ, số thửa từ chuỗi tìm kiếm
func InferParcelSearch(s string) ParcelSearchInfo {
	info := ParcelSearchInfo{}

	// Chuẩn hóa chuỗi trước khi parse (loại bỏ dấu để dễ bắt pattern)
	normalized := strings.ToLower(RemoveAccents(s))

	// Pattern: tờ [số] thửa [số] hoặc to [số] thua [số]
	// Hỗ trợ các biến thể: "tờ 12 thửa 34", "to 12 thua 34", "tờ:12 thửa:34", "to12 thua34"
	mapReg := regexp.MustCompile(`to\s*[:\s]*(\d+)`)
	landReg := regexp.MustCompile(`thua\s*[:\s]*(\d+)`)

	mapMatch := mapReg.FindStringSubmatch(normalized)
	if len(mapMatch) > 1 {
		info.MapNumber = mapMatch[1]
	}

	landMatch := landReg.FindStringSubmatch(normalized)
	if len(landMatch) > 1 {
		info.LandNumber = landMatch[1]
	}

	// Clean text để search address như bình thường
	info.CleanedText = CleanSearchText(s)

	return info
}
