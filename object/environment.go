package object

// NewEnclosedEnvironment makes creating such an enclosed environment easy. The Get method has also been changed.
// It checks the enclosing environment for the given name.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewEnvironment returns a new Environment
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get is an Environment getter
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set is an Environment setter
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
