package driver

import (
	"fmt"

	"github.com/connormckelvey/tmplrun/evaluator"
	"github.com/dop251/goja"
)

type gojaDriver struct{}

func NewGoja() *gojaDriver {
	return new(gojaDriver)
}

func (gd *gojaDriver) CreateContext(env *evaluator.Environment) (evaluator.DriverContext, error) {
	vm := &gojaContext{goja.New()}
	for k, v := range env.Props() {
		if err := vm.runtime.Set(k, v); err != nil {
			return nil, err
		}
	}
	if err := vm.runtime.Set("include", env.Include); err != nil {
		return nil, err
	}
	if err := vm.runtime.Set("template", env.Template); err != nil {
		return nil, err
	}
	if err := vm.runtime.Set("log", func(args ...any) {
		fmt.Println(args)
	}); err != nil {
		return nil, err
	}
	return vm, nil
}

type gojaContext struct {
	runtime *goja.Runtime
}

func (gc *gojaContext) Eval(code string) (string, error) {
	v, err := gc.runtime.RunString(code)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(v.Export()), nil
}
