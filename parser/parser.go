/**
 * Recursive descent parser
 */

package parser

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jellycat-io/gero/ast"
	"github.com/jellycat-io/gero/lexer"
	"github.com/jellycat-io/gero/token"
)

type Parser struct {
	l         *lexer.Lexer
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	tok, err := p.l.NextToken()
	if err != nil {
		log.Fatal(err)
	}
	p.peekToken = tok
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) Program() *ast.Program {
	program := &ast.Program{}
	program.Body = []ast.Expression{}
	for !p.isAtEnd() {
		program.Body = append(program.Body, p.Literal())
	}
	return program
}

func (p *Parser) Literal() ast.Expression {
	switch p.peekToken.Type {
	case token.INT:
		return p.IntegerLiteral()
	case token.STRING:
		return p.StringLiteral()
	default:
		msg := fmt.Sprintf("Unexpected literal: %q", p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) IntegerLiteral() *ast.IntegerLiteral {
	tok := p.eat(token.INT)
	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", tok.Literal)
		p.errors = append(p.errors, msg)
	}

	return &ast.IntegerLiteral{Token: tok, Value: int64(value)}
}

func (p *Parser) StringLiteral() *ast.StringLiteral {
	tok := p.eat(token.STRING)
	return &ast.StringLiteral{Token: tok, Value: tok.Literal[1 : len(tok.Literal)-1]}
}

func (p *Parser) eat(tokenType token.TokenType) token.Token {
	token := p.peekToken

	if token.Type != tokenType {
		msg := fmt.Sprintf("Unexpected token %q, expected %q", token.Literal, tokenType)
		p.errors = append(p.errors, msg)
	}

	tok, err := p.l.NextToken()
	if err != nil {
		log.Fatal(err)
	}
	p.peekToken = tok

	return token
}

func (p *Parser) isAtEnd() bool {
	return p.peekToken.Type == token.EOF
}
