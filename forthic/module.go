package forthic

// Module - Container for words, variables, and imported modules
//
// Modules provide namespacing and code organization in Forthic.
// Each module maintains its own dictionary of words, variables, and imported modules.
//
// Features:
// - Word and variable management
// - Module importing with optional prefixes
// - Exportable word lists for controlled visibility
// - Module duplication for isolated execution contexts
type Module struct {
	words          []Word
	exportable     []string
	variables      map[string]*Variable
	modules        map[string]*Module
	modulePrefixes map[string]map[string]bool // module_name -> set of prefixes
	name           string
	forthicCode    string
	interp         *Interpreter
}

// NewModule creates a new Module
func NewModule(name string, forthicCode ...string) *Module {
	code := ""
	if len(forthicCode) > 0 {
		code = forthicCode[0]
	}

	return &Module{
		words:          make([]Word, 0),
		exportable:     make([]string, 0),
		variables:      make(map[string]*Variable),
		modules:        make(map[string]*Module),
		modulePrefixes: make(map[string]map[string]bool),
		name:           name,
		forthicCode:    code,
		interp:         nil,
	}
}

// GetName returns the module's name
func (m *Module) GetName() string {
	return m.name
}

// SetInterp sets the interpreter for this module
func (m *Module) SetInterp(interp *Interpreter) {
	m.interp = interp
}

// GetInterp returns the interpreter for this module
func (m *Module) GetInterp() (*Interpreter, error) {
	if m.interp == nil {
		return nil, NewModuleError(m.name, "Module has no interpreter")
	}
	return m.interp, nil
}

// ============================================================================
// Duplication Methods
// ============================================================================

// Dup creates a duplicate of the module
func (m *Module) Dup() *Module {
	result := NewModule(m.name, m.forthicCode)

	// Copy words slice
	result.words = make([]Word, len(m.words))
	copy(result.words, m.words)

	// Copy exportable slice
	result.exportable = make([]string, len(m.exportable))
	copy(result.exportable, m.exportable)

	// Copy variables
	for key, variable := range m.variables {
		result.variables[key] = variable.Dup()
	}

	// Copy module references (shallow copy)
	for key, module := range m.modules {
		result.modules[key] = module
	}

	return result
}

// Copy creates a copy of the module with restored prefixes
func (m *Module) Copy(interp *Interpreter) *Module {
	result := m.Dup()

	// Restore module_prefixes
	for moduleName, prefixes := range m.modulePrefixes {
		module := m.modules[moduleName]
		for prefix := range prefixes {
			result.ImportModule(prefix, module, interp)
		}
	}

	return result
}

// ============================================================================
// Module Management
// ============================================================================

// FindModule finds a module by name
func (m *Module) FindModule(name string) *Module {
	return m.modules[name]
}

// RegisterModule registers a module with a prefix
func (m *Module) RegisterModule(moduleName string, prefix string, module *Module) {
	m.modules[moduleName] = module

	if m.modulePrefixes[moduleName] == nil {
		m.modulePrefixes[moduleName] = make(map[string]bool)
	}
	m.modulePrefixes[moduleName][prefix] = true
}

// ImportModule imports a module with optional prefix
// prefix can be "" for unprefixed imports
func (m *Module) ImportModule(prefix string, module *Module, interp *Interpreter) {
	newModule := module.Dup()

	words := newModule.ExportableWords()
	for _, word := range words {
		if prefix == "" {
			// Unprefixed import - add word directly
			m.AddWord(word)
		} else {
			// Prefixed import - create ExecuteWord
			prefixedWord := NewExecuteWord(prefix+"."+word.GetName(), word)
			m.AddWord(prefixedWord)
		}
	}

	m.RegisterModule(module.name, prefix, newModule)
}

// ============================================================================
// Word Management
// ============================================================================

// AddWord adds a word to the module
func (m *Module) AddWord(word Word) {
	m.words = append(m.words, word)
}

// AddMemoWords adds memo word and refresh variants
func (m *Module) AddMemoWords(word Word) *ModuleMemoWord {
	memoWord := NewModuleMemoWord(word)
	m.words = append(m.words, memoWord)
	m.words = append(m.words, NewModuleMemoBangWord(memoWord))
	m.words = append(m.words, NewModuleMemoBangAtWord(memoWord))
	return memoWord
}

// AddExportable adds word names to the exportable list
func (m *Module) AddExportable(names []string) {
	m.exportable = append(m.exportable, names...)
}

// AddExportableWord adds a word and marks it as exportable
func (m *Module) AddExportableWord(word Word) {
	m.words = append(m.words, word)
	m.exportable = append(m.exportable, word.GetName())
}

// AddModuleWord creates a ModuleWord and marks it as exportable
func (m *Module) AddModuleWord(wordName string, handler func(*Interpreter) error) {
	word := NewModuleWord(wordName, handler)
	m.AddExportableWord(word)
}

