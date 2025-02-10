package lexer

import (
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
; コメント行をスキップする
(define (square x)
  (* x x))
'(1 2 "three" 4.0)
`

	lexer := NewLexer(strings.NewReader(input))

	// 期待するトークン列を定義
	expectedTokens := []Token{
		{Type: TokenLParen, Literal: "("},
		{Type: TokenIdentifier, Literal: "define"},
		{Type: TokenLParen, Literal: "("},
		{Type: TokenIdentifier, Literal: "square"},
		{Type: TokenIdentifier, Literal: "x"},
		{Type: TokenRParen, Literal: ")"},
		{Type: TokenLParen, Literal: "("},
		{Type: TokenIdentifier, Literal: "*"},
		{Type: TokenIdentifier, Literal: "x"},
		{Type: TokenIdentifier, Literal: "x"},
		{Type: TokenRParen, Literal: ")"},
		{Type: TokenRParen, Literal: ")"},
		{Type: TokenQuote, Literal: "'"},
		{Type: TokenLParen, Literal: "("},
		{Type: TokenNumber, Literal: "1"},
		{Type: TokenNumber, Literal: "2"},
		{Type: TokenString, Literal: "three"},
		{Type: TokenNumber, Literal: "4.0"},
		{Type: TokenRParen, Literal: ")"},
		{Type: TokenEOF, Literal: ""},
	}

	for i, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.Type || token.Literal != expected.Literal {
			t.Errorf("Token %d: expected (%s, %q), got (%s, %q)",
				i, expected.Type, expected.Literal, token.Type, token.Literal)
		}
	}
}