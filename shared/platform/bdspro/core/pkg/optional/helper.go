package optional

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func FromStringWrapper(v *wrapperspb.StringValue) String {
	if v == nil {
		return AbsentString()
	}
	return NewString(v.Value)
}

func FromUint64Wrapper(v *wrapperspb.UInt64Value) Uint64 {
	if v == nil {
		return AbsentUint64()
	}
	return NewUint64(v.Value)
}

// Tương tự cho các kiểu khác...
