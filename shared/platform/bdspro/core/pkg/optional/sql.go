// pkg/optional/sql.go
package optional

// func (o Optional[T]) SQLValue() (driver.Value, error) {
// 	if !o.IsPresent() {
// 		return nil, nil
// 	}
// 	return o.value, nil
// }

// func (o *Optional[T]) Scan(src interface{}) error {
// 	if src == nil {
// 		o.state = StateNull
// 		o.value = *new(T)
// 		return nil
// 	}
// 	// Tuỳ T mà xử lý, có thể dùng reflection hoặc type switch.
// 	// Đơn giản: giả sử T là kiểu cơ bản, dùng convert.
// 	var v T
// 	switch dst := any(&v).(type) {
// 	case *string:
// 		if s, ok := src.(string); ok {
// 			*dst = s
// 		} else if b, ok := src.([]byte); ok {
// 			*dst = string(b)
// 		} else {
// 			return fmt.Errorf("cannot scan %T into optional.String", src)
// 		}
// 	case *int64:
// 		// ...
// 	default:
// 		return fmt.Errorf("unsupported type for scan")
// 	}
// 	o.state = StatePresent
// 	o.value = v
// 	return nil
// }
