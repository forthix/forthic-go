package modules

import (
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupArrayInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()

	// Import array module
	arrayMod := NewArrayModule()
	interp.ImportModule(arrayMod.Module, "")

	// Import math module for numeric operations in MAP/SELECT/REDUCE
	mathMod := NewMathModule()
	interp.ImportModule(mathMod.Module, "")

	// Import boolean module for comparison operations in SELECT
	boolMod := NewBooleanModule()
	interp.ImportModule(boolMod.Module, "")

	return interp
}

// ========================================
// Basic Operations
// ========================================

func TestArray_Append(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3] 4 APPEND`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 4 {
		t.Fatalf("Expected 4 elements, got %d", len(result))
	}
	if result[3].(int64) != 4 {
		t.Errorf("Expected last element to be 4, got %v", result[3])
	}
}

func TestArray_AppendToNil(t *testing.T) {
	interp := setupArrayInterpreter()
	interp.StackPush(nil)
	interp.StackPush(1)
	err := interp.Run(`APPEND`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(result))
	}
}

func TestArray_Reverse(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3] REVERSE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(int64) != 3 || result[1].(int64) != 2 || result[2].(int64) != 1 {
		t.Errorf("Expected [3 2 1], got %v", result)
	}
}

func TestArray_Unique(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 2 3 1 4] UNIQUE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 4 {
		t.Fatalf("Expected 4 unique elements, got %d", len(result))
	}
	// Check that we have 1, 2, 3, 4 (order preserved)
	if result[0].(int64) != 1 || result[1].(int64) != 2 || result[2].(int64) != 3 || result[3].(int64) != 4 {
		t.Errorf("Expected [1 2 3 4], got %v", result)
	}
}

func TestArray_Length(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] LENGTH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int)
	if result != 5 {
		t.Errorf("Expected 5, got %d", result)
	}
}

func TestArray_LengthNil(t *testing.T) {
	interp := setupArrayInterpreter()
	interp.StackPush(nil)
	err := interp.Run(`LENGTH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int)
	if result != 0 {
		t.Errorf("Expected 0 for nil, got %d", result)
	}
}

// ========================================
// Access Operations
// ========================================

func TestArray_Nth(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40] 2 NTH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int64)
	if result != 30 {
		t.Errorf("Expected 30, got %d", result)
	}
}

func TestArray_NthOutOfBounds(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30] 5 NTH`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for out of bounds, got %v", result)
	}
}

func TestArray_Last(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40] LAST`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int64)
	if result != 40 {
		t.Errorf("Expected 40, got %d", result)
	}
}

func TestArray_LastEmpty(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[] LAST`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for empty array, got %v", result)
	}
}

func TestArray_Slice(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40 50] 1 3 SLICE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(int64) != 20 || result[1].(int64) != 30 || result[2].(int64) != 40 {
		t.Errorf("Expected [20 30 40], got %v", result)
	}
}

func TestArray_SliceNegativeIndices(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40 50] -2 -1 SLICE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(result))
	}
	if result[0].(int64) != 40 || result[1].(int64) != 50 {
		t.Errorf("Expected [40 50], got %v", result)
	}
}

func TestArray_SliceReverse(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40 50] 3 1 SLICE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(int64) != 40 || result[1].(int64) != 30 || result[2].(int64) != 20 {
		t.Errorf("Expected [40 30 20] (reverse), got %v", result)
	}
}

func TestArray_Take(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40 50] 3 TAKE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(int64) != 10 || result[1].(int64) != 20 || result[2].(int64) != 30 {
		t.Errorf("Expected [10 20 30], got %v", result)
	}
}

func TestArray_TakeMoreThanLength(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20] 5 TAKE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(result))
	}
}

func TestArray_Drop(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20 30 40 50] 2 DROP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(int64) != 30 || result[1].(int64) != 40 || result[2].(int64) != 50 {
		t.Errorf("Expected [30 40 50], got %v", result)
	}
}

func TestArray_DropMoreThanLength(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[10 20] 5 DROP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 0 {
		t.Fatalf("Expected 0 elements, got %d", len(result))
	}
}

// ========================================
// Set Operations
// ========================================

func TestArray_Difference(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] [2 4] DIFFERENCE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	// Should have [1 3 5]
	if result[0].(int64) != 1 || result[1].(int64) != 3 || result[2].(int64) != 5 {
		t.Errorf("Expected [1 3 5], got %v", result)
	}
}

func TestArray_Intersection(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] [2 4 6] INTERSECTION`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(result))
	}
	// Should have [2 4]
	if result[0].(int64) != 2 || result[1].(int64) != 4 {
		t.Errorf("Expected [2 4], got %v", result)
	}
}

