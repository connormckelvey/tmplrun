package lexer

import (
	"regexp"

	"github.com/connormckelvey/tmplrun/token"
)

type Rule interface {
	Test(*sourceReader) (tokenType token.TokenType, ok bool)
	Tokenize(sr *sourceReader) (length int, more bool)
}

type EOFRule struct{}

func (to *EOFRule) Test(sr *sourceReader) (token.TokenType, bool) {
	return token.EOF, isEOF(sr.Char())
}

func (to *EOFRule) Tokenize(sr *sourceReader) (length int, more bool) {
	return 0, false
}

type TextRule struct{}

func (r *TextRule) Test(sr *sourceReader) (token.TokenType, bool) {
	return token.TEXT, true
}

func (r *TextRule) Tokenize(sr *sourceReader) (n int, more bool) {
	return 1, true
}

type regexRule struct {
	pattern  *regexp.Regexp
	peekSize int
	matches  map[int]int
}

func newRegexRule(pattern *regexp.Regexp, peekSize int) *regexRule {
	return &regexRule{
		pattern:  pattern,
		peekSize: peekSize,
		matches:  make(map[int]int),
	}
}

func (r *regexRule) Tokenize(sr *sourceReader) (n int, more bool) {
	n = r.matches[sr.Pos()]
	delete(r.matches, sr.Pos())
	return n, false
}

func (r *regexRule) Test(sr *sourceReader) (token.TokenType, bool) {
	peek, _ := sr.PeekString(r.peekSize)
	test := string(sr.Char()) + peek

	match := r.pattern.FindStringSubmatch(test)
	if match == nil {
		return "", false
	}

	r.matches[sr.Pos()] = len(match[0])
	return token.UNKNOWN, true
}

var (
	CloseTagPattern = regexp.MustCompile(`^([0-9a-zA-Z]{2,16})?(%{1,16})(>)`)
	OpenTagPattern  = regexp.MustCompile(`^(<)(%{1,16})([0-9a-zA-Z]{2,16})?`)
	tagPeekChars    = 32
)

type OpenTagRule struct {
	*regexRule
}

func NewOpenTagRule() *OpenTagRule {
	return &OpenTagRule{
		regexRule: newRegexRule(OpenTagPattern, tagPeekChars),
	}
}

func (r *OpenTagRule) Test(sr *sourceReader) (token.TokenType, bool) {
	peek, _ := sr.Peek()
	if sr.Char() == '<' && peek == '%' {
		_, ok := r.regexRule.Test(sr)
		return token.OPEN, ok
	}
	return "", false
}

type CloseTagRule struct {
	*regexRule
}

func NewCloseTagRule() *CloseTagRule {
	return &CloseTagRule{
		regexRule: newRegexRule(CloseTagPattern, tagPeekChars),
	}
}

func (r *CloseTagRule) Test(sr *sourceReader) (token.TokenType, bool) {
	if isLetter(sr.Char()) || sr.Char() == '%' {
		_, ok := r.regexRule.Test(sr)
		return token.CLOSE, ok
	}
	return "", false
}
