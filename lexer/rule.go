package lexer

import (
	"regexp"

	"github.com/connormckelvey/tmplrun/token"
)

// Rule defines the interface for lexer rules.
type Rule interface {
	// Test tests if the rule matches the current input.
	Test(*sourceReader) (tokenType token.TokenType, ok bool)
	// Tokenize generates a token from the input.
	Tokenize(*sourceReader) (length int, more bool)
}

// EOFRule represents a rule for end-of-file.
type EOFRule struct{}

// Test tests if the end-of-file rule matches the current input.
func (to *EOFRule) Test(sr *sourceReader) (token.TokenType, bool) {
	return token.EOF, isEOF(sr.Char())
}

// Tokenize generates an end-of-file token.
func (to *EOFRule) Tokenize(sr *sourceReader) (length int, more bool) {
	return 0, false
}

// TextRule represents a rule for text.
type TextRule struct{}

// Test always returns true for text rule.
func (r *TextRule) Test(sr *sourceReader) (token.TokenType, bool) {
	return token.TEXT, true
}

// Tokenize generates a text token.
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
	// CloseTagPattern represents the regular expression pattern for close tags.
	CloseTagPattern = regexp.MustCompile(`^([0-9a-zA-Z]{2,16})?(%{1,16})(>)`)
	// OpenTagPattern represents the regular expression pattern for open tags.
	OpenTagPattern = regexp.MustCompile(`^(<)(%{1,16})([0-9a-zA-Z]{2,16})?`)
	tagPeekChars   = 32
)

// OpenTagRule represents a rule for open tags.
type OpenTagRule struct {
	*regexRule
}

// NewOpenTagRule creates a new instance of OpenTagRule.
func NewOpenTagRule() *OpenTagRule {
	return &OpenTagRule{
		regexRule: newRegexRule(OpenTagPattern, tagPeekChars),
	}
}

// Test tests if the input matches the pattern for open tags.
func (r *OpenTagRule) Test(sr *sourceReader) (token.TokenType, bool) {
	peek, _ := sr.Peek()
	if sr.Char() == '<' && peek == '%' {
		_, ok := r.regexRule.Test(sr)
		return token.OPEN, ok
	}
	return "", false
}

// CloseTagRule represents a rule for close tags.
type CloseTagRule struct {
	*regexRule
}

// NewCloseTagRule creates a new instance of CloseTagRule.
func NewCloseTagRule() *CloseTagRule {
	return &CloseTagRule{
		regexRule: newRegexRule(CloseTagPattern, tagPeekChars),
	}
}

// Test tests if the input matches the pattern for close tags.
func (r *CloseTagRule) Test(sr *sourceReader) (token.TokenType, bool) {
	if isLetter(sr.Char()) || sr.Char() == '%' {
		_, ok := r.regexRule.Test(sr)
		return token.CLOSE, ok
	}
	return "", false
}
