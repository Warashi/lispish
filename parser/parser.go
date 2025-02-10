package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/Warashi/lispish/lexer"
)

// Expr は Scheme の式を表すインターフェースです。
type Expr interface{}

// Symbol は Scheme のシンボル（識別子）を表します。
type Symbol string

// Integer は整数リテラルを表します。
type Integer int64

// Float は浮動小数点数リテラルを表します。
type Float float64

// String は文字列リテラルを表します。
type String string

// List は Scheme のリスト（S式）を表します。
type List []Expr

// Parser は lexer からのトークンをもとに Scheme の式を構文解析します。
type Parser struct {
	l        *lexer.Lexer
	curToken lexer.Token
}

// NewParser は入力リーダーからパーサを初期化して返します。
func NewParser(r io.Reader) *Parser {
	p := &Parser{
		l: lexer.NewLexer(r),
	}
	p.nextToken() // 最初のトークンを取得
	return p
}

// nextToken は次のトークンを取得します（コメントはスキップ）。
func (p *Parser) nextToken() {
	tok := p.l.NextToken()
	// コメントトークンは読み飛ばす
	for tok.Type == lexer.TokenComment {
		tok = p.l.NextToken()
	}
	p.curToken = tok
}

// ParseExpr は1つの Scheme 式をパースして返します。
func (p *Parser) ParseExpr() (Expr, error) {
	switch p.curToken.Type {
	case lexer.TokenEOF:
		return nil, io.EOF
	case lexer.TokenInteger:
		// 整数リテラルをパース
		val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer literal: %s", p.curToken.Literal)
		}
		expr := Integer(val)
		p.nextToken()
		return expr, nil
	case lexer.TokenFloat:
		// 浮動小数点数リテラルをパース
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float literal: %s", p.curToken.Literal)
		}
		expr := Float(val)
		p.nextToken()
		return expr, nil
	case lexer.TokenString:
		// 文字列リテラル
		expr := String(p.curToken.Literal)
		p.nextToken()
		return expr, nil
	case lexer.TokenIdentifier:
		// 識別子はシンボルとして扱う
		expr := Symbol(p.curToken.Literal)
		p.nextToken()
		return expr, nil
	case lexer.TokenLParen:
		return p.parseList()
	case lexer.TokenQuote:
		return p.parseQuote()
	case lexer.TokenRParen:
		return nil, fmt.Errorf("unexpected ')'")
	default:
		return nil, fmt.Errorf("unexpected token: %v", p.curToken)
	}
}

// parseList はリスト式をパースします。
func (p *Parser) parseList() (Expr, error) {
	// 現在のトークンは '(' なので、これを消費
	p.nextToken()
	var list List
	// ')' が現れるまで式を読み込む
	for p.curToken.Type != lexer.TokenRParen {
		if p.curToken.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF while reading list")
		}
		expr, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
	}
	// 終了括弧 ')' を消費
	p.nextToken()
	return list, nil
}

// parseQuote は引用式をパースします。
// 例: 'expr  → (quote expr)
func (p *Parser) parseQuote() (Expr, error) {
	// クォートトークンを消費
	p.nextToken()
	expr, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}
	// (quote <expr>) として返す
	return List{Symbol("quote"), expr}, nil
}

// ParseAll は入力全体から式を読み込み、式のスライスを返します。
func (p *Parser) ParseAll() ([]Expr, error) {
	var exprs []Expr
	for p.curToken.Type != lexer.TokenEOF {
		expr, err := p.ParseExpr()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	return exprs, nil
}
