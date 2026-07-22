// pkg/optional/json.go
package optional

// import "encoding/json"

// func (o Optional[T]) MarshalJSON() ([]byte, error) {
// 	switch o.state {
// 	case StateAbsent:
// 		return []byte("null"), nil // hoặc omitempty? Tuỳ logic.
// 	case StateNull:
// 		return []byte("null"), nil
// 	case StatePresent:
// 		return json.Marshal(o.value)
// 	}
// 	return []byte("null"), nil
// }

// func (o *Optional[T]) UnmarshalJSON(data []byte) error {
// 	if string(data) == "null" {
// 		o.state = StateNull
// 		o.value = *new(T)
// 		return nil
// 	}
// 	var v T
// 	if err := json.Unmarshal(data, &v); err != nil {
// 		return err
// 	}
// 	o.state = StatePresent
// 	o.value = v
// 	return nil
// }
