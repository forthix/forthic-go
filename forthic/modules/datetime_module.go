package modules

import (
	"strings"
	"time"

	"github.com/forthix/forthic-go/forthic"
)

// DateTimeModule provides date and time operations
type DateTimeModule struct {
	*forthic.Module
}

// NewDateTimeModule creates a new datetime module
func NewDateTimeModule() *DateTimeModule {
	m := &DateTimeModule{
		Module: forthic.NewModule("datetime", ""),
	}
	m.registerWords()
	return m
}

func (m *DateTimeModule) registerWords() {
	// Current date/time
	m.AddModuleWord("TODAY", m.today)
	m.AddModuleWord("NOW", m.now)

	// Conversion TO datetime types
	m.AddModuleWord(">TIME", m.toTime)
	m.AddModuleWord(">DATE", m.toDate)
	m.AddModuleWord(">DATETIME", m.toDateTime)
	m.AddModuleWord("AT", m.at)

	// Conversion FROM datetime types
	m.AddModuleWord("TIME>STR", m.timeToStr)
	m.AddModuleWord("DATE>STR", m.dateToStr)
	m.AddModuleWord("DATE>INT", m.dateToInt)

	// Timestamps
	m.AddModuleWord(">TIMESTAMP", m.toTimestamp)
	m.AddModuleWord("TIMESTAMP>DATETIME", m.timestampToDatetime)

	// Date math
	m.AddModuleWord("ADD-DAYS", m.addDays)
	m.AddModuleWord("SUBTRACT-DATES", m.subtractDates)

	// Time adjustment
	m.AddModuleWord("AM", m.am)
	m.AddModuleWord("PM", m.pm)
}

// ========================================
// Current
// ========================================

func (m *DateTimeModule) today(interp *forthic.Interpreter) error {
	now := time.Now().UTC()
	// Return date with time set to 00:00:00
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	interp.StackPush(today)
	return nil
}

func (m *DateTimeModule) now(interp *forthic.Interpreter) error {
	interp.StackPush(time.Now().UTC())
	return nil
}

// ========================================
// Conversion TO datetime types
// ========================================

func (m *DateTimeModule) toTime(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	// If already a time.Time, extract just the time part
	if t, ok := item.(time.Time); ok {
		// Return time with year 0, month 1, day 1 (time-only representation)
		timeOnly := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
		interp.StackPush(timeOnly)
		return nil
	}

	// Parse as string
	str, ok := item.(string)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	str = strings.TrimSpace(str)

	// Try parsing various time formats
	formats := []string{
		"15:04",
		"15:04:05",
		"3:04 PM",
		"3:04PM",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			// Return time-only (year 0, month 1, day 1)
			timeOnly := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
			interp.StackPush(timeOnly)
			return nil
		}
	}

	interp.StackPush(nil)
	return nil
}

func (m *DateTimeModule) toDate(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	// If already a time.Time, return date part only
	if t, ok := item.(time.Time); ok {
		dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		interp.StackPush(dateOnly)
		return nil
	}

	// Parse as string
	str, ok := item.(string)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	str = strings.TrimSpace(str)

	// Try parsing date formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"Jan 2, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			// Return date only (time set to 00:00:00)
			dateOnly := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			interp.StackPush(dateOnly)
			return nil
		}
	}

	interp.StackPush(nil)
	return nil
}

func (m *DateTimeModule) toDateTime(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	// If already a time.Time, return it
	if t, ok := item.(time.Time); ok {
		interp.StackPush(t)
		return nil
	}

	// If it's a number, treat as Unix timestamp (seconds)
	if num, err := toNumber(item); err == nil {
		dt := time.Unix(int64(num), 0).UTC()
		interp.StackPush(dt)
		return nil
	}

	// Parse as string
	str, ok := item.(string)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	str = strings.TrimSpace(str)

	// Try parsing datetime formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, str); err == nil {
			interp.StackPush(t.UTC())
			return nil
		}
	}

	interp.StackPush(nil)
	return nil
}

