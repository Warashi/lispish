package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			input: "(define x 10)",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "define", Line: 1, Column: 2},
				{Type: Identifier, Literal: "x", Line: 1, Column: 8},
				{Type: Literal, Literal: "10", Line: 1, Column: 10},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 12},
			},
		},
		{
			input: "(+ 1 2)",
			expected: []Token{
				{Type: Delimiter, Literal: "(", Line: 1, Column: 1},
				{Type: Identifier, Literal: "+", Line: 1, Column: 2},
				{Type: Literal, Literal: "1", Line: 1, Column: 4},
				{Type: Literal, Literal: "2", Line: 1, Column: 6},
				{Type: Delimiter, Literal: ")", Line: 1, Column: 7},
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
			input: "; this is a comment\n(define y 20)",
			expected: []Token{
				{Type: Comment, Literal: "; this is a comment", Line: 1, Column: 1},
				{Type: Delimiter, Literal: "(", Line: 2, Column: 1},
				{Type: Identifier, Literal: "define", Line: 2, Column: 2},
				{Type: Identifier, Literal: "y", Line: 2, Column: 9},
				{Type: Literal, Literal: "20", Line: 2, Column: 11},
				{Type: Delimiter, Literal: ")", Line: 2, Column: 13},
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
