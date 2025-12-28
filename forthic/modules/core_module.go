package modules

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/forthix/forthic-go/forthic"
)

// CoreModule provides essential interpreter operations
type CoreModule struct {
	*forthic.Module
}

// NewCoreModule creates a new core module
func NewCoreModule() *CoreModule {
	m := &CoreModule{
		Module: forthic.NewModule("core", ""),
	}
	m.registerWords()
	return m
}

func (m *CoreModule) registerWords() {
	// Stack operations
	m.AddModuleWord("POP", m.pop)
	m.AddModuleWord("DUP", m.dup)
	m.AddModuleWord("SWAP", m.swap)

	// Variable operations
	m.AddModuleWord("VARIABLES", m.variables)
	m.AddModuleWord("!", m.set)
	m.AddModuleWord("@", m.get)
	m.AddModuleWord("!@", m.setGet)

	// Module operations
	m.AddModuleWord("EXPORT", m.export_word)
	m.AddModuleWord("USE-MODULES", m.useModules)

	// Execution
	m.AddModuleWord("INTERPRET", m.interpret)

	// Control flow
	m.AddModuleWord("IDENTITY", m.identity)
	m.AddModuleWord("NOP", m.nop)
	m.AddModuleWord("NULL", m.null)
	m.AddModuleWord("ARRAY?", m.arrayCheck)
	m.AddModuleWord("DEFAULT", m.default_word)
	m.AddModuleWord("*DEFAULT", m.defaultStar)

	// Options
	m.AddModuleWord("~>", m.toOptions)

	// Profiling
	m.AddModuleWord("PROFILE-START", m.profileStart)
	m.AddModuleWord("PROFILE-END", m.profileEnd)
	m.AddModuleWord("PROFILE-TIMESTAMP", m.profileTimestamp)
	m.AddModuleWord("PROFILE-DATA", m.profileData)

	// Logging
	m.AddModuleWord("START-LOG", m.startLog)
	m.AddModuleWord("END-LOG", m.endLog)

	// String operations
	m.AddModuleWord("INTERPOLATE", m.interpolate)
	m.AddModuleWord("PRINT", m.print)

	// Debug
	m.AddModuleWord("PEEK!", m.peek)
	m.AddModuleWord("STACK!", m.stackDebug)
}

// getOrCreateVariable gets or creates a variable, validating the name
func getOrCreateVariable(interp *forthic.Interpreter, name string) (*forthic.Variable, error) {
	// Validate variable name - no __ prefix allowed
	if strings.HasPrefix(name, "__") {
		return nil, forthic.NewInvalidVariableNameError(name)
	}

	curModule := interp.CurModule()

	// Check if variable already exists
	variable := curModule.GetVariable(name)

	// Create it if it doesn't exist
	if variable == nil {
		curModule.AddVariable(name, nil)
		variable = curModule.GetVariable(name)
	}

	return variable, nil
}

// ========================================
// Stack Operations
// ========================================

func (m *CoreModule) pop(interp *forthic.Interpreter) error {
	interp.StackPop()
	return nil
}

func (m *CoreModule) dup(interp *forthic.Interpreter) error {
	a := interp.StackPop()
	interp.StackPush(a)
	interp.StackPush(a)
	return nil
}

func (m *CoreModule) swap(interp *forthic.Interpreter) error {
	b := interp.StackPop()
	a := interp.StackPop()
	interp.StackPush(b)
	interp.StackPush(a)
	return nil
}

// ========================================
// Variable Operations
// ========================================

func (m *CoreModule) variables(interp *forthic.Interpreter) error {
	varnames := interp.StackPop()
	curModule := interp.CurModule()

	if arr, ok := varnames.([]interface{}); ok {
		for _, v := range arr {
			if varName, ok := v.(string); ok {
				// Validate variable name
				if strings.HasPrefix(varName, "__") {
					return forthic.NewInvalidVariableNameError(varName)
				}
				curModule.AddVariable(varName, nil)
			}
		}
	}
	return nil
}

func (m *CoreModule) set(interp *forthic.Interpreter) error {
	variable := interp.StackPop()
	value := interp.StackPop()

	var varObj *forthic.Variable
	var err error

	if varName, ok := variable.(string); ok {
		// Auto-create variable if string name
		varObj, err = getOrCreateVariable(interp, varName)
		if err != nil {
			return err
		}
	} else {
		// Use existing variable object
		varObj = variable.(*forthic.Variable)
	}

	varObj.SetValue(value)
	return nil
}

