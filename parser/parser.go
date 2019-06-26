package parser

import (
	"github.com/dudewhocode/interpreter/ast"
	"github.com/dudewhocode/interpreter/lexer"
	"github.com/dudewhocode/interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  *token.Token
	peekToken *token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokesn, so curToken and peekToken are set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	return program
}
