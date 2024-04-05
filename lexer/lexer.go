package lexer

import (
	"io"
	"strings"
	"sync"

	"github.com/connormckelvey/tmplrun/token"
)

type Lexer struct {
	source *sourceReader
	rules  []Rule
	once   *sync.Once
}

func New(source io.Reader) *Lexer {
	return NewWithRules(
		source,
		new(EOFRule),
		NewOpenTagRule(),
		NewCloseTagRule(),
		new(TextRule),
	)
}

func NewWithRules(source io.Reader, rules ...Rule) *Lexer {
	l := &Lexer{
		source: newSourceReader(source),
		rules:  rules,
		once:   new(sync.Once),
	}
	return l
}

func (l *Lexer) init() (err error) {
	l.once.Do(func() {
		err = l.source.nextChar()
	})
	return err
}

func (l *Lexer) NextToken() (*token.Token, error) {
	if err := l.init(); err != nil {
		return nil, err
	}

	rule, idx, tokenType := l.findRule(l.rules)
	if rule == nil {
		err := errRuleNotFound.create(l.source.Pos())
		return nil, err
	}

	var buf strings.Builder
	for {
		length, more := rule.Tokenize(l.source)
		for i := 0; i < length; i++ {
			buf.WriteByte(l.source.Char())
			err := l.source.nextChar()
			if err != nil {
				return nil, err
			}
		}
		if !more {
			break
		}

		// check for higher priority rule before continuing
		rule, _, _ := l.findRule(l.rules[0:idx])
		if rule != nil {
			break
		}
	}

	return &token.Token{
		Type:    tokenType,
		Literal: buf.String(),
	}, nil
}

func (l *Lexer) findRule(rules []Rule) (rule Rule, idx int, tokenType token.TokenType) {
	for i, rule := range rules {
		if tokenType, ok := rule.Test(l.source); ok {
			return rule, i, tokenType
		}
	}
	return nil, -1, ""
}
