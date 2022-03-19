package glox

import "fmt"

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func NewEnvironment() Environment {
	return Environment{
		values:    map[string]interface{}{},
		enclosing: nil,
	}
}

func NewEnvironmentWithEnclosing(enclosing *Environment) Environment {
	return Environment{
		values:    map[string]interface{}{},
		enclosing: enclosing,
	}
}

func (env *Environment) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Environment) Get(name string) (interface{}, error) {
	val, ok := env.values[name]
	if !ok {
		if env.enclosing != nil {
			return env.enclosing.Get(name)
		}
		return nil, fmt.Errorf("undefined variable: %s", name)
	}
	return val, nil
}

func (env *Environment) GetAt(distance int, name string) (interface{}, error) {
	value, ok := env.ancestor(distance).values[name]
	if !ok {
		return nil, fmt.Errorf("variable: %v not found at distance: %v", name, distance)
	}
	return value, nil
}

func (env *Environment) Assign(name string, value interface{}) error {
	_, ok := env.values[name]
	if !ok {
		if env.enclosing != nil {
			return env.enclosing.Assign(name, value)
		}
		return fmt.Errorf("undefined variable: %s", name)
	}
	env.values[name] = value
	return nil
}

func (env *Environment) AssignAt(distance int, name string, value interface{}) error {
	ancestorEnvironment := env.ancestor(distance)
	ancestorEnvironment.values[name] = value
	return nil
}

func (env *Environment) ancestor(distance int) *Environment {
	environment := env
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}
