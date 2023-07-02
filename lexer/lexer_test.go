package lexer

import (
	"testing"

	"github.com/jellycat-io/gero/token"
)

func TestNextToken(t *testing.T) {
	input := `
		42
		"Hello"
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NEWLINE, ""},
		{token.INT, "42"},
		{token.NEWLINE, ""},
		{token.STRING, "Hello"},
		{token.NEWLINE, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - Wrong token type. Expected = %q, got = %q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Wrong token literal. Expected = %q, got = %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
