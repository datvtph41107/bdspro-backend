package utils

import (
	"strconv"
	"strings"
)

func ParseUintCSV(value string) []uint64 {
	parts := strings.Split(value, ",")
	result := make([]uint64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		parsed, err := strconv.ParseUint(part, 10, 64)
		if err == nil && parsed > 0 {
			result = append(result, parsed)
		}
	}
	return result
}

func MergePositiveUint64(primary uint64, values []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(values)+1)
	result := make([]uint64, 0, len(values)+1)
	if primary > 0 {
		seen[primary] = struct{}{}
		result = append(result, primary)
	}
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
