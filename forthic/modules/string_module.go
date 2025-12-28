package modules

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/forthix/forthic-go/forthic"
)

// StringModule provides string manipulation operations
type StringModule struct {
	*forthic.Module
}

// NewStringModule creates a new string module
func NewStringModule() *StringModule {
	m := &StringModule{
		Module: forthic.NewModule("string", ""),
	}
	m.registerWords()
	return m
}

func (m *StringModule) registerWords() {
	// Conversion
	m.AddModuleWord(">STR", m.toStr)
	m.AddModuleWord("URL-ENCODE", m.urlEncode)
	m.AddModuleWord("URL-DECODE", m.urlDecode)

	// Transform
	m.AddModuleWord("LOWERCASE", m.lowercase)
	m.AddModuleWord("UPPERCASE", m.uppercase)
	m.AddModuleWord("STRIP", m.strip)
	m.AddModuleWord("ASCII", m.ascii)

	// Split/Join
	m.AddModuleWord("SPLIT", m.split)
	m.AddModuleWord("JOIN", m.join)
	m.AddModuleWord("CONCAT", m.concat)

	// Pattern
	m.AddModuleWord("REPLACE", m.replace)
	m.AddModuleWord("RE-MATCH", m.reMatch)
	m.AddModuleWord("RE-MATCH-ALL", m.reMatchAll)
	m.AddModuleWord("RE-MATCH-GROUP", m.reMatchGroup)

	// Constants
	m.AddModuleWord("/N", m.slashN)
	m.AddModuleWord("/R", m.slashR)
	m.AddModuleWord("/T", m.slashT)
}

// ========================================
// Conversion
// ========================================

func (m *StringModule) toStr(interp *forthic.Interpreter) error {
	item := interp.StackPop()
	interp.StackPush(fmt.Sprintf("%v", item))
	return nil
}

func (m *StringModule) urlEncode(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil || str == "" {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		interp.StackPush(url.QueryEscape(s))
	} else {
		interp.StackPush("")
	}
	return nil
}

func (m *StringModule) urlDecode(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil || str == "" {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		decoded, err := url.QueryUnescape(s)
		if err != nil {
			interp.StackPush("")
		} else {
			interp.StackPush(decoded)
		}
	} else {
		interp.StackPush("")
	}
	return nil
}

// ========================================
// Transform
// ========================================

func (m *StringModule) lowercase(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		interp.StackPush(strings.ToLower(s))
	} else {
		interp.StackPush("")
	}
	return nil
}

func (m *StringModule) uppercase(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		interp.StackPush(strings.ToUpper(s))
	} else {
		interp.StackPush("")
	}
	return nil
}

func (m *StringModule) strip(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		interp.StackPush(strings.TrimSpace(s))
	} else {
		interp.StackPush("")
	}
	return nil
}

func (m *StringModule) ascii(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil {
		interp.StackPush("")
		return nil
	}
	if s, ok := str.(string); ok {
		result := ""
		for _, ch := range s {
			if ch < 256 {
				result += string(ch)
			}
		}
		interp.StackPush(result)
	} else {
		interp.StackPush("")
	}
	return nil
}

// ========================================
// Split/Join
// ========================================

func (m *StringModule) split(interp *forthic.Interpreter) error {
	sep := interp.StackPop()
	str := interp.StackPop()

	if str == nil {
		str = ""
	}

	s, ok1 := str.(string)
	sepStr, ok2 := sep.(string)

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	parts := strings.Split(s, sepStr)
	result := make([]interface{}, len(parts))
	for i, part := range parts {
		result[i] = part
	}
	interp.StackPush(result)
	return nil
}

func (m *StringModule) join(interp *forthic.Interpreter) error {
	sep := interp.StackPop()
	arr := interp.StackPop()

	if arr == nil {
		interp.StackPush("")
		return nil
	}

	slice, ok1 := arr.([]interface{})
	sepStr, ok2 := sep.(string)

	if !ok1 || !ok2 {
		interp.StackPush("")
		return nil
	}

	parts := make([]string, len(slice))
	for i, item := range slice {
		if s, ok := item.(string); ok {
			parts[i] = s
		} else {
			parts[i] = fmt.Sprintf("%v", item)
		}
	}

	interp.StackPush(strings.Join(parts, sepStr))
	return nil
}

