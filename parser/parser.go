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

	tok, ok := p.l.NextToken().(token.Token)
	if !ok {
		return nil
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
	return ast.NewProgram(p.StatementList(token.EOF))
}

/**
 * StatementList
 * 	: Statement
 * 	| StatementList Statement -> Statement*
 * 	;
 */
func (p *Parser) StatementList(stopTokenType token.TokenType) []ast.Statement {
	statementList := []ast.Statement{p.Statement()}

	for !p.match(stopTokenType) {
		statementList = append(statementList, p.Statement())
	}

	return statementList
}

/**
 * Statement
 * 	: ExpressionStatement
 * 	| BlockStatement
 * 	;
 */
func (p *Parser) Statement() ast.Statement {
	switch p.peekToken.Type {
	case token.LBRACE:
		return p.BlockStatement()
	default:
		return p.ExpressionStatement()
	}
}

/**
 * BlockStatement
 * 	: '{' OptStatementList '}'
 * 	;
 */
func (p *Parser) BlockStatement() *ast.BlockStatement {
	var body []ast.Statement
	p.eat(token.LBRACE)

	if p.peekToken.Type != token.RBRACE {
		body = p.StatementList(token.RBRACE)
	} else {
		body = []ast.Statement{}
	}

	p.eat(token.RBRACE)

	return ast.NewBlockStatement(p.peekToken, body)
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
	return p.AdditiveExpression()
}

/**
 * AdditiveExpression
 * 	: MultiplicativeExpression
 * 	| AdditiveExpression ADDITIVE_OPERATOR MultiplicativeExpression -> MultiplicativeExpression ADDITIVE_OPERATOR MultiplicativeExpression ADDITIVE_OPERATOR MultiplicativeExpression
 * 	;
 */
func (p *Parser) AdditiveExpression() ast.Expression {
	return p.BinaryExpression(p.MultiplicativeExpression, token.PLUS, token.MINUS)
}

/**
 * MultiplicativeExpression
 * 	: PrimaryExpression
 * 	| MultiplicativeExpression MULTIPLICATIVE_OPERATOR PrimaryExpression -> PrimaryExpression MULTIPLICATIVE_OPERATOR PrimaryExpression MULTIPLICATIVE_OPERATOR PrimaryExpression
 * 	;
 */
func (p *Parser) MultiplicativeExpression() ast.Expression {
	return p.BinaryExpression(p.PrimaryExpression, token.ASTERISK, token.SLASH, token.PERCENT)
}

func (p *Parser) BinaryExpression(builder func() ast.Expression, ops ...token.TokenType) ast.Expression {
	left := builder()

	for _, op := range ops {
		for p.match(op) {
			var operator token.Token
			operator = p.eat(op).(token.Token)

			right := builder()

			left = ast.NewBinaryExpression(operator.Literal, left, right)
		}
	}

	return left
}

/**
 * PrimaryExpression
 * 	: Literal
 *	| GroupedExpression
 * 	;
 */
func (p *Parser) PrimaryExpression() ast.Expression {
	switch p.peekToken.Type {
	case token.LPAREN:
		return p.GroupedExpression()
	default:
		return p.Literal()
	}
}

/**
 * GroupedExpression
 * 	: Literal
 * 	;
 */
func (p *Parser) GroupedExpression() ast.Expression {
	p.eat(token.LPAREN)
	exp := p.Expression()
	p.eat(token.RPAREN)

	return exp
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
		//geroError.PrintError(geroError.PARSER_ERROR, msg)
		return nil
	}
}

func (p *Parser) IntegerLiteral() *ast.IntegerLiteral {
	tok, ok := p.eat(token.INT).(token.Token)
	if !ok {
		return nil
	}
	value, err := strconv.ParseInt(tok.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer at line %d", tok.Literal, tok.Line)
		p.errors = append(p.errors, msg)
		//geroError.PrintError(geroError.PARSER_ERROR, msg)
	}

	return ast.NewIntegerLiteral(tok, int64(value))
}

func (p *Parser) StringLiteral() *ast.StringLiteral {
	tok, ok := p.eat(token.STRING).(token.Token)
	if !ok {
		return nil
	}
	return ast.NewStringLiteral(tok, tok.Literal[1:len(tok.Literal)-1])
}

func (p *Parser) eat(tokenType token.TokenType) interface{} {
	curToken := p.peekToken

	if curToken.Type != tokenType {
		msg := fmt.Sprintf("Unexpected token %q, expected %q at line %d", curToken.Type, tokenType, curToken.Line)
		p.errors = append(p.errors, msg)
		//geroError.PrintError(geroError.PARSER_ERROR, msg)
		return nil
	}

	tok, ok := p.l.NextToken().(token.Token)
	if !ok {
		return nil
	}

	p.peekToken = tok

	return curToken
}

func (p *Parser) isAtEnd() bool {
	return p.peekToken.Type == token.EOF
}

func (p *Parser) match(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}
