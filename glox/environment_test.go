package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalEnvironment(t *testing.T) {
	env := NewEnvironment()

	val, err := env.Get("non_existent_var")
	assert.Nil(t, val)
	assert.EqualError(t, err, "undefined variable: non_existent_var")

	env.Define("existent_var", true)
	val, err = env.Get("existent_var")
	assert.Equal(t, val, true)
	assert.Nil(t, err)

	err = env.Assign("non_existent_var", true)
	assert.EqualError(t, err, "undefined variable: non_existent_var")

	err = env.Assign("existent_var", "test")
	assert.Nil(t, err)
	val, err = env.Get("existent_var")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}

func TestEnclosingEnvironment(t *testing.T) {
	globalEnv := NewEnvironment()
	globalEnv.Define("global_var", true)
	globalEnv.Define("global_var2", true)

	localEnv := NewEnvironmentWithEnclosing(&globalEnv)
	// define a local variable that shadows a global variable
	localEnv.Define("global_var", "test")
	// define a local variable
	localEnv.Define("local_var", true)

	// get local variable that shadows a global variable
	val, err := localEnv.Get("global_var")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)

	// get enclosing environment variable
	val, err = localEnv.Get("global_var2")
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	// get local variable
	val, err = localEnv.Get("local_var")
	assert.Nil(t, err)
	assert.Equal(t, true, val)

	// assign global variable on local environment
	localEnv.Assign("global_var2", "test")
	val, err = localEnv.Get("global_var2")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)

	// assign local variable on local environment
	localEnv.Assign("local_var", "test")
	val, err = localEnv.Get("local_var")
	assert.Nil(t, err)
	assert.Equal(t, "test", val)
}
