package forthic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"true", "TRUE", true},
		{"false", "FALSE", false},
		{"invalid", "true", nil},
		{"invalid", "True", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ToBool(tt.input)
			_ = ok
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"positive int", "42", int64(42)},
		{"negative int", "-10", int64(-10)},
		{"zero", "0", int64(0)},
		{"large int", "1000000", int64(1000000)},
		{"float should fail", "3.14", nil},
		{"invalid", "abc", nil},
		{"partial number", "42abc", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ToInt(tt.input)
			_ = ok
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"simple float", "3.14", 3.14},
		{"negative float", "-2.5", -2.5},
		{"zero float", "0.0", 0.0},
		{"no decimal should fail", "42", nil},
		{"invalid", "abc", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ToFloat(tt.input)
			_ = ok
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.InDelta(t, tt.expected, result, 0.0001)
			}
		})
	}
}

func TestToTime(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedHour int
		expectedMin  int
		shouldPass   bool
	}{
		{"simple time", "9:00", 9, 0, true},
		{"afternoon time", "14:30", 14, 30, true},
		{"PM time", "2:30 PM", 14, 30, true},
		{"AM time", "9:00 AM", 9, 0, true},
		{"noon", "12:00 PM", 12, 0, true},
		{"midnight", "12:00 AM", 0, 0, true},
		{"invalid", "25:00", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := ToTime(tt.input)
			_ = ok

			if tt.shouldPass {
				assert.NotNil(t, result)
				tm := result.(time.Time)
				assert.Equal(t, tt.expectedHour, tm.Hour())
				assert.Equal(t, tt.expectedMin, tm.Minute())
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestToLiteralDate(t *testing.T) {
	loc := time.UTC
	handler := ToLiteralDate(loc)

	tests := []struct {
		name       string
		input      string
		shouldPass bool
	}{
		{"valid date", "2020-06-05", true},
		{"year wildcard", "YYYY-06-05", true},
		{"month wildcard", "2020-MM-05", true},
		{"day wildcard", "2020-06-DD", true},
		{"all wildcards", "YYYY-MM-DD", true},
		{"invalid format", "2020/06/05", false},
		{"invalid", "not-a-date", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := handler(tt.input)
			_ = ok

			if tt.shouldPass {
				assert.NotNil(t, result)
				tm := result.(time.Time)
				// Just verify we got a time.Time
				assert.NotZero(t, tm)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestToZonedDateTime(t *testing.T) {
	loc := time.UTC
	handler := ToZonedDateTime(loc)

	tests := []struct {
		name       string
		input      string
		shouldPass bool
	}{
		{"UTC datetime", "2025-05-24T10:15:00Z", true},
		{"offset datetime", "2025-05-24T10:15:00-05:00", true},
		{"plain datetime", "2025-05-24T10:15:00", true},
		{"IANA timezone", "2025-05-20T08:00:00[America/Los_Angeles]", true},
		{"offset with IANA", "2025-05-20T08:00:00-07:00[America/Los_Angeles]", true},
		{"invalid", "not-a-datetime", false},
		{"no T separator", "2025-05-24 10:15:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := handler(tt.input)
			_ = ok

			if tt.shouldPass {
				assert.NotNil(t, result)
				tm := result.(time.Time)
				// Just verify we got a time.Time
				assert.NotZero(t, tm)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
