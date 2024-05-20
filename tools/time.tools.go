package tools

import (
	"fmt"
	"strconv"
	"time"
)

// parseTime tries to parse various formats of time strings or timestamps
func ParseTime(input string, customFormats ...string) (time.Time, error) {
	// Try to parse Unix timestamp
	if timestamp, err := strconv.ParseInt(input, 10, 64); err == nil {
		// If the input is a timestamp in milliseconds
		if len(input) == 13 {
			return time.Unix(0, timestamp*int64(time.Millisecond)), nil
		}
		// If the input is a timestamp in seconds
		return time.Unix(timestamp, 0), nil
	}

	// Define default date formats
	defaultFormats := []string{
		time.RFC3339,           // Example: "2006-01-02T15:04:05Z07:00"
		"2006-01-02 15:04:05",  // Example: "2006-01-02 15:04:05"
		"2006-01-02",           // Example: "2006-01-02"
		"02 Jan 2006 15:04:05", // Example: "02 Jan 2006 15:04:05"
		"02-Jan-2006",          // Example: "02-Jan-2006"
	}

	// Combine default formats with custom formats
	formats := append(defaultFormats, customFormats...)

	// Try to parse date strings using the formats
	for _, format := range formats {
		if t, err := time.Parse(format, input); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time input: %s", input)
}
