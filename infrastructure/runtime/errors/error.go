package _errors

import (
	"reflect"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func InvalidRequest(details ...protoadapt.MessageV1) error {
	st := status.New(codes.InvalidArgument, "invalid request")
	if len(details) > 0 {
		st, _ = st.WithDetails(details...)
	}
	return st.Err()
}

func ThrowError(code int32, err error, details ...protoadapt.MessageV1) error {
	st := status.New(codes.Code(code), err.Error())

	if len(details) > 0 {
		stWithDetails, errDetails := st.WithDetails(details...)
		if errDetails != nil {
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return st.Err()
}

// ReturnError is the one caller-facing application-error entry point.
//
// During repository migration it accepts both:
//   - ReturnError(Spec, Option...)
//   - the historical ReturnError(number, message)
//
// The legacy shape is preserved only for protected/unmigrated callers and must
// be removed after their explicit compatibility gate. New code must pass Spec.
func ReturnError(value interface{}, arguments ...interface{}) error {
	if spec, ok := value.(Spec); ok {
		options := make([]Option, 0, len(arguments))
		for _, argument := range arguments {
			option, ok := argument.(Option)
			if !ok {
				return status.Error(codes.Internal, "invalid canonical error option")
			}
			options = append(options, option)
		}
		return newError(spec, options...)
	}

	code, ok := legacyNumericCode(value)
	if !ok || len(arguments) != 1 {
		return status.Error(codes.Internal, "invalid legacy error contract")
	}
	message, ok := arguments[0].(string)
	if !ok {
		return status.Error(codes.Internal, "invalid legacy error message")
	}
	return legacyReturnError(code, message)
}

func legacyNumericCode(value interface{}) (int32, bool) {
	if value == nil {
		return 0, false
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int32(reflected.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int32(reflected.Uint()), true
	default:
		return 0, false
	}
}

func legacyReturnError(code int32, message string) error {
	st := status.New(codes.Internal, message)
	return withLegacyErrorDetail(st, code, message, nil).Err()
}
