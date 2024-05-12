package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonMap[T any] map[string]interface{}

func (j *JsonMap[T]) Scan(src interface{}) error {
	var source []byte

	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	case nil:
		source = nil
	default:
		return fmt.Errorf("cannot convert %T to JsonB", src)
	}

	json.Unmarshal(source, j)
	return nil
}

func (j JsonMap[T]) Value() (driver.Value, error) {
	return json.Marshal(j)
}
