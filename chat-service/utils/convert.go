package utils

import "strconv"

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