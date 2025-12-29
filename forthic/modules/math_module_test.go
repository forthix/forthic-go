package modules

import (
	"math"
	"testing"

	"github.com/forthix/forthic-go/forthic"
)

func setupMathInterpreter() *forthic.Interpreter {
	interp := forthic.NewInterpreter()
	mathMod := NewMathModule()
	interp.ImportModule(mathMod.Module, "")
	return interp
}

func TestMath_Arithmetic(t *testing.T) {
	interp := setupMathInterpreter()

	err := interp.Run(`
		2 4 +
		2 4 -
		2 4 *
		2 4 /
		5 3 MOD
		2.51 ROUND
		[1 2 3] +
		[2 3 4] *
	`)

	if err != nil {
		t.Fatalf("Error running code: %v", err)
	}

	stack := interp.GetStack()
	items := stack.Items()

	if len(items) != 8 {
		t.Fatalf("Expected 8 items on stack, got %d", len(items))
	}

	// 2 + 4 = 6
	if items[0].(float64) != 6.0 {
		t.Errorf("Expected 6.0, got %v", items[0])
	}

	// 2 - 4 = -2
	if items[1].(float64) != -2.0 {
		t.Errorf("Expected -2.0, got %v", items[1])
	}

	// 2 * 4 = 8
	if items[2].(float64) != 8.0 {
		t.Errorf("Expected 8.0, got %v", items[2])
	}

	// 2 / 4 = 0.5
	if items[3].(float64) != 0.5 {
		t.Errorf("Expected 0.5, got %v", items[3])
	}

	// 5 MOD 3 = 2
	if items[4].(int) != 2 {
		t.Errorf("Expected 2, got %v", items[4])
	}

	// 2.51 ROUND = 3
	if items[5].(float64) != 3.0 {
		t.Errorf("Expected 3.0, got %v", items[5])
	}

	// [1 2 3] + = 6
	if items[6].(float64) != 6.0 {
		t.Errorf("Expected 6.0, got %v", items[6])
	}

	// [2 3 4] * = 24
	if items[7].(float64) != 24.0 {
		t.Errorf("Expected 24.0, got %v", items[7])
	}
}

func TestMath_Divide(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(10)
	interp.StackPush(2)
	err := interp.Run("DIVIDE")
	if err != nil {
		t.Fatalf("Error running DIVIDE: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 5.0 {
		t.Errorf("Expected 5.0, got %v", result)
	}
}

func TestMath_Mean(t *testing.T) {
	interp := setupMathInterpreter()

	// Mean of [1,2,3,4,5] = 3
	err := interp.Run("[1 2 3 4 5] MEAN")
	if err != nil {
		t.Fatalf("Error running MEAN: %v", err)
	}
	result := interp.StackPop()
	if result.(float64) != 3.0 {
		t.Errorf("Expected 3.0, got %v", result)
	}

	// Mean of [4] = 4
	err = interp.Run("[4] MEAN")
	if err != nil {
		t.Fatalf("Error running MEAN: %v", err)
	}
	result = interp.StackPop()
	if result.(float64) != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
	}

	// Mean of [] = 0
	err = interp.Run("[] MEAN")
	if err != nil {
		t.Fatalf("Error running MEAN: %v", err)
	}
	result = interp.StackPop()
	if result.(float64) != 0.0 {
		t.Errorf("Expected 0.0, got %v", result)
	}
}

