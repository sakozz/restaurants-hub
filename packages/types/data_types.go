package types

import (
	"database/sql"
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

type NullTime struct {
	sql.NullTime
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.NullTime.Time)
	}
	return json.Marshal(nil)
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	var t *sql.NullTime
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	if t != nil {
		nt.Valid = true
		nt.NullTime = *t
	} else {
		nt.Valid = false
	}
	return nil
}

type NullInt struct {
	sql.NullInt64
}

func (nInt NullInt) MarshalJSON() ([]byte, error) {
	if nInt.Valid {
		return json.Marshal(nInt.NullInt64.Int64)
	}
	return json.Marshal(nil)
}

func (nInt *NullInt) UnmarshalJSON(data []byte) error {
	var i *sql.NullInt64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	if i != nil {
		nInt.Valid = true
		nInt.NullInt64 = *i
	} else {
		nInt.Valid = false
	}
	return nil
}
