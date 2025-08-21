package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FlexibleDate is a custom time type that can parse multiple date formats
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface
func (fd *FlexibleDate) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	str := strings.Trim(string(data), "\"")

	// List of supported date formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339 with timezone
		"2006-01-02T15:04:05Z",      // RFC3339 UTC
		"2006-01-02T15:04:05",       // ISO8601 without timezone
		"2006-01-02 15:04:05",       // SQL datetime format
		"2006-01-02",                // Date only (YYYY-MM-DD)
	}

	var err error
	for _, format := range formats {
		fd.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to parse date '%s', supported formats: YYYY-MM-DD, YYYY-MM-DDTHH:MM:SS, or RFC3339", str)
}

// MarshalJSON implements json.Marshaler interface
func (fd FlexibleDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.Time.Format(time.RFC3339))
}

// String returns the string representation
func (fd FlexibleDate) String() string {
	return fd.Time.Format("2006-01-02T15:04:05Z07:00")
}