// ExportableWords returns all exportable words
func (m *Module) ExportableWords() []Word {
	result := make([]Word, 0)
	exportableMap := make(map[string]bool)
	for _, name := range m.exportable {
		exportableMap[name] = true
	}

	for _, word := range m.words {
		if exportableMap[word.GetName()] {
			result = append(result, word)
		}
	}

	return result
}

// FindWord finds a word by name (checks words then variables)
func (m *Module) FindWord(name string) Word {
	// Check dictionary words first
	word := m.FindDictionaryWord(name)
	if word != nil {
		return word
	}

	// Check variables
	word = m.FindVariable(name)
	return word
}

// FindDictionaryWord finds a word in the word dictionary
// Searches from end to beginning (last added word wins)
func (m *Module) FindDictionaryWord(wordName string) Word {
	for i := len(m.words) - 1; i >= 0; i-- {
		w := m.words[i]
		if w.GetName() == wordName {
			return w
		}
	}
	return nil
}

// FindVariable finds a variable and returns it as a PushValueWord
func (m *Module) FindVariable(varName string) Word {
	variable, ok := m.variables[varName]
	if ok {
		return NewPushValueWord(varName, variable)
	}
	return nil
}

// ============================================================================
// Variable Management
// ============================================================================

// AddVariable adds a variable to the module
func (m *Module) AddVariable(name string, value interface{}) {
	if m.variables[name] == nil {
		m.variables[name] = NewVariable(name, value)
	}
}

// GetVariable returns a variable by name
func (m *Module) GetVariable(name string) *Variable {
	return m.variables[name]
}

// ============================================================================
// Additional Word Types for Module System
// ============================================================================

// ExecuteWord - Wrapper word that executes another word
// Used for prefixed module imports (e.g., prefix.word)
type ExecuteWord struct {
	*BaseWord
	targetWord Word
}

// NewExecuteWord creates a new ExecuteWord
func NewExecuteWord(name string, targetWord Word) *ExecuteWord {
	return &ExecuteWord{
		BaseWord:   NewBaseWord(name),
		targetWord: targetWord,
	}
}

func (w *ExecuteWord) Execute(interp *Interpreter) error {
	return w.targetWord.Execute(interp)
}

func (w *ExecuteWord) GetRuntimeInfo() *RuntimeInfo {
	return w.targetWord.GetRuntimeInfo()
}

// ModuleMemoWord - Memoized word that caches its result
type ModuleMemoWord struct {
	*BaseWord
	word     Word
	hasValue bool
	value    interface{}
}

// NewModuleMemoWord creates a new ModuleMemoWord
func NewModuleMemoWord(word Word) *ModuleMemoWord {
	return &ModuleMemoWord{
		BaseWord: NewBaseWord(word.GetName()),
		word:     word,
		hasValue: false,
		value:    nil,
	}
}

func (w *ModuleMemoWord) Refresh(interp *Interpreter) error {
	err := w.word.Execute(interp)
	if err != nil {
		return err
	}
	w.value = interp.StackPop()
	w.hasValue = true
	return nil
}

func (w *ModuleMemoWord) Execute(interp *Interpreter) error {
	if !w.hasValue {
		err := w.Refresh(interp)
		if err != nil {
			return err
		}
	}
	interp.StackPush(w.value)
	return nil
}

// ModuleMemoBangWord - Forces refresh of a memoized word
type ModuleMemoBangWord struct {
	*BaseWord
	memoWord *ModuleMemoWord
}

// NewModuleMemoBangWord creates a new ModuleMemoBangWord
func NewModuleMemoBangWord(memoWord *ModuleMemoWord) *ModuleMemoBangWord {
	return &ModuleMemoBangWord{
		BaseWord: NewBaseWord(memoWord.GetName() + "!"),
		memoWord: memoWord,
	}
}

func (w *ModuleMemoBangWord) Execute(interp *Interpreter) error {
	return w.memoWord.Refresh(interp)
}

// ModuleMemoBangAtWord - Refreshes a memoized word and returns its value
type ModuleMemoBangAtWord struct {
	*BaseWord
	memoWord *ModuleMemoWord
}

// NewModuleMemoBangAtWord creates a new ModuleMemoBangAtWord
func NewModuleMemoBangAtWord(memoWord *ModuleMemoWord) *ModuleMemoBangAtWord {
	return &ModuleMemoBangAtWord{
		BaseWord: NewBaseWord(memoWord.GetName() + "!@"),
		memoWord: memoWord,
	}
}

func (w *ModuleMemoBangAtWord) Execute(interp *Interpreter) error {
	err := w.memoWord.Refresh(interp)
	if err != nil {
		return err
	}
	interp.StackPush(w.memoWord.value)
	return nil
}
