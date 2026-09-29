package utils

import (
	"fmt"
	"strconv"
)

func Int32SliceToUint64Slice(members []int32) []uint64 {
	converted := make([]uint64, len(members))
	for i, member := range members {
		converted[i] = uint64(member)
	}
	return converted
}

func Int32SliceToStringSlice(numbers []int32) []string {
	converted := make([]string, len(numbers))
	for i, num := range numbers {
		converted[i] = strconv.Itoa(int(num))
	}
	return converted
}

// ParseUint64 safely converts interface{} to uint64
// Supports: float64, int, int32, int64, uint, uint32, uint64, string
func ParseUint64(value interface{}) (uint64, bool) {
	switch v := value.(type) {
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int32:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case uint:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

// ParseUint64Slice safely converts []interface{} to []uint64
func ParseUint64Slice(values []interface{}) ([]uint64, error) {
	result := make([]uint64, 0, len(values))
	for i, v := range values {
		parsed, ok := ParseUint64(v)
		if !ok {
			return nil, fmt.Errorf("failed to parse element at index %d: %v", i, v)
		}
		result = append(result, parsed)
	}
	return result, nil
}

// ParseUint64WithDefault safely converts interface{} to uint64 with default value
func ParseUint64WithDefault(value interface{}, defaultValue uint64) uint64 {
	if parsed, ok := ParseUint64(value); ok {
		return parsed
	}
	return defaultValue
}
