package main

import (
	"context"
	"os"

	"github.com/connormckelvey/tmplrun/internal/cmd"
	"github.com/urfave/cli/v2"

	"github.com/spf13/afero"
)

func main() {
	env := &cmd.Environment{
		Reader:     os.Stdin,
		Writer:     os.Stdout,
		ErrWriter:  os.Stderr,
		FileSystem: afero.NewOsFs(),
		Args:       os.Args,
		Exit: func(err error) {
			if err != nil {
				os.Exit(1)
			}
		},
		Clock: &cmd.SystemClock{},
	}

	cliApp := cmd.NewCLIApp(
		cmd.UseEnvironment(env),
		func(a *cli.App) {
			a.DefaultCommand = "render"
			a.Usage = "Bring your own language templating engine"
		},
		cmd.UseCommands(
			newRenderCommand(env),
		),
		cmd.UseDefaultCommand("render"),
	)

	app := cmd.NewApplication(env, cliApp)
	app.Run(context.Background())
}
