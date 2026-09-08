package utils

import (
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	_utils "common/utils"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var Rules = []_utils.Rule{
	{
		Field:   "bedroom",
		Pattern: regexp.MustCompile(`((\d+|một|hai|ba|bốn|năm|sáu|bảy|tám|chín|mười)[\s-]*)(phòng\s*ngủ|ngủ|chỗ\s*ngủ)`),
		GetValue: func(match []string) interface{} {
			return _utils.ParseNumber(match[1])
		},
	},
	{
		Field:   "bathroom",
		Pattern: regexp.MustCompile(`((\d+|một|hai|ba|bốn|năm|sáu|bảy|tám|chín|mười)[\s-]*)(phòng\s*tắm|tắm|chỗ\s*tắm)`),
		GetValue: func(match []string) interface{} {
			return _utils.ParseNumber(match[1])
		},
	},
	{
		Field:   "toilet",
		Pattern: regexp.MustCompile(`((\d+|một|hai|ba|bốn|năm|sáu|bảy|tám|chín|mười)?[\s-]*)(vệ\s*sinh|nhà\s*vệ\s*sinh)`),
		GetValue: func(match []string) interface{} {
			if n := _utils.ParseNumber(match[1]); n != nil {
				return *n
			}
			return 0
		},
	},
	{
		Field:   "kitchen",
		Pattern: regexp.MustCompile(`((\d+|một|hai|ba|bốn|năm|sáu|bảy|tám|chín|mười)?[\s-]*)(nhà\s*bếp|bếp|chỗ\s*nấu)`),
		GetValue: func(match []string) interface{} {
			if n := _utils.ParseNumber(match[1]); n != nil {
				return *n
			}
			return 0
		},
	},
	{
		Field:   "park",
		Pattern: regexp.MustCompile(`(?i)(?:(\d+|một|hai|ba|bốn|năm|sáu|bảy|tám|chín|mười)\s*)?chỗ\s*để\s*(?:xe|ôtô|ô tô|xe\s*máy)`),
		GetValue: func(match []string) interface{} {
			if n := _utils.ParseNumber(match[1]); n != nil {
				return *n
			}
			return nil
		},
	},
	{
		Field: "priceSuggest",
		// Pattern: regexp.MustCompile(`(?i)(\d{1,3}(?:[\.,]?\d{3})+|\d+(?:[\.,]\d+)?|[\p{L}\s]+)\s*tỷ`),
		Pattern: regexp.MustCompile(`(?i)(\d+(?:[\.,]\d+)?)\s*tỷ`),

		// Pattern: regexp.MustCompile(`(?i)(\d+(?:[\.,]\d+)?|[\p{L}\s]+)\s*tỷ`),
		GetValue: func(match []string) interface{} {
			// val, _ := strconv.ParseFloat(match[1], 64)
			// val, _ := _utils.ExtractPriceInSentence(match[1])
			s := strings.ReplaceAll(match[1], ",", ".")
			val, _ := strconv.ParseFloat(s, 64)
			return int64(math.Round(val * 1_000_000_000))
		},
	},
	{
		Field:   "area",
		Pattern: regexp.MustCompile(`([\d\.]+)\s*(m2|m²|mét vuông|m\s*vuông|m\b)`),
		GetValue: func(match []string) interface{} {
			// val, _ := _utils.ExtractPriceInSentence(match[1])
			val, _ := strconv.ParseFloat(match[1], 64)
			return val
		},
	},
	{
		Field:   "frontage",
		Pattern: regexp.MustCompile(`mặt\s*tiền\s*([\d\.]+)`),
		GetValue: func(match []string) interface{} {
			val, _ := strconv.ParseFloat(match[1], 64)
			return val
		},
	},
	{
		Field:   "address",
		Pattern: regexp.MustCompile(`(?i)(?:quận|gần|khu\s*vực|khu|chỗ|phía|đằng|ở|tại|thuộc)?\s*([\pL\s\d]{2,})`),
		GetValue: func(match []string) interface{} {
			return strings.ToLower(strings.TrimSpace(match[1]))
		},
	},
	{
		Field:   "transactionType",
		Pattern: regexp.MustCompile(`(?i)\b(bán|cho\s*thuê)\b`),
		GetValue: func(match []string) interface{} {
			switch strings.ToLower(strings.ReplaceAll(match[1], " ", "")) {
			case "bán":
				return 10
			case "chothue":
				return 20
			default:
				return nil
			}
		},
	},
}

func InferSearchText(text string) (dto.TotalSearchParser, error) {
	raw := _utils.ParseWithRules(text, Rules)

	dataParser := dto.TotalSearchParser{}
	err := mapstructure.Decode(raw, &dataParser)
	return dataParser, err

}

func InferSearchText2(text string) (_dto.ProductV3DTO, error) {
	raw := _utils.ParseWithRules(text, Rules)

	dataParser := _dto.ProductV3DTO{}
	err := mapstructure.Decode(raw, &dataParser)
	return dataParser, err

}
