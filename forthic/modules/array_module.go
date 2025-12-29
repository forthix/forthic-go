package modules

import (
	"sort"

	"github.com/forthix/forthic-go/forthic"
)

// ArrayModule provides array manipulation operations
type ArrayModule struct {
	*forthic.Module
}

// NewArrayModule creates a new array module
func NewArrayModule() *ArrayModule {
	m := &ArrayModule{
		Module: forthic.NewModule("array", ""),
	}
	m.registerWords()
	return m
}

func (m *ArrayModule) registerWords() {
	// Basic operations
	m.AddModuleWord("APPEND", m.append_word)
	m.AddModuleWord("REVERSE", m.reverse)
	m.AddModuleWord("UNIQUE", m.unique)
	m.AddModuleWord("LENGTH", m.length)

	// Access operations
	m.AddModuleWord("NTH", m.nth)
	m.AddModuleWord("LAST", m.last)
	m.AddModuleWord("SLICE", m.slice)
	m.AddModuleWord("TAKE", m.take)
	m.AddModuleWord("DROP", m.drop)
	m.AddModuleWord("KEY-OF", m.keyOf)

	// Set operations
	m.AddModuleWord("DIFFERENCE", m.difference)
	m.AddModuleWord("INTERSECTION", m.intersection)
	m.AddModuleWord("UNION", m.union)

	// Sort and shuffle
	m.AddModuleWord("SORT", m.sortArray)
	m.AddModuleWord("SHUFFLE", m.shuffle)
	m.AddModuleWord("ROTATE", m.rotate)

	// Combine
	m.AddModuleWord("ZIP", m.zip)
	m.AddModuleWord("ZIP-WITH", m.zipWith)
	m.AddModuleWord("FLATTEN", m.flatten)
	m.AddModuleWord("UNPACK", m.unpack)

	// Group and index
	m.AddModuleWord("INDEX", m.index)
	m.AddModuleWord("BY-FIELD", m.byField)
	m.AddModuleWord("GROUP-BY-FIELD", m.groupByField)
	m.AddModuleWord("GROUP-BY", m.groupBy)
	m.AddModuleWord("GROUPS-OF", m.groupsOf)

	// Transform
	m.AddModuleWord("MAP", m.mapArray)
	m.AddModuleWord("SELECT", m.selectArray)
	m.AddModuleWord("REDUCE", m.reduce)
	m.AddModuleWord("FOREACH", m.foreach)
	m.AddModuleWord("<REPEAT", m.repeat)
}

// ========================================
// Basic Operations
// ========================================

func (m *ArrayModule) append_word(interp *forthic.Interpreter) error {
	item := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		container = []interface{}{}
	}

	if arr, ok := container.([]interface{}); ok {
		result := append(arr, item)
		interp.StackPush(result)
	} else {
		// For records/maps, append [key, value] pair
		if rec, ok := container.(map[string]interface{}); ok {
			if pair, ok := item.([]interface{}); ok && len(pair) == 2 {
				if key, ok := pair[0].(string); ok {
					rec[key] = pair[1]
				}
			}
			interp.StackPush(rec)
		} else {
			interp.StackPush(container)
		}
	}
	return nil
}

func (m *ArrayModule) reverse(interp *forthic.Interpreter) error {
	container := interp.StackPop()

	if arr, ok := container.([]interface{}); ok {
		result := make([]interface{}, len(arr))
		for i, v := range arr {
			result[len(arr)-1-i] = v
		}
		interp.StackPush(result)
	} else {
		interp.StackPush(container)
	}
	return nil
}

func (m *ArrayModule) unique(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	if slice, ok := arr.([]interface{}); ok {
		seen := make(map[interface{}]bool)
		result := []interface{}{}

		for _, item := range slice {
			// Use string representation as key for complex types
			key := item
			if _, exists := seen[key]; !exists {
				seen[key] = true
				result = append(result, item)
			}
		}
		interp.StackPush(result)
	} else {
		interp.StackPush(arr)
	}
	return nil
}

func (m *ArrayModule) length(interp *forthic.Interpreter) error {
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(0)
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		interp.StackPush(len(arr))
	} else if rec, ok := container.(map[string]interface{}); ok {
		interp.StackPush(len(rec))
	} else {
		interp.StackPush(0)
	}
	return nil
}

// ========================================
// Access Operations
// ========================================

func (m *ArrayModule) nth(interp *forthic.Interpreter) error {
	n := interp.StackPop()
	container := interp.StackPop()

	if container == nil || n == nil {
		interp.StackPush(nil)
		return nil
	}

	var index int
	switch v := n.(type) {
	case int:
		index = v
	case int64:
		index = int(v)
	case float64:
		index = int(v)
	default:
		interp.StackPush(nil)
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		if index < 0 || index >= len(arr) {
			interp.StackPush(nil)
			return nil
		}
		interp.StackPush(arr[index])
	} else {
		interp.StackPush(nil)
	}
	return nil
}

