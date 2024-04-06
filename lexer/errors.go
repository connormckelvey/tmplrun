package lexer

import "fmt"

// Error represents an error that occurred during lexing.
type Error struct {
	// Code is the error code.
	Code errorCode
	// Pos is the position where the error occurred.
	Pos int
}

// Error returns the error message for the lexer error.
func (le Error) Error() string {
	return fmt.Sprintf(
		"%v: pos: %v",
		errorMessages[le.Code],
		le.Pos,
	)
}

type errorCode int

func (code errorCode) create(pos int) Error {
	return Error{
		Code: code,
		Pos:  pos,
	}
}

const (
	errRuleNotFound errorCode = iota + 1
)

var errorMessages = []string{
	"unknown error",
	"no matching rule",
}
