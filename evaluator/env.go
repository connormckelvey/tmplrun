package evaluator

import (
	"fmt"
	"io/fs"
)

type Environment struct {
	props  map[string]any
	fs     fs.FS
	errors []error
	hooks  HooksAPI
}

func NewEnvironment(fsys fs.FS, props map[string]any, hooksAPI HooksAPI) *Environment {
	return &Environment{
		props: props,
		fs:    fsys,
		hooks: hooksAPI,
	}
}

func (env *Environment) Props() map[string]any {
	return env.props
}

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

func (env *Environment) Template(name string, props map[string]any) string {
	file := env.Include(name)
	if file == "" || env.hasErrors() {
		return ""
	}

	res, err := env.hooks.Render(file, props)
	if err != nil {
		env.errors = append(env.errors, err)
		return ""
	}

	return res
}
