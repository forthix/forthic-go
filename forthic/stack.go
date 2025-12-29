package forthic

import (
	"encoding/json"
	"fmt"
)

// Stack - Wrapper for the interpreter's data stack
//
// Provides LIFO stack operations for Forthic interpreter.
// All stack values are stored as interface{} to allow dynamic typing.
type Stack struct {
	items []interface{}
}

// NewStack creates a new Stack with optional initial items
func NewStack(items ...interface{}) *Stack {
	if items == nil {
		return &Stack{
			items: make([]interface{}, 0),
		}
	}
	return &Stack{
		items: items,
	}
}

// Push adds a value to the top of the stack
func (s *Stack) Push(val interface{}) {
	s.items = append(s.items, val)
}

// Pop removes and returns the top value from the stack
// Returns error if stack is empty
func (s *Stack) Pop() (interface{}, error) {
	if len(s.items) == 0 {
		return nil, NewStackUnderflowError()
	}
	val := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return val, nil
}

// Peek returns the top value without removing it
// Returns error if stack is empty
func (s *Stack) Peek() (interface{}, error) {
	if len(s.items) == 0 {
		return nil, NewStackUnderflowError()
	}
	return s.items[len(s.items)-1], nil
}

// Length returns the number of items on the stack
func (s *Stack) Length() int {
	return len(s.items)
}

// Clear removes all items from the stack
func (s *Stack) Clear() {
	s.items = make([]interface{}, 0)
}

// Items returns a copy of the stack items
func (s *Stack) Items() []interface{} {
	result := make([]interface{}, len(s.items))
	copy(result, s.items)
	return result
}

// RawItems returns the internal items slice (for internal use only)
func (s *Stack) RawItems() []interface{} {
	return s.items
}

// Get retrieves the item at the specified index (0 = bottom, length-1 = top)
// Returns error if index is out of bounds
func (s *Stack) Get(index int) (interface{}, error) {
	if index < 0 || index >= len(s.items) {
		return nil, fmt.Errorf("Stack index out of bounds: %d (length: %d)", index, len(s.items))
	}
	return s.items[index], nil
}

// Set sets the item at the specified index (0 = bottom, length-1 = top)
// Returns error if index is out of bounds
func (s *Stack) Set(index int, val interface{}) error {
	if index < 0 || index >= len(s.items) {
		return fmt.Errorf("Stack index out of bounds: %d (length: %d)", index, len(s.items))
	}
	s.items[index] = val
	return nil
}

// String returns a formatted string representation for debugging
func (s *Stack) String() string {
	if len(s.items) == 0 {
		return "Stack[]"
	}

	itemStrs := make([]string, len(s.items))
	for i, item := range s.items {
		if b, err := json.Marshal(item); err == nil {
			itemStrs[i] = string(b)
		} else {
			itemStrs[i] = fmt.Sprintf("%v", item)
		}
	}

	return fmt.Sprintf("Stack[%d items]", len(s.items))
}

// ToJSON returns the stack items as a JSON array string
func (s *Stack) ToJSON() (string, error) {
	b, err := json.Marshal(s.items)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
