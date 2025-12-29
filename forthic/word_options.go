package forthic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// WordOptions - Type-safe options container for module words
//
// Overview:
// WordOptions provides a structured way for Forthic words to accept optional
// configuration parameters without requiring fixed parameter positions. This
// enables flexible, extensible APIs similar to keyword arguments in other languages.
//
// Usage in Forthic:
//   [.option_name value ...] ~> WORD
//
// Example in Forthic code:
//   [1 2 3] '2 *' [.with_key TRUE] ~> MAP
//   [10 20 30] [.comparator "-1 *"] ~> SORT
//   [[[1 2]]] [.depth 1] ~> FLATTEN
//
// Internal Representation:
// Created from flat array: [.key1 val1 .key2 val2]
// Stored as map internally for efficient lookup
//
// Note: Dot-symbols in Forthic have the leading '.' already stripped,
// so keys come in as "key1", "key2", etc.
type WordOptions struct {
	options map[string]interface{}
}

// NewWordOptions creates a new WordOptions from a flat array of key-value pairs
// flatArray must be []interface{} with even length: [key1, val1, key2, val2, ...]
// Keys must be strings (dot-symbols with . already stripped)
func NewWordOptions(flatArray interface{}) (*WordOptions, error) {
	// Check if it's an array
	arr, ok := flatArray.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Options must be an array")
	}

	// Check even length
	if len(arr)%2 != 0 {
		return nil, fmt.Errorf("Options must be key-value pairs (even length). Got %d elements", len(arr))
	}

	opts := &WordOptions{
		options: make(map[string]interface{}),
	}

	// Parse key-value pairs
	for i := 0; i < len(arr); i += 2 {
		key := arr[i]
		value := arr[i+1]

		// Key must be a string
		keyStr, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("Option key must be a string (dot-symbol). Got: %T", key)
		}

		opts.options[keyStr] = value
	}

	return opts, nil
}

// Get retrieves an option value with optional default
func (wo *WordOptions) Get(key string, defaultValue ...interface{}) interface{} {
	if val, ok := wo.options[key]; ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return nil
}

// Has checks if an option key exists
func (wo *WordOptions) Has(key string) bool {
	_, ok := wo.options[key]
	return ok
}

// ToRecord converts all options to a plain map
func (wo *WordOptions) ToRecord() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range wo.options {
		result[k] = v
	}
	return result
}

// Keys returns all option keys
func (wo *WordOptions) Keys() []string {
	keys := make([]string, 0, len(wo.options))
	for k := range wo.options {
		keys = append(keys, k)
	}
	return keys
}

// String returns a formatted string representation for debugging
func (wo *WordOptions) String() string {
	if len(wo.options) == 0 {
		return "<WordOptions: >"
	}

	pairs := make([]string, 0, len(wo.options))
	for k, v := range wo.options {
		// Try to JSON encode the value for display
		var valStr string
		if b, err := json.Marshal(v); err == nil {
			valStr = string(b)
		} else {
			valStr = fmt.Sprintf("%v", v)
		}
		pairs = append(pairs, fmt.Sprintf(".%s %s", k, valStr))
	}

	return fmt.Sprintf("<WordOptions: %s>", strings.Join(pairs, " "))
}
