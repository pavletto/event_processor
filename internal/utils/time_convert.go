package utils

import (
	"strings"
	"time"
)

var knownFormats = []string{
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.DateTime,
	time.DateOnly,
	time.TimeOnly,
}

func parseAnyDate(value string) (string, bool) {
	for _, format := range knownFormats {
		if parsedTime, err := time.Parse(format, value); err == nil {
			return parsedTime.Format(time.RFC3339Nano), true
		}
	}

	if strings.Contains(value, " ") && !strings.Contains(value, "T") {
		value = strings.Replace(value, " ", "T", 1)
		for _, format := range knownFormats {
			if parsedTime, err := time.Parse(format, value); err == nil {
				return parsedTime.Format(time.RFC3339Nano), true
			}
		}
	}

	return value, false
}

func TransformDates(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			v[key] = TransformDates(value)
		}
	case []interface{}:
		for i, value := range v {
			v[i] = TransformDates(value)
		}
	case string:
		if newDate, ok := parseAnyDate(v); ok {
			return newDate
		}
	}
	return data
}
