package parser

import (
	"reflect"
	"strings"
	"testing"
)

// TestParser_SimpleExpressions tests the parsing of simple expressions.
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

	// Test 1: (define (square x) (* x x))
	expr1, ok := exprs[0].(List)
	if !ok {
		t.Fatalf("expected first expression to be a List, got %T", exprs[0])
	}
	if len(expr1) != 3 {
		t.Fatalf("expected first list to have 3 elements, got %d", len(expr1))
	}
	// Check that the first element is 'define'
	if sym, ok := expr1[0].(Symbol); !ok || sym != "define" {
		t.Errorf("expected first element to be Symbol 'define', got %v", expr1[0])
	}
	// Check that the second element is (square x)
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

	// Check that the third element is (* x x)
	innerList2, ok := expr1[2].(List)
	if !ok || len(innerList2) != 3 {
		t.Errorf("expected third element to be a List of length 3, got %v", expr1[2])
	} else {
		if sym, ok := innerList2[0].(Symbol); !ok || sym != "*" {
			t.Errorf("expected first element of inner list to be Symbol '*', got %v", innerList2[0])
		}
	}

	// Test 2: '(1 2 "three" 4.0)
	expr2, ok := exprs[1].(List)
	if !ok || len(expr2) != 2 {
		t.Fatalf("expected second expression to be a List of length 2, got %v", exprs[1])
	}
	// Since it's a quoted expression, the first element should be 'quote'
	if sym, ok := expr2[0].(Symbol); !ok || sym != "quote" {
		t.Errorf("expected first element of quoted expression to be 'quote', got %v", expr2[0])
	}
	// The second element should be the list (1 2 "three" 4.0)
	quotedList, ok := expr2[1].(List)
	if !ok || len(quotedList) != 4 {
		t.Fatalf("expected quoted list to be a List of length 4, got %v", expr2[1])
	}
	// Verify each element
	expected := []Expr{Integer(1), Integer(2), String("three"), Float(4.0)}
	for i, exp := range expected {
		if !reflect.DeepEqual(quotedList[i], exp) {
			t.Errorf("at index %d, expected %v, got %v", i, exp, quotedList[i])
		}
	}
}

// TestParser_QuoteExpression tests the parsing of quoted expressions.
func TestParser_QuoteExpression(t *testing.T) {
	input := "'(a b c)"
	p := NewParser(strings.NewReader(input))
	expr, err := p.ParseExpr()
	if err != nil {
		t.Fatalf("ParseExpr error: %v", err)
	}
	// expr should be represented as (quote (a b c))
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

// TestParser_Comments tests the parsing of comments.
func TestParser_Comments(t *testing.T) {
	input := `
    ; This is a comment
    (define x 42) ; Another comment
    ; Comment before quoted expression
    '(1 2 3) ; Comment after quoted expression
    `
	p := NewParser(strings.NewReader(input))
	exprs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("ParseAll error: %v", err)
	}

	if len(exprs) != 6 {
		t.Fatalf("expected 6 expressions, got %d", len(exprs))
	}

	// Test 1: ; This is a comment
	if comment, ok := exprs[0].(Comment); !ok || comment != "; This is a comment" {
		t.Errorf("expected first expression to be a Comment, got %v", exprs[0])
	}

	// Test 2: (define x 42)
	expr1, ok := exprs[1].(List)
	if !ok {
		t.Fatalf("expected second expression to be a List, got %T", exprs[1])
	}
	if len(expr1) != 3 {
		t.Fatalf("expected second list to have 3 elements, got %d", len(expr1))
	}
	// Check that the first element is 'define'
	if sym, ok := expr1[0].(Symbol); !ok || sym != "define" {
		t.Errorf("expected first element to be Symbol 'define', got %v", expr1[0])
	}
	// Check that the second element is 'x'
	if sym, ok := expr1[1].(Symbol); !ok || sym != "x" {
		t.Errorf("expected second element to be Symbol 'x', got %v", expr1[1])
	}
	// Check that the third element is 42
	if num, ok := expr1[2].(Integer); !ok || num != 42 {
		t.Errorf("expected third element to be Integer 42, got %v", expr1[2])
	}

	// Test 3: ; Another comment
	if comment, ok := exprs[2].(Comment); !ok || comment != "; Another comment" {
		t.Errorf("expected third expression to be a Comment, got %v", exprs[2])
	}

	// Test 4: ; Comment before quoted expression
	if comment, ok := exprs[3].(Comment); !ok || comment != "; Comment before quoted expression" {
		t.Errorf("expected fourth expression to be a Comment, got %v", exprs[3])
	}

	// Test 5: '(1 2 3)
	expr2, ok := exprs[4].(List)
	if !ok || len(expr2) != 2 {
		t.Fatalf("expected fifth expression to be a List of length 2, got %v", exprs[4])
	}
	// Since it's a quoted expression, the first element should be 'quote'
	if sym, ok := expr2[0].(Symbol); !ok || sym != "quote" {
		t.Errorf("expected first element of quoted expression to be 'quote', got %v", expr2[0])
	}
	// The second element should be the list (1 2 3)
	quotedList, ok := expr2[1].(List)
	if !ok || len(quotedList) != 3 {
		t.Fatalf("expected quoted list to be a List of length 3, got %v", expr2[1])
	}
	// Verify each element
	expected := []Expr{Integer(1), Integer(2), Integer(3)}
	for i, exp := range expected {
		if !reflect.DeepEqual(quotedList[i], exp) {
			t.Errorf("at index %d, expected %v, got %v", i, exp, quotedList[i])
		}
	}

	// Test 6: ; Comment after quoted expression
	if comment, ok := exprs[5].(Comment); !ok || comment != "; Comment after quoted expression" {
		t.Errorf("expected sixth expression to be a Comment, got %v", exprs[5])
	}
}
