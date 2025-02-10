package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_SimpleExpressions(t *testing.T) {
	input := `
    (define (square x)
       (* x x))
    '(1 2 "three" 4.0)
    `
	p := NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}

	if len(exprs) != 2 {
		t.Fatalf("expected 2 expressions, got %d", len(exprs))
	}

	// テスト1: (define (square x) (* x x))
	expr1, ok := exprs[0].(List)
	if !ok {
		t.Fatalf("expected first expression to be a List, got %T", exprs[0])
	}
	if len(expr1) != 3 {
		t.Fatalf("expected first list to have 3 elements, got %d", len(expr1))
	}
	// 最初の要素が 'define' であることを確認
	if sym, ok := expr1[0].(Symbol); !ok || sym != "define" {
		t.Errorf("expected first element to be Symbol 'define', got %v", expr1[0])
	}
	// 2番目の要素は (square x) であることを確認
	innerList, ok := expr1[1].(List)
	if !ok || len(innerList) != 2 {
		t.Errorf("expected second element to be a List of length 2, got %v", expr1[1])
	} else {
		if sym, ok := innerList[0].(Symbol); !ok || sym != "square" {
			t.Errorf("expected first element of inner list to be Symbol 'square', got %v", innerList[0])
		}
		if sym, ok := innerList[1].(Symbol); !ok || sym != "x" {
			t.Errorf("expected second element of inner list to be Symbol 'x', got %v", innerList[1])
		}
	}

	// 3番目の要素は (* x x) であることを確認
	innerList2, ok := expr1[2].(List)
	if !ok || len(innerList2) != 3 {
		t.Errorf("expected third element to be a List of length 3, got %v", expr1[2])
	} else {
		if sym, ok := innerList2[0].(Symbol); !ok || sym != "*" {
			t.Errorf("expected first element of inner list to be Symbol '*', got %v", innerList2[0])
		}
	}

	// テスト2: '(1 2 "three" 4.0)
	expr2, ok := exprs[1].(List)
	if !ok || len(expr2) != 2 {
		t.Fatalf("expected second expression to be a List of length 2, got %v", exprs[1])
	}
	// クォート式なので、最初の要素は 'quote' であるはず
	if sym, ok := expr2[0].(Symbol); !ok || sym != "quote" {
		t.Errorf("expected first element of quoted expression to be 'quote', got %v", expr2[0])
	}
	// 2番目の要素は (1 2 "three" 4.0) のリスト
	quotedList, ok := expr2[1].(List)
	if !ok || len(quotedList) != 4 {
		t.Fatalf("expected quoted list to be a List of length 4, got %v", expr2[1])
	}
	// 各要素を検証
	expected := []Expr{Integer(1), Integer(2), String("three"), Float(4.0)}
	for i, exp := range expected {
		if !reflect.DeepEqual(quotedList[i], exp) {
			t.Errorf("at index %d, expected %v, got %v", i, exp, quotedList[i])
		}
	}
}

func TestParser_QuoteExpression(t *testing.T) {
	input := "'(a b c)"
	p := NewParser(strings.NewReader(input))
	expr, err := p.ParseExpr()
	if err != nil {
		t.Fatalf("ParseExpr error: %v", err)
	}
	// expr は (quote (a b c)) として表現されるはず
	listExpr, ok := expr.(List)
	if !ok || len(listExpr) != 2 {
		t.Fatalf("expected a List of length 2 for quote expression, got %v", expr)
	}
	if sym, ok := listExpr[0].(Symbol); !ok || sym != "quote" {
		t.Errorf("expected first element to be 'quote', got %v", listExpr[0])
	}
	quoted, ok := listExpr[1].(List)
	if !ok || len(quoted) != 3 {
		t.Fatalf("expected quoted part to be a List of length 3, got %v", listExpr[1])
	}
	expectedSymbols := []Symbol{"a", "b", "c"}
	for i, exp := range expectedSymbols {
		if sym, ok := quoted[i].(Symbol); !ok || sym != exp {
			t.Errorf("expected element %d to be %v, got %v", i, exp, quoted[i])
		}
	}
}
