package driver

import (
	"fmt"

	"github.com/connormckelvey/tmplrun/evaluator"
	"github.com/dop251/goja"
)

// gojaDriver represents a driver implementation using the Goja JavaScript engine.
type gojaDriver struct{}

// NewGoja creates a new instance of the Goja driver.
func NewGoja() *gojaDriver {
	return new(gojaDriver)
}

// CreateContext creates a new evaluation context using Goja.
func (gd *gojaDriver) CreateContext(env *evaluator.Environment) (evaluator.DriverContext, error) {
	vm := &gojaContext{goja.New()}
	// Set properties from the environment.
	for k, v := range env.Props() {
		if err := vm.runtime.Set(k, v); err != nil {
			return nil, err
		}
	}

	err := env.RegisterHooks(func(name string, value any) error {
		return vm.runtime.Set(name, value)
	})

	if err != nil {
		return nil, err
	}
	// Set built-in functions.
	if err := vm.runtime.Set("include", env.Include); err != nil {
		return nil, err
	}
	if err := vm.runtime.Set("template", env.Template); err != nil {
		return nil, err
	}
	if err := vm.runtime.Set("log", func(args ...interface{}) {
		fmt.Println(args...)
	}); err != nil {
		return nil, err
	}

	return vm, nil
}

type gojaContext struct {
	runtime *goja.Runtime
}

// Eval evaluates code within the Goja context.
func (gc *gojaContext) Eval(code string) (string, error) {
	// Run the code and export the result.
	v, err := gc.runtime.RunString(code)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(v.Export()), nil
}
