package models

import (
	"fmt"
	"time"
)

// Custom Date type to enforce "YYYY-MM-DD" format
type Date time.Time

// MarshalJSON formats the Date field as "YYYY-MM-DD"
func (d Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(d).Format("2006-01-02"))
	return []byte(formatted), nil
}

// UnmarshalJSON parses incoming JSON date strings as "YYYY-MM-DD"
func (d *Date) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse("\"2006-01-02\"", string(data))
	if err != nil {
		return err
	}
	*d = Date(parsedTime)
	return nil
}

// String function to return "YYYY-MM-DD" format
func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}
