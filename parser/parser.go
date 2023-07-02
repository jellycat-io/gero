/**
 * Recursive descent parser
 */

package parser

import (
	"fmt"
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
		fmt.Printf(err.Error())
	}
	p.peekToken = tok
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

/**
 * Main entry point.
 *
 * Program
 * 	: StatementList
 * 	;
 */
func (p *Parser) Program() *ast.Program {
	return ast.NewProgram(p.StatementList())
}

/**
 * StatementList
 * 	: Statement
 * 	| StatementList Statement -> Statement*
 * 	;
 */
func (p *Parser) StatementList() []ast.Statement {
	statementList := []ast.Statement{p.Statement()}

	for !p.isAtEnd() {
		statementList = append(statementList, p.Statement())
	}

	return statementList
}

/**
 * Statement
 * 	: ExpressionStatement
 * 	;
 */
func (p *Parser) Statement() ast.Statement {
	return p.ExpressionStatement()
}

/**
 * ExpressionStatement
 * 	: Expression ';'
 * 	;
 */
func (p *Parser) ExpressionStatement() *ast.ExpressionStatement {
	exp := p.Expression()

	p.eat(token.SEMI)

	return ast.NewExpressionStatement(p.peekToken, exp)
}

/**
 * Expression
 * 	: Literal
 * 	;
 */
func (p *Parser) Expression() ast.Expression {
	return p.Literal()
}

/**
 * Literal
 * 	: IntegerLiteral
 * 	| FloatLiteral
 * 	| StringLiteral
 * 	;
 */
func (p *Parser) Literal() ast.Expression {
	switch p.peekToken.Type {
	case token.INT:
		return p.IntegerLiteral()
	case token.STRING:
		return p.StringLiteral()
	default:
		msg := fmt.Sprintf("Unexpected literal %q at line %d", p.peekToken.Type, p.peekToken.Line)
		p.errors = append(p.errors, msg)
		return nil
	}
}

func (p *Parser) IntegerLiteral() *ast.IntegerLiteral {
	tok, ok := p.eat(token.INT)
	if !ok {
		return nil
	}
	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer at line %d", tok.Literal, tok.Line)
		p.errors = append(p.errors, msg)
	}

	return ast.NewIntegerLiteral(tok, int64(value))
}

func (p *Parser) StringLiteral() *ast.StringLiteral {
	tok, ok := p.eat(token.STRING)
	if !ok {
		return nil
	}
	return ast.NewStringLiteral(tok, tok.Literal[1:len(tok.Literal)-1])
}

func (p *Parser) eat(tokenType token.TokenType) (t token.Token, ok bool) {
	token := p.peekToken

	if token.Type != tokenType {
		msg := fmt.Sprintf("Unexpected token %q, expected %q at line %d", token.Literal, tokenType, token.Line)
		p.errors = append(p.errors, msg)
		return token, false
	}

	tok, err := p.l.NextToken()
	if err != nil {
		fmt.Printf(err.Error())
	}
	p.peekToken = tok

	return token, true
}

func (p *Parser) isAtEnd() bool {
	return p.peekToken.Type == token.EOF
}