func TestMath_MaxTwoNumbers(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(4.0)
	interp.StackPush(18.0)
	err := interp.Run("MAX")
	if err != nil {
		t.Fatalf("Error running MAX: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 18.0 {
		t.Errorf("Expected 18.0, got %v", result)
	}
}

func TestMath_MaxArray(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush([]interface{}{14.0, 8.0, 55.0, 4.0, 5.0})
	err := interp.Run("MAX")
	if err != nil {
		t.Fatalf("Error running MAX: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 55.0 {
		t.Errorf("Expected 55.0, got %v", result)
	}
}

func TestMath_MinTwoNumbers(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(4.0)
	interp.StackPush(18.0)
	err := interp.Run("MIN")
	if err != nil {
		t.Fatalf("Error running MIN: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
	}
}

func TestMath_MinArray(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush([]interface{}{14.0, 8.0, 55.0, 4.0, 5.0})
	err := interp.Run("MIN")
	if err != nil {
		t.Fatalf("Error running MIN: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
	}
}

func TestMath_MeanArray(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush([]interface{}{1.0, 2.0, 3.0, 4.0, 5.0})
	err := interp.Run("MEAN")
	if err != nil {
		t.Fatalf("Error running MEAN: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 3.0 {
		t.Errorf("Expected 3.0, got %v", result)
	}
}

func TestMath_MeanStringArray(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush([]interface{}{"a", "a", "b", "c"})
	err := interp.Run("MEAN")
	if err != nil {
		t.Fatalf("Error running MEAN: %v", err)
	}

	result := interp.StackPop()
	resultMap := result.(map[string]float64)

	if resultMap["a"] != 0.5 {
		t.Errorf("Expected a=0.5, got %v", resultMap["a"])
	}
	if resultMap["b"] != 0.25 {
		t.Errorf("Expected b=0.25, got %v", resultMap["b"])
	}
	if resultMap["c"] != 0.25 {
		t.Errorf("Expected c=0.25, got %v", resultMap["c"])
	}
}

func TestMath_Abs(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(-5.0)
	err := interp.Run("ABS")
	if err != nil {
		t.Fatalf("Error running ABS: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 5.0 {
		t.Errorf("Expected 5.0, got %v", result)
	}
}

func TestMath_Sqrt(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(16.0)
	err := interp.Run("SQRT")
	if err != nil {
		t.Fatalf("Error running SQRT: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
	}
}

func TestMath_Floor(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(3.7)
	err := interp.Run("FLOOR")
	if err != nil {
		t.Fatalf("Error running FLOOR: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 3.0 {
		t.Errorf("Expected 3.0, got %v", result)
	}
}

func TestMath_Ceil(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(3.2)
	err := interp.Run("CEIL")
	if err != nil {
		t.Fatalf("Error running CEIL: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 4.0 {
		t.Errorf("Expected 4.0, got %v", result)
	}
}

func TestMath_Clamp(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(5.0) // value
	interp.StackPush(0.0) // min
	interp.StackPush(10.0) // max
	err := interp.Run("CLAMP")
	if err != nil {
		t.Fatalf("Error running CLAMP: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 5.0 {
		t.Errorf("Expected 5.0, got %v", result)
	}

	// Test lower bound
	interp.StackPush(-5.0) // value
	interp.StackPush(0.0)  // min
	interp.StackPush(10.0) // max
	err = interp.Run("CLAMP")
	if err != nil {
		t.Fatalf("Error running CLAMP: %v", err)
	}

	result = interp.StackPop()
	if result.(float64) != 0.0 {
		t.Errorf("Expected 0.0, got %v", result)
	}

	// Test upper bound
	interp.StackPush(15.0) // value
	interp.StackPush(0.0)  // min
	interp.StackPush(10.0) // max
	err = interp.Run("CLAMP")
	if err != nil {
		t.Fatalf("Error running CLAMP: %v", err)
	}

	result = interp.StackPop()
	if result.(float64) != 10.0 {
		t.Errorf("Expected 10.0, got %v", result)
	}
}

func TestMath_Infinity(t *testing.T) {
	interp := setupMathInterpreter()
	err := interp.Run("INFINITY")
	if err != nil {
		t.Fatalf("Error running INFINITY: %v", err)
	}

	result := interp.StackPop()
	if !math.IsInf(result.(float64), 1) {
		t.Errorf("Expected +Inf, got %v", result)
	}
}

func TestMath_ToInt(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(3.7)
	err := interp.Run(">INT")
	if err != nil {
		t.Fatalf("Error running >INT: %v", err)
	}

	result := interp.StackPop()
	if result.(int) != 3 {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestMath_ToFloat(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(3)
	err := interp.Run(">FLOAT")
	if err != nil {
		t.Fatalf("Error running >FLOAT: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 3.0 {
		t.Errorf("Expected 3.0, got %v", result)
	}
}

func TestMath_ToFixed(t *testing.T) {
	interp := setupMathInterpreter()
	interp.StackPush(3.14159)
	interp.StackPush(2.0)
	err := interp.Run(">FIXED")
	if err != nil {
		t.Fatalf("Error running >FIXED: %v", err)
	}

	result := interp.StackPop()
	if result.(float64) != 3.14 {
		t.Errorf("Expected 3.14, got %v", result)
	}
}
