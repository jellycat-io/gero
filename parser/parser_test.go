package parser

import (
	"testing"

	"github.com/jellycat-io/gero/ast"
	"github.com/jellycat-io/gero/lexer"
)

func TestParsingIntegerLiteral(t *testing.T) {
	input := `5`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()

	if len(program.Body) != 1 {
		t.Fatalf("Program has wrong number of expressions. Expected=%d, got=%d", 1, len(program.Body))
	}
	literal, ok := program.Body[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("program.Body[0] is not ast.IntegerLiteral, got=%T", program.Body[0])
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not %d, got=%d", 5, literal.Value)
	}
}

func TestParsingStringLiteral(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()

	if len(program.Body) != 1 {
		t.Fatalf("Program has wrong number of expressions. Expected=%d, got=%d", 1, len(program.Body))
	}
	literal, ok := program.Body[0].(*ast.StringLiteral)
	if !ok {
		t.Fatalf("program.Body[0] is not ast.StringLiteral, got=%T", program.Body[0])
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %s, got=%s", "hello world", literal.Value)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