func (m *StringModule) concat(interp *forthic.Interpreter) error {
	str2 := interp.StackPop()

	// Case 1: Array on top of stack
	if arr, ok := str2.([]interface{}); ok {
		parts := make([]string, len(arr))
		for i, item := range arr {
			if s, ok := item.(string); ok {
				parts[i] = s
			} else {
				parts[i] = fmt.Sprintf("%v", item)
			}
		}
		interp.StackPush(strings.Join(parts, ""))
		return nil
	}

	// Case 2: Two strings
	str1 := interp.StackPop()

	s1, ok1 := str1.(string)
	s2, ok2 := str2.(string)

	if !ok1 {
		s1 = fmt.Sprintf("%v", str1)
	}
	if !ok2 {
		s2 = fmt.Sprintf("%v", str2)
	}

	interp.StackPush(s1 + s2)
	return nil
}

// ========================================
// Pattern
// ========================================

func (m *StringModule) replace(interp *forthic.Interpreter) error {
	replaceStr := interp.StackPop()
	text := interp.StackPop()
	str := interp.StackPop()

	if str == nil {
		interp.StackPush("")
		return nil
	}

	s, ok1 := str.(string)
	t, ok2 := text.(string)
	r, ok3 := replaceStr.(string)

	if !ok1 || !ok2 || !ok3 {
		interp.StackPush(s)
		return nil
	}

	// Treat text as regex pattern (standard Forthic behavior)
	re, err := regexp.Compile(t)
	if err != nil {
		// If regex is invalid, return original string
		interp.StackPush(s)
		return nil
	}
	result := re.ReplaceAllString(s, r)
	interp.StackPush(result)
	return nil
}

func (m *StringModule) reMatch(interp *forthic.Interpreter) error {
	pattern := interp.StackPop()
	str := interp.StackPop()

	if str == nil {
		interp.StackPush(false)
		return nil
	}

	s, ok1 := str.(string)
	p, ok2 := pattern.(string)

	if !ok1 || !ok2 {
		interp.StackPush(false)
		return nil
	}

	re, err := regexp.Compile(p)
	if err != nil {
		interp.StackPush(false)
		return nil
	}

	matches := re.FindStringSubmatch(s)
	if matches == nil {
		interp.StackPush(false)
		return nil
	}

	// Convert to []interface{} for consistency
	result := make([]interface{}, len(matches))
	for i, match := range matches {
		result[i] = match
	}
	interp.StackPush(result)
	return nil
}

func (m *StringModule) reMatchAll(interp *forthic.Interpreter) error {
	pattern := interp.StackPop()
	str := interp.StackPop()

	if str == nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	s, ok1 := str.(string)
	p, ok2 := pattern.(string)

	if !ok1 || !ok2 {
		interp.StackPush([]interface{}{})
		return nil
	}

	re, err := regexp.Compile(p)
	if err != nil {
		interp.StackPush([]interface{}{})
		return nil
	}

	allMatches := re.FindAllStringSubmatch(s, -1)
	result := make([]interface{}, 0)

	// Extract first capture group from each match (matches TypeScript behavior)
	for _, matches := range allMatches {
		if len(matches) > 1 {
			result = append(result, matches[1])
		}
	}

	interp.StackPush(result)
	return nil
}

func (m *StringModule) reMatchGroup(interp *forthic.Interpreter) error {
	num := interp.StackPop()
	match := interp.StackPop()

	if match == nil || match == false {
		interp.StackPush(nil)
		return nil
	}

	arr, ok1 := match.([]interface{})
	idx, ok2 := num.(int64)

	if !ok1 {
		interp.StackPush(nil)
		return nil
	}

	if !ok2 {
		// Try int
		if i, ok := num.(int); ok {
			idx = int64(i)
		} else {
			interp.StackPush(nil)
			return nil
		}
	}

	if idx < 0 || idx >= int64(len(arr)) {
		interp.StackPush(nil)
		return nil
	}

	interp.StackPush(arr[idx])
	return nil
}

// ========================================
// Constants
// ========================================

func (m *StringModule) slashN(interp *forthic.Interpreter) error {
	interp.StackPush("\n")
	return nil
}

func (m *StringModule) slashR(interp *forthic.Interpreter) error {
	interp.StackPush("\r")
	return nil
}

func (m *StringModule) slashT(interp *forthic.Interpreter) error {
	interp.StackPush("\t")
	return nil
}
