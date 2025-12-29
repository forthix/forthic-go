package modules

import (
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupRecordInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()

	// Import record module
	recMod := NewRecordModule()
	interp.ImportModule(recMod.Module, "")

	// Import array module for some tests
	arrayMod := NewArrayModule()
	interp.ImportModule(arrayMod.Module, "")

	return interp
}

// ========================================
// Creation
// ========================================

func TestRecord_CreateRecord(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	if result["age"].(int64) != 30 {
		t.Errorf("Expected age=30, got %v", result["age"])
	}
}

func TestRecord_CreateRecordEmpty(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[] REC`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if len(result) != 0 {
		t.Errorf("Expected empty record, got %v", result)
	}
}

func TestRecord_CreateRecordInvalidPairs(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["incomplete"]] REC`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	// Should skip invalid pairs
	if len(result) != 1 {
		t.Errorf("Expected 1 valid key, got %v", result)
	}
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
}

func TestRecord_SetRecordValue(t *testing.T) {
	interp := setupRecordInterpreter()
	// Stack signature: ( rec value field -- rec )
	err := interp.Run(`[["name" "Alice"]] REC 30 "age" <REC!`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	if result["age"].(int64) != 30 {
		t.Errorf("Expected age=30, got %v", result["age"])
	}
}

func TestRecord_SetRecordValueOverwrite(t *testing.T) {
	interp := setupRecordInterpreter()
	// Stack signature: ( rec value field -- rec )
	err := interp.Run(`[["name" "Alice"]] REC "Bob" "name" <REC!`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Bob" {
		t.Errorf("Expected name='Bob' (overwritten), got %v", result["name"])
	}
}

func TestRecord_SetRecordNestedPath(t *testing.T) {
	interp := setupRecordInterpreter()
	// Stack signature: ( rec value field -- rec )
	err := interp.Run(`[["user" [["name" "Alice"]] REC]] REC 30 ["user" "age"] <REC!`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	user := result["user"].(map[string]interface{})
	if user["name"].(string) != "Alice" {
		t.Errorf("Expected user.name='Alice', got %v", user["name"])
	}
	if user["age"].(int64) != 30 {
		t.Errorf("Expected user.age=30, got %v", user["age"])
	}
}

// ========================================
// Access
// ========================================

func TestRecord_GetRecordValue(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC "name" REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "Alice" {
		t.Errorf("Expected 'Alice', got %v", result)
	}
}

func TestRecord_GetRecordValueMissing(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"]] REC "age" REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for missing key, got %v", result)
	}
}

func TestRecord_GetRecordValueNilRecord(t *testing.T) {
	interp := setupRecordInterpreter()
	interp.StackPush(nil)
	interp.StackPush("key")
	err := interp.Run(`REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for nil record, got %v", result)
	}
}

func TestRecord_GetRecordNestedPath(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["user" [["name" "Alice"] ["age" 30]] REC]] REC ["user" "name"] REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "Alice" {
		t.Errorf("Expected 'Alice', got %v", result)
	}
}

func TestRecord_PipeRecAt(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[[["name" "Alice"]] REC [["name" "Bob"]] REC] "name" |REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result))
	}
	if result[0].(string) != "Alice" {
		t.Errorf("Expected 'Alice', got %v", result[0])
	}
	if result[1].(string) != "Bob" {
		t.Errorf("Expected 'Bob', got %v", result[1])
	}
}

func TestRecord_Keys(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30] ["city" "NYC"]] REC KEYS`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 keys, got %d", len(result))
	}
	// Check that keys are present (order is not guaranteed for maps)
	keys := make(map[string]bool)
	for _, k := range result {
		keys[k.(string)] = true
	}
	if !keys["name"] || !keys["age"] || !keys["city"] {
		t.Errorf("Expected keys [name age city], got %v", result)
	}
}

func TestRecord_KeysEmpty(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[] REC KEYS`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 0 {
		t.Errorf("Expected empty array, got %v", result)
	}
}

func TestRecord_Values(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["a" 1] ["b" 2] ["c" 3]] REC VALUES`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 values, got %d", len(result))
	}
	// Check that values are present (order is not guaranteed for maps)
	values := make(map[int64]bool)
	for _, v := range result {
		values[v.(int64)] = true
	}
	if !values[1] || !values[2] || !values[3] {
		t.Errorf("Expected values [1 2 3], got %v", result)
	}
}

func TestRecord_ValuesEmpty(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[] REC VALUES`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 0 {
		t.Errorf("Expected empty array, got %v", result)
	}
}

// ========================================
// Transform
// ========================================

func TestRecord_InvertKeys(t *testing.T) {
	interp := setupRecordInterpreter()
	// Two-level nested structure
	err := interp.Run(`[["A" [["X" 1] ["Y" 2]] REC] ["B" [["X" 3] ["Y" 4]] REC]] REC INVERT-KEYS`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 keys (X, Y), got %d", len(result))
	}

	xRec := result["X"].(map[string]interface{})
	yRec := result["Y"].(map[string]interface{})

	if xRec["A"].(int64) != 1 {
		t.Errorf("Expected X.A=1, got %v", xRec["A"])
	}
	if xRec["B"].(int64) != 3 {
		t.Errorf("Expected X.B=3, got %v", xRec["B"])
	}
	if yRec["A"].(int64) != 2 {
		t.Errorf("Expected Y.A=2, got %v", yRec["A"])
	}
	if yRec["B"].(int64) != 4 {
		t.Errorf("Expected Y.B=4, got %v", yRec["B"])
	}
}

func TestRecord_Relabel(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["a" 1] ["b" 2] ["c" 3]] REC ["a" "b"] ["x" "y"] RELABEL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(result))
	}
	if result["x"].(int64) != 1 {
		t.Errorf("Expected x=1, got %v", result["x"])
	}
	if result["y"].(int64) != 2 {
		t.Errorf("Expected y=2, got %v", result["y"])
	}
}

func TestRecord_RecDefaults(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" ""]] REC [["age" 25] ["city" "NYC"]] REC-DEFAULTS`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	if result["age"].(int64) != 25 {
		t.Errorf("Expected age=25 (default applied), got %v", result["age"])
	}
	if result["city"].(string) != "NYC" {
		t.Errorf("Expected city='NYC' (default applied), got %v", result["city"])
	}
}

func TestRecord_Del(t *testing.T) {
	interp := setupRecordInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30] ["city" "NYC"]] REC "age" <DEL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if len(result) != 2 {
		t.Fatalf("Expected 2 keys after delete, got %d", len(result))
	}
	if _, exists := result["age"]; exists {
		t.Errorf("Expected 'age' to be deleted")
	}
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
}

// ========================================
// Integration Tests
// ========================================

func TestRecord_CreateGetSet(t *testing.T) {
	interp := setupRecordInterpreter()

	// Create record and get name
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC "name" REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	name := interp.StackPop().(string)
	if name != "Alice" {
		t.Errorf("Expected name='Alice', got %v", name)
	}

	// Create record, set age, and get it
	err = interp.Run(`[["name" "Alice"] ["age" 30]] REC 31 "age" <REC! "age" REC@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	age := interp.StackPop().(int64)
	if age != 31 {
		t.Errorf("Expected age=31, got %v", age)
	}
}
