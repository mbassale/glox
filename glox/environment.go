package glox

import "fmt"

type Environment struct {
	values map[string]interface{}
}

func NewEnvironment() Environment {
	return Environment{
		values: map[string]interface{}{},
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Get(name string) (interface{}, error) {
	val, exists := env.values[name]
	if !exists {
		return 0, fmt.Errorf("undefined variable: %s", name)
	}
	return val, nil
}
