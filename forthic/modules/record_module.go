package modules

import (
	"github.com/forthix/forthic-go/forthic"
)

// RecordModule provides record/map manipulation operations
type RecordModule struct {
	*forthic.Module
}

// NewRecordModule creates a new record module
func NewRecordModule() *RecordModule {
	m := &RecordModule{
		Module: forthic.NewModule("record", ""),
	}
	m.registerWords()
	return m
}

func (m *RecordModule) registerWords() {
	// Creation
	m.AddModuleWord("REC", m.createRecord)
	m.AddModuleWord("<REC!", m.setRecordValue)

	// Access
	m.AddModuleWord("REC@", m.getRecordValue)
	m.AddModuleWord("|REC@", m.pipeRecAt)
	m.AddModuleWord("KEYS", m.keys)
	m.AddModuleWord("VALUES", m.values)

	// Transform
	m.AddModuleWord("RELABEL", m.relabel)
	m.AddModuleWord("INVERT-KEYS", m.invertKeys)
	m.AddModuleWord("REC-DEFAULTS", m.recDefaults)
	m.AddModuleWord("<DEL", m.del)
}

// ========================================
// Creation
// ========================================

func (m *RecordModule) createRecord(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	// Build record from [[key, val], ...] pairs
	result := make(map[string]interface{})
	for _, item := range slice {
		pair, ok := item.([]interface{})
		if !ok || len(pair) < 2 {
			// Skip invalid pairs
			continue
		}
		key, ok := pair[0].(string)
		if !ok {
			// Skip non-string keys
			continue
		}
		result[key] = pair[1]
	}

	interp.StackPush(result)
	return nil
}

func (m *RecordModule) setRecordValue(interp *forthic.Interpreter) error {
	// Stack signature: ( rec value field -- rec )
	field := interp.StackPop()
	value := interp.StackPop()
	record := interp.StackPop()

	if record == nil {
		record = make(map[string]interface{})
	}

	rec, ok := record.(map[string]interface{})
	if !ok {
		interp.StackPush(record)
		return nil
	}

	// Create a copy to avoid modifying original
	result := make(map[string]interface{})
	for k, v := range rec {
		result[k] = v
	}

	// Support both string and array of field names
	var fields []string
	if fieldStr, ok := field.(string); ok {
		fields = []string{fieldStr}
	} else if fieldArr, ok := field.([]interface{}); ok {
		fields = make([]string, len(fieldArr))
		for i, f := range fieldArr {
			if fStr, ok := f.(string); ok {
				fields[i] = fStr
			} else {
				// Invalid field path
				interp.StackPush(result)
				return nil
			}
		}
	} else {
		interp.StackPush(result)
		return nil
	}

	// Drill down to set value at path
	if len(fields) == 0 {
		interp.StackPush(result)
		return nil
	}

	// Drill down, creating nested maps as needed
	curRec := result
	for i := 0; i < len(fields)-1; i++ {
		fieldName := fields[i]
		if existing, ok := curRec[fieldName].(map[string]interface{}); ok {
			// Copy the existing nested record
			newRec := make(map[string]interface{})
			for k, v := range existing {
				newRec[k] = v
			}
			curRec[fieldName] = newRec
			curRec = newRec
		} else {
			// Create new nested record
			newRec := make(map[string]interface{})
			curRec[fieldName] = newRec
			curRec = newRec
		}
	}

	// Set the final value
	curRec[fields[len(fields)-1]] = value

	interp.StackPush(result)
	return nil
}

// ========================================
// Access
// ========================================

func (m *RecordModule) getRecordValue(interp *forthic.Interpreter) error {
	field := interp.StackPop()
	record := interp.StackPop()

	if record == nil {
		interp.StackPush(nil)
		return nil
	}

	rec, ok := record.(map[string]interface{})
	if !ok {
		interp.StackPush(nil)
		return nil
	}

	// Support both string and array of field names
	var fields []string
	if fieldStr, ok := field.(string); ok {
		fields = []string{fieldStr}
	} else if fieldArr, ok := field.([]interface{}); ok {
		fields = make([]string, len(fieldArr))
		for i, f := range fieldArr {
			if fStr, ok := f.(string); ok {
				fields[i] = fStr
			} else {
				interp.StackPush(nil)
				return nil
			}
		}
	} else {
		interp.StackPush(nil)
		return nil
	}

	// Drill down through nested records
	result := drillForValue(rec, fields)
	interp.StackPush(result)
	return nil
}