func (m *ArrayModule) last(interp *forthic.Interpreter) error {
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(nil)
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		if len(arr) == 0 {
			interp.StackPush(nil)
			return nil
		}
		interp.StackPush(arr[len(arr)-1])
	} else {
		interp.StackPush(nil)
	}
	return nil
}

func (m *ArrayModule) slice(interp *forthic.Interpreter) error {
	endVal := interp.StackPop()
	startVal := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	arr, ok := container.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	start := toInt(startVal)
	end := toInt(endVal)
	length := len(arr)

	// Normalize negative indices
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// Handle reverse slice
	if start > end {
		result := []interface{}{}
		for i := start; i >= end && i >= 0 && i < length; i-- {
			result = append(result, arr[i])
		}
		interp.StackPush(result)
		return nil
	}

	// Forward slice
	if start < 0 || start >= length {
		interp.StackPush([]interface{}{})
		return nil
	}

	if end >= length {
		end = length - 1
	}

	result := []interface{}{}
	for i := start; i <= end && i < length; i++ {
		result = append(result, arr[i])
	}
	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) take(interp *forthic.Interpreter) error {
	n := interp.StackPop()
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	count := toInt(n)
	if count <= 0 {
		interp.StackPush([]interface{}{})
		return nil
	}

	if count >= len(slice) {
		interp.StackPush(slice)
		return nil
	}

	result := slice[:count]
	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) drop(interp *forthic.Interpreter) error {
	n := interp.StackPop()
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	count := toInt(n)
	if count <= 0 {
		interp.StackPush(slice)
		return nil
	}

	if count >= len(slice) {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := slice[count:]
	interp.StackPush(result)
	return nil
}

// ========================================
// Set Operations
// ========================================

func (m *ArrayModule) difference(interp *forthic.Interpreter) error {
	arr2 := interp.StackPop()
	arr1 := interp.StackPop()

	slice1, ok1 := arr1.([]interface{})
	slice2, ok2 := arr2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	// Create set from arr2
	set2 := make(map[interface{}]bool)
	for _, item := range slice2 {
		set2[item] = true
	}

	// Find items in arr1 not in arr2
	result := []interface{}{}
	for _, item := range slice1 {
		if !set2[item] {
			result = append(result, item)
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) intersection(interp *forthic.Interpreter) error {
	arr2 := interp.StackPop()
	arr1 := interp.StackPop()

	slice1, ok1 := arr1.([]interface{})
	slice2, ok2 := arr2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	// Create set from arr2
	set2 := make(map[interface{}]bool)
	for _, item := range slice2 {
		set2[item] = true
	}

	// Find items in both arrays
	seen := make(map[interface{}]bool)
	result := []interface{}{}
	for _, item := range slice1 {
		if set2[item] && !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) union(interp *forthic.Interpreter) error {
	arr2 := interp.StackPop()
	arr1 := interp.StackPop()

	slice1, ok1 := arr1.([]interface{})
	slice2, ok2 := arr2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	// Combine and deduplicate
	seen := make(map[interface{}]bool)
	result := []interface{}{}

	for _, item := range slice1 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	for _, item := range slice2 {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	interp.StackPush(result)
	return nil
}

// ========================================
// Sort
// ========================================

func (m *ArrayModule) sortArray(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush(arr)
		return nil
	}

	// Create a copy to avoid modifying original
	result := make([]interface{}, len(slice))
	copy(result, slice)

	// Simple numeric/string sort
	sort.Slice(result, func(i, j int) bool {
		return compareValues(result[i], result[j]) < 0
	})

	interp.StackPush(result)
	return nil
}

// ========================================
// Combine
// ========================================

func (m *ArrayModule) zip(interp *forthic.Interpreter) error {
	arr2 := interp.StackPop()
	arr1 := interp.StackPop()

	slice1, ok1 := arr1.([]interface{})
	slice2, ok2 := arr2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	minLen := len(slice1)
	if len(slice2) < minLen {
		minLen = len(slice2)
	}

	result := make([]interface{}, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = []interface{}{slice1[i], slice2[i]}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) flatten(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush(arr)
		return nil
	}

	// Fully flatten by default (depth = -1 means infinite depth)
	// TODO: Support depth option via ~> operator when implemented
	result := flattenArray(slice, -1)
	interp.StackPush(result)
	return nil
}

// ========================================
// Transform
// ========================================

func (m *ArrayModule) mapArray(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	arr := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := make([]interface{}, len(slice))
	for i, item := range slice {
		interp.StackPush(item)
		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		result[i] = interp.StackPop()
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) selectArray(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	arr := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := []interface{}{}
	for _, item := range slice {
		interp.StackPush(item)
		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		keep := interp.StackPop()
		if isTruthy(keep) {
			result = append(result, item)
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) reduce(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	initial := interp.StackPop()
	arr := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush(initial)
		return nil
	}

	slice, ok := arr.([]interface{})
	if !ok {
		interp.StackPush(initial)
		return nil
	}

	accumulator := initial
	for _, item := range slice {
		interp.StackPush(accumulator)
		interp.StackPush(item)
		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		accumulator = interp.StackPop()
	}

	interp.StackPush(accumulator)
	return nil
}

// ========================================
// Sort and Shuffle
// ========================================

func (m *ArrayModule) shuffle(interp *forthic.Interpreter) error {
	arr := interp.StackPop()

	slice, ok := arr.([]interface{})
	if !ok || len(slice) == 0 {
		interp.StackPush(arr)
		return nil
	}

	// Create a copy to avoid modifying original
	result := make([]interface{}, len(slice))
	copy(result, slice)

	// Fisher-Yates shuffle
	for i := len(result) - 1; i > 0; i-- {
		j := randInt(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) rotate(interp *forthic.Interpreter) error {
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(container)
		return nil
	}

	arr, ok := container.([]interface{})
	if !ok || len(arr) == 0 {
		interp.StackPush(container)
		return nil
	}

	// Rotate: move last element to front
	result := make([]interface{}, len(arr))
	result[0] = arr[len(arr)-1]
	copy(result[1:], arr[:len(arr)-1])

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) keyOf(interp *forthic.Interpreter) error {
	value := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(nil)
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		for i, item := range arr {
			if areEqual(item, value) {
				interp.StackPush(i)
				return nil
			}
		}
		interp.StackPush(nil)
	} else if rec, ok := container.(map[string]interface{}); ok {
		for key, val := range rec {
			if areEqual(val, value) {
				interp.StackPush(key)
				return nil
			}
		}
		interp.StackPush(nil)
	} else {
		interp.StackPush(nil)
	}

	return nil
}

// ========================================
// Combine Operations
// ========================================

func (m *ArrayModule) zipWith(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	arr2 := interp.StackPop()
	arr1 := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush([]interface{}{})
		return nil
	}

	slice1, ok1 := arr1.([]interface{})
	slice2, ok2 := arr2.([]interface{})

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	result := []interface{}{}
	for i := 0; i < len(slice1); i++ {
		var value2 interface{} = nil
		if i < len(slice2) {
			value2 = slice2[i]
		}
		interp.StackPush(slice1[i])
		interp.StackPush(value2)
		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		result = append(result, interp.StackPop())
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) unpack(interp *forthic.Interpreter) error {
	container := interp.StackPop()

	if container == nil {
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		for _, item := range arr {
			interp.StackPush(item)
		}
	} else if rec, ok := container.(map[string]interface{}); ok {
		// Get sorted keys for consistent order
		keys := make([]string, 0, len(rec))
		for k := range rec {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			interp.StackPush(rec[k])
		}
	}

	return nil
}

// ========================================
// Group and Index Operations
// ========================================

func (m *ArrayModule) index(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	items := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	slice, ok := items.([]interface{})
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	result := make(map[string]interface{})
	for _, item := range slice {
		interp.StackPush(item)
		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		keys := interp.StackPop()
		if keyArr, ok := keys.([]interface{}); ok {
			for _, k := range keyArr {
				keyStr := toLowerCase(k)
				if existing, ok := result[keyStr].([]interface{}); ok {
					result[keyStr] = append(existing, item)
				} else {
					result[keyStr] = []interface{}{item}
				}
			}
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) byField(interp *forthic.Interpreter) error {
	field := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	fieldStr, ok := field.(string)
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	var values []interface{}
	if arr, ok := container.([]interface{}); ok {
		values = arr
	} else if rec, ok := container.(map[string]interface{}); ok {
		values = []interface{}{}
		for _, v := range rec {
			values = append(values, v)
		}
	} else {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	result := make(map[string]interface{})
	for _, v := range values {
		if rec, ok := v.(map[string]interface{}); ok {
			if fieldVal, exists := rec[fieldStr]; exists {
				if fieldValStr, ok := fieldVal.(string); ok {
					result[fieldValStr] = v
				} else {
					keyStr := toString(fieldVal)
					result[keyStr] = v
				}
			}
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) groupByField(interp *forthic.Interpreter) error {
	field := interp.StackPop()
	container := interp.StackPop()

	if container == nil {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	fieldStr, ok := field.(string)
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	var values []interface{}
	if arr, ok := container.([]interface{}); ok {
		values = arr
	} else if rec, ok := container.(map[string]interface{}); ok {
		values = []interface{}{}
		for _, v := range rec {
			values = append(values, v)
		}
	} else {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	result := make(map[string]interface{})
	for _, v := range values {
		if rec, ok := v.(map[string]interface{}); ok {
			if fieldVal, exists := rec[fieldStr]; exists {
				// Handle field value that is an array
				if fieldArr, ok := fieldVal.([]interface{}); ok {
					for _, fv := range fieldArr {
						keyStr := toString(fv)
						if existing, ok := result[keyStr].([]interface{}); ok {
							result[keyStr] = append(existing, v)
						} else {
							result[keyStr] = []interface{}{v}
						}
					}
				} else {
					keyStr := toString(fieldVal)
					if existing, ok := result[keyStr].([]interface{}); ok {
						result[keyStr] = append(existing, v)
					} else {
						result[keyStr] = []interface{}{v}
					}
				}
			}
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) groupBy(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	items := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	if items == nil {
		interp.StackPush(map[string]interface{}{})
		return nil
	}

	result := make(map[string]interface{})

	if arr, ok := items.([]interface{}); ok {
		for _, item := range arr {
			interp.StackPush(item)
			err := interp.Run(codeStr)
			if err != nil {
				return err
			}
			groupKey := toString(interp.StackPop())
			if existing, ok := result[groupKey].([]interface{}); ok {
				result[groupKey] = append(existing, item)
			} else {
				result[groupKey] = []interface{}{item}
			}
		}
	} else if rec, ok := items.(map[string]interface{}); ok {
		for _, item := range rec {
			interp.StackPush(item)
			err := interp.Run(codeStr)
			if err != nil {
				return err
			}
			groupKey := toString(interp.StackPop())
			if existing, ok := result[groupKey].([]interface{}); ok {
				result[groupKey] = append(existing, item)
			} else {
				result[groupKey] = []interface{}{item}
			}
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *ArrayModule) groupsOf(interp *forthic.Interpreter) error {
	n := interp.StackPop()
	container := interp.StackPop()

	groupSize := toInt(n)
	if groupSize <= 0 {
		return forthicError("GROUPS-OF requires group size > 0")
	}

	if container == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	if arr, ok := container.([]interface{}); ok {
		numGroups := (len(arr) + groupSize - 1) / groupSize
		result := make([]interface{}, numGroups)
		for i := 0; i < numGroups; i++ {
			start := i * groupSize
			end := start + groupSize
			if end > len(arr) {
				end = len(arr)
			}
			result[i] = arr[start:end]
		}
		interp.StackPush(result)
	} else {
		interp.StackPush([]interface{}{})
	}

	return nil
}

// ========================================
// Iteration Operations
// ========================================

func (m *ArrayModule) foreach(interp *forthic.Interpreter) error {
	forthicCode := interp.StackPop()
	items := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		return nil
	}

	if items == nil {
		return nil
	}

	if arr, ok := items.([]interface{}); ok {
		for _, item := range arr {
			interp.StackPush(item)
			err := interp.Run(codeStr)
			if err != nil {
				return err
			}
		}
	} else if rec, ok := items.(map[string]interface{}); ok {
		for _, item := range rec {
			interp.StackPush(item)
			err := interp.Run(codeStr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *ArrayModule) repeat(interp *forthic.Interpreter) error {
	numTimes := interp.StackPop()
	forthicCode := interp.StackPop()

	codeStr, ok := forthicCode.(string)
	if !ok {
		return nil
	}

	count := toInt(numTimes)
	for i := 0; i < count; i++ {
		// Store item so we can push it back later
		item := interp.StackPop()
		interp.StackPush(item)

		err := interp.Run(codeStr)
		if err != nil {
			return err
		}
		res := interp.StackPop()

		// Push original item and result
		interp.StackPush(item)
		interp.StackPush(res)
	}

	return nil
}

// ========================================
// Helper Functions
// ========================================

func compareValues(a, b interface{}) int {
	// Try numeric comparison
	aNum, aOk := toNumericValue(a)
	bNum, bOk := toNumericValue(b)
	if aOk && bOk {
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

	return 0
}

func toNumericValue(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func flattenArray(arr []interface{}, depth int) []interface{} {
	// depth = -1 means fully flatten (infinite depth)
	// depth = 0 means don't flatten
	// depth > 0 means flatten that many levels
	if depth == 0 {
		return arr
	}

	result := []interface{}{}
	for _, item := range arr {
		if subArr, ok := item.([]interface{}); ok {
			// For infinite depth (-1), keep depth at -1
			// For finite depth, decrement
			nextDepth := depth
			if depth > 0 {
				nextDepth = depth - 1
			}
			flattened := flattenArray(subArr, nextDepth)
			result = append(result, flattened...)
		} else {
			result = append(result, item)
		}
	}
	return result
}
