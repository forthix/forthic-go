package modules

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/forthix/forthic-go/forthic"
)

// MathModule provides mathematical operations
type MathModule struct {
	*forthic.Module
}

// NewMathModule creates a new math module
func NewMathModule() *MathModule {
	m := &MathModule{
		Module: forthic.NewModule("math", ""),
	}
	m.registerWords()
	return m
}

func (m *MathModule) registerWords() {
	// Arithmetic Operations
	m.AddModuleWord("+", m.plus)
	m.AddModuleWord("ADD", m.plus)
	m.AddModuleWord("-", m.minus)
	m.AddModuleWord("SUBTRACT", m.minus)
	m.AddModuleWord("*", m.times)
	m.AddModuleWord("MULTIPLY", m.times)
	m.AddModuleWord("/", m.divide)
	m.AddModuleWord("DIVIDE", m.divide)
	m.AddModuleWord("MOD", m.mod)

	// Aggregates
	m.AddModuleWord("SUM", m.sum)
	m.AddModuleWord("MEAN", m.mean)
	m.AddModuleWord("MAX", m.max_word)
	m.AddModuleWord("MIN", m.min_word)

	// Type conversion
	m.AddModuleWord(">INT", m.toInt)
	m.AddModuleWord(">FLOAT", m.toFloat)
	m.AddModuleWord("ROUND", m.round)
	m.AddModuleWord(">FIXED", m.toFixed)

	// Math functions
	m.AddModuleWord("ABS", m.abs)
	m.AddModuleWord("SQRT", m.sqrt)
	m.AddModuleWord("FLOOR", m.floor)
	m.AddModuleWord("CEIL", m.ceil)
	m.AddModuleWord("CLAMP", m.clamp)

	// Special values
	m.AddModuleWord("INFINITY", m.infinity)
	m.AddModuleWord("UNIFORM-RANDOM", m.uniformRandom)
}

// ========================================
// Arithmetic Operations
// ========================================

func (m *MathModule) plus(interp *forthic.Interpreter) error {
	b := interp.StackPop()

	// Case 1: Array on top of stack
	if arr, ok := b.([]interface{}); ok {
		result := 0.0
		for _, val := range arr {
			if val != nil {
				if num, err := toNumber(val); err == nil {
					result += num
				}
			}
		}
		interp.StackPush(result)
		return nil
	}

	// Case 2: Two numbers
	a := interp.StackPop()
	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(0.0)
		return nil
	}

	interp.StackPush(numA + numB)
	return nil
}

func (m *MathModule) minus(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()

	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	interp.StackPush(numA - numB)
	return nil
}

func (m *MathModule) times(interp *forthic.Interpreter) error {
	b := interp.StackPop()

	// Case 1: Array on top of stack
	if arr, ok := b.([]interface{}); ok {
		result := 1.0
		for _, val := range arr {
			if val == nil {
				interp.StackPush(nil)
				return nil
			}
			if num, err := toNumber(val); err == nil {
				result *= num
			} else {
				interp.StackPush(nil)
				return nil
			}
		}
		interp.StackPush(result)
		return nil
	}

	// Case 2: Two numbers
	a := interp.StackPop()
	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	interp.StackPush(numA * numB)
	return nil
}

func (m *MathModule) divide(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()

	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	if numB == 0 {
		interp.StackPush(math.Inf(1))
		return nil
	}

	interp.StackPush(numA / numB)
	return nil
}

func (m *MathModule) mod(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()

	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	interp.StackPush(int(numA) % int(numB))
	return nil
}

// ========================================
// Aggregates
// ========================================

func (m *MathModule) sum(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush(0.0)
		return nil
	}

	if slice, ok := arr.([]interface{}); ok {
		result := 0.0
		for _, val := range slice {
			if val != nil {
				if num, err := toNumber(val); err == nil {
					result += num
				}
			}
		}
		interp.StackPush(result)
		return nil
	}

	interp.StackPush(0.0)
	return nil
}

func (m *MathModule) mean(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush(0.0)
		return nil
	}

	if slice, ok := arr.([]interface{}); ok {
		if len(slice) == 0 {
			interp.StackPush(0.0)
			return nil
		}

		// Check if it's an array of strings
		if isStringArray(slice) {
			result := computeStringMean(slice)
			interp.StackPush(result)
			return nil
		}

		// Numeric mean
		sum := 0.0
		count := 0
		for _, val := range slice {
			if val != nil {
				if num, err := toNumber(val); err == nil {
					sum += num
					count++
				}
			}
		}

		if count == 0 {
			interp.StackPush(0.0)
			return nil
		}

		interp.StackPush(sum / float64(count))
		return nil
	}

	interp.StackPush(0.0)
	return nil
}

