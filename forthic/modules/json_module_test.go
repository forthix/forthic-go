package modules

import (
	"strings"
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupJSONInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()

	// Import JSON module
	jsonMod := NewJSONModule()
	interp.ImportModule(jsonMod.Module, "")

	// Import record module for creating test data
	recMod := NewRecordModule()
	interp.ImportModule(recMod.Module, "")

	return interp
}

// ========================================
// Encoding
// ========================================

func TestJSON_ToJSONString(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"hello" >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != `"hello"` {
		t.Errorf("Expected '\"hello\"', got %v", result)
	}
}

func TestJSON_ToJSONNumber(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`42 >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "42" {
		t.Errorf("Expected '42', got %v", result)
	}
}

func TestJSON_ToJSONBoolean(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`TRUE >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "true" {
		t.Errorf("Expected 'true', got %v", result)
	}
}

func TestJSON_ToJSONNull(t *testing.T) {
	interp := setupJSONInterpreter()
	interp.StackPush(nil)
	err := interp.Run(`>JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "null" {
		t.Errorf("Expected 'null', got %v", result)
	}
}

func TestJSON_ToJSONArray(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[1 2 3] >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "[1,2,3]" {
		t.Errorf("Expected '[1,2,3]', got %v", result)
	}
}

func TestJSON_ToJSONRecord(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	// JSON key order is not guaranteed, check both possibilities
	if result != `{"age":30,"name":"Alice"}` && result != `{"name":"Alice","age":30}` {
		t.Errorf("Expected JSON object, got %v", result)
	}
}

func TestJSON_ToJSONNestedStructure(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[["users" [[["name" "Alice"]] REC [["name" "Bob"]] REC]]] REC >JSON`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	// Should contain nested structure
	if !strings.Contains(result, "users") || !strings.Contains(result, "Alice") {
		t.Errorf("Expected nested JSON structure, got %v", result)
	}
}

// ========================================
// Prettify
// ========================================

func TestJSON_Prettify(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC JSON-PRETTIFY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	// Should contain newlines and indentation
	if !strings.Contains(result, "\n") || !strings.Contains(result, "  ") {
		t.Errorf("Expected prettified JSON with newlines, got %v", result)
	}
}

func TestJSON_PrettifyArray(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[1 2 3 4 5] JSON-PRETTIFY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	// Array should be formatted
	if !strings.Contains(result, "\n") {
		t.Errorf("Expected prettified array, got %v", result)
	}
}

func TestJSON_PrettifyNull(t *testing.T) {
	interp := setupJSONInterpreter()
	interp.StackPush(nil)
	err := interp.Run(`JSON-PRETTIFY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "null" {
		t.Errorf("Expected 'null', got %v", result)
	}
}

// ========================================
// Decoding
// ========================================

func TestJSON_FromJSONString(t *testing.T) {
	interp := setupJSONInterpreter()
	// Push JSON string directly to avoid escaping issues
	interp.StackPush(`"hello"`)
	err := interp.Run(`JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}
}

func TestJSON_FromJSONNumber(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"42" JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(float64)
	if result != 42.0 {
		t.Errorf("Expected 42.0, got %v", result)
	}
}

func TestJSON_FromJSONBoolean(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"true" JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(bool)
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestJSON_FromJSONNull(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"null" JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestJSON_FromJSONArray(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"[1,2,3]" JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
	if result[0].(float64) != 1.0 || result[2].(float64) != 3.0 {
		t.Errorf("Expected [1 2 3], got %v", result)
	}
}

func TestJSON_FromJSONObject(t *testing.T) {
	interp := setupJSONInterpreter()
	// Push JSON string directly to avoid escaping issues
	interp.StackPush(`{"name":"Alice","age":30}`)
	err := interp.Run(`JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	if result["age"].(float64) != 30.0 {
		t.Errorf("Expected age=30, got %v", result["age"])
	}
}

func TestJSON_FromJSONNestedStructure(t *testing.T) {
	interp := setupJSONInterpreter()
	// Push JSON string directly to avoid escaping issues
	interp.StackPush(`{"users":[{"name":"Alice"},{"name":"Bob"}]}`)
	err := interp.Run(`JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	users := result["users"].([]interface{})
	if len(users) != 2 {
		t.Fatalf("Expected 2 users, got %d", len(users))
	}
	firstUser := users[0].(map[string]interface{})
	if firstUser["name"].(string) != "Alice" {
		t.Errorf("Expected first user name='Alice', got %v", firstUser["name"])
	}
}

func TestJSON_FromJSONInvalid(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"invalid json{" JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for invalid JSON, got %v", result)
	}
}

func TestJSON_FromJSONNil(t *testing.T) {
	interp := setupJSONInterpreter()
	interp.StackPush(nil)
	err := interp.Run(`JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

// ========================================
// Round-trip Tests
// ========================================

func TestJSON_RoundTripString(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`"hello world" >JSON JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(string)
	if result != "hello world" {
		t.Errorf("Expected 'hello world', got %v", result)
	}
}

func TestJSON_RoundTripNumber(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`42 >JSON JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(float64)
	if result != 42.0 {
		t.Errorf("Expected 42.0, got %v", result)
	}
}

func TestJSON_RoundTripArray(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[1 2 3] >JSON JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().([]interface{})
	if len(result) != 3 {
		t.Fatalf("Expected 3 elements, got %d", len(result))
	}
}

func TestJSON_RoundTripRecord(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[["name" "Alice"] ["age" 30]] REC >JSON JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice', got %v", result["name"])
	}
	if result["age"].(float64) != 30.0 {
		t.Errorf("Expected age=30.0, got %v", result["age"])
	}
}

func TestJSON_PrettifyThenParse(t *testing.T) {
	interp := setupJSONInterpreter()
	err := interp.Run(`[["name" "Alice"]] REC JSON-PRETTIFY JSON>`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	result := interp.StackPop().(map[string]interface{})
	if result["name"].(string) != "Alice" {
		t.Errorf("Expected name='Alice' after prettify+parse, got %v", result["name"])
	}
}
