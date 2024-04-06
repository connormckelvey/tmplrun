package parser

import (
	"errors"

	"github.com/connormckelvey/tmplrun/ast"
	"github.com/connormckelvey/tmplrun/internal/containers"
	"github.com/connormckelvey/tmplrun/lexer"
	"github.com/connormckelvey/tmplrun/token"
)

// Parser represents a template parser.
type Parser struct {
	l         *lexer.Lexer
	errors    []string
	curToken  *token.Token
	peekToken *token.Token
}

// New creates a new instance of Parser with the given lexer.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	return p
}

func (p *Parser) init() error {
	if err := p.nextToken(); err != nil {
		return err
	}
	if err := p.nextToken(); err != nil {
		return err
	}
	return nil
}

// Parse parses the input template and returns the abstract syntax tree (AST) document.
func (p *Parser) Parse() (*ast.Document, error) {
	if err := p.init(); err != nil {
		return nil, err
	}

	var document ast.Document
	var stack containers.Stack[ast.Node]
	stack.Push(&document)

	for !p.curTokenIs(token.EOF) {
		currentNode := stack.Peek()
		switch p.curToken.Type {
		case token.TEXT:
			currentNode.Append(&ast.TextNode{
				Token: p.curToken,
			})
		case token.OPEN:
			templateNode := ast.NewTemplateNode(p.curToken)
			currentNode.Append(templateNode)
			stack.Push(templateNode)
		case token.CLOSE:
			if templateNode, ok := currentNode.(*ast.TemplateNode); !ok {
				return nil, errors.New("unexpected node")
			} else if isClosingTag(templateNode.Token, p.curToken) {
				templateNode.Closed = true
				stack.Pop()
			}
		}
		err := p.nextToken()
		if err != nil {
			return nil, err
		}
	}
	return &document, nil
}

func (p *Parser) nextToken() (err error) {
	next, err := p.l.NextToken()
	if err != nil {
		return err
	}
	p.curToken = p.peekToken
	p.peekToken = next

	return nil
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
