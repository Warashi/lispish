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
	TokenIdentifier  // 識別子
	TokenInteger     // 整数
	TokenFloat       // 浮動小数点数
	TokenString      // 文字列リテラル
	TokenComment     // コメント
)

// String は TokenType の文字列表現を返します。
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
	case TokenInteger:
		return "Integer"
	case TokenFloat:
		return "Float"
	case TokenString:
		return "String"
	case TokenComment:
		return "Comment"
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
	// モードを設定：識別子、文字列、整数、浮動小数点を認識
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanInts | scanner.ScanFloats
	// デフォルトの Whitespace には改行('\n')も含まれるため、コメント終了検出のために改行は除外する
	s.Whitespace = scanner.GoWhitespace &^ (1 << '\n')
	// Scheme では識別子に記号などが使われることがあるため、IsIdentRune を上書き
	s.IsIdentRune = func(ch rune, i int) bool {
		// '#' はどこでも許容（例: #t, #f など）
		if ch == '#' {
			return true
		}
		// Scheme の識別子に使われる記号を許容
		switch ch {
		case '!', '$', '%', '&', '*', '+', '-', '/', ':', '<', '=', '>', '?', '^', '_', '~':
			return true
		}
		// 2文字目以降なら数字も許容
		if i > 0 && unicode.IsDigit(ch) {
			return true
		}
		// それ以外は Unicode の文字（アルファベット）を許容
		return unicode.IsLetter(ch)
	}
	return &Lexer{s: s}
}

// NextToken は入力から次のトークンを返します。
// Scheme のコメント（';' から行末まで）は読み飛ばします。
func (l *Lexer) NextToken() Token {
	for {
		tok := l.s.Scan()
		if tok == scanner.EOF {
			return Token{Type: TokenEOF, Literal: ""}
		}
		text := l.s.TokenText()

		// セミコロン ';' で始まる場合、コメント行として改行まで読み飛ばす
		if text == ";" {
			// Save the current whitespace flag
			originalWhitespace := l.s.Whitespace
			// Include all whitespace characters during comment processing
			l.s.Whitespace = 0

			commentText := ";"
			for {
				tok = l.s.Scan()
				if tok == '\n' || tok == scanner.EOF {
					break
				}
				commentText += l.s.TokenText()
			}

			// Restore the original whitespace flag
			l.s.Whitespace = originalWhitespace
			return Token{Type: TokenComment, Literal: commentText}
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
		case scanner.Int:
			return Token{Type: TokenInteger, Literal: text}
		case scanner.Float:
			return Token{Type: TokenFloat, Literal: text}
		case scanner.Ident:
			return Token{Type: TokenIdentifier, Literal: text}
		default:
			// 改行、タブ、スペースなどはスキップ
			if tok == '\n' || tok == '\r' || tok == '\t' || tok == ' ' {
				continue
			}
			// 上記以外は識別子として扱う
			return Token{Type: TokenIdentifier, Literal: text}
		}
	}
}
