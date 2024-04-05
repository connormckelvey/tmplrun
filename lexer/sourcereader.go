package lexer

import (
	"bufio"
	"errors"
	"io"
)

type sourceReader struct {
	r   bufio.Reader
	pos int
	ch  byte
}

const maxPeekSize = 4096

func newSourceReader(r io.Reader) *sourceReader {
	return newSourceReaderSize(r, maxPeekSize)
}

func newSourceReaderSize(r io.Reader, size int) *sourceReader {
	return &sourceReader{
		r:   *bufio.NewReaderSize(r, size),
		pos: -1,
	}
}

func (sr *sourceReader) Char() byte {
	return sr.ch
}

func (sr *sourceReader) Pos() int {
	return sr.pos
}

func (sr *sourceReader) Peek() (byte, error) {
	b, err := sr.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (sr *sourceReader) PeekString(length int) (string, error) {
	b, err := sr.r.Peek(length)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}

func (sr *sourceReader) nextChar() error {
	var err error
	sr.ch, err = sr.r.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	sr.pos++
	return nil
}
