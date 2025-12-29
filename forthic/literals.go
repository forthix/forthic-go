package forthic

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// Boolean Literals
// ============================================================================

// ToBool parses boolean literals: TRUE, FALSE
func ToBool(str string) (interface{}, bool) {
	if str == "TRUE" {
		return true, true
	}
	if str == "FALSE" {
		return false, true
	}
	return nil, false
}

// ============================================================================
// Numeric Literals
// ============================================================================

// ToFloat parses float literals: 3.14, -2.5, 0.0
// Must contain a decimal point
func ToFloat(str string) (interface{}, bool) {
	if !strings.Contains(str, ".") {
		return nil, false
	}
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, false
	}
	return result, true
}

// ToInt parses integer literals: 42, -10, 0
// Must not contain a decimal point
func ToInt(str string) (interface{}, bool) {
	if strings.Contains(str, ".") {
		return nil, false
	}
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, false
	}
	// Verify it's actually an integer string (not "42abc")
	if strconv.FormatInt(result, 10) != str {
		return nil, false
	}
	return result, true
}

// ============================================================================
// Time Literals
// ============================================================================

// ToTime parses time literals: 9:00, 11:30 PM, 22:15
func ToTime(str string) (interface{}, bool) {
	// Pattern: HH:MM or HH:MM AM/PM
	re := regexp.MustCompile(`^(\d{1,2}):(\d{2})(?:\s*(AM|PM))?$`)
	match := re.FindStringSubmatch(str)
	if match == nil {
		return nil, false
	}

	hours, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, false
	}
	minutes, err := strconv.Atoi(match[2])
	if err != nil {
		return nil, false
	}
	meridiem := match[3]

	// Adjust for AM/PM
	if meridiem == "PM" && hours < 12 {
		hours += 12
	} else if meridiem == "AM" && hours == 12 {
		hours = 0
	} else if meridiem == "AM" && hours > 12 {
		// Handle invalid cases like "22:15 AM"
		hours -= 12
	}

	if hours > 23 || minutes >= 60 {
		return nil, false
	}

	// Return a time.Time with year 0, month 1, day 1 (time-only representation)
	return time.Date(0, 1, 1, hours, minutes, 0, 0, time.UTC), true
}

// ============================================================================
// Date Literals
// ============================================================================

// ToLiteralDate creates a date literal handler
// Parses: 2020-06-05, YYYY-MM-DD (with wildcards)
func ToLiteralDate(timezone *time.Location) LiteralHandler {
	return func(str string) (interface{}, bool) {
		// Pattern: YYYY-MM-DD or wildcards (YYYY, MM, DD)
		re := regexp.MustCompile(`^(\d{4}|YYYY)-(\d{2}|MM)-(\d{2}|DD)$`)
		match := re.FindStringSubmatch(str)
		if match == nil {
			return nil, false
		}

		now := time.Now().In(timezone)
		year := now.Year()
		month := int(now.Month())
		day := now.Day()

		if match[1] != "YYYY" {
			y, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, false
			}
			year = y
		}

		if match[2] != "MM" {
			m, err := strconv.Atoi(match[2])
			if err != nil {
				return nil, false
			}
			month = m
		}

		if match[3] != "DD" {
			d, err := strconv.Atoi(match[3])
			if err != nil {
				return nil, false
			}
			day = d
		}

		result := time.Date(year, time.Month(month), day, 0, 0, 0, 0, timezone)
		return result, true
	}
}

// ============================================================================
// ZonedDateTime Literals
// ============================================================================

// ToZonedDateTime creates a zoned datetime literal handler
// Parses:
// - 2025-05-24T10:15:00[America/Los_Angeles] (IANA named timezone, RFC 9557)
// - 2025-05-24T10:15:00-07:00[America/Los_Angeles] (offset + IANA timezone)
// - 2025-05-24T10:15:00Z (UTC)
// - 2025-05-24T10:15:00-05:00 (offset timezone)
// - 2025-05-24T10:15:00 (uses interpreter's timezone)
func ToZonedDateTime(timezone *time.Location) LiteralHandler {
	return func(str string) (interface{}, bool) {
		if !strings.Contains(str, "T") {
			return nil, false
		}

		// Handle IANA named timezone in bracket notation (RFC 9557)
		// Examples: 2025-05-20T08:00:00[America/Los_Angeles]
		//           2025-05-20T08:00:00-07:00[America/Los_Angeles]
		if strings.Contains(str, "[") && strings.HasSuffix(str, "]") {
			// Extract timezone name from brackets
			bracketStart := strings.Index(str, "[")
			bracketEnd := strings.Index(str, "]")
			tzName := str[bracketStart+1 : bracketEnd]

			// Load the timezone
			loc, err := time.LoadLocation(tzName)
			if err != nil {
				return nil, false
			}

			// Parse the datetime part (before the bracket)
			dtStr := str[:bracketStart]

			// Try parsing with offset first (2025-05-20T08:00:00-07:00)
			if strings.Contains(dtStr, "+") || strings.LastIndex(dtStr, "-") > 10 {
				// Has offset, parse as RFC3339
				t, err := time.Parse(time.RFC3339, dtStr)
				if err != nil {
					return nil, false
				}
				// Convert to the named timezone
				return t.In(loc), true
			}

			// No offset, parse as plain datetime and assign timezone
			t, err := time.Parse("2006-01-02T15:04:05", dtStr)
			if err != nil {
				return nil, false
			}
			return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc), true
		}

		// Handle explicit UTC (Z suffix)
		if strings.HasSuffix(str, "Z") {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil, false
			}
			return t.UTC(), true
		}

		// Handle explicit timezone offset (+05:00, -05:00)
		offsetRe := regexp.MustCompile(`[+-]\d{2}:\d{2}$`)
		if offsetRe.MatchString(str) {
			t, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil, false
			}
			// Convert to UTC for canonical storage
			return t.UTC(), true
		}

		// No timezone specified, use interpreter's timezone
		t, err := time.Parse("2006-01-02T15:04:05", str)
		if err != nil {
			return nil, false
		}
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), timezone), true
	}
}
