package modules

import (
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupBooleanInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()
	boolMod := NewBooleanModule()
	interp.ImportModule(boolMod.Module, "")
	return interp
}

// ========================================
// Comparison Operations
// ========================================

func TestBoolean_Comparison(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`
		2 4 ==
		2 4 !=
		2 4 <
		2 4 <=
		2 4 >
		2 4 >=
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 6 {
		t.Fatalf("Expected 6 items on stack, got %d", len(items))
	}

	// 2 == 4 -> false
	if items[0].(bool) != false {
		t.Errorf("Expected false, got %v", items[0])
	}

	// 2 != 4 -> true
	if items[1].(bool) != true {
		t.Errorf("Expected true, got %v", items[1])
	}

	// 2 < 4 -> true
	if items[2].(bool) != true {
		t.Errorf("Expected true, got %v", items[2])
	}

	// 2 <= 4 -> true
	if items[3].(bool) != true {
		t.Errorf("Expected true, got %v", items[3])
	}

	// 2 > 4 -> false
	if items[4].(bool) != false {
		t.Errorf("Expected false, got %v", items[4])
	}

	// 2 >= 4 -> false
	if items[5].(bool) != false {
		t.Errorf("Expected false, got %v", items[5])
	}
}

func TestBoolean_EqualityDifferentTypes(t *testing.T) {
	interp := setupBooleanInterpreter()

	// Numbers
	err := interp.Run(`2 2 ==`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected 2 == 2 to be true")
	}

	// Strings
	err = interp.Run(`"hello" "hello" ==`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected strings to be equal")
	}

	// Booleans
	err = interp.Run(`TRUE TRUE ==`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected TRUE == TRUE to be true")
	}
}

// ========================================
// Logic Operations
// ========================================

func TestBoolean_Logic(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`
		FALSE FALSE OR
		[FALSE FALSE TRUE FALSE] OR
		FALSE TRUE AND
		[FALSE FALSE TRUE FALSE] AND
		FALSE NOT
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 5 {
		t.Fatalf("Expected 5 items on stack, got %d", len(items))
	}

	// FALSE OR FALSE -> false
	if items[0].(bool) != false {
		t.Errorf("Expected false, got %v", items[0])
	}

	// [FALSE FALSE TRUE FALSE] OR -> true
	if items[1].(bool) != true {
		t.Errorf("Expected true, got %v", items[1])
	}

	// FALSE AND TRUE -> false
	if items[2].(bool) != false {
		t.Errorf("Expected false, got %v", items[2])
	}

	// [FALSE FALSE TRUE FALSE] AND -> false
	if items[3].(bool) != false {
		t.Errorf("Expected false, got %v", items[3])
	}

	// NOT FALSE -> true
	if items[4].(bool) != true {
		t.Errorf("Expected true, got %v", items[4])
	}
}

func TestBoolean_ORTwoValues(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"TRUE FALSE OR", true},
		{"FALSE FALSE OR", false},
		{"TRUE TRUE OR", true},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

func TestBoolean_ORArray(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"[FALSE FALSE FALSE] OR", false},
		{"[TRUE FALSE FALSE] OR", true},
		{"[FALSE TRUE FALSE] OR", true},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

func TestBoolean_ANDTwoValues(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"TRUE TRUE AND", true},
		{"TRUE FALSE AND", false},
		{"FALSE FALSE AND", false},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

func TestBoolean_ANDArray(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"[TRUE TRUE TRUE] AND", true},
		{"[TRUE FALSE TRUE] AND", false},
		{"[FALSE FALSE FALSE] AND", false},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

