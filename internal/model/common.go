package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type HttpCode int

type TimeOnly struct {
	time.Time
}

// Implements custom JSON marshalling
func (timeOnly TimeOnly) MarshalJSON() ([]byte, error) {
	time := fmt.Sprintf(`"%s"`, timeOnly.Format("15:04"))
	return []byte(time), nil
}

// Implements the sql.Scanner interface
func (timeOnly *TimeOnly) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into TimeOnly", value)
	}
	timeOnly.Time = t
	return nil
}

// Implements the driver.Valuer interface
func (timeOnly TimeOnly) Value() (driver.Value, error) {
	return timeOnly.Time, nil
}
