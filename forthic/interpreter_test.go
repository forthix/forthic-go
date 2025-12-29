package forthic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Interpreter Core Tests (17 tests from Phase 3 plan)
// ============================================================================

func TestInterpreter_InitialState(t *testing.T) {
	interp := NewInterpreter()
	assert.Equal(t, 0, interp.GetStack().Length())
	assert.Equal(t, "", interp.CurModule().GetName())
}

func TestInterpreter_PushString(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`"hello"`)
	assert.NoError(t, err)
	assert.Equal(t, 1, interp.GetStack().Length())
	assert.Equal(t, "hello", interp.StackPop())
}

func TestInterpreter_Comment(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("# This is a comment")
	assert.NoError(t, err)
	assert.Equal(t, 0, interp.GetStack().Length())

	// Test comment with code
	interp2 := NewInterpreter()
	err = interp2.Run(`"before" # This is a comment`)
	assert.NoError(t, err)
	assert.Equal(t, 1, interp2.GetStack().Length())
}

func TestInterpreter_EmptyArray(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("[]")
	assert.NoError(t, err)
	assert.Equal(t, 1, interp.GetStack().Length())

	result := interp.StackPop()
	arr, ok := result.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 0, len(arr))
}

func TestInterpreter_ArrayWithItems(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`[1 2 3]`)
	assert.NoError(t, err)
	assert.Equal(t, 1, interp.GetStack().Length())

	result := interp.StackPop()
	arr, ok := result.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 3, len(arr))
	assert.Equal(t, int64(1), arr[0])
	assert.Equal(t, int64(2), arr[1])
	assert.Equal(t, int64(3), arr[2])
}

func TestInterpreter_StartModule(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("{")
	assert.NoError(t, err)
	// Module stack should have 2 modules (app + pushed app)
	assert.Equal(t, 2, len(interp.moduleStack))
}

func TestInterpreter_ModuleNested(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("{mymodule")
	assert.NoError(t, err)
	assert.Equal(t, "mymodule", interp.CurModule().GetName())
}

func TestInterpreter_ModuleClosure(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("{mymodule }")
	assert.NoError(t, err)
	// Back to app module
	assert.Equal(t, "", interp.CurModule().GetName())
}

func TestInterpreter_Definition(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`: PUSH_42 42 ;`)
	assert.NoError(t, err)

	// Word should be defined
	word := interp.CurModule().FindDictionaryWord("PUSH_42")
	assert.NotNil(t, word)
}

func TestInterpreter_DefinitionExecution(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`: PUSH_42 42 ; PUSH_42`)
	assert.NoError(t, err)

	// Should have 42 on stack
	assert.Equal(t, 1, interp.GetStack().Length())
	assert.Equal(t, int64(42), interp.StackPop())
}

func TestInterpreter_Memo(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`@: CONSTANT 42 ;`)
	assert.NoError(t, err)

	// Should have created the memo word and its variants
	memoWord := interp.CurModule().FindDictionaryWord("CONSTANT")
	assert.NotNil(t, memoWord)

	refreshWord := interp.CurModule().FindDictionaryWord("CONSTANT!")
	assert.NotNil(t, refreshWord)

	refreshAtWord := interp.CurModule().FindDictionaryWord("CONSTANT!@")
	assert.NotNil(t, refreshAtWord)
}

func TestInterpreter_Literals(t *testing.T) {
	interp := NewInterpreter()

	// Test boolean
	err := interp.Run("TRUE FALSE")
	assert.NoError(t, err)
	assert.Equal(t, false, interp.StackPop())
	assert.Equal(t, true, interp.StackPop())

	// Test integer
	err = interp.Run("42")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), interp.StackPop())

	// Test float
	err = interp.Run("3.14")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, interp.StackPop())
}

func TestInterpreter_UnknownWord(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run("UNKNOWN_WORD")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown word")
}

func TestInterpreter_StackUnderflow(t *testing.T) {
	interp := NewInterpreter()
	assert.Panics(t, func() {
		interp.StackPop()
	})
}

func TestInterpreter_MissingSemicolon(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`: WORD`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Missing semicolon")
}

func TestInterpreter_ExtraSemicolon(t *testing.T) {
	interp := NewInterpreter()
	err := interp.Run(`;`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Extra semicolon")
}
