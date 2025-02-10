package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `(define x 42)
; これはコメントです
(define y "hello world")
(quote (a b c))
#t #f
'(-3.14 foo)`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		// (define x 42)
		{LPAREN, "("},
		{SYMBOL, "define"},
		{SYMBOL, "x"},
		{NUMBER, "42"},
		{RPAREN, ")"},

		// コメント行はスキップされるのでテスト対象外

		// (define y "hello world")
		{LPAREN, "("},
		{SYMBOL, "define"},
		{SYMBOL, "y"},
		{STRING, "hello world"},
		{RPAREN, ")"},

		// (quote (a b c))
		{LPAREN, "("},
		{SYMBOL, "quote"},
		{LPAREN, "("},
		{SYMBOL, "a"},
		{SYMBOL, "b"},
		{SYMBOL, "c"},
		{RPAREN, ")"},
		{RPAREN, ")"},

		// #t #f
		{BOOLEAN, "#t"},
		{BOOLEAN, "#f"},

		// '(-3.14 foo)
		{LPAREN, "("},
		{QUOTE, "'"},
		{LPAREN, "("},
		{NUMBER, "-3.14"},
		{SYMBOL, "foo"},
		{RPAREN, ")"},
		{RPAREN, ")"},

		// 入力終了
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - トークンタイプが不正です。期待値=%q, 実際=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - リテラルが不正です。期待値=%q, 実際=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}