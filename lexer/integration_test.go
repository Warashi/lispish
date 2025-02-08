package lexer

import (
	"testing"
)

func TestLexerIntegration(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "(define (square x) (* x x))",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "define", Line: 1, Column: 2},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 9},
				{Type: Identifier, Literal: "square", Line: 1, Column: 10},
				{Type: Identifier, Literal: "x", Line: 1, Column: 17},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 18},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 20},
				{Type: Identifier, Literal: "*", Line: 1, Column: 21},
				{Type: Identifier, Literal: "x", Line: 1, Column: 23},
				{Type: Identifier, Literal: "x", Line: 1, Column: 25},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 26},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 27},
			},
		},
		{
			input: "(let ((x 10) (y 20)) (+ x y))",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "let", Line: 1, Column: 2},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 6},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 7},
				{Type: Identifier, Literal: "x", Line: 1, Column: 8},
				{Type: Literal, Literal: "10", Line: 1, Column: 10},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 12},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 14},
				{Type: Identifier, Literal: "y", Line: 1, Column: 15},
				{Type: Literal, Literal: "20", Line: 1, Column: 17},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 19},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 20},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 22},
				{Type: Identifier, Literal: "+", Line: 1, Column: 23},
				{Type: Identifier, Literal: "x", Line: 1, Column: 25},
				{Type: Identifier, Literal: "y", Line: 1, Column: 27},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 28},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 29},
			},
		},
		{
			input: "(if (> x 10) (display \"x is greater than 10\") (display \"x is less than or equal to 10\"))",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "if", Line: 1, Column: 2},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 5},
				{Type: Identifier, Literal: ">", Line: 1, Column: 6},
				{Type: Identifier, Literal: "x", Line: 1, Column: 8},
				{Type: Literal, Literal: "10", Line: 1, Column: 10},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 12},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 14},
				{Type: Identifier, Literal: "display", Line: 1, Column: 15},
				{Type: Literal, Literal: "\"x is greater than 10\"", Line: 1, Column: 23},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 44},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 46},
				{Type: Identifier, Literal: "display", Line: 1, Column: 47},
				{Type: Literal, Literal: "\"x is less than or equal to 10\"", Line: 1, Column: 55},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 84},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 85},
			},
		},
		{
			input: "(begin (define r 10) (* pi (* r r)))",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "begin", Line: 1, Column: 2},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 8},
				{Type: Identifier, Literal: "define", Line: 1, Column: 9},
				{Type: Identifier, Literal: "r", Line: 1, Column: 16},
				{Type: Literal, Literal: "10", Line: 1, Column: 18},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 20},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 22},
				{Type: Identifier, Literal: "*", Line: 1, Column: 23},
				{Type: Identifier, Literal: "pi", Line: 1, Column: 25},
				{Type: Delimiter, Literal: "(", Line: 1, Column: 28},
				{Type: Identifier, Literal: "*", Line: 1, Column: 29},
				{Type: Identifier, Literal: "r", Line: 1, Column: 31},
				{Type: Identifier, Literal: "r", Line: 1, Column: 33},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 34},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 35},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 36},
			},
		},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		for i, expectedToken := range tt.expected {
			tok, err := l.NextToken()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tok != expectedToken {
				t.Errorf("test[%d] - token wrong. expected=%v, got=%v", i, expectedToken, tok)
			}
		}
	}
}
