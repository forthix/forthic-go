package modules

import (
	"github.com/forthix/forthic-go/forthic"
)

// BooleanModule provides boolean and comparison operations
type BooleanModule struct {
	*forthic.Module
}

// NewBooleanModule creates a new boolean module
func NewBooleanModule() *BooleanModule {
	m := &BooleanModule{
		Module: forthic.NewModule("boolean", ""),
	}
	m.registerWords()
	return m
}

func (m *BooleanModule) registerWords() {
	// Comparison operations
	m.AddModuleWord("==", m.equals)
	m.AddModuleWord("!=", m.notEquals)
	m.AddModuleWord("<", m.lessThan)
	m.AddModuleWord("<=", m.lessThanOrEqual)
	m.AddModuleWord(">", m.greaterThan)
	m.AddModuleWord(">=", m.greaterThanOrEqual)

	// Logic operations
	m.AddModuleWord("OR", m.or)
	m.AddModuleWord("AND", m.and)
	m.AddModuleWord("NOT", m.not)
	m.AddModuleWord("XOR", m.xor)
	m.AddModuleWord("NAND", m.nand)

	// Membership operations
	m.AddModuleWord("IN", m.in)
	m.AddModuleWord("ANY", m.any)
	m.AddModuleWord("ALL", m.all)

	// Type conversion
	m.AddModuleWord(">BOOL", m.toBool)
}

// ========================================
// Comparison Operations
// ========================================

func (m *BooleanModule) equals(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(areEqual(a, b))
	return nil
}

func (m *BooleanModule) notEquals(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(!areEqual(a, b))
	return nil
}

func (m *BooleanModule) lessThan(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(compare(a, b) < 0)
	return nil
}

func (m *BooleanModule) lessThanOrEqual(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(compare(a, b) <= 0)
	return nil
}

func (m *BooleanModule) greaterThan(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(compare(a, b) > 0)
	return nil
}

func (m *BooleanModule) greaterThanOrEqual(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(compare(a, b) >= 0)
	return nil
}

// ========================================
// Logic Operations
// ========================================

func (m *BooleanModule) or(interp *forthic.Interpreter) error {
	b := interp.StackPop()

	// Case 1: Array on top of stack
	if arr, ok := b.([]interface{}); ok {
		for _, val := range arr {
			if isTruthy(val) {
				interp.StackPush(true)
				return nil
			}
		}
		interp.StackPush(false)
		return nil
	}

	// Case 2: Two values
	a := interp.StackPop()
	interp.StackPush(isTruthy(a) || isTruthy(b))
	return nil
}

func (m *BooleanModule) and(interp *forthic.Interpreter) error {
	b := interp.StackPop()

	// Case 1: Array on top of stack
	if arr, ok := b.([]interface{}); ok {
		for _, val := range arr {
			if !isTruthy(val) {
				interp.StackPush(false)
				return nil
			}
		}
		interp.StackPush(true)
		return nil
	}

	// Case 2: Two values
	a := interp.StackPop()
	interp.StackPush(isTruthy(a) && isTruthy(b))
	return nil
}

func (m *BooleanModule) not(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	interp.StackPush(!isTruthy(val))
	return nil
}

func (m *BooleanModule) xor(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	aBool := isTruthy(a)
	bBool := isTruthy(b)
	interp.StackPush((aBool || bBool) && !(aBool && bBool))
	return nil
}

func (m *BooleanModule) nand(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(!(isTruthy(a) && isTruthy(b)))
	return nil
}

// ========================================
// Membership Operations
// ========================================

func (m *BooleanModule) in(interp *forthic.Interpreter) error {
	arr := interp.StackPop()
	item := interp.StackPop()

	if slice, ok := arr.([]interface{}); ok {
		for _, val := range slice {
			if areEqual(val, item) {
				interp.StackPush(true)
				return nil
			}
		}
	}

	interp.StackPush(false)
	return nil
}

func (m *BooleanModule) any(interp *forthic.Interpreter) error {
	items2 := interp.StackPop()
	items1 := interp.StackPop()

	slice1, ok1 := items1.([]interface{})
	slice2, ok2 := items2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush(false)
		return nil
	}

	// If items2 is empty, return true
	if len(slice2) == 0 {
		interp.StackPush(true)
		return nil
	}

	// Check if any item from items1 is in items2
	for _, item1 := range slice1 {
		for _, item2 := range slice2 {
			if areEqual(item1, item2) {
				interp.StackPush(true)
				return nil
			}
		}
	}

	interp.StackPush(false)
	return nil
}

func (m *BooleanModule) all(interp *forthic.Interpreter) error {
	items2 := interp.StackPop()
	items1 := interp.StackPop()

	slice1, ok1 := items1.([]interface{})
	slice2, ok2 := items2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush(false)
		return nil
	}

	// If items2 is empty, return true
	if len(slice2) == 0 {
		interp.StackPush(true)
		return nil
	}

	// Check if all items from items2 are in items1
	for _, item2 := range slice2 {
		found := false
		for _, item1 := range slice1 {
			if areEqual(item1, item2) {
				found = true
				break
			}
		}
		if !found {
			interp.StackPush(false)
			return nil
		}
	}

	interp.StackPush(true)
	return nil
}

// ========================================
// Type Conversion
// ========================================

func (m *BooleanModule) toBool(interp *forthic.Interpreter) error {
	val := interp.StackPop()
	interp.StackPush(isTruthy(val))
	return nil
}

// ========================================
// Helper Functions
// ========================================

func compare(a, b interface{}) int {
	// Try numeric comparison first
	aNum, aErr := toNumber(a)
	bNum, bErr := toNumber(b)
	if aErr == nil && bErr == nil {
		if aNum < bNum {
			return -1
		} else if aNum > bNum {
			return 1
		}
		return 0
	}

	// Try string comparison
	aStr, aOk := a.(string)
	bStr, bOk := b.(string)
	if aOk && bOk {
		if aStr < bStr {
			return -1
		} else if aStr > bStr {
			return 1
		}
		return 0
	}

	// Default: equal
	return 0
}
