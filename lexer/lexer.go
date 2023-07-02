package lexer

import (
	"strings"

	"github.com/jellycat-io/gero/token"
)

// Lazily pulls a token from a stream.
type Lexer struct {
	input        string
	position     int
	readPosition int
	line         int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return '0'
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '\n':
		l.line++
		tok.Type = token.NEWLINE
		tok.Literal = ""
		l.readChar()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			if isFloat(tok.Literal) {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.input[l.readPosition]) {
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == '0' {
			break
		}
	}
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isFloat(n string) bool {
	return strings.Contains(n, ".")
}
