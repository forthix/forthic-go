package forthic

// Word - Base class for all executable words in Forthic
//
// A word is the fundamental unit of execution in Forthic. When interpreted,
// it performs an action (typically manipulating the stack or control flow).
// All concrete word types must override the Execute method.
type Word interface {
	Execute(interp *Interpreter) error
	GetName() string
	GetString() string
	GetLocation() *CodeLocation
	SetLocation(location *CodeLocation)
	AddErrorHandler(handler WordErrorHandler)
	RemoveErrorHandler(handler WordErrorHandler)
	ClearErrorHandlers()
	GetErrorHandlers() []WordErrorHandler
	GetRuntimeInfo() *RuntimeInfo
}

// WordErrorHandler is a function that handles errors during word execution
// Returns nil if error was handled, or returns error if it should propagate
type WordErrorHandler func(error, Word, *Interpreter) error

// RuntimeInfo - Metadata about where and how a word can execute
//
// Used by the ExecutionPlanner to batch remote word execution efficiently.
// Standard library words can execute in any runtime, while runtime-specific
// words (like RemoteWord) can only execute in their designated runtime.
type RuntimeInfo struct {
	Runtime     string   // "local" | "python" | "ruby" | "rust" | etc.
	IsRemote    bool     // True if this word requires remote execution
	IsStandard  bool     // True if this is a standard library word (available in all runtimes)
	AvailableIn []string // List of runtimes where this word is available
}

// BaseWord provides default implementation of Word interface
type BaseWord struct {
	name          string
	str           string
	location      *CodeLocation
	errorHandlers []WordErrorHandler
}

// NewBaseWord creates a new BaseWord
func NewBaseWord(name string) *BaseWord {
	return &BaseWord{
		name:          name,
		str:           name,
		location:      nil,
		errorHandlers: make([]WordErrorHandler, 0),
	}
}

func (w *BaseWord) Execute(interp *Interpreter) error {
	return NewForthicError("Must override Word.Execute")
}

func (w *BaseWord) GetName() string {
	return w.name
}

func (w *BaseWord) GetString() string {
	return w.str
}

func (w *BaseWord) GetLocation() *CodeLocation {
	return w.location
}

func (w *BaseWord) SetLocation(location *CodeLocation) {
	w.location = location
}

func (w *BaseWord) AddErrorHandler(handler WordErrorHandler) {
	w.errorHandlers = append(w.errorHandlers, handler)
}

func (w *BaseWord) RemoveErrorHandler(handler WordErrorHandler) {
	// In Go, we can't directly compare function pointers for equality
	// This is a limitation - we'll need to use a different approach
	// For now, this is a no-op. Users should use ClearErrorHandlers instead.
	// TODO: Consider using a handle/ID system for removable handlers
}

func (w *BaseWord) ClearErrorHandlers() {
	w.errorHandlers = make([]WordErrorHandler, 0)
}

func (w *BaseWord) GetErrorHandlers() []WordErrorHandler {
	// Return a copy
	result := make([]WordErrorHandler, len(w.errorHandlers))
	copy(result, w.errorHandlers)
	return result
}

// TryErrorHandlers tries error handlers in order
// Returns nil if error was handled, otherwise returns error
func (w *BaseWord) TryErrorHandlers(err error, word Word, interp *Interpreter) error {
	// Check if error is IntentionalStopError - if so, bypass handlers
	if _, ok := err.(*IntentionalStopError); ok {
		return err
	}

	for _, handler := range w.errorHandlers {
		handlerErr := handler(err, word, interp)
		if handlerErr == nil {
			// Handler succeeded, error is handled
			return nil
		}
		// Handler failed, try next one
	}
	// No handler succeeded
	return err
}

func (w *BaseWord) GetRuntimeInfo() *RuntimeInfo {
	return &RuntimeInfo{
		Runtime:     "local",
		IsRemote:    false,
		IsStandard:  false,
		AvailableIn: []string{"go"},
	}
}

// ============================================================================
// Concrete Word Types
// ============================================================================

// PushValueWord - Word that pushes a value onto the stack
type PushValueWord struct {
	*BaseWord
	value interface{}
}

// NewPushValueWord creates a new PushValueWord
func NewPushValueWord(name string, value interface{}) *PushValueWord {
	return &PushValueWord{
		BaseWord: NewBaseWord(name),
		value:    value,
	}
}

func (w *PushValueWord) Execute(interp *Interpreter) error {
	interp.StackPush(w.value)
	return nil
}

// ModuleWord - Word that wraps a function with error handler support
type ModuleWord struct {
	*BaseWord
	handler func(*Interpreter) error
}

// NewModuleWord creates a new ModuleWord
func NewModuleWord(name string, handler func(*Interpreter) error) *ModuleWord {
	return &ModuleWord{
		BaseWord: NewBaseWord(name),
		handler:  handler,
	}
}

func (w *ModuleWord) Execute(interp *Interpreter) error {
	err := w.handler(interp)
	if err != nil {
		// Try error handlers
		if handledErr := w.TryErrorHandlers(err, w, interp); handledErr == nil {
			return nil
		}
		return err
	}
	return nil
}

// DefinitionWord - Word defined by a sequence of other words
type DefinitionWord struct {
	*BaseWord
	words []Word
}

// NewDefinitionWord creates a new DefinitionWord
func NewDefinitionWord(name string, words []Word) *DefinitionWord {
	return &DefinitionWord{
		BaseWord: NewBaseWord(name),
		words:    words,
	}
}

func (w *DefinitionWord) Execute(interp *Interpreter) error {
	for _, word := range w.words {
		err := word.Execute(interp)
		if err != nil {
			// Try error handlers
			if handledErr := w.TryErrorHandlers(err, w, interp); handledErr == nil {
				continue
			}
			return err
		}
	}
	return nil
}

func (w *DefinitionWord) GetWords() []Word {
	return w.words
}
