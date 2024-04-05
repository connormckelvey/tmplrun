package parser

import (
	"github.com/connormckelvey/tmplrun/lexer"
	"github.com/connormckelvey/tmplrun/token"
)

func isClosingTag(open *token.Token, close *token.Token) bool {
	openMatch := lexer.OpenTagPattern.FindStringSubmatch(open.Literal)
	openPadding, openIdent := openMatch[2], openMatch[3]

	closeMatch := lexer.CloseTagPattern.FindStringSubmatch(close.Literal)
	closePadding, closeIdent := closeMatch[2], closeMatch[1]

	return openPadding == closePadding && openIdent == closeIdent
}
