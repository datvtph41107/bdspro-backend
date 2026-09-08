package utils

// Uint32SliceToInt32Slice converts a slice of uint32 to a slice of int32
func Uint32SliceToInt32Slice(slice []uint32) []int32 {
	result := make([]int32, len(slice))
	for i, v := range slice {
		result[i] = int32(v)
	}
	return result
}
