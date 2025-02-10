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
; コメント中のホワイトスペースを含む
; コメント中の	タブを含む
`

	lexer := NewLexer(strings.NewReader(input))

	// 期待するトークン列
	expectedTokens := []Token{
		{Type: TokenComment, Literal: "; コメント行をスキップする"},
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
		{Type: TokenInteger, Literal: "1"},
		{Type: TokenInteger, Literal: "2"},
		{Type: TokenString, Literal: "three"},
		{Type: TokenFloat, Literal: "4.0"},
		{Type: TokenRParen, Literal: ")"},
		{Type: TokenComment, Literal: "; コメント中のホワイトスペースを含む"},
		{Type: TokenComment, Literal: "; コメント中の	タブを含む"},
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

func TestLexerWithAdditionalCases(t *testing.T) {
	input := `
123
456.789
`

	lexer := NewLexer(strings.NewReader(input))

	// 期待するトークン列
	expectedTokens := []Token{
		{Type: TokenInteger, Literal: "123"},
		{Type: TokenFloat, Literal: "456.789"},
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
