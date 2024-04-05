package cmd

import (
	"context"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

type HandlerContext[P any] struct {
	cCtx   *cli.Context
	params *P
	Reader io.Reader
	Log    *zerolog.Logger
}

func newHandlerContext[P any](cCtx *cli.Context, log *zerolog.Logger) *HandlerContext[P] {
	ctx := &HandlerContext[P]{
		cCtx:   cCtx,
		params: new(P),
		Reader: cCtx.App.Reader,
		Log:    log,
	}
	cCtx.Context = context.WithValue(
		cCtx.Context,
		handlerContextKey_params,
		ctx.params,
	)
	return ctx
}

func loadHandlerContext[P any](cCtx *cli.Context, log *zerolog.Logger) *HandlerContext[P] {
	ctx := &HandlerContext[P]{
		cCtx:   cCtx,
		params: cCtx.Context.Value(handlerContextKey_params).(*P),
		Reader: cCtx.App.Reader,
		Log:    log,
	}
	return ctx
}

func (ctx *HandlerContext[P]) Deadline() (deadline time.Time, ok bool) {
	return ctx.cCtx.Context.Deadline()
}
func (ctx *HandlerContext[P]) Done() <-chan struct{} {
	return ctx.cCtx.Context.Done()
}
func (ctx *HandlerContext[P]) Err() error {
	return ctx.cCtx.Context.Err()

}
func (ctx *HandlerContext[P]) Value(key any) any {
	return ctx.cCtx.Context.Value(key)
}
