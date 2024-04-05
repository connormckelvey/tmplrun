package token

type TokenType string

const (
	UNKNOWN TokenType = "UNKNOWN"
	TEXT    TokenType = "TEXT"
	OPEN    TokenType = "<%"
	CLOSE   TokenType = "%>"
	EOF     TokenType = "EOF"
)

type Token struct {
	Type     TokenType
	Literal  string
	Position int
}
