package modules

import (
	"testing"
	"time"

	"github.com/forthix/forthic-go/forthic"
)

func setupDateTimeInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()

	// Import datetime module
	dtMod := NewDateTimeModule()
	interp.ImportModule(dtMod.Module, "")

	return interp
}

// ========================================
// Current
// ========================================

func TestDateTime_Today(t *testing.T) {
	interp := setupDateTimeInterpreter()
	before := time.Now().UTC()
	err := interp.Run(`TODAY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)

	// TODAY should be today's date with time 00:00:00
	if result.Year() != before.Year() || result.Month() != before.Month() || result.Day() != before.Day() {
		t.Errorf("Expected today's date, got %v", result)
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Errorf("Expected time 00:00:00, got %02d:%02d:%02d", result.Hour(), result.Minute(), result.Second())
	}
}

func TestDateTime_Now(t *testing.T) {
	interp := setupDateTimeInterpreter()
	before := time.Now().UTC()
	err := interp.Run(`NOW`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	after := time.Now().UTC()

	// NOW should be between before and after
	if result.Before(before) || result.After(after) {
		t.Errorf("Expected NOW to be current time, got %v", result)
	}
}

// ========================================
// Conversion TO datetime types
// ========================================

func TestDateTime_ToTime(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"14:30" >TIME`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Hour() != 14 || result.Minute() != 30 {
		t.Errorf("Expected 14:30, got %02d:%02d", result.Hour(), result.Minute())
	}
}

func TestDateTime_ToDate(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Year() != 2023 || result.Month() != 6 || result.Day() != 15 {
		t.Errorf("Expected 2023-06-15, got %v", result)
	}
}

func TestDateTime_ToDateTimeFromString(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15T14:30:45" >DATETIME`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Year() != 2023 || result.Month() != 6 || result.Day() != 15 ||
		result.Hour() != 14 || result.Minute() != 30 || result.Second() != 45 {
		t.Errorf("Expected 2023-06-15 14:30:45, got %v", result)
	}
}

func TestDateTime_ToDateTimeFromTimestamp(t *testing.T) {
	interp := setupDateTimeInterpreter()
	// 1672531200 = 2023-01-01 00:00:00 UTC
	err := interp.Run(`1672531200 >DATETIME`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Year() != 2023 || result.Month() != 1 || result.Day() != 1 {
		t.Errorf("Expected 2023-01-01, got %v", result)
	}
}

func TestDateTime_AT(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE "14:30" >TIME AT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Year() != 2023 || result.Month() != 6 || result.Day() != 15 ||
		result.Hour() != 14 || result.Minute() != 30 {
		t.Errorf("Expected 2023-06-15 14:30, got %v", result)
	}
}

// ========================================
// Conversion FROM datetime types
// ========================================

func TestDateTime_TimeToStr(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"14:30" >TIME TIME>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "14:30" {
		t.Errorf("Expected '14:30', got '%s'", result)
	}
}

func TestDateTime_DateToStr(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE DATE>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "2023-06-15" {
		t.Errorf("Expected '2023-06-15', got '%s'", result)
	}
}

func TestDateTime_DateToInt(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE DATE>INT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int)
	if result != 20230615 {
		t.Errorf("Expected 20230615, got %d", result)
	}
}

// ========================================
// Timestamps
// ========================================

func TestDateTime_ToTimestamp(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-01-01T00:00:00" >DATETIME >TIMESTAMP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int64)
	// 2023-01-01 00:00:00 UTC is 1672531200
	if result != 1672531200 {
		t.Errorf("Expected timestamp 1672531200, got %v", result)
	}
}

func TestDateTime_TimestampToDatetime(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`1672531200 TIMESTAMP>DATETIME`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(time.Time)
	if result.Year() != 2023 || result.Month() != 1 || result.Day() != 1 {
		t.Errorf("Expected 2023-01-01, got %v", result)
	}
}

// ========================================
// Date math
// ========================================

func TestDateTime_AddDays(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-01-01" >DATE 30 ADD-DAYS DATE>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "2023-01-31" {
		t.Errorf("Expected '2023-01-31', got '%s'", result)
	}
}

func TestDateTime_AddDaysNegative(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE -5 ADD-DAYS DATE>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "2023-06-10" {
		t.Errorf("Expected '2023-06-10', got '%s'", result)
	}
}

func TestDateTime_SubtractDates(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-25" >DATE "2023-06-15" >DATE SUBTRACT-DATES`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int)
	if result != 10 {
		t.Errorf("Expected 10 days, got %v", result)
	}
}

func TestDateTime_SubtractDatesNegative(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"2023-06-15" >DATE "2023-06-25" >DATE SUBTRACT-DATES`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int)
	if result != -10 {
		t.Errorf("Expected -10 days, got %v", result)
	}
}

// ========================================
// Time adjustment
// ========================================

func TestDateTime_AM(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"14:30" >TIME AM TIME>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "02:30" {
		t.Errorf("Expected '02:30', got '%s'", result)
	}
}

func TestDateTime_PM(t *testing.T) {
	interp := setupDateTimeInterpreter()
	err := interp.Run(`"10:30" >TIME PM TIME>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "22:30" {
		t.Errorf("Expected '22:30', got '%s'", result)
	}
}

// ========================================
// Integration Tests with Date Literals
// ========================================

func TestDateTime_DateLiteralsWithOperations(t *testing.T) {
	interp := setupDateTimeInterpreter()

	// Test: 2023-01-01 + 30 days → 2023-01-31
	err := interp.Run(`2023-01-01 30 ADD-DAYS DATE>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result := interp.StackPop().(string)
	if result != "2023-01-31" {
		t.Errorf("Expected '2023-01-31', got '%s'", result)
	}
}

func TestDateTime_DateTimeLiteralsWithOperations(t *testing.T) {
	interp := setupDateTimeInterpreter()

	// Test: 2023-06-15T14:30:45 → timestamp
	err := interp.Run(`2023-06-15T14:30:45 >TIMESTAMP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	timestamp := interp.StackPop().(int64)
	// Verify it's a reasonable timestamp (June 2023)
	if timestamp < 1686000000 || timestamp > 1687000000 {
		t.Errorf("Expected timestamp around June 2023, got %d", timestamp)
	}
}

func TestDateTime_CombineDateAndTime(t *testing.T) {
	interp := setupDateTimeInterpreter()

	// Test: TODAY + "14:30" → datetime
	err := interp.Run(`TODAY "14:30" >TIME AT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result := interp.StackPop().(time.Time)
	if result.Hour() != 14 || result.Minute() != 30 {
		t.Errorf("Expected time 14:30, got %02d:%02d", result.Hour(), result.Minute())
	}
}

func TestDateTime_RoundTrip(t *testing.T) {
	interp := setupDateTimeInterpreter()

	// Test: date → string → date
	err := interp.Run(`"2023-06-15" >DATE DATE>STR >DATE DATE>STR`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result := interp.StackPop().(string)
	if result != "2023-06-15" {
		t.Errorf("Expected '2023-06-15', got '%s'", result)
	}
}
