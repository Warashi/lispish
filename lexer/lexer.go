package lexer

import (
	"io"
	"strconv"
	"text/scanner"
	"unicode"
)

// TokenType はトークンの種類を表します。
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenLParen      // (
	TokenRParen      // )
	TokenQuote       // '
	TokenIdentifier  // 変数名など
	TokenNumber      // 数値（整数・浮動小数点）
	TokenString      // 文字列リテラル
)

// String は TokenType のデバッグ用文字列表現を返します。
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenLParen:
		return "LParen"
	case TokenRParen:
		return "RParen"
	case TokenQuote:
		return "Quote"
	case TokenIdentifier:
		return "Identifier"
	case TokenNumber:
		return "Number"
	case TokenString:
		return "String"
	default:
		return "Unknown"
	}
}

// Token は字句解析された1単位（トークン）を表します。
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer は Scheme の入力を走査する字句解析器です。
type Lexer struct {
	s scanner.Scanner
}

// NewLexer は io.Reader から入力を受け取り、Lexer を初期化して返します。
func NewLexer(r io.Reader) *Lexer {
	var s scanner.Scanner
	s.Init(r)
	// 識別子、文字列、整数、浮動小数点数をスキャンするように設定
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanInts | scanner.ScanFloats

	// Scheme では識別子に一部記号が含まれることがあるため、識別子のルールを拡張
	s.IsIdentRune = func(ch rune, i int) bool {
		// '#' を含む（例: #t, #f など）
		if ch == '#' {
			return true
		}
		// Scheme の識別子として使われる記号を許容
		if ch == '!' || ch == '$' || ch == '%' || ch == '&' ||
			ch == '*' || ch == '+' || ch == '-' || ch == '/' ||
			ch == ':' || ch == '<' || ch == '=' || ch == '>' ||
			ch == '?' || ch == '^' || ch == '_' || ch == '~' {
			return true
		}
		// 2文字目以降であれば数字も許容
		if i > 0 && unicode.IsDigit(ch) {
			return true
		}
		// その他は Unicode の文字（アルファベット）を許容
		return unicode.IsLetter(ch)
	}
	return &Lexer{s: s}
}

// NextToken は入力から次のトークンを返します。
// Scheme のコメント（セミコロン ';' から行末）は読み飛ばします。
func (l *Lexer) NextToken() Token {
	for {
		tok := l.s.Scan()
		// 入力終端なら EOF トークンを返す
		if tok == scanner.EOF {
			return Token{Type: TokenEOF, Literal: ""}
		}

		text := l.s.TokenText()

		// ';' で始まる場合、行末までコメントとして読み飛ばす
		if tok == ';' {
			for tok != '\n' && tok != scanner.EOF {
				tok = l.s.Scan()
			}
			continue
		}

		switch tok {
		case '(':
			return Token{Type: TokenLParen, Literal: text}
		case ')':
			return Token{Type: TokenRParen, Literal: text}
		case '\'':
			return Token{Type: TokenQuote, Literal: text}
		case scanner.String:
			// 文字列リテラルの場合、囲みのクォートを除去
			unquoted, err := strconv.Unquote(text)
			if err != nil {
				unquoted = text
			}
			return Token{Type: TokenString, Literal: unquoted}
		case scanner.Int, scanner.Float:
			return Token{Type: TokenNumber, Literal: text}
		case scanner.Ident:
			return Token{Type: TokenIdentifier, Literal: text}
		default:
			// 空白や改行などはスキップ
			if tok == '\n' || tok == '\r' || tok == '\t' || tok == ' ' {
				continue
			}
			// その他は識別子として扱う
			return Token{Type: TokenIdentifier, Literal: text}
		}
	}
}