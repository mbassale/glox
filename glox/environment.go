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
	val, ok := env.values[name]
	if !ok {
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return val, nil
}

func (env *Environment) Assign(name string, value interface{}) error {
	_, ok := env.values[name]
	if !ok {
		return fmt.Errorf("undefined variable: %s", name)
	}
	env.values[name] = value
	return nil
}
