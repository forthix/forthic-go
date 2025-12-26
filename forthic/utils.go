package forthic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// Type Checking Utilities
// ============================================================================

// IsInt checks if a value can be treated as an integer
func IsInt(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsFloat checks if a value can be treated as a float
func IsFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	default:
		return false
	}
}

// IsString checks if a value is a string
func IsString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

// IsBool checks if a value is a boolean
func IsBool(v interface{}) bool {
	_, ok := v.(bool)
	return ok
}

// IsArray checks if a value is a slice/array
func IsArray(v interface{}) bool {
	switch v.(type) {
	case []interface{}:
		return true
	default:
		return false
	}
}

// IsRecord checks if a value is a map/record
func IsRecord(v interface{}) bool {
	switch v.(type) {
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

// ConvertToInt attempts to convert a value to int64
func ConvertToInt(v interface{}) (int64, error) {
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// ConvertToFloat attempts to convert a value to float64
func ConvertToFloat(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float32:
		return float64(val), nil
	case float64:
		return val, nil
	case int, int8, int16, int32, int64:
		i, _ := ConvertToInt(val)
		return float64(i), nil
	case uint, uint8, uint16, uint32, uint64:
		i, _ := ConvertToInt(val)
		return float64(i), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float", v)
	}
}

// ConvertToString attempts to convert a value to string
func ConvertToString(v interface{}) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ============================================================================
// String Utilities
// ============================================================================

// Trim removes leading and trailing whitespace
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// Split splits a string by a separator
func Split(s, sep string) []string {
	if sep == "" {
		// Split into individual characters
		chars := []string{}
		for _, r := range s {
			chars = append(chars, string(r))
		}
		return chars
	}
	return strings.Split(s, sep)
}

// Join joins strings with a separator
func Join(parts []string, sep string) string {
	return strings.Join(parts, sep)
}

// Replace replaces all occurrences of old with new in s
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// ============================================================================
// Date/Time Utilities
// ============================================================================

// ParseDate parses a date string in YYYY-MM-DD format
// Supports wildcards: YYYY-**-**, ****-MM-**, ****-**-DD
func ParseDate(s string) (time.Time, error) {
	// Check for wildcards
	if strings.Contains(s, "*") {
		// Replace wildcards with current date values
		now := time.Now()
		year := now.Year()
		month := int(now.Month())
		day := now.Day()

		parts := strings.Split(s, "-")
		if len(parts) != 3 {
			return time.Time{}, fmt.Errorf("invalid date format: %s", s)
		}

		if parts[0] != "****" {
			y, err := strconv.Atoi(parts[0])
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid year: %s", parts[0])
			}
			year = y
		}

		if parts[1] != "**" {
			m, err := strconv.Atoi(parts[1])
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid month: %s", parts[1])
			}
			month = m
		}

		if parts[2] != "**" {
			d, err := strconv.Atoi(parts[2])
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid day: %s", parts[2])
			}
			day = d
		}

		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
	}

	// Standard parsing
	return time.Parse("2006-01-02", s)
}

// ParseTime parses a time string in HH:MM or HH:MM:SS format
// Also supports 12-hour format with AM/PM (e.g., "2:30 PM")
func ParseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// Check for AM/PM format
	ampmRegex := regexp.MustCompile(`^(\d{1,2}):(\d{2})\s*(AM|PM)$`)
	if matches := ampmRegex.FindStringSubmatch(s); matches != nil {
		hour, _ := strconv.Atoi(matches[1])
		minute, _ := strconv.Atoi(matches[2])
		meridiem := matches[3]

		if meridiem == "PM" && hour < 12 {
			hour += 12
		} else if meridiem == "AM" && hour == 12 {
			hour = 0
		}

		return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC), nil
	}

	// Try HH:MM:SS format
	if t, err := time.Parse("15:04:05", s); err == nil {
		return t, nil
	}

	// Try HH:MM format
	return time.Parse("15:04", s)
}

// FormatDate formats a time as YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTime formats a time as HH:MM
func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

// FormatDateTime formats a time as RFC3339
func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseDateTime parses an RFC3339 datetime string
func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// DateToInt converts a date to YYYYMMDD integer format
func DateToInt(t time.Time) int64 {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()
	return int64(year*10000 + month*100 + day)
}

// IntToDate converts a YYYYMMDD integer to a date
func IntToDate(n int64) time.Time {
	year := int(n / 10000)
	month := int((n % 10000) / 100)
	day := int(n % 100)
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