func TestBoolean_NOT(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`TRUE NOT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected NOT TRUE to be false")
	}

	err = interp.Run(`FALSE NOT`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected NOT FALSE to be true")
	}
}

func TestBoolean_XOR(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"TRUE TRUE XOR", false},
		{"TRUE FALSE XOR", true},
		{"FALSE TRUE XOR", true},
		{"FALSE FALSE XOR", false},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

func TestBoolean_NAND(t *testing.T) {
	interp := setupBooleanInterpreter()

	tests := []struct {
		code     string
		expected bool
	}{
		{"TRUE TRUE NAND", false},
		{"TRUE FALSE NAND", true},
		{"FALSE TRUE NAND", true},
		{"FALSE FALSE NAND", true},
	}

	for _, tt := range tests {
		err := interp.Run(tt.code)
		if err != nil {
			t.Fatalf("Error running %s: %v", tt.code, err)
		}
		result := interp.StackPop().(bool)
		if result != tt.expected {
			t.Errorf("%s: expected %v, got %v", tt.code, tt.expected, result)
		}
	}
}

// ========================================
// Membership Operations
// ========================================

func TestBoolean_IN(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`
		"alpha" ["beta" "gamma"] IN
		"alpha" ["beta" "gamma" "alpha"] IN
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 2 {
		t.Fatalf("Expected 2 items on stack, got %d", len(items))
	}

	if items[0].(bool) != false {
		t.Errorf("Expected false, got %v", items[0])
	}

	if items[1].(bool) != true {
		t.Errorf("Expected true, got %v", items[1])
	}
}

func TestBoolean_INNumbers(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`5 [1 2 3 4 5] IN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected 5 to be in array")
	}

	err = interp.Run(`10 [1 2 3 4 5] IN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected 10 not to be in array")
	}
}

func TestBoolean_INEmptyArray(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`"test" [] IN`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected false for empty array")
	}
}

func TestBoolean_ANY(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`
		["alpha" "beta"] ["beta" "gamma"] ANY
		["delta" "beta"] ["gamma" "alpha"] ANY
		["alpha" "beta"] [] ANY
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 3 {
		t.Fatalf("Expected 3 items on stack, got %d", len(items))
	}

	if items[0].(bool) != true {
		t.Errorf("Expected true, got %v", items[0])
	}

	if items[1].(bool) != false {
		t.Errorf("Expected false, got %v", items[1])
	}

	if items[2].(bool) != true {
		t.Errorf("Expected true (empty array), got %v", items[2])
	}
}

func TestBoolean_ANYNumbers(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`[1 2 3] [3 4 5] ANY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected true (3 is in both)")
	}

	err = interp.Run(`[1 2 3] [4 5 6] ANY`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected false (no overlap)")
	}
}

func TestBoolean_ALL(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`
		["alpha" "beta"] ["beta" "gamma"] ALL
		["delta" "beta"] ["beta"] ALL
		["alpha" "beta"] [] ALL
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 3 {
		t.Fatalf("Expected 3 items on stack, got %d", len(items))
	}

	if items[0].(bool) != false {
		t.Errorf("Expected false, got %v", items[0])
	}

	if items[1].(bool) != true {
		t.Errorf("Expected true, got %v", items[1])
	}

	if items[2].(bool) != true {
		t.Errorf("Expected true (empty array), got %v", items[2])
	}
}

func TestBoolean_ALLNumbers(t *testing.T) {
	interp := setupBooleanInterpreter()

	err := interp.Run(`[1 2 3 4 5] [2 3 4] ALL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected true (all of [2,3,4] are in [1,2,3,4,5])")
	}

	err = interp.Run(`[1 2 3] [2 3 4] ALL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected false (4 not in [1,2,3])")
	}
}

// ========================================
// Type Conversion
// ========================================

func TestBoolean_ToBool(t *testing.T) {
	interp := setupBooleanInterpreter()

	// NULL -> false (push nil directly)
	interp.StackPush(nil)
	err := interp.Run(">BOOL")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected nil to be false")
	}

	// 0 -> false
	err = interp.Run("0 >BOOL")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected 0 to be false")
	}

	// 1 -> true
	err = interp.Run("1 >BOOL")
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected 1 to be true")
	}

	// "" -> false
	err = interp.Run(`"" >BOOL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != false {
		t.Error("Expected empty string to be false")
	}

	// "Hi" -> true
	err = interp.Run(`"Hi" >BOOL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected non-empty string to be true")
	}
}

func TestBoolean_ToBoolArrays(t *testing.T) {
	interp := setupBooleanInterpreter()

	// Empty arrays are truthy (JavaScript behavior)
	err := interp.Run(`[] >BOOL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected empty array to be truthy")
	}

	err = interp.Run(`[1] >BOOL`)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if interp.StackPop().(bool) != true {
		t.Error("Expected non-empty array to be truthy")
	}
}
