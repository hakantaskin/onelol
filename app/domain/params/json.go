package params

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSON struct {
	json.RawMessage
}

// Value get value of Jsonb
func (j JSON) Value() (driver.Value, error) {
	if len(j.RawMessage) == 0 {
		return nil, nil
	}
	return j.MarshalJSON()
}

// Scan scan value into Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	if len(bytes) <= 0 {
		return nil
	}

	return json.Unmarshal(bytes, j)
}
