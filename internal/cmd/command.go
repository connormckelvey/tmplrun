package cmd

import (
	"github.com/urfave/cli/v2"
)

func NewCommand(command *cli.Command, options ...CommandOption) *cli.Command {
	for _, apply := range options {
		apply(command)
	}
	return command
}

type CommandOption func(*cli.Command)

func UseHandler[P any](handler CommandHandler[P], opts ...HandlerOption[P]) CommandOption {
	return func(c *cli.Command) {
		h := NewHandler(handler, opts...)
		for _, apply := range opts {
			apply(h)(c)
		}
		c.Before = h.Before
		c.Action = h.Action
		c.After = h.After
		c.OnUsageError = h.OnUsageError
	}
}
