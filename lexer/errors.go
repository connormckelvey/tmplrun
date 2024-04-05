package lexer

import "fmt"

type Error struct {
	Code errorCode
	Pos  int
}

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
