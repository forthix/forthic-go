package modules

import (
	"strings"
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupCoreInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()
	coreMod := NewCoreModule()
	mathMod := NewMathModule()
	interp.ImportModule(coreMod.Module, "")
	interp.ImportModule(mathMod.Module, "")
	return interp
}

// ========================================
// Stack Operations
// ========================================

func TestCore_POP(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("1 2 3 POP")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 2 {
		t.Fatalf("Expected 2 items on stack, got %d", len(items))
	}

	if items[1].(int64) != 2 {
		t.Errorf("Expected top to be 2, got %v", items[1])
	}
}

func TestCore_DUP(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("42 DUP")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 2 {
		t.Fatalf("Expected 2 items on stack, got %d", len(items))
	}

	if items[0].(int64) != 42 || items[1].(int64) != 42 {
		t.Errorf("Expected both items to be 42, got %v and %v", items[0], items[1])
	}
}

func TestCore_SWAP(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("1 2 SWAP")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 2 {
		t.Fatalf("Expected 2 items on stack, got %d", len(items))
	}

	if items[0].(int64) != 2 {
		t.Errorf("Expected bottom to be 2, got %v", items[0])
	}
	if items[1].(int64) != 1 {
		t.Errorf("Expected top to be 1, got %v", items[1])
	}
}

// ========================================
// Variable Operations
// ========================================

func TestCore_VARIABLES(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`["x" "y"] VARIABLES`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	appModule := interp.GetAppModule()
	xVar := appModule.GetVariable("x")
	yVar := appModule.GetVariable("y")

	if xVar == nil {
		t.Error("Expected variable x to be created")
	}
	if yVar == nil {
		t.Error("Expected variable y to be created")
	}
}

func TestCore_InvalidVariableName(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`["__test"] VARIABLES`)
	if err == nil {
		t.Error("Expected error for invalid variable name")
	}
}

func TestCore_SetGetVariables(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`["x"] VARIABLES 24 x ! x @`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(int64) != 24 {
		t.Errorf("Expected 24, got %v", result)
	}
}

func TestCore_BangAt(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`["x"] VARIABLES 42 x !@`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(int64) != 42 {
		t.Errorf("Expected 42 on stack, got %v", result)
	}

	appModule := interp.GetAppModule()
	xVar := appModule.GetVariable("x")
	if xVar.GetValue().(int64) != 42 {
		t.Errorf("Expected variable to be 42, got %v", xVar.GetValue())
	}
}

