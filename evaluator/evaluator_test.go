package evaluator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Warashi/lispish/parser"
)

// TestEvaluatorArithmetic は基本的な算術演算 (+) の評価結果をテストします。
func TestEvaluatorArithmetic(t *testing.T) {
	input := "(+ 1 2 3)"
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(6)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorMultiplication は (*) の評価結果をテストします。
func TestEvaluatorMultiplication(t *testing.T) {
	input := "(* 2 3 4)"
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(24)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorQuote はクォート式の評価をテストします。
// クォート式は引数を評価せずにそのまま返すため、 '(1 2 3) はリスト [1, 2, 3] として返されるはずです。
func TestEvaluatorQuote(t *testing.T) {
	input := "'(1 2 3)"
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.List{parser.Integer(1), parser.Integer(2), parser.Integer(3)}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorDefine は define を使った変数定義の評価をテストします。
func TestEvaluatorDefine(t *testing.T) {
	input := `
	(define x 42)
	x
	`
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(42)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorLambda は lambda（および define を使った関数定義）の評価をテストします。
func TestEvaluatorLambda(t *testing.T) {
	input := `
	(define (square x) (* x x))
	(square 5)
	`
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(25)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorNestedExpressions は複数の定義および関数適用を含む式群の評価をテストします。
func TestEvaluatorNestedExpressions(t *testing.T) {
	input := `
	(define (square x) (* x x))
	(define (add a b) (+ a b))
	(add (square 3) (square 4))
	`
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(25)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// TestEvaluatorComments はコメントを含む入力の評価結果が正しいことをテストします。
// コメントは評価対象そのものとしては扱われますが、最終的な評価結果は後続の式に依存します。
func TestEvaluatorComments(t *testing.T) {
	input := `
	; This is a comment
	(define x 10) ; Another comment
	; Comment before variable reference
	x
	`
	p := parser.NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}
	env := NewGlobalEnv()
	// 複数の式を評価した場合、 EvalAll は最後の式の評価結果を返すので x の値が返るはずです。
	result, err := EvalAll(exprs, env)
	if err != nil {
		t.Fatalf("EvalAll error: %v", err)
	}
	expected := parser.Integer(10)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
