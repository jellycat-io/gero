package ast

import (
	"bytes"

	"github.com/jellycat-io/gero/token"
)

type NodeType string

type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Type       string
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func NewProgram(stmts []Statement) *Program {
	return &Program{
		Type:       "Program",
		Statements: stmts,
	}
}

type ExpressionStatement struct {
	Type       string
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
func NewExpressionStatement(t token.Token, e Expression) *ExpressionStatement {
	return &ExpressionStatement{
		Type:       "ExpressionStatement",
		Token:      t,
		Expression: e,
	}
}

type BlockStatement struct {
	Type  string
	Token token.Token // the first token of the expression
	Body  []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Body {
		out.WriteString(s.String())
	}

	return out.String()
}
func NewBlockStatement(t token.Token, stmts []Statement) *BlockStatement {
	return &BlockStatement{
		Type:  "BlockStatement",
		Token: t,
		Body:  stmts,
	}
}

type BinaryExpression struct {
	Type     string
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}
func NewBinaryExpression(o string, l Expression, r Expression) *BinaryExpression {
	return &BinaryExpression{
		Type:     "BinaryExpression",
		Left:     l,
		Operator: o,
		Right:    r,
	}
}

type IntegerLiteral struct {
	Type  string
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string  { return il.Token.Literal }
func NewIntegerLiteral(t token.Token, value int64) *IntegerLiteral {
	return &IntegerLiteral{
		Type:  "IntegerLiteral",
		Token: t,
		Value: value,
	}
}

type FloatLiteral struct {
	Type  string
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string  { return fl.Token.Literal }
func NewFloatLiteral(t token.Token, value float64) *FloatLiteral {
	return &FloatLiteral{
		Type:  "FloatLiteral",
		Token: t,
		Value: value,
	}
}

type StringLiteral struct {
	Type  string
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string  { return sl.Token.Literal }
func NewStringLiteral(t token.Token, value string) *StringLiteral {
	return &StringLiteral{
		Type:  "StringLiteral",
		Token: t,
		Value: value,
	}
}
