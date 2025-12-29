package modules

import (
	"encoding/json"

	"github.com/forthix/forthic-go/forthic"
)

// JSONModule provides JSON encoding/decoding operations
type JSONModule struct {
	*forthic.Module
}

// NewJSONModule creates a new JSON module
func NewJSONModule() *JSONModule {
	m := &JSONModule{
		Module: forthic.NewModule("json", ""),
	}
	m.registerWords()
	return m
}

func (m *JSONModule) registerWords() {
	// Encoding
	m.AddModuleWord(">JSON", m.toJSON)
	m.AddModuleWord("JSON-PRETTIFY", m.jsonPrettify)

	// Decoding
	m.AddModuleWord("JSON>", m.fromJSON)
}

// ========================================
// Encoding
// ========================================

func (m *JSONModule) toJSON(interp *forthic.Interpreter) error {
	value := interp.StackPop()

	if value == nil {
		interp.StackPush("null")
		return nil
	}

	// Convert to JSON
	bytes, err := json.Marshal(value)
	if err != nil {
		interp.StackPush("")
		return nil
	}

	interp.StackPush(string(bytes))
	return nil
}

func (m *JSONModule) jsonPrettify(interp *forthic.Interpreter) error {
	value := interp.StackPop()

	if value == nil {
		interp.StackPush("null")
		return nil
	}

	// Pretty print with 2-space indentation
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		interp.StackPush("")
		return nil
	}

	interp.StackPush(string(bytes))
	return nil
}

// ========================================
// Decoding
// ========================================

func (m *JSONModule) fromJSON(interp *forthic.Interpreter) error {
	jsonStr := interp.StackPop()

	if jsonStr == nil {
		interp.StackPush(nil)
		return nil
	}

	str, ok := jsonStr.(string)
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	// Parse JSON
	var result interface{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		interp.StackPush(nil)
		return nil
	}

	// Convert parsed result to Forthic types
	interp.StackPush(normalizeJSONValue(result))
	return nil
}

// ========================================
// Helper Functions
// ========================================

// normalizeJSONValue converts JSON parsed values to Forthic types
func normalizeJSONValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case map[string]interface{}:
		// Recursively normalize map values
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = normalizeJSONValue(value)
		}
		return result

	case []interface{}:
		// Recursively normalize array elements
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = normalizeJSONValue(item)
		}
		return result

	case float64:
		// JSON numbers are float64 by default
		// Keep as float64 for consistency
		return v

	case bool:
		return v

	case string:
		return v

	default:
		return v
	}
}
