package evaluator

import (
	"fmt"
	"io/fs"
)

// Environment represents the environment for template evaluation.
type Environment struct {
	props  map[string]any
	fs     fs.FS
	errors []error
	hooks  HooksAPI
}

// NewEnvironment creates a new instance of Environment with the given parameters.
func NewEnvironment(fsys fs.FS, props map[string]any, hooksAPI HooksAPI) *Environment {
	return &Environment{
		props: props,
		fs:    fsys,
		hooks: hooksAPI,
	}
}

// Props returns the properties associated with the environment.
func (env *Environment) Props() map[string]any {
	return env.props
}

// Include includes the template specified by name.
func (env *Environment) Include(name string) string {
	file, err := env.hooks.Include(name)
	if err != nil {
		env.errors = append(env.errors, fmt.Errorf("no file %s", name))
		return ""
	}
	return string(file)
}

func (env *Environment) hasErrors() bool {
	return len(env.errors) > 0
}

// Template renders the template specified by name with the given properties.
func (env *Environment) Template(name string, props map[string]any) string {
	res, err := env.hooks.Render(name, props)
	if err != nil {
		env.errors = append(env.errors, err)
		return ""
	}

	return res
}

func (env *Environment) RegisterHooks(register func(name string, value any) error) error {
	reg, ok := env.hooks.(RegisterHooks)
	if !ok {
		return nil
	}
	if err := reg.Register(register); err != nil {
		return err
	}
	return nil
}