func TestArray_Union(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3] [3 4 5] UNION`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 5 {
		t.Fatalf("Expected 5 unique elements, got %d", len(result))
	}
	// Should have [1 2 3 4 5] (deduplicated, order preserved)
	if result[0].(int64) != 1 || result[1].(int64) != 2 || result[2].(int64) != 3 || result[3].(int64) != 4 || result[4].(int64) != 5 {
		t.Errorf("Expected [1 2 3 4 5], got %v", result)
	}
}

// ========================================
// Sort
// ========================================

func TestArray_SortNumbers(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[3 1 4 1 5 9 2 6] SORT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 8 {
		t.Fatalf("Expected 8 elements, got %d", len(result))
	}
	// Check first and last elements
	if result[0].(int64) != 1 || result[7].(int64) != 9 {
		t.Errorf("Expected sorted array starting with 1 and ending with 9, got %v", result)
	}
}

func TestArray_SortStrings(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`["zebra" "apple" "mango" "banana"] SORT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 4 {
		t.Fatalf("Expected 4 elements, got %d", len(result))
	}
	if result[0].(string) != "apple" || result[3].(string) != "zebra" {
		t.Errorf("Expected sorted strings, got %v", result)
	}
}

// ========================================
// Combine
// ========================================

func TestArray_Zip(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3] ["a" "b" "c"] ZIP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 pairs, got %d", len(result))
	}

	// Check first pair
	pair0 := result[0].([]interface{})
	if pair0[0].(int64) != 1 || pair0[1].(string) != "a" {
		t.Errorf("Expected [1 'a'], got %v", pair0)
	}

	// Check last pair
	pair2 := result[2].([]interface{})
	if pair2[0].(int64) != 3 || pair2[1].(string) != "c" {
		t.Errorf("Expected [3 'c'], got %v", pair2)
	}
}

func TestArray_ZipDifferentLengths(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] ["a" "b"] ZIP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	// Should only zip up to length of shorter array
	if len(result) != 2 {
		t.Fatalf("Expected 2 pairs (min length), got %d", len(result))
	}
}

func TestArray_Flatten(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[[1 2] [3 4] [5 6]] FLATTEN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 6 {
		t.Fatalf("Expected 6 elements, got %d", len(result))
	}
	if result[0].(int64) != 1 || result[5].(int64) != 6 {
		t.Errorf("Expected [1 2 3 4 5 6], got %v", result)
	}
}

func TestArray_FlattenMixed(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[[1 2] 3 [4 5]] FLATTEN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 5 {
		t.Fatalf("Expected 5 elements, got %d", len(result))
	}
	if result[0].(int64) != 1 || result[2].(int64) != 3 || result[4].(int64) != 5 {
		t.Errorf("Expected [1 2 3 4 5], got %v", result)
	}
}

// ========================================
// Transform
// ========================================

func TestArray_Map(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] "2 *" MAP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 5 {
		t.Fatalf("Expected 5 elements, got %d", len(result))
	}
	// Check that values are doubled
	if result[0].(float64) != 2.0 || result[4].(float64) != 10.0 {
		t.Errorf("Expected [2 4 6 8 10], got %v", result)
	}
}

func TestArray_MapWithStrings(t *testing.T) {
	interp := setupArrayInterpreter()
	// Need to import string module for UPPERCASE
	stringMod := NewStringModule()
	interp.ImportModule(stringMod.Module, "")

	err := interp.Run(`["hello" "world"] "UPPERCASE" MAP`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 elements, got %d", len(result))
	}
	if result[0].(string) != "HELLO" || result[1].(string) != "WORLD" {
		t.Errorf("Expected ['HELLO' 'WORLD'], got %v", result)
	}
}

func TestArray_Select(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5 6] "2 MOD 0 ==" SELECT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 even numbers, got %d", len(result))
	}
	if result[0].(int64) != 2 || result[1].(int64) != 4 || result[2].(int64) != 6 {
		t.Errorf("Expected [2 4 6], got %v", result)
	}
}

func TestArray_SelectWithComparison(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] "3 >" SELECT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 elements > 3, got %d", len(result))
	}
	if result[0].(int64) != 4 || result[1].(int64) != 5 {
		t.Errorf("Expected [4 5], got %v", result)
	}
}

func TestArray_Reduce(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] 0 "+" REDUCE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(float64)
	if result != 15.0 {
		t.Errorf("Expected 15 (sum), got %v", result)
	}
}

func TestArray_ReduceProduct(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[1 2 3 4 5] 1 "*" REDUCE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(float64)
	if result != 120.0 {
		t.Errorf("Expected 120 (factorial), got %v", result)
	}
}

func TestArray_ReduceEmpty(t *testing.T) {
	interp := setupArrayInterpreter()
	err := interp.Run(`[] 42 "+" REDUCE`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(int64)
	if result != 42 {
		t.Errorf("Expected 42 (initial value), got %v", result)
	}
}