func (m *CoreModule) get(interp *forthic.Interpreter) error {
	variable := interp.StackPop()

	var varObj *forthic.Variable
	var err error

	if varName, ok := variable.(string); ok {
		// Auto-create variable if string name
		varObj, err = getOrCreateVariable(interp, varName)
		if err != nil {
			return err
		}
	} else {
		// Use existing variable object
		varObj = variable.(*forthic.Variable)
	}

	interp.StackPush(varObj.GetValue())
	return nil
}

func (m *CoreModule) setGet(interp *forthic.Interpreter) error {
	variable := interp.StackPop()
	value := interp.StackPop()

	var varObj *forthic.Variable
	var err error

	if varName, ok := variable.(string); ok {
		// Auto-create variable if string name
		varObj, err = getOrCreateVariable(interp, varName)
		if err != nil {
			return err
		}
	} else {
		// Use existing variable object
		varObj = variable.(*forthic.Variable)
	}

	varObj.SetValue(value)
	interp.StackPush(varObj.GetValue())
	return nil
}

// ========================================
// Module Operations
// ========================================

func (m *CoreModule) export_word(interp *forthic.Interpreter) error {
	names := interp.StackPop()
	if arr, ok := names.([]interface{}); ok {
		strNames := make([]string, 0, len(arr))
		for _, name := range arr {
			if str, ok := name.(string); ok {
				strNames = append(strNames, str)
			}
		}
		interp.CurModule().AddExportable(strNames)
	}
	return nil
}

func (m *CoreModule) useModules(interp *forthic.Interpreter) error {
	names := interp.StackPop()
	if names == nil {
		return nil
	}
	if arr, ok := names.([]interface{}); ok {
		return interp.UseModules(arr)
	}
	return nil
}

// ========================================
// Execution
// ========================================

func (m *CoreModule) interpret(interp *forthic.Interpreter) error {
	str := interp.StackPop()
	if str == nil {
		return nil
	}
	if code, ok := str.(string); ok {
		return interp.Run(code)
	}
	return nil
}

// ========================================
// Control Flow
// ========================================

func (m *CoreModule) identity(interp *forthic.Interpreter) error {
	// No-op
	return nil
}

func (m *CoreModule) nop(interp *forthic.Interpreter) error {
	// No-op
	return nil
}

func (m *CoreModule) null(interp *forthic.Interpreter) error {
	interp.StackPush(nil)
	return nil
}

func (m *CoreModule) arrayCheck(interp *forthic.Interpreter) error {
	value := interp.StackPop()
	_, isArray := value.([]interface{})
	interp.StackPush(isArray)
	return nil
}

func (m *CoreModule) default_word(interp *forthic.Interpreter) error {
	defaultValue := interp.StackPop()
	value := interp.StackPop()

	if value == nil || value == "" {
		interp.StackPush(defaultValue)
	} else {
		interp.StackPush(value)
	}
	return nil
}

func (m *CoreModule) defaultStar(interp *forthic.Interpreter) error {
	defaultForthic := interp.StackPop()
	value := interp.StackPop()

	if value == nil || value == "" {
		if code, ok := defaultForthic.(string); ok {
			err := interp.Run(code)
			if err != nil {
				return err
			}
			result := interp.StackPop()
			interp.StackPush(result)
			return nil
		}
	}
	interp.StackPush(value)
	return nil
}

// ========================================
// Options
// ========================================

func (m *CoreModule) toOptions(interp *forthic.Interpreter) error {
	array := interp.StackPop()
	opts, err := forthic.NewWordOptions(array)
	if err != nil {
		return err
	}
	interp.StackPush(opts)
	return nil
}

// ========================================
// Profiling (Placeholder implementations)
// ========================================

func (m *CoreModule) profileStart(interp *forthic.Interpreter) error {
	// TODO: Implement profiling in interpreter
	return nil
}

func (m *CoreModule) profileEnd(interp *forthic.Interpreter) error {
	// TODO: Implement profiling in interpreter
	return nil
}

func (m *CoreModule) profileTimestamp(interp *forthic.Interpreter) error {
	label := interp.StackPop()
	_ = label // TODO: Implement profiling
	return nil
}

func (m *CoreModule) profileData(interp *forthic.Interpreter) error {
	// TODO: Implement profiling
	result := map[string]interface{}{
		"word_counts": []interface{}{},
		"timestamps":  []interface{}{},
	}
	interp.StackPush(result)
	return nil
}

// ========================================
// Logging (Placeholder implementations)
// ========================================

func (m *CoreModule) startLog(interp *forthic.Interpreter) error {
	// TODO: Implement logging in interpreter
	return nil
}

func (m *CoreModule) endLog(interp *forthic.Interpreter) error {
	// TODO: Implement logging in interpreter
	return nil
}

