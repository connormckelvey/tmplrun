package token

// TokenType represents the type of a token.
type TokenType string

// Token types
const (
	// UNKNOWN represents an unknown token type.
	UNKNOWN TokenType = "UNKNOWN"
	// TEXT represents a text token type.
	TEXT TokenType = "TEXT"
	// OPEN represents an open tag token type.
	OPEN TokenType = "<%"
	// CLOSE represents a close tag token type.
	CLOSE TokenType = "%>"
	// EOF represents an end-of-file token type.
	EOF TokenType = "EOF"
)

// Token represents a token in the template.
type Token struct {
	// Type is the type of the token.
	Type TokenType
	// Literal is the literal value of the token.
	Literal string
	// Position is the position of the token in the input.
	Position int
}
