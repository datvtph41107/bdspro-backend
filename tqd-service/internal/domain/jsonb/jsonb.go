package jsonb

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONB - Custom type cho PostgreSQL JSONB
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(map[string]interface{}(j))
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*j = nil
		return nil
	}

	// Thử parse như object
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		*j = obj
		return nil
	}

	// Thử parse như string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			*j = nil
			return nil
		}
		if err := json.Unmarshal([]byte(str), &obj); err != nil {
			return err
		}
		*j = obj
		return nil
	}

	return nil
}