// ========================================
// String Operations
// ========================================

func (m *CoreModule) interpolate(interp *forthic.Interpreter) error {
	// Pop options if present
	topVal := interp.StackPop()
	var str string
	var opts *forthic.WordOptions

	// Check if we have options
	if optsVal, ok := topVal.(*forthic.WordOptions); ok {
		opts = optsVal
		str, _ = interp.StackPop().(string)
	} else {
		str, _ = topVal.(string)
		opts, _ = forthic.NewWordOptions([]interface{}{})
	}

	separator, _ := opts.Get("separator", ", ").(string)
	nullText, _ := opts.Get("null_text", "null").(string)
	useJSON, _ := opts.Get("json", false).(bool)

	result := interpolateString(interp, str, separator, nullText, useJSON)
	interp.StackPush(result)
	return nil
}

func (m *CoreModule) print(interp *forthic.Interpreter) error {
	// Pop options if present
	topVal := interp.StackPop()
	var value interface{}
	var opts *forthic.WordOptions

	// Check if we have options
	if optsVal, ok := topVal.(*forthic.WordOptions); ok {
		opts = optsVal
		value = interp.StackPop()
	} else {
		value = topVal
		opts, _ = forthic.NewWordOptions([]interface{}{})
	}

	separator, _ := opts.Get("separator", ", ").(string)
	nullText, _ := opts.Get("null_text", "null").(string)
	useJSON, _ := opts.Get("json", false).(bool)

	var result string
	if str, ok := value.(string); ok {
		// String: interpolate variables
		result = interpolateString(interp, str, separator, nullText, useJSON)
	} else {
		// Non-string: format directly
		result = valueToString(value, separator, nullText, useJSON)
	}

	fmt.Println(result)
	return nil
}

func interpolateString(interp *forthic.Interpreter, str string, separator string, nullText string, useJSON bool) string {
	if str == "" {
		return ""
	}

	// Handle escape sequences by replacing \. with a temporary placeholder
	escaped := strings.ReplaceAll(str, "\\.", "\x00ESCAPED_DOT\x00")

	// Replace whitespace-preceded or start-of-string .variable patterns
	re := regexp.MustCompile(`(^|\s)\.([a-zA-Z_][a-zA-Z0-9_-]*)`)
	interpolated := re.ReplaceAllStringFunc(escaped, func(match string) string {
		// Extract variable name (skip leading space if present and the dot)
		trimmed := strings.TrimSpace(match)
		varName := trimmed[1:] // Remove leading .

		variable, err := getOrCreateVariable(interp, varName)
		if err != nil {
			return match // Return original if error
		}
		value := variable.GetValue()

		// Preserve leading whitespace if it was there
		if strings.HasPrefix(match, " ") || strings.HasPrefix(match, "\t") {
			return string(match[0]) + valueToString(value, separator, nullText, useJSON)
		}
		return valueToString(value, separator, nullText, useJSON)
	})

	// Restore escaped dots
	return strings.ReplaceAll(interpolated, "\x00ESCAPED_DOT\x00", ".")
}

func valueToString(value interface{}, separator string, nullText string, useJSON bool) string {
	if value == nil {
		return nullText
	}
	if useJSON {
		bytes, _ := json.Marshal(value)
		return string(bytes)
	}
	if arr, ok := value.([]interface{}); ok {
		strs := make([]string, len(arr))
		for i, v := range arr {
			strs[i] = valueToString(v, separator, nullText, false)
		}
		return strings.Join(strs, separator)
	}
	if _, ok := value.(map[string]interface{}); ok {
		bytes, _ := json.Marshal(value)
		return string(bytes)
	}
	return fmt.Sprintf("%v", value)
}

// ========================================
// Debug Operations
// ========================================

func (m *CoreModule) peek(interp *forthic.Interpreter) error {
	stack := interp.GetStack()
	items := stack.Items()
	if len(items) > 0 {
		fmt.Println(items[len(items)-1])
	} else {
		fmt.Println("<STACK EMPTY>")
	}
	return forthic.NewIntentionalStopError("PEEK!")
}

func (m *CoreModule) stackDebug(interp *forthic.Interpreter) error {
	stack := interp.GetStack()
	items := stack.Items()

	// Reverse the items
	reversed := make([]interface{}, len(items))
	for i := 0; i < len(items); i++ {
		reversed[i] = items[len(items)-1-i]
	}

	bytes, _ := json.MarshalIndent(reversed, "", "  ")
	fmt.Println(string(bytes))
	return forthic.NewIntentionalStopError("STACK!")
}
