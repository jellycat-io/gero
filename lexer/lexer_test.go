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
		{ "hello"; }
		2 + 2;
		2 - 2;
		2 * 2;
		2 / 2;
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "42"},
		{token.FLOAT, "3.14"},
		{token.STRING, `"hello"`},
		{token.STRING, `""`},
		{token.LBRACE, `{`},
		{token.STRING, `"hello"`},
		{token.SEMI, `;`},
		{token.RBRACE, `}`},
		{token.INT, `2`},
		{token.PLUS, `+`},
		{token.INT, `2`},
		{token.SEMI, `;`},
		{token.INT, `2`},
		{token.MINUS, `-`},
		{token.INT, `2`},
		{token.SEMI, `;`},
		{token.INT, `2`},
		{token.ASTERISK, `*`},
		{token.INT, `2`},
		{token.SEMI, `;`},
		{token.INT, `2`},
		{token.SLASH, `/`},
		{token.INT, `2`},
		{token.SEMI, `;`},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok, ok := l.NextToken().(token.Token)
		if !ok {
			return
		}

		if tok.Type != tt.expectedType {
			t.Fatalf("Tests[%d] - Wrong token type. Expected = %q, got = %q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("Tests[%d] - Wrong token literal. Expected = %q, got = %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
