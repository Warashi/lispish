package lexer

// TokenType は各トークンの種類を表す文字列です。
type TokenType string

const (
	ILLEGAL = "ILLEGAL" // 不正なトークン
	EOF     = "EOF"     // 入力の終端
	LPAREN  = "("       // 左括弧
	RPAREN  = ")"       // 右括弧
	DOT     = "."       // ドット（リストのドット記法用）
	QUOTE   = "'"       // クォート
	NUMBER  = "NUMBER"  // 数値リテラル
	STRING  = "STRING"  // 文字列リテラル
	SYMBOL  = "SYMBOL"  // シンボル
	BOOLEAN = "BOOLEAN" // #t, #f
)

// Token は字句解析器が認識したトークンを表現します。
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer は入力文字列を保持し、読み出し位置などの状態を管理します。
type Lexer struct {
	input   string
	pos     int  // 現在の位置（現在の文字を指す）
	readPos int  // 次の文字を読むための位置
	ch      byte // 現在注目している文字
}

// New は与えられた入力文字列から Lexer を初期化して返します。
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar は次の文字を読み進め、現在の文字（l.ch）を更新します。
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // 入力終了を 0 とする
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// peekChar は次に読む文字を返します（位置は更新しません）。
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// skipWhitespace は空白文字（スペース、タブ、改行等）を読み飛ばします。
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment はセミコロンで始まるコメント行（改行まで）を読み飛ばします。
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readString はダブルクォートで囲まれた文字列を読み取ります。
func (l *Lexer) readString() string {
	// 現在 l.ch は '"' と仮定し、これを読み飛ばす
	l.readChar()
	start := l.pos
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[start:l.pos]
	l.readChar() // 終了するダブルクォートを読み飛ばす
	return str
}

// readNumber は符号付きおよび小数部を持つ数値を読み取ります。
func (l *Lexer) readNumber() string {
	start := l.pos
	// マイナス符号の場合
	if l.ch == '-' {
		l.readChar()
	}
	for isDigit(l.ch) {
		l.readChar()
	}
	// 小数点付き数値の場合
	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[start:l.pos]
}

// readSymbol はシンボル（アルファベットや数字、記号の組み合わせ）を読み取ります。
func (l *Lexer) readSymbol() string {
	start := l.pos
	for isSymbolChar(l.ch) {
		l.readChar()
	}
	return l.input[start:l.pos]
}

// NextToken は入力から次のトークンを切り出して返します。
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token

	switch l.ch {
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch)}
	case '\'':
		tok = Token{Type: QUOTE, Literal: string(l.ch)}
	case '.':
		tok = Token{Type: DOT, Literal: string(l.ch)}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case ';':
		// コメント行の場合、コメントを読み飛ばして次のトークンを返す
		l.skipComment()
		return l.NextToken()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		return tok
	default:
		// 数字またはマイナス記号＋数字の場合は数値
		if isDigit(l.ch) || (l.ch == '-' && isDigit(l.peekChar())) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else if isInitialSymbol(l.ch) {
			tok.Literal = l.readSymbol()
			// Scheme では #t, #f はブール値として扱う
			if tok.Literal == "#t" || tok.Literal == "#f" {
				tok.Type = BOOLEAN
			} else {
				tok.Type = SYMBOL
			}
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// isDigit は ch が数字かどうかを判定します。
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter は ch がアルファベットかどうかを判定します。
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// isAllowedSymbolChar はシンボル中に使える記号かどうかを判定します。
func isAllowedSymbolChar(ch byte) bool {
	switch ch {
	case '!', '$', '%', '&', '*', '+', '-', '.', '/', ':', '<', '=', '>', '?', '@', '^', '_', '~', '#':
		return true
	}
	return false
}

// isSymbolChar はシンボルの文字として有効かどうかを判定します。
func isSymbolChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || isAllowedSymbolChar(ch)
}

// isInitialSymbol はシンボルの最初の文字として有効かどうかを判定します。
// 数字で始まる場合は数値と判断するため、通常は文字または許可された記号で始まる必要があります。
func isInitialSymbol(ch byte) bool {
	return isLetter(ch) || isAllowedSymbolChar(ch)
}