func (m *DateTimeModule) at(interp *forthic.Interpreter) error {
	timeVal := interp.StackPop()
	dateVal := interp.StackPop()

	if dateVal == nil || timeVal == nil {
		interp.StackPush(nil)
		return nil
	}

	date, ok1 := dateVal.(time.Time)
	timeOnly, ok2 := timeVal.(time.Time)

	if !ok1 || !ok2 {
		interp.StackPush(nil)
		return nil
	}

	// Combine date and time components
	result := time.Date(
		date.Year(), date.Month(), date.Day(),
		timeOnly.Hour(), timeOnly.Minute(), timeOnly.Second(),
		0, time.UTC,
	)

	interp.StackPush(result)
	return nil
}

// ========================================
// Conversion FROM datetime types
// ========================================

func (m *DateTimeModule) timeToStr(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush("")
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush("")
		return nil
	}

	// Format as HH:MM
	result := t.Format("15:04")
	interp.StackPush(result)
	return nil
}

func (m *DateTimeModule) dateToStr(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush("")
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush("")
		return nil
	}

	// Format as YYYY-MM-DD
	result := t.Format("2006-01-02")
	interp.StackPush(result)
	return nil
}

func (m *DateTimeModule) dateToInt(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	// Convert to YYYYMMDD integer
	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	result := year*10000 + month*100 + day
	interp.StackPush(result)
	return nil
}

// ========================================
// Timestamps
// ========================================

func (m *DateTimeModule) toTimestamp(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	// Return Unix timestamp in seconds
	interp.StackPush(t.Unix())
	return nil
}

func (m *DateTimeModule) timestampToDatetime(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	// Convert to number
	timestamp, err := toNumber(item)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}

	// Convert Unix timestamp (seconds) to datetime
	dt := time.Unix(int64(timestamp), 0).UTC()
	interp.StackPush(dt)
	return nil
}

// ========================================
// Date math
// ========================================

func (m *DateTimeModule) addDays(interp *forthic.Interpreter) error {
	numDays := interp.StackPop()
	date := interp.StackPop()

	if date == nil || numDays == nil {
		interp.StackPush(nil)
		return nil
	}

	t, ok := date.(time.Time)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	days := toInt(numDays)
	result := t.AddDate(0, 0, days)

	interp.StackPush(result)
	return nil
}

func (m *DateTimeModule) subtractDates(interp *forthic.Interpreter) error {
	date2 := interp.StackPop()
	date1 := interp.StackPop()

	if date1 == nil || date2 == nil {
		interp.StackPush(nil)
		return nil
	}

	t1, ok1 := date1.(time.Time)
	t2, ok2 := date2.(time.Time)

	if !ok1 || !ok2 {
		interp.StackPush(nil)
		return nil
	}

	// Calculate difference: date1 - date2
	diff := t1.Sub(t2)
	days := int(diff.Hours() / 24)

	interp.StackPush(days)
	return nil
}

// ========================================
// Time adjustment
// ========================================

func (m *DateTimeModule) am(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush(item)
		return nil
	}

	// If hour is >= 12, subtract 12
	hour := t.Hour()
	if hour >= 12 {
		result := time.Date(t.Year(), t.Month(), t.Day(), hour-12, t.Minute(), t.Second(), 0, time.UTC)
		interp.StackPush(result)
	} else {
		interp.StackPush(t)
	}

	return nil
}

func (m *DateTimeModule) pm(interp *forthic.Interpreter) error {
	item := interp.StackPop()

	if item == nil {
		interp.StackPush(nil)
		return nil
	}

	t, ok := item.(time.Time)
	if !ok {
		interp.StackPush(item)
		return nil
	}

	// If hour is < 12, add 12
	hour := t.Hour()
	if hour < 12 {
		result := time.Date(t.Year(), t.Month(), t.Day(), hour+12, t.Minute(), t.Second(), 0, time.UTC)
		interp.StackPush(result)
	} else {
		interp.StackPush(t)
	}

	return nil
}
