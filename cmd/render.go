package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/connormckelvey/tmplrun"
	"github.com/connormckelvey/tmplrun/internal/cmd"
	"github.com/connormckelvey/tmplrun/internal/fsys"
	"github.com/connormckelvey/tmplrun/internal/prompt"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type RenderParams struct {
	Entrypoint string
	PropsFile  string
	Include    []string
	Output     string
	Overwrite  bool
}

type RenderHandler struct {
	env *cmd.Environment
}

func newRenderHandler(env *cmd.Environment) *RenderHandler {
	return &RenderHandler{
		env: env,
	}
}

func (h *RenderHandler) loadProps(filename string) (map[string]any, error) {
	fsys := afero.NewIOFS(h.env.FileSystem)
	var props map[string]any
	if filename == "" {
		return nil, nil
	}
	binProps, err := fs.ReadFile(fsys, filepath.Clean(filename))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(binProps, &props); err != nil {
		return nil, err
	}
	return props, nil
}

const overwritePromptFmt = "Output file '%s' already exists. Overwrite? (y/N): "

func newOverwritePrompt(r io.Reader, w io.Writer) *prompt.Confirm {
	return &prompt.Confirm{
		Reader: r,
		Writer: w,
		Labelf: overwritePromptFmt,
	}
}

func (h *RenderHandler) prepareOutput(ctx *cmd.HandlerContext[RenderParams], entrypointFile string, outputFile string, overwrite bool) (string, fs.FileMode, error) {
	outputInfo, err := h.env.FileSystem.Stat(outputFile)
	if err != nil && !os.IsNotExist(err) {
		return "", 0, err
	}

	if outputInfo != nil {
		if outputInfo.IsDir() {
			entrypointBase := filepath.Base(entrypointFile)
			return h.prepareOutput(ctx, entrypointFile, filepath.Join(outputFile, entrypointBase), overwrite)
		}

		if !overwrite {
			confirm := newOverwritePrompt(ctx.Reader, h.env.Writer)
			ok, err := confirm.Run(context.Background(), outputFile)
			if err != nil {
				return "", 0, err
			}
			if !ok {
				return "", 0, errors.New("canceled")
			}
		}
	}

	entryInfo, err := h.env.FileSystem.Stat(entrypointFile)
	if err != nil && !os.IsNotExist(err) {
		return "", 0, err
	}
	return outputFile, entryInfo.Mode().Perm(), nil

}
func (h *RenderHandler) Action(ctx *cmd.HandlerContext[RenderParams], params *RenderParams, output io.Writer) error {
	props, err := h.loadProps(params.PropsFile)
	if err != nil {
		return err
	}

	readOnlyFS := afero.NewReadOnlyFs(h.env.FileSystem)
	limitedFS, err := fsys.NewLimitedFS(
		afero.NewIOFS(readOnlyFS),
		fsys.WithGlobs(params.Include),
	)
	if err != nil {
		return err
	}

	f, err := limitedFS.Open(filepath.Clean(params.Entrypoint))
	if err != nil {
		return fmt.Errorf("unable to open entrypoint: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			ctx.Log.Err(err).Msgf("error closing file: %s", params.Entrypoint)
		}
	}()

	tmpl := tmplrun.New(limitedFS)
	result, err := tmpl.Run(f, props)
	if err != nil {
		return err
	}

	if params.Output == "" {
		_, err := fmt.Fprintln(output, result)
		return err
	}

	outputPath := filepath.Dir(params.Output)
	outputName := filepath.Base(params.Output)
	tmpName := fmt.Sprintf("%s.tmp.%d", outputName, h.env.Clock.Now().UnixNano())
	tmpOutput := filepath.Join(outputPath, tmpName)

	tmpPath, err := filepath.Abs(tmpOutput)
	if err != nil {
		return err
	}

	outputFile, _, err := h.prepareOutput(ctx, params.Entrypoint, params.Output, params.Overwrite)
	if err != nil {
		return err
	}

	file, err := h.env.FileSystem.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	if _, err := file.WriteString(result); err != nil {
		return err
	}

	if err := h.env.FileSystem.Rename(tmpOutput, outputFile); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(output, "Rendered to %s\n", outputFile); err != nil {
		return err
	}
	return nil
}

func newRenderCommand(env *cmd.Environment) *cli.Command {
	return cmd.NewCommand(
		&cli.Command{
			Name:    "render",
			Aliases: []string{"r"},
			Usage:   "todo",
		},
		cmd.UseHandler(
			newRenderHandler(env),
			cmd.UseArguments(func(p *RenderParams, s []string) {
				if len(s) > 0 {
					p.Entrypoint = s[0]
				}
			}),
			cmd.UseStringFlag(
				&cli.StringFlag{Name: "props", Aliases: []string{"p"}, Usage: "TODO"},
				func(p *RenderParams, s string) {
					p.PropsFile = s
				}),
			cmd.UseStringSliceFlag(
				&cli.StringSliceFlag{Name: "include", Aliases: []string{"i"}, Usage: "TODO"},
				func(p *RenderParams, s []string) {
					p.Include = s
				}),
			cmd.UseStringFlag(
				&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "TODO"},
				func(p *RenderParams, s string) {
					p.Output = s
				}),
			cmd.UseBoolFlag(
				&cli.BoolFlag{Name: "yes", Aliases: []string{"y"}, Usage: "TODO"},
				func(p *RenderParams, b bool) {
					p.Overwrite = b
				}),
		),
	)
}
