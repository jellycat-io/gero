package parser

import (
	"testing"

	"github.com/jellycat-io/gero/ast"
	"github.com/jellycat-io/gero/lexer"
)

func TestParsingStatementList(t *testing.T) {
	input := `
		5;
		"hello";
	`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()
	checkParserErrors(t, p)

	if len(program.Body) != 2 {
		t.Fatalf("Program has wrong number of statements. Expected=%d, got=%d", 2, len(program.Body))
	}
}

func TestParsingExpressionStatement(t *testing.T) {
	input := `
		5;
		"hello";
	`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()
	checkParserErrors(t, p)

	_, ok := program.Body[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Body[0] is not ast.ExpressionStatement, got=%T", program.Body[0])
	}

	_, ok = program.Body[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Body[1] is not ast.ExpressionStatement, got=%T", program.Body[1])
	}
}

func TestParsingBlockStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{`{ 5; "hello world"; }`, 2},
		{`{}`, 0},
		{`{ 5; { "hello world"; 10; } }`, 2},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Program()
		checkParserErrors(t, p)

		block, ok := program.Body[0].(*ast.BlockStatement)
		if !ok {
			t.Fatalf("program.Body[0] is not ast.BlockStatement, got=%T", program.Body[0])
		}

		if len(block.Body) != tt.expected {
			t.Fatalf("Program has wrong number of statements. Expected=%d, got=%d", tt.expected, len(block.Body))
		}

		if len(block.Body) > 0 {
			stmt := block.Body[0].(*ast.ExpressionStatement)
			testIntegerLiteral(t, stmt.Expression, 5)

			nested, ok := block.Body[1].(*ast.BlockStatement)
			if ok {
				if len(nested.Body) != tt.expected {
					t.Fatalf("Program has wrong number of statements. Expected=%d, got=%d", tt.expected, len(block.Body))
				}
			} else {
				stmt = block.Body[1].(*ast.ExpressionStatement)
				testStringLiteral(t, stmt.Expression, "hello world")
			}

		}
	}
}

func TestParsingIntegerLiteral(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()
	checkParserErrors(t, p)

	stmt := program.Body[0].(*ast.ExpressionStatement)

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestParsingStringLiteral(t *testing.T) {
	input := `
		"hello world";
		'hello world';
	`

	l := lexer.New(input)
	p := New(l)
	program := p.Program()
	checkParserErrors(t, p)

	stmt := program.Body[0].(*ast.ExpressionStatement)
	testStringLiteral(t, stmt.Expression, "hello world")
	stmt = program.Body[1].(*ast.ExpressionStatement)
	testStringLiteral(t, stmt.Expression, "hello world")
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

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	lit, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Literal is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if lit.Value != value {
		t.Errorf("Literal.Value not %d. got=%d", value, lit.Value)
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, il ast.Expression, value string) bool {
	lit, ok := il.(*ast.StringLiteral)
	if !ok {
		t.Errorf("Literal is not *ast.StringLiteral. got=%T", il)
		return false
	}

	if lit.Value != value {
		t.Errorf("Literal.Value not %s. got=%s", value, lit.Value)
		return false
	}

	return true
}