func (m *RecordModule) keys(interp *forthic.Interpreter) error {
	record := interp.StackPop()

	if record == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	rec, ok := record.(map[string]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := []interface{}{}
	for key := range rec {
		result = append(result, key)
	}

	interp.StackPush(result)
	return nil
}

func (m *RecordModule) values(interp *forthic.Interpreter) error {
	record := interp.StackPop()

	if record == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	rec, ok := record.(map[string]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := []interface{}{}
	for _, val := range rec {
		result = append(result, val)
	}

	interp.StackPush(result)
	return nil
}

// ========================================
// Transform
// ========================================


func (m *RecordModule) invertKeys(interp *forthic.Interpreter) error {
	record := interp.StackPop()

	if record == nil {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	rec, ok := record.(map[string]interface{})
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	// Invert two-level nested record structure
	// Input:  {A: {X: 1, Y: 2}, B: {X: 3, Y: 4}}
	// Output: {X: {A: 1, B: 3}, Y: {A: 2, B: 4}}
	result := make(map[string]interface{})

	for firstKey, subRecordVal := range rec {
		subRecord, ok := subRecordVal.(map[string]interface{})
		if !ok {
			continue
		}

		for secondKey, value := range subRecord {
			if _, exists := result[secondKey]; !exists {
				result[secondKey] = make(map[string]interface{})
			}
			result[secondKey].(map[string]interface{})[firstKey] = value
		}
	}

	interp.StackPush(result)
	return nil
}


// ========================================
// Additional Words
// ========================================

func (m *RecordModule) pipeRecAt(interp *forthic.Interpreter) error {
	field := interp.StackPop()
	records := interp.StackPop()

	if records == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice, ok := records.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	// Map REC@ over array of records
	result := make([]interface{}, len(slice))
	for i, record := range slice {
		if rec, ok := record.(map[string]interface{}); ok {
			// Convert field to fields array
			var fields []string
			if fieldStr, ok := field.(string); ok {
				fields = []string{fieldStr}
			} else if fieldArr, ok := field.([]interface{}); ok {
				fields = make([]string, len(fieldArr))
				for j, f := range fieldArr {
					if fStr, ok := f.(string); ok {
						fields[j] = fStr
					}
				}
			}
			result[i] = drillForValue(rec, fields)
		} else {
			result[i] = nil
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *RecordModule) relabel(interp *forthic.Interpreter) error {
	newKeys := interp.StackPop()
	oldKeys := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(container)
		return nil
	}

	oldKeyArr, ok1 := oldKeys.([]interface{})
	newKeyArr, ok2 := newKeys.([]interface{})

	if !ok1 || !ok2 || len(oldKeyArr) != len(newKeyArr) {
		interp.StackPush(container)
		return nil
	}

	// Build mapping from new keys to old keys
	newToOld := make(map[string]string)
	for i := 0; i < len(oldKeyArr); i++ {
		oldKey, ok1 := oldKeyArr[i].(string)
		newKey, ok2 := newKeyArr[i].(string)
		if ok1 && ok2 {
			newToOld[newKey] = oldKey
		}
	}

	// Apply relabeling
	if rec, ok := container.(map[string]interface{}); ok {
		result := make(map[string]interface{})
		for newKey, oldKey := range newToOld {
			if val, exists := rec[oldKey]; exists {
				result[newKey] = val
			}
		}
		interp.StackPush(result)
	} else {
		interp.StackPush(container)
	}

	return nil
}

func (m *RecordModule) recDefaults(interp *forthic.Interpreter) error {
	keyVals := interp.StackPop()
	record := interp.StackPop()

	if record == nil {
		record = make(map[string]interface{})
	}

	rec, ok1 := record.(map[string]interface{})
	keyValArr, ok2 := keyVals.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush(record)
		return nil
	}

	// Create copy
	result := make(map[string]interface{})
	for k, v := range rec {
		result[k] = v
	}

	// Set defaults for missing/empty fields
	for _, item := range keyValArr {
		pair, ok := item.([]interface{})
		if !ok || len(pair) < 2 {
			continue
		}
		key, ok := pair[0].(string)
		if !ok {
			continue
		}

		// Set default if missing, null, or empty string
		if val, exists := result[key]; !exists || val == nil || val == "" {
			result[key] = pair[1]
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *RecordModule) del(interp *forthic.Interpreter) error {
	key := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(container)
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		// Delete from array by index
		idx := toInt(key)
		if idx < 0 || idx >= len(arr) {
			interp.StackPush(arr)
			return nil
		}
		result := append(arr[:idx], arr[idx+1:]...)
		interp.StackPush(result)
	} else if rec, ok := container.(map[string]interface{}); ok {
		// Delete from record by key
		keyStr, ok := key.(string)
		if !ok {
			interp.StackPush(rec)
			return nil
		}
		result := make(map[string]interface{})
		for k, v := range rec {
			if k != keyStr {
				result[k] = v
			}
		}
		interp.StackPush(result)
	} else {
		interp.StackPush(container)
	}

	return nil
}

// ========================================
// Helper Functions
// ========================================

func drillForValue(record map[string]interface{}, fields []string) interface{} {
	var result interface{} = record
	for _, field := range fields {
		if result == nil {
			return nil
		}
		if rec, ok := result.(map[string]interface{}); ok {
			if val, exists := rec[field]; exists {
				result = val
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return result
}
