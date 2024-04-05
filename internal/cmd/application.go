package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type Application struct {
	env *Environment
	cli *cli.App
	log zerolog.Logger
}

func NewApplication(env *Environment, cli *cli.App) *Application {
	return &Application{
		env: env,
		cli: cli,
		log: zerolog.New(env.ErrWriter),
	}
}

func (app *Application) Run(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		sig := <-sigs
		cancel(fmt.Errorf("signal received %s", sig.String()))
	}()

	err := app.cli.RunContext(ctx, app.env.Args)
	switch {
	case err != nil:
		app.log.Err(err).Msg("")
		app.env.Exit(err)
	case ctx.Err() != nil:
		app.log.Err(err).Msg("")
		app.env.Exit(ctx.Err())
	default:
		app.env.Exit(nil)
	}
}

type Environment struct {
	Reader     io.Reader
	Writer     io.Writer
	ErrWriter  io.Writer
	FileSystem afero.Fs
	Args       []string
	Exit       func(error)
	Clock      Clock
}

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (st *SystemClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	nanos int64
}

func (st *FakeClock) Now() time.Time {
	return time.Unix(0, st.nanos)
}
