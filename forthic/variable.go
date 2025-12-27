package forthic

// Variable - Named mutable value container
//
// Represents a variable that can store and retrieve values within a module scope.
// Variables are accessed by name and can be set to any value type.
type Variable struct {
	name  string
	value interface{}
}

// NewVariable creates a new Variable
func NewVariable(name string, value interface{}) *Variable {
	return &Variable{
		name:  name,
		value: value,
	}
}

// GetName returns the variable's name
func (v *Variable) GetName() string {
	return v.name
}

// SetValue sets the variable's value
func (v *Variable) SetValue(val interface{}) {
	v.value = val
}

// GetValue returns the variable's value
func (v *Variable) GetValue() interface{} {
	return v.value
}

// Dup creates a duplicate of the variable
func (v *Variable) Dup() *Variable {
	return NewVariable(v.name, v.value)
}
