package cmd

import (
	"github.com/urfave/cli/v2"
)

type CLIAppOption func(*cli.App)

func NewCLIApp(options ...CLIAppOption) *cli.App {
	app := cli.NewApp()
	for _, apply := range options {
		apply(app)
	}
	return app
}

func UseEnvironment(env *Environment) CLIAppOption {
	return func(a *cli.App) {
		a.Reader = env.Reader
		a.Writer = env.Writer
		a.ErrWriter = env.ErrWriter
	}
}

func UseCommands(commands ...*cli.Command) CLIAppOption {
	return func(a *cli.App) {
		a.Commands = append(a.Commands, commands...)
	}
}

func UseDefaultCommand(command string) CLIAppOption {
	return func(a *cli.App) {
		a.DefaultCommand = command
	}
}