func (m *MathModule) max_word(interp *forthic.Interpreter) error {
	val := interp.StackPop()

	// Case 1: Array
	if arr, ok := val.([]interface{}); ok {
		if len(arr) == 0 {
			interp.StackPush(nil)
			return nil
		}

		max, err := toNumber(arr[0])
		if err != nil {
			interp.StackPush(nil)
			return nil
		}

		for i := 1; i < len(arr); i++ {
			if num, err := toNumber(arr[i]); err == nil {
				if num > max {
					max = num
				}
			}
		}

		interp.StackPush(max)
		return nil
	}

	// Case 2: Two numbers on stack
	b := val
	a := interp.StackPop()

	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	if numA > numB {
		interp.StackPush(numA)
	} else {
		interp.StackPush(numB)
	}
	return nil
}

func (m *MathModule) min_word(interp *forthic.Interpreter) error {
	val := interp.StackPop()

	// Case 1: Array
	if arr, ok := val.([]interface{}); ok {
		if len(arr) == 0 {
			interp.StackPush(nil)
			return nil
		}

		min, err := toNumber(arr[0])
		if err != nil {
			interp.StackPush(nil)
			return nil
		}

		for i := 1; i < len(arr); i++ {
			if num, err := toNumber(arr[i]); err == nil {
				if num < min {
					min = num
				}
			}
		}

		interp.StackPush(min)
		return nil
	}

	// Case 2: Two numbers on stack
	b := val
	a := interp.StackPop()

	numA, errA := toNumber(a)
	numB, errB := toNumber(b)

	if errA != nil || errB != nil {
		interp.StackPush(nil)
		return nil
	}

	if numA < numB {
		interp.StackPush(numA)
	} else {
		interp.StackPush(numB)
	}
	return nil
}

// ========================================
// Type Conversion
// ========================================

func (m *MathModule) toInt(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(0)
		return nil
	}
	interp.StackPush(int(num))
	return nil
}

func (m *MathModule) toFloat(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(0.0)
		return nil
	}
	interp.StackPush(num)
	return nil
}

func (m *MathModule) round(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}
	interp.StackPush(math.Round(num))
	return nil
}

func (m *MathModule) toFixed(interp *forthic.Interpreter) error {
	decimals := interp.StackPop()
	val := interp.StackPop()

	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}

	dec, err := toNumber(decimals)
	if err != nil {
		dec = 0
	}

	multiplier := math.Pow(10, dec)
	result := math.Round(num*multiplier) / multiplier
	interp.StackPush(result)
	return nil
}

// ========================================
// Math Functions
// ========================================

func (m *MathModule) abs(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}
	interp.StackPush(math.Abs(num))
	return nil
}

func (m *MathModule) sqrt(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}
	interp.StackPush(math.Sqrt(num))
	return nil
}

func (m *MathModule) floor(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}
	interp.StackPush(math.Floor(num))
	return nil
}

func (m *MathModule) ceil(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	num, err := toNumber(val)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}
	interp.StackPush(math.Ceil(num))
	return nil
}

func (m *MathModule) clamp(interp *forthic.Interpreter) error {
	max := interp.StackPop()
	min := interp.StackPop()
	val := interp.StackPop()

	numVal, err1 := toNumber(val)
	numMin, err2 := toNumber(min)
	numMax, err3 := toNumber(max)

	if err1 != nil || err2 != nil || err3 != nil {
		interp.StackPush(nil)
		return nil
	}

	if numVal < numMin {
		interp.StackPush(numMin)
	} else if numVal > numMax {
		interp.StackPush(numMax)
	} else {
		interp.StackPush(numVal)
	}
	return nil
}

// ========================================
// Special Values
// ========================================

func (m *MathModule) infinity(interp *forthic.Interpreter) error {
	interp.StackPush(math.Inf(1))
	return nil
}

func (m *MathModule) uniformRandom(interp *forthic.Interpreter) error {
	max := interp.StackPop()
	min := interp.StackPop()

	numMin, err1 := toNumber(min)
	numMax, err2 := toNumber(max)

	if err1 != nil || err2 != nil {
		interp.StackPush(0.0)
		return nil
	}

	result := numMin + rand.Float64()*(numMax-numMin)
	interp.StackPush(result)
	return nil
}

// ========================================
// Helper Functions
// ========================================

func toNumber(val interface{}) (float64, error) {
	if val == nil {
		return 0, fmt.Errorf("nil value")
	}

	switch v := val.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("not a number: %T", val)
	}
}

func isStringArray(arr []interface{}) bool {
	for _, val := range arr {
		if val == nil {
			continue
		}
		if _, ok := val.(string); !ok {
			return false
		}
	}
	return true
}

func computeStringMean(arr []interface{}) map[string]float64 {
	counts := make(map[string]int)
	total := 0

	for _, val := range arr {
		if str, ok := val.(string); ok {
			counts[str]++
			total++
		}
	}

	result := make(map[string]float64)
	for key, count := range counts {
		result[key] = float64(count) / float64(total)
	}

	return result
}
