package lexer

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/TwiN/go-color"
	"github.com/jellycat-io/gero/token"
)

var specs = map[string]token.TokenType{
	//-----------------------------------
	// Skipped
	"^\\s+":                   token.WHITESPACE,
	"^\\t+":                   token.WHITESPACE,
	"^\\n":                    token.NEWLINE,
	"^\\/\\/.*":               token.COMMENT,
	"^\\/\\*[\\s\\S]*?\\*\\/": token.COMMENT,
	//-----------------------------------
	// Symbols, delimiters
	"^;": token.SEMI,
	"^{": token.LBRACE,
	"^}": token.RBRACE,
	//-----------------------------------
	// Numbers
	"^[0-9]*(\\.[0-9]+)": token.FLOAT,
	"^\\d+":              token.INT,
	//-----------------------------------
	// Strings
	`^"[^"]*"`: token.STRING,
	`^'[^']*'`: token.STRING,
}

// Lazily pulls a token from a stream.
type Lexer struct {
	input  string
	line   int
	cursor int
}

func New(input string) *Lexer {
	l := &Lexer{input: input, cursor: 0, line: 1}
	return l
}

func (l *Lexer) NextToken() (token.Token, error) {
	if !l.hasMoreTokens() {
		return l.newToken(token.EOF, ""), nil
	}

	s := l.input[l.cursor:len(l.input)]

	for regex, tokenType := range specs {
		value, ok := l.match(regex, s)

		if !ok {
			continue
		}

		if tokenType == token.NEWLINE {
			// TODO: Fix line increment
			l.line += 1
			return l.NextToken()
		}

		if tokenType == token.WHITESPACE || tokenType == token.COMMENT {
			return l.NextToken()
		}

		return l.newToken(tokenType, value), nil
	}

	return l.newToken(token.ILLEGAL, ""), errors.New(color.InRed(fmt.Sprintf(
		`Syntax error: Unexpected token "%s" at line %d`,
		string(s[0]),
		l.line,
	)))
}

func (l *Lexer) hasMoreTokens() bool {
	return l.cursor < len(l.input)
}

func (l *Lexer) match(regex string, input string) (s string, ok bool) {
	re := regexp.MustCompile(regex)
	matched := re.FindString(input)

	if matched == "" {
		return "", false
	}

	l.cursor += len(matched)
	return matched, true
}

func (l *Lexer) newToken(tokenType token.TokenType, value string) token.Token {
	return token.Token{Type: tokenType, Literal: string(value), Line: l.line}
}
