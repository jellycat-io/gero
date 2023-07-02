package lexer

import (
	"testing"

	"github.com/jellycat-io/gero/token"
)

func TestNextToken(t *testing.T) {
	input := `
		// This is a comment
		42
		3.14
		/**
		 * Another comment
		 */
		"hello"
		""
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "42"},
		{token.FLOAT, "3.14"},
		{token.STRING, `"hello"`},
		{token.STRING, `""`},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Error(err.Error())
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - Wrong token type. Expected = %q, got = %q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Wrong token literal. Expected = %q, got = %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
