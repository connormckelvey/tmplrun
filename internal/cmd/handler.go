package cmd

import (
	"io"

	"github.com/urfave/cli/v2"
)

type BeforeHandler[P any] interface {
	Before(ctx *HandlerContext[P], params *P, w io.Writer) error
}

type AfterHandler[P any] interface {
	After(ctx *HandlerContext[P], params *P, w io.Writer) error
}

type OnUsageErrorHandler[P any] interface {
	OnUsageError(ctx *HandlerContext[P], err error, isSubcommand bool, params *P, w io.Writer) error
}

type CommandHandler[P any] interface {
	Action(ctx *HandlerContext[P], params *P, w io.Writer) error
}

type Handler[P any] struct {
	CommandHandler[P]
	paramMappers []func(*cli.Context, *P)
}

func NewHandler[P any](commandHandler CommandHandler[P], opts ...HandlerOption[P]) *Handler[P] {
	h := &Handler[P]{
		CommandHandler: commandHandler,
	}
	for _, apply := range opts {
		apply(h)
	}
	return h
}

type HandlerOption[P any] func(*Handler[P]) CommandOption

func UseArguments[P any](mapping func(*P, []string)) HandlerOption[P] {
	return func(h *Handler[P]) CommandOption {
		h.paramMappers = append(h.paramMappers, func(ctx *cli.Context, p *P) {
			v := ctx.Args().Slice()
			mapping(p, v)
		})
		return func(c *cli.Command) {
			c.Args = true
		}
	}
}

func UseStringFlag[P any](flag *cli.StringFlag, mapping func(*P, string)) HandlerOption[P] {
	return func(h *Handler[P]) CommandOption {
		h.paramMappers = append(h.paramMappers, func(ctx *cli.Context, p *P) {
			v := ctx.String(flag.Name)
			mapping(p, v)
		})
		return func(c *cli.Command) {
			c.Flags = append(c.Flags, flag)
		}
	}
}

func UseStringSliceFlag[P any](flag *cli.StringSliceFlag, mapping func(*P, []string)) HandlerOption[P] {
	return func(h *Handler[P]) CommandOption {
		h.paramMappers = append(h.paramMappers, func(ctx *cli.Context, p *P) {
			v := ctx.StringSlice(flag.Name)
			mapping(p, v)
		})
		return func(c *cli.Command) {
			c.Flags = append(c.Flags, flag)
		}
	}
}

func UseBoolFlag[P any](flag *cli.BoolFlag, mapping func(*P, bool)) HandlerOption[P] {
	return func(h *Handler[P]) CommandOption {
		h.paramMappers = append(h.paramMappers, func(ctx *cli.Context, p *P) {
			v := ctx.Bool(flag.Name)
			mapping(p, v)
		})
		return func(c *cli.Command) {
			c.Flags = append(c.Flags, flag)
		}
	}
}

func (h *Handler[P]) Before(cCtx *cli.Context) error {
	ctx := newHandlerContext[P](cCtx, nil)
	for _, mapper := range h.paramMappers {
		mapper(cCtx, ctx.params)
	}
	if bh, ok := h.CommandHandler.(BeforeHandler[P]); ok {
		return bh.Before(ctx, ctx.params, ctx.cCtx.App.Writer)
	}
	return nil
}

func (h *Handler[P]) After(cCtx *cli.Context) error {
	ah, ok := h.CommandHandler.(AfterHandler[P])
	if !ok {
		return nil
	}
	ctx := loadHandlerContext[P](cCtx, nil)
	return ah.After(ctx, ctx.params, ctx.cCtx.App.Writer)
}

func (h *Handler[P]) Action(cCtx *cli.Context) error {
	ctx := loadHandlerContext[P](cCtx, nil)
	return h.CommandHandler.Action(ctx, ctx.params, cCtx.App.Writer)
}

func (h *Handler[P]) OnUsageError(cCtx *cli.Context, err error, isSubcommand bool) error {
	ueh, ok := h.CommandHandler.(OnUsageErrorHandler[P])
	if !ok {
		return nil
	}
	ctx := loadHandlerContext[P](cCtx, nil)
	return ueh.OnUsageError(ctx, err, isSubcommand, ctx.params, cCtx.App.Writer)
}

type handlerContextKey string

const (
	handlerContextKey_params handlerContextKey = "params"
)
