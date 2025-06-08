package utils

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func GetPageFromRequest(page *wrapperspb.Int32Value) int {
	if page == nil {
		return 0
	}
	return int(page.Value)
}

func GetSizeFromRequest(size *wrapperspb.Int32Value) int {
	if size == nil {
		return 10
	}
	return int(size.Value)
} 