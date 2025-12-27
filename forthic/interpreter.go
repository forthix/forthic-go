package forthic

import (
	"fmt"
	"time"
)

// LiteralHandler tries to parse a string as a literal value
// Returns value and true if successful, nil and false otherwise
type LiteralHandler func(string) (interface{}, bool)

// Interpreter - Core Forthic interpreter
//
// Core interpreter that tokenizes and executes Forthic code.
// Manages the data stack, module stack, and execution context.
type Interpreter struct {
	stack           *Stack
	appModule       *Module
	moduleStack     []*Module
	registeredMods  map[string]*Module
	tokenizerStack  []*Tokenizer
	previousToken   *Token
	isCompiling     bool
	isMemoDefinition bool
	curDefinition   *DefinitionWord
	literalHandlers []LiteralHandler
	timezone        string
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(modules ...*Module) *Interpreter {
	interp := &Interpreter{
		stack:           NewStack(),
		appModule:       NewModule(""),
		moduleStack:     make([]*Module, 0),
		registeredMods:  make(map[string]*Module),
		tokenizerStack:  make([]*Tokenizer, 0),
		previousToken:   nil,
		isCompiling:     false,
		isMemoDefinition: false,
		curDefinition:   nil,
		literalHandlers: make([]LiteralHandler, 0),
		timezone:        "UTC",
	}

	// Set app module's interpreter
	interp.appModule.SetInterp(interp)

	// Initialize module stack with app module
	interp.moduleStack = append(interp.moduleStack, interp.appModule)

	// Register standard literal handlers
	interp.registerStandardLiterals()

	// Import provided modules (unprefixed)
	for _, module := range modules {
		interp.ImportModule(module, "")
	}

	return interp
}

// ============================================================================
// Stack Operations
// ============================================================================

// StackPush pushes a value onto the stack
func (i *Interpreter) StackPush(val interface{}) {
	i.stack.Push(val)
}

// StackPop pops a value from the stack
// Throws StackUnderflowError if stack is empty
func (i *Interpreter) StackPop() interface{} {
	val, err := i.stack.Pop()
	if err != nil {
		// Get token location if available
		var loc *CodeLocation
		if len(i.tokenizerStack) > 0 {
			tokenizer := i.GetTokenizer()
			loc = tokenizer.getTokenLocation()
		}
		panic(NewStackUnderflowError().WithLocation(loc))
	}
	return val
}

// StackPeek peeks at the top of the stack without removing it
func (i *Interpreter) StackPeek() interface{} {
	val, err := i.stack.Peek()
	if err != nil {
		// Get token location if available
		var loc *CodeLocation
		if len(i.tokenizerStack) > 0 {
			tokenizer := i.GetTokenizer()
			loc = tokenizer.getTokenLocation()
		}
		panic(NewStackUnderflowError().WithLocation(loc))
	}
	return val
}

// GetStack returns the stack
func (i *Interpreter) GetStack() *Stack {
	return i.stack
}

// ============================================================================
// Module Operations
// ============================================================================

// GetAppModule returns the app module
func (i *Interpreter) GetAppModule() *Module {
	return i.appModule
}

// CurModule returns the current module (top of module stack)
func (i *Interpreter) CurModule() *Module {
	return i.moduleStack[len(i.moduleStack)-1]
}

// ModuleStackPush pushes a module onto the module stack
func (i *Interpreter) ModuleStackPush(module *Module) {
	i.moduleStack = append(i.moduleStack, module)
}

// ModuleStackPop pops a module from the module stack
func (i *Interpreter) ModuleStackPop() *Module {
	if len(i.moduleStack) <= 1 {
		panic(NewForthicError("Cannot pop app module from module stack"))
	}
	module := i.moduleStack[len(i.moduleStack)-1]
	i.moduleStack = i.moduleStack[:len(i.moduleStack)-1]
	return module
}

// RegisterModule registers a module with the interpreter
func (i *Interpreter) RegisterModule(module *Module) {
	i.registeredMods[module.name] = module
	module.SetInterp(i)
}

// FindModule finds a registered module by name
func (i *Interpreter) FindModule(name string) (*Module, error) {
	module, ok := i.registeredMods[name]
	if !ok {
		return nil, NewUnknownModuleError(name)
	}
	return module, nil
}

// UseModules imports modules into the app module
// names can be strings or [string, string] pairs (module_name, prefix)
func (i *Interpreter) UseModules(names []interface{}) error {
	for _, name := range names {
		moduleName := ""
		prefix := ""

		// Check if it's an array [module_name, prefix]
		if arr, ok := name.([]interface{}); ok {
			if len(arr) >= 1 {
				moduleName = arr[0].(string)
			}
			if len(arr) >= 2 {
				prefix = arr[1].(string)
			}
		} else {
			// Simple string name
			moduleName = name.(string)
		}

		module, err := i.FindModule(moduleName)
		if err != nil {
			return err
		}

		i.appModule.ImportModule(prefix, module, i)
	}
	return nil
}

// ImportModule registers and imports a module
func (i *Interpreter) ImportModule(module *Module, prefix string) {
	i.RegisterModule(module)
	i.appModule.ImportModule(prefix, module, i)
}

// ============================================================================
// Tokenizer Operations
// ============================================================================

// GetTokenizer returns the current tokenizer
func (i *Interpreter) GetTokenizer() *Tokenizer {
	return i.tokenizerStack[len(i.tokenizerStack)-1]
}

// ============================================================================
// Literal Handlers
// ============================================================================

// registerStandardLiterals registers the standard literal handlers
func (i *Interpreter) registerStandardLiterals() {
	// Load timezone
	loc, err := time.LoadLocation(i.timezone)
	if err != nil {
		loc = time.UTC // Fallback to UTC if timezone is invalid
	}

	// Order matters: more specific handlers first
	i.literalHandlers = []LiteralHandler{
		ToBool,
		ToFloat,
		ToZonedDateTime(loc),
		ToLiteralDate(loc),
		ToTime,
		ToInt,
	}
}

// RegisterLiteralHandler adds a custom literal handler
func (i *Interpreter) RegisterLiteralHandler(handler LiteralHandler) {
	// Add to front so it can override existing handlers
	i.literalHandlers = append([]LiteralHandler{handler}, i.literalHandlers...)
}

// findLiteralWord tries to parse a string as a literal
func (i *Interpreter) findLiteralWord(name string) Word {
	for _, handler := range i.literalHandlers {
		value, ok := handler(name)
		if ok {
			return NewPushValueWord(name, value)
		}
	}
	return nil
}

// ============================================================================
// Find Word
// ============================================================================

// FindWord finds a word by name
// Searches module stack from top to bottom, then checks literal handlers
func (i *Interpreter) FindWord(name string) (Word, error) {
	// 1. Check module stack (from top to bottom)
	for j := len(i.moduleStack) - 1; j >= 0; j-- {
		module := i.moduleStack[j]
		word := module.FindWord(name)
		if word != nil {
			return word, nil
		}
	}

	// 2. Check literal handlers
	word := i.findLiteralWord(name)
	if word != nil {
		return word, nil
	}

	// 3. Not found
	return nil, NewUnknownWordError(name)
}

// ============================================================================
// Main Execution
// ============================================================================

// Run executes Forthic code
func (i *Interpreter) Run(code string) error {
	tokenizer := NewTokenizer(code, nil, false)
	i.tokenizerStack = append(i.tokenizerStack, tokenizer)

	err := i.runWithTokenizer(tokenizer)

	i.tokenizerStack = i.tokenizerStack[:len(i.tokenizerStack)-1]
	return err
}

// runWithTokenizer executes code using the given tokenizer
func (i *Interpreter) runWithTokenizer(tokenizer *Tokenizer) error {
	for {
		token, err := tokenizer.NextToken()
		if err != nil {
			return err
		}

		err = i.handleToken(token)
		if err != nil {
			return err
		}

		if token.Type == TOKEN_EOS {
			break
		}

		i.previousToken = token
	}
	return nil
}

// ============================================================================
// Token Handling
// ============================================================================

// handleToken dispatches token to appropriate handler
func (i *Interpreter) handleToken(token *Token) error {
	switch token.Type {
	case TOKEN_STRING:
		return i.handleStringToken(token)
	case TOKEN_COMMENT:
		return i.handleCommentToken(token)
	case TOKEN_START_ARRAY:
		return i.handleStartArrayToken(token)
	case TOKEN_END_ARRAY:
		return i.handleEndArrayToken(token)
	case TOKEN_START_MODULE:
		return i.handleStartModuleToken(token)
	case TOKEN_END_MODULE:
		return i.handleEndModuleToken(token)
	case TOKEN_START_DEF:
		return i.handleStartDefinitionToken(token)
	case TOKEN_START_MEMO:
		return i.handleStartMemoToken(token)
	case TOKEN_END_DEF:
		return i.handleEndDefinitionToken(token)
	case TOKEN_DOT_SYMBOL:
		return i.handleDotSymbolToken(token)
	case TOKEN_WORD:
		return i.handleWordToken(token)
	case TOKEN_EOS:
		if i.isCompiling {
			if i.previousToken != nil {
				return NewMissingSemicolonError().WithLocation(i.previousToken.Location)
			}
			return NewMissingSemicolonError()
		}
		return nil
	default:
		return NewForthicError(fmt.Sprintf("Unknown token type: %v", token.Type))
	}
}

// handleStringToken handles string literals
func (i *Interpreter) handleStringToken(token *Token) error {
	word := NewPushValueWord("<string>", token.String)
	return i.handleWord(word, token.Location)
}

// handleDotSymbolToken handles dot symbols
func (i *Interpreter) handleDotSymbolToken(token *Token) error {
	word := NewPushValueWord("<dot-symbol>", token.String)
	return i.handleWord(word, token.Location)
}

// handleCommentToken handles comments (no-op)
func (i *Interpreter) handleCommentToken(token *Token) error {
	return nil
}

// handleStartArrayToken handles [
func (i *Interpreter) handleStartArrayToken(token *Token) error {
	word := NewPushValueWord("<start_array_token>", token)
	return i.handleWord(word, token.Location)
}

// handleEndArrayToken handles ]
func (i *Interpreter) handleEndArrayToken(token *Token) error {
	word := NewEndArrayWord()
	return i.handleWord(word, token.Location)
}

// handleStartModuleToken handles {
func (i *Interpreter) handleStartModuleToken(token *Token) error {
	word := NewStartModuleWord(token.String)

	// Module words are immediate (execute during compilation) and also compiled
	if i.isCompiling {
		i.curDefinition.words = append(i.curDefinition.words, word)
	}

	return word.Execute(i)
}

// handleEndModuleToken handles }
func (i *Interpreter) handleEndModuleToken(token *Token) error {
	word := NewEndModuleWord()

	// Module words are immediate (execute during compilation) and also compiled
	if i.isCompiling {
		i.curDefinition.words = append(i.curDefinition.words, word)
	}

	return word.Execute(i)
}

// handleStartDefinitionToken handles :
func (i *Interpreter) handleStartDefinitionToken(token *Token) error {
	if i.isCompiling {
		return NewMissingSemicolonError().WithLocation(i.previousToken.Location)
	}
	i.curDefinition = NewDefinitionWord(token.String, nil)
	i.isCompiling = true
	i.isMemoDefinition = false
	return nil
}

// handleStartMemoToken handles @:
func (i *Interpreter) handleStartMemoToken(token *Token) error {
	if i.isCompiling {
		return NewMissingSemicolonError().WithLocation(i.previousToken.Location)
	}
	i.curDefinition = NewDefinitionWord(token.String, nil)
	i.isCompiling = true
	i.isMemoDefinition = true
	return nil
}

// handleEndDefinitionToken handles ;
func (i *Interpreter) handleEndDefinitionToken(token *Token) error {
	if !i.isCompiling || i.curDefinition == nil {
		return NewExtraSemicolonError().WithLocation(token.Location)
	}

	if i.isMemoDefinition {
		i.CurModule().AddMemoWords(i.curDefinition)
	} else {
		i.CurModule().AddWord(i.curDefinition)
	}

	i.isCompiling = false
	return nil
}

// handleWordToken handles word tokens
func (i *Interpreter) handleWordToken(token *Token) error {
	word, err := i.FindWord(token.String)
	if err != nil {
		return err
	}
	return i.handleWord(word, token.Location)
}

// handleWord executes or compiles a word
func (i *Interpreter) handleWord(word Word, location *CodeLocation) error {
	if i.isCompiling {
		word.SetLocation(location)
		i.curDefinition.words = append(i.curDefinition.words, word)
		return nil
	} else {
		return word.Execute(i)
	}
}

// ============================================================================
// Special Word Types
// ============================================================================

// StartModuleWord handles module creation and switching
type StartModuleWord struct {
	*BaseWord
}

// NewStartModuleWord creates a new StartModuleWord
func NewStartModuleWord(name string) *StartModuleWord {
	return &StartModuleWord{
		BaseWord: NewBaseWord(name),
	}
}

func (w *StartModuleWord) Execute(interp *Interpreter) error {
	// Empty name refers to app module
	if w.name == "" {
		interp.ModuleStackPush(interp.GetAppModule())
		return nil
	}

	// Check if module exists in current module
	module := interp.CurModule().FindModule(w.name)
	if module == nil {
		// Create new module
		module = NewModule(w.name)
		interp.CurModule().RegisterModule(w.name, w.name, module)

		// If we're at app module, also register with interpreter
		if interp.CurModule().name == "" {
			interp.RegisterModule(module)
		}
	}

	interp.ModuleStackPush(module)
	return nil
}

// EndModuleWord pops the current module
type EndModuleWord struct {
	*BaseWord
}

// NewEndModuleWord creates a new EndModuleWord
func NewEndModuleWord() *EndModuleWord {
	return &EndModuleWord{
		BaseWord: NewBaseWord("}"),
	}
}

func (w *EndModuleWord) Execute(interp *Interpreter) error {
	interp.ModuleStackPop()
	return nil
}

// EndArrayWord collects items into an array
type EndArrayWord struct {
	*BaseWord
}

// NewEndArrayWord creates a new EndArrayWord
func NewEndArrayWord() *EndArrayWord {
	return &EndArrayWord{
		BaseWord: NewBaseWord("]"),
	}
}

func (w *EndArrayWord) Execute(interp *Interpreter) error {
	items := make([]interface{}, 0)
	for {
		item := interp.StackPop()

		// Check if it's a START_ARRAY token
		if token, ok := item.(*Token); ok && token.Type == TOKEN_START_ARRAY {
			break
		}

		items = append(items, item)
	}

	// Reverse the items
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	interp.StackPush(items)
	return nil
}
