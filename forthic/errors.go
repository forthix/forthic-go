package forthic

import (
	"fmt"
	"strings"
)

// CodeLocation represents a location in Forthic source code
type CodeLocation struct {
	Source   string
	File     string
	Line     int
	Column   int
	StartPos int
	EndPos   int
}

func (l CodeLocation) String() string {
	if l.File == "" {
		return fmt.Sprintf("line %d, col %d", l.Line, l.Column)
	}
	return fmt.Sprintf("%s:%d:%d", l.File, l.Line, l.Column)
}

// ForthicError is the base error type for all Forthic errors
type ForthicError struct {
	Message  string
	Forthic  string
	Location *CodeLocation
	Cause    error
}

func (e *ForthicError) Error() string {
	var parts []string

	parts = append(parts, e.Message)

	if e.Location != nil {
		parts = append(parts, fmt.Sprintf("at %s", e.Location))
	}

	if e.Forthic != "" {
		parts = append(parts, fmt.Sprintf("in: %s", e.Forthic))
	}

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("caused by: %v", e.Cause))
	}

	return strings.Join(parts, "\n  ")
}

func (e *ForthicError) Unwrap() error {
	return e.Cause
}

// NewForthicError creates a new ForthicError
func NewForthicError(message string) *ForthicError {
	return &ForthicError{
		Message: message,
	}
}

// WithLocation adds location information to the error
func (e *ForthicError) WithLocation(loc *CodeLocation) *ForthicError {
	e.Location = loc
	return e
}

// WithForthic adds the Forthic code snippet to the error
func (e *ForthicError) WithForthic(forthic string) *ForthicError {
	e.Forthic = forthic
	return e
}

// WithCause adds a causal error
func (e *ForthicError) WithCause(cause error) *ForthicError {
	e.Cause = cause
	return e
}

// UnknownWordError represents an attempt to execute an unknown word
type UnknownWordError struct {
	*ForthicError
	Word string
}

func NewUnknownWordError(word string) *UnknownWordError {
	return &UnknownWordError{
		ForthicError: NewForthicError(fmt.Sprintf("Unknown word: %s", word)),
		Word:         word,
	}
}

// UnknownModuleError represents an attempt to use an unknown module
type UnknownModuleError struct {
	*ForthicError
	Module string
}

func NewUnknownModuleError(module string) *UnknownModuleError {
	return &UnknownModuleError{
		ForthicError: NewForthicError(fmt.Sprintf("Unknown module: %s", module)),
		Module:       module,
	}
}

// StackUnderflowError represents an attempt to pop from an empty stack
type StackUnderflowError struct {
	*ForthicError
}

func NewStackUnderflowError() *StackUnderflowError {
	return &StackUnderflowError{
		ForthicError: NewForthicError("Stack underflow"),
	}
}

// WordExecutionError represents an error during word execution
type WordExecutionError struct {
	*ForthicError
	Word string
}

func NewWordExecutionError(word string, err error) *WordExecutionError {
	return &WordExecutionError{
		ForthicError: NewForthicError(fmt.Sprintf("Error executing word: %s", word)).WithCause(err),
		Word:         word,
	}
}

// MissingSemicolonError represents a missing semicolon in a definition
type MissingSemicolonError struct {
	*ForthicError
}

func NewMissingSemicolonError() *MissingSemicolonError {
	return &MissingSemicolonError{
		ForthicError: NewForthicError("Missing semicolon (;) to end definition"),
	}
}

// ExtraSemicolonError represents an extra semicolon outside a definition
type ExtraSemicolonError struct {
	*ForthicError
}

func NewExtraSemicolonError() *ExtraSemicolonError {
	return &ExtraSemicolonError{
		ForthicError: NewForthicError("Extra semicolon (;) outside of definition"),
	}
}

// ModuleError represents an error related to module operations
type ModuleError struct {
	*ForthicError
	Module string
}

func NewModuleError(module string, message string) *ModuleError {
	return &ModuleError{
		ForthicError: NewForthicError(fmt.Sprintf("Module error in %s: %s", module, message)),
		Module:       module,
	}
}

// IntentionalStopError represents an intentional stop (not a real error)
type IntentionalStopError struct {
	*ForthicError
}

func NewIntentionalStopError(message string) *IntentionalStopError {
	return &IntentionalStopError{
		ForthicError: NewForthicError(message),
	}
}

// InvalidVariableNameError represents an invalid variable name
type InvalidVariableNameError struct {
	*ForthicError
	VarName string
}

func NewInvalidVariableNameError(varName string) *InvalidVariableNameError {
	return &InvalidVariableNameError{
		ForthicError: NewForthicError(fmt.Sprintf("Invalid variable name: %s", varName)),
		VarName:      varName,
	}
}
