package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	EOF TokenType = iota
	Identifier
	Keyword
	Literal
	Operator
	Delimiter
	Comment
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer represents a lexical scanner.
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// NewLexer initializes a new instance of Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar reads the next character in the input and advances the positions.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() (Token, error) {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
	case '(':
		tok = Token{Type: Delimiter, Literal: "(", Line: l.line, Column: l.column}
	case ')':
		tok = Token{Type: Delimiter, Literal: ")", Line: l.line, Column: l.column}
	case ';':
		tok = l.readComment()
	case '"':
		return l.readString()
	default:
		if isLetter(l.ch) {
			return l.readIdentifier()
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			return Token{}, fmt.Errorf("invalid character '%c' at line %d, column %d", l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok, nil
}

// skipWhitespace skips over whitespace characters.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readComment reads a comment token.
func (l *Lexer) readComment() Token {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return Token{Type: Comment, Literal: l.input[position:l.position], Line: l.line, Column: l.column}
}

// readString reads a string literal token.
func (l *Lexer) readString() (Token, error) {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
		}
	}
	if l.ch == 0 {
		return Token{}, fmt.Errorf("unterminated string literal at line %d, column %d", l.line, l.column)
	}
	literal := l.input[position:l.position]
	return Token{Type: Literal, Literal: literal, Line: l.line, Column: l.column}, nil
}

// readIdentifier reads an identifier token.
func (l *Lexer) readIdentifier() (Token, error) {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	literal := l.input[position:l.position]
	return Token{Type: Identifier, Literal: literal, Line: l.line, Column: l.column}, nil
}

// readNumber reads a numeric literal token.
func (l *Lexer) readNumber() (Token, error) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	literal := l.input[position:l.position]
	return Token{Type: Literal, Literal: literal, Line: l.line, Column: l.column}, nil
}

// isLetter checks if a character is a letter.
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
