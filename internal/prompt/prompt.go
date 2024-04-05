package prompt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type Prompt interface {
	Run(a ...any) error
}

type Confirm struct {
	Reader io.Reader
	Writer io.Writer
	Labelf string
}

type CharReader struct {
	f        *os.File
	rr       *bufio.Reader
	oldState *term.State
}

func NewCharReader(r io.Reader) *CharReader {
	tr := &CharReader{
		rr: bufio.NewReader(r),
	}

	f, ok := r.(*os.File)
	if !ok {
		return tr
	}

	// check if file is pipe:
	info, err := f.Stat()
	if err != nil {
		return tr
	}
	if info.Mode()&os.ModeCharDevice == 0 {
		return tr
	}

	// enable reading one byte at a time from stdin
	oldState, err := term.MakeRaw(int(f.Fd()))
	if err != nil {
		return tr
	}

	tr.rr = bufio.NewReader(f)
	tr.f = f
	tr.oldState = oldState

	return tr
}

func (cr *CharReader) ReadChar() (string, error) {
	r, _, err := cr.rr.ReadRune()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// does not close the underlying reader
func (cr *CharReader) Close() error {
	if cr.oldState == nil {
		return nil
	}
	err := term.Restore(int(cr.f.Fd()), cr.oldState)
	if err != nil {
		return err
	}
	return nil
}

func (c *Confirm) Run(ctx context.Context, a ...any) (bool, error) {
	if c.Reader == nil {
		c.Reader = os.Stdin
	}
	if c.Writer == nil {
		c.Writer = os.Stdout
	}

	input := make(chan string, 1)
	result := make(chan bool, 1)
	errs := make(chan error, 1)

	fmt.Fprintf(c.Writer, c.Labelf, a...)

	go func() {
		defer close(errs)
		defer close(input)

		reader := NewCharReader(c.Reader)
		char, readErr := reader.ReadChar()
		if closeErr := reader.Close(); closeErr != nil {
			errs <- closeErr
			return
		}
		if readErr != nil {
			errs <- readErr
			return
		}
		if _, err := fmt.Fprintln(c.Writer, char); err != nil {
			errs <- err
			return
		}

		input <- char
	}()

	go func() {
		defer close(result)

		select {
		case <-ctx.Done():
			return
		case text := <-input:
			switch strings.ToLower(strings.TrimSpace(text)) {
			case "y":
				result <- true
			}
		}
	}()

	return <-result, <-errs
}