func TestCore_AutoCreateVariables(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test ! with string variable name (auto-creates variable)
	err := interp.Run(`"hello" "autovar1" !`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	err = interp.Run(`autovar1 @`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result := interp.StackPop()
	if result.(string) != "hello" {
		t.Errorf("Expected 'hello', got %v", result)
	}

	// Verify variable was created
	appModule := interp.GetAppModule()
	autovar1 := appModule.GetVariable("autovar1")
	if autovar1 == nil {
		t.Error("Expected autovar1 to be created")
	}

	// Test @ with string variable name (auto-creates with null)
	err = interp.Run(`"autovar2" @`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result = interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	// Test !@ with string variable name
	err = interp.Run(`"world" "autovar3" !@`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	result = interp.StackPop()
	if result.(string) != "world" {
		t.Errorf("Expected 'world', got %v", result)
	}
}

func TestCore_AutoCreateVariablesValidation(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test that __ prefix variables are rejected
	err := interp.Run(`"value" "__invalid" !`)
	if err == nil {
		t.Error("Expected error for __ prefix variable")
	}

	// Test that validation works for @ as well
	err = interp.Run(`"__invalid2" @`)
	if err == nil {
		t.Error("Expected error for __ prefix variable")
	}

	// Test that validation works for !@ as well
	err = interp.Run(`"value" "__invalid3" !@`)
	if err == nil {
		t.Error("Expected error for __ prefix variable")
	}
}

// ========================================
// Module Operations
// ========================================

func TestCore_EXPORT(t *testing.T) {
	interp := setupCoreInterpreter()

	// Create a module and export some words
	err := interp.Run(`["POP" "DUP"] EXPORT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	// Note: actual export testing would require checking the module's exportable list
	// This is a basic smoke test
}

func TestCore_INTERPRET(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`"5 10 +" INTERPRET`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 15.0 {
		t.Errorf("Expected 15, got %v", result)
	}
}

// ========================================
// Control Flow
// ========================================

func TestCore_IDENTITY(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("42 IDENTITY")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	// Stack should still have 42
	result := interp.StackPop()
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestCore_NOP(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("NOP")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	// Stack should be empty
	stack := interp.GetStack()
	if stack.Length() != 0 {
		t.Errorf("Expected empty stack, got %d items", stack.Length())
	}
}

func TestCore_NULL(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run("NULL")
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestCore_ARRAY_Check(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`[1 2 3] ARRAY?`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(bool) != true {
		t.Error("Expected true for array check")
	}

	err = interp.Run(`42 ARRAY?`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result = interp.StackPop()
	if result.(bool) != false {
		t.Error("Expected false for non-array")
	}
}

func TestCore_DEFAULT(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test with null
	err := interp.Run(`NULL 42 DEFAULT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	// Test with non-null value
	err = interp.Run(`10 42 DEFAULT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result = interp.StackPop()
	if result.(int64) != 10 {
		t.Errorf("Expected 10, got %v", result)
	}

	// Test with empty string
	err = interp.Run(`"" 42 DEFAULT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result = interp.StackPop()
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

func TestCore_DefaultStar(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test with null - should execute Forthic
	err := interp.Run(`NULL "10 20 +" *DEFAULT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 30.0 {
		t.Errorf("Expected 30, got %v", result)
	}

	// Test with non-null value - should not execute Forthic
	err = interp.Run(`42 "10 20 +" *DEFAULT`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result = interp.StackPop()
	if result.(int64) != 42 {
		t.Errorf("Expected 42, got %v", result)
	}
}

// ========================================
// Options
// ========================================

func TestCore_ToOptions(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`[.key1 "value1" .key2 42] ~>`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	opts, ok := result.(*forthic.WordOptions)
	if !ok {
		t.Fatalf("Expected WordOptions, got %T", result)
	}

	val1 := opts.Get("key1", nil)
	if val1.(string) != "value1" {
		t.Errorf("Expected 'value1', got %v", val1)
	}

	val2 := opts.Get("key2", nil)
	if val2.(int64) != 42 {
		t.Errorf("Expected 42, got %v", val2)
	}
}

// ========================================
// String Operations
// ========================================

func TestCore_INTERPOLATE_Basic(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`5 .count ! "Count: .count" INTERPOLATE`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(string) != "Count: 5" {
		t.Errorf("Expected 'Count: 5', got %v", result)
	}
}

func TestCore_INTERPOLATE_WithOptions(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`[1 2 3] .items ! "Items: .items" [.separator " | "] ~> INTERPOLATE`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(string) != "Items: 1 | 2 | 3" {
		t.Errorf("Expected 'Items: 1 | 2 | 3', got %v", result)
	}
}

func TestCore_INTERPOLATE_EscapedDots(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`"Test \\. escaped" INTERPOLATE`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if !strings.Contains(result.(string), ".") {
		t.Errorf("Expected escaped dot to be preserved, got %v", result)
	}
}

func TestCore_INTERPOLATE_NullText(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`NULL .value ! "Value: .value" [.null_text "<empty>"] ~> INTERPOLATE`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(string) != "Value: <empty>" {
		t.Errorf("Expected 'Value: <empty>', got %v", result)
	}
}

// ========================================
// Profiling (Placeholder tests)
// ========================================

func TestCore_Profiling(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test basic profiling operations
	err := interp.Run(`PROFILE-START PROFILE-END`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	err = interp.Run(`"test" PROFILE-TIMESTAMP`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	err = interp.Run(`PROFILE-DATA`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result == nil {
		t.Error("Expected profiling data to be returned")
	}
}

// ========================================
// Logging (Placeholder tests)
// ========================================

func TestCore_Logging(t *testing.T) {
	interp := setupCoreInterpreter()

	// Test basic logging operations
	err := interp.Run(`START-LOG END-LOG`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}
}

// ========================================
// Integration Tests
// ========================================

func TestCore_VariableIntegration(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`
		["x" "y"] VARIABLES
		10 x !
		20 y !
		x @ y @ +
	`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 30.0 {
		t.Errorf("Expected 30, got %v", result)
	}
}

func TestCore_StackManipulation(t *testing.T) {
	interp := setupCoreInterpreter()

	err := interp.Run(`
		1 2 3
		DUP    # Stack: 1 2 3 3
		POP    # Stack: 1 2 3
		SWAP   # Stack: 1 3 2
	`)
	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(items))
	}

	if items[2].(int64) != 2 {
		t.Errorf("Expected top to be 2, got %v", items[2])
	}
	if items[1].(int64) != 3 {
		t.Errorf("Expected middle to be 3, got %v", items[1])
	}
	if items[0].(int64) != 1 {
		t.Errorf("Expected bottom to be 1, got %v", items[0])
	}
}
