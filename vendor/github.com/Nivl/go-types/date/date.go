package date

import (
	"database/sql/driver"
	"strings"
	"time"
)

// DATE is a time.Time layout for the a date (no time)
const DATE = "2006-01-02"

// Date represents a time.Time that uses DATE for json input/output
// instead of RFC3339
type Date struct {
	time.Time
}

// Today returns the current local day.
func Today() *Date {
	var day, year int
	var month time.Month
	year, month, day = time.Now().Date()
	return &Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// New accepts "year-month" or "year-month-day"
func New(date string) (*Date, error) {
	// If we only have year-month, then we add "-day"
	if strings.Count(date, "-") == 1 {
		date += "-01"
	}

	t, err := time.Parse(DATE, date)
	if err != nil {
		return nil, err
	}
	return &Date{Time: t}, nil
}

// Value - Implementation of valuer for database/sql
func (t *Date) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	return t.Format(DATE), nil
}

// Scan - Implement the database/sql scanner interface
func (t *Date) Scan(value interface{}) error {
	if value != nil {
		t.Time = value.(time.Time)
	}
	return nil
}

// String implements the fmt.Stringer interface
func (t Date) String() string {
	return t.Format(DATE)
}

// ScanString implements the go-params Scanner interface
func (t *Date) ScanString(date string) error {
	if strings.Count(date, "-") == 1 {
		date += "-01"
	}

	var err error
	t.Time, err = time.Parse(DATE, date)
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(DATE) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *Date) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(`"`+DATE+`"`, string(data))
	return
}

// Equal checks if the given date is equal to the current one
func (t *Date) Equal(u *Date) bool {
	return (t.Time.Year() == u.Time.Year()) &&
		(t.Time.Month() == u.Time.Month()) &&
		(t.Time.Day() == u.Time.Day())
}

// IsBefore checks if the current date is before the given one
func (t *Date) IsBefore(u *Date) bool {
	if t.Time.Year() != u.Time.Year() {
		return t.Time.Year() < u.Time.Year()
	}
	if t.Time.Month() != u.Time.Month() {
		return t.Time.Month() < u.Time.Month()
	}
	if t.Time.Day() != u.Time.Day() {
		return t.Time.Day() < u.Time.Day()
	}
	return false
}

// IsAfter checks if the current date is after the given one
func (t *Date) IsAfter(u *Date) bool {
	if t.Time.Year() != u.Time.Year() {
		return t.Time.Year() > u.Time.Year()
	}
	if t.Time.Month() != u.Time.Month() {
		return t.Time.Month() > u.Time.Month()
	}
	if t.Time.Day() != u.Time.Day() {
		return t.Time.Day() > u.Time.Day()
	}
	return false
}
