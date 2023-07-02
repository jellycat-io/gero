package ast

import "github.com/jellycat-io/gero/token"

type NodeType string

type Node interface{}

type Program struct {
	Type string `default:"Program"`
	Body []Statement
}

func NewProgram(b []Statement) *Program {
	return &Program{
		Type: "Program",
		Body: b,
	}
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type ExpressionStatement struct {
	Type       string
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func NewExpressionStatement(t token.Token, e Expression) *ExpressionStatement {
	return &ExpressionStatement{
		Type:       "ExpressionStatement",
		Token:      t,
		Expression: e,
	}
}

type IntegerLiteral struct {
	Type  string
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
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
func NewStringLiteral(t token.Token, value string) *StringLiteral {
	return &StringLiteral{
		Type:  "StringLiteral",
		Token: t,
		Value: value,
	}
}
