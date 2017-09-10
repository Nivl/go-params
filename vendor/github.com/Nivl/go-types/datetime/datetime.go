package datetime

import (
	"database/sql/driver"
	"time"
)

// ISO8601 is a time.Time layout for the ISO8601 format
const ISO8601 = "2006-01-02T15:04:05-0700"

// DateTime represents a time.Time that uses ISO8601 for json input/output
// instead of RFC3339
type DateTime struct {
	time.Time
}

// Now returns the current local time.
func Now() *DateTime {
	return &DateTime{Time: time.Now()}
}

// Value - Implementation of valuer for database/sql
func (t *DateTime) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	return t.Format(ISO8601), nil
}

// Scan - Implement the database/sql scanner interface
func (t *DateTime) Scan(value interface{}) error {
	if value != nil {
		t.Time = value.(time.Time)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(ISO8601) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(`"`+ISO8601+`"`, string(data))
	return
}

// Equal check if the given date is equal to the current one
func (t *DateTime) Equal(u *DateTime) bool {
	return t.Time.Equal(u.Time)
}

// AddDate returns the time corresponding to adding the given number of years, months, and days to t.
func (t *DateTime) AddDate(years int, months int, days int) *DateTime {
	return &DateTime{Time: t.Time.AddDate(years, months, days)}
}
