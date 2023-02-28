package evaluator

type Environment struct {
	store map[string]int
}

func NewEnvironment() *Environment {
	s := make(map[string]int)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (int, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val int) int {
	e.store[name] = val
	return val
}
