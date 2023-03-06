package evaluator

type Value interface{}

type Environment struct {
	store map[string]Value
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Value)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Value, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}
