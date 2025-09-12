package cmd

import (
	"encoding/json"
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type CustomMarshaler struct{}

type ListResponse struct {
	Data          interface{} `json:"data"`
	TotalElements int         `json:"totalElements"`
}

func (c *CustomMarshaler) Marshal(v interface{}) ([]byte, error) {
	// Convert v thành map nếu có thể
	m, _ := v.(proto.Message)
	// log.Println("m", v, ok, m)
	val := reflect.ValueOf(m)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Lấy trường theo tên
	fieldVal := val.FieldByName("Data")
	if !fieldVal.IsValid() || !fieldVal.CanInterface() || val.Kind() == reflect.Array {
		return json.Marshal(map[string]interface{}{
			"code":    0,
			"message": "success",
			"data":    v,
		})
	}

	var data, totalElements interface{}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldVal := val.Field(i)

		switch field.Name {
		case "Data":
			if fieldVal.CanInterface() {
				data = fieldVal.Interface()
			}
		case "TotalElements":
			if fieldVal.CanInterface() {
				totalElements = fieldVal.Interface()
			}
		case "Total":
			if fieldVal.CanInterface() {
				totalElements = fieldVal.Interface()
			}
		}
	}

	// fallback default
	return json.Marshal(map[string]interface{}{
		"code":          0,
		"message":       "success",
		"data":          data,
		"totalElements": totalElements,
		"total":         totalElements,
	})
}

func (c *CustomMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (c *CustomMarshaler) NewDecoder(r io.Reader) runtime.Decoder {
	return json.NewDecoder(r)
}

func (c *CustomMarshaler) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

func (c *CustomMarshaler) ContentType(_ interface{}) string {
	return "application/json"
}
