package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Custom type for date-only handling
type DateOnly time.Time

// MarshalJSON serializes the date in `YYYY-MM-DD` format
func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

// UnmarshalJSON deserializes a `YYYY-MM-DD` string into a DateOnly
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var dateString string
	if err := json.Unmarshal(b, &dateString); err != nil {
		return err
	}

	parsedTime, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return err
	}

	*d = DateOnly(parsedTime)
	return nil
}

// ToTime converts DateOnly to time.Time
func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}

// Value converts DateOnly to a SQL-compatible value
func (d DateOnly) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

// Scan assigns a value from a database driver
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02", v)
		if err != nil {
			return fmt.Errorf("cannot parse date: %v", err)
		}
		*d = DateOnly(parsedTime)
		return nil
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}
}
