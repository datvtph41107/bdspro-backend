package _utils

import (
	sharepb "pb/types/shared"
	"sort"
)

func EnumMapToList[T ~uint32 | ~uint | ~int](enumMap map[T]string) []*sharepb.EnumItem {
	keys := make([]T, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}

	// Sắp xếp keys theo giá trị tăng dần
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := make([]*sharepb.EnumItem, 0, len(keys))
	for _, k := range keys {
		result = append(result, &sharepb.EnumItem{
			Value: uint32(k),
			Name:  enumMap[k],
		})
	}
	return result
}
