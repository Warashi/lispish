package evaluator

import (
	"fmt"

	"github.com/Warashi/lispish/parser"
)

// Env は変数とその値の束縛を保持する環境です。
// outer があればネストした環境（静的スコープ）を実現します。
type Env struct {
	vars  map[parser.Symbol]parser.Expr
	outer *Env
}

// NewEnv は新しい環境を生成します。
func NewEnv(outer *Env) *Env {
	return &Env{
		vars:  make(map[parser.Symbol]parser.Expr),
		outer: outer,
	}
}

// Get はシンボルに束縛された値を探索します。
func (env *Env) Get(sym parser.Symbol) (parser.Expr, bool) {
	if val, ok := env.vars[sym]; ok {
		return val, true
	}
	if env.outer != nil {
		return env.outer.Get(sym)
	}
	return nil, false
}

// Set はシンボルと値の束縛を設定します。
func (env *Env) Set(sym parser.Symbol, val parser.Expr) {
	env.vars[sym] = val
}

// BuiltinFunc は組み込み関数の型です。
type BuiltinFunc func(args []parser.Expr) (parser.Expr, error)

// Builtin は組み込み関数のラッパーです。
type Builtin struct {
	Fn BuiltinFunc
}

// Closure はユーザ定義の関数（ラムダ式）のクロージャです。
type Closure struct {
	params []parser.Symbol
	body   parser.Expr
	env    *Env
}

// Eval はAST（parser.Expr）を評価し、その結果の値を返します。
func Eval(expr parser.Expr, env *Env) (parser.Expr, error) {
	switch exp := expr.(type) {
	// リテラルはそのまま返す
	case parser.Integer, parser.Float, parser.String:
		return exp, nil

	// シンボルは環境から値を取得
	case parser.Symbol:
		val, ok := env.Get(exp)
		if !ok {
			return nil, fmt.Errorf("undefined symbol: %s", exp)
		}
		return val, nil

	// リストは特殊フォームもしくは関数適用として評価する
	case parser.List:
		if len(exp) == 0 {
			return nil, fmt.Errorf("cannot evaluate empty list")
		}

		// 最初の要素がシンボルの場合、特殊フォームの可能性をチェック
		if firstSym, ok := exp[0].(parser.Symbol); ok {
			switch firstSym {
			case "quote":
				// (quote expr) → expr を評価せずに返す
				if len(exp) != 2 {
					return nil, fmt.Errorf("quote: wrong number of arguments")
				}
				return exp[1], nil

			case "define":
				// (define var expr) あるいは関数定義の短縮形
				if len(exp) < 3 {
					return nil, fmt.Errorf("define: too few arguments")
				}
				// 関数定義の形: (define (fun arg1 arg2 ...) body ...)
				if list, ok := exp[1].(parser.List); ok {
					if len(list) == 0 {
						return nil, fmt.Errorf("define: invalid function definition")
					}
					funName, ok := list[0].(parser.Symbol)
					if !ok {
						return nil, fmt.Errorf("define: function name must be a symbol")
					}
					var params []parser.Symbol
					for _, param := range list[1:] {
						s, ok := param.(parser.Symbol)
						if !ok {
							return nil, fmt.Errorf("define: function parameters must be symbols")
						}
						params = append(params, s)
					}
					// 複数の式があれば順次評価し、最後の値を返す（ここでは簡単のため List としてまとめる）
					var body parser.Expr
					if len(exp) == 3 {
						body = exp[2]
					} else {
						body = parser.List(exp[2:])
					}
					closure := &Closure{
						params: params,
						body:   body,
						env:    env,
					}
					env.Set(funName, closure)
					return funName, nil
				} else {
					// 変数定義: (define var expr)
					varName, ok := exp[1].(parser.Symbol)
					if !ok {
						return nil, fmt.Errorf("define: first argument must be a symbol")
					}
					value, err := Eval(exp[2], env)
					if err != nil {
						return nil, err
					}
					env.Set(varName, value)
					return varName, nil
				}

			case "lambda":
				// (lambda (params...) body...) → 関数（クロージャ）を返す
				if len(exp) < 3 {
					return nil, fmt.Errorf("lambda: too few arguments")
				}
				paramList, ok := exp[1].(parser.List)
				if !ok {
					return nil, fmt.Errorf("lambda: first argument must be a list of parameters")
				}
				var params []parser.Symbol
				for _, param := range paramList {
					s, ok := param.(parser.Symbol)
					if !ok {
						return nil, fmt.Errorf("lambda: parameters must be symbols")
					}
					params = append(params, s)
				}
				var body parser.Expr
				if len(exp) == 3 {
					body = exp[2]
				} else {
					body = parser.List(exp[2:])
				}
				return &Closure{
					params: params,
					body:   body,
					env:    env,
				}, nil
			}
		}

		// 関数適用の場合
		op, err := Eval(exp[0], env)
		if err != nil {
			return nil, err
		}
		// 引数を評価
		var args []parser.Expr
		for _, arg := range exp[1:] {
			evaluatedArg, err := Eval(arg, env)
			if err != nil {
				return nil, err
			}
			args = append(args, evaluatedArg)
		}
		// 関数適用
		switch fn := op.(type) {
		case *Builtin:
			return fn.Fn(args)
		case *Closure:
			if len(args) != len(fn.params) {
				return nil, fmt.Errorf("expected %d arguments, got %d", len(fn.params), len(args))
			}
			newEnv := NewEnv(fn.env)
			for i, param := range fn.params {
				newEnv.Set(param, args[i])
			}
			return Eval(fn.body, newEnv)
		default:
			return nil, fmt.Errorf("not a function: %v", op)
		}

	// コメントはそのまま返す（評価対象にならない）
	case parser.Comment:
		return exp, nil

	default:
		return nil, fmt.Errorf("cannot evaluate expression: %v", expr)
	}
}

// EvalAll は複数の式を順次評価し、最後の評価結果を返します。
func EvalAll(exprs []parser.Expr, env *Env) (parser.Expr, error) {
	var result parser.Expr
	var err error
	for _, expr := range exprs {
		result, err = Eval(expr, env)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// --- 組み込み関数 ---

// builtinMul は "*" を実装します。
// 引数が整数・浮動小数点の場合に乗算を行い、どれかが浮動小数点なら結果も浮動小数点となります。
func builtinMul(args []parser.Expr) (parser.Expr, error) {
	if len(args) == 0 {
		return parser.Integer(1), nil
	}
	isFloat := false
	productInt := int64(1)
	productFloat := 1.0
	for _, arg := range args {
		switch v := arg.(type) {
		case parser.Integer:
			productInt *= int64(v)
			productFloat *= float64(v)
		case parser.Float:
			isFloat = true
			productFloat *= float64(v)
		default:
			return nil, fmt.Errorf("*: invalid argument type %T", arg)
		}
	}
	if isFloat {
		return parser.Float(productFloat), nil
	}
	return parser.Integer(productInt), nil
}

// builtinAdd は "+" を実装します。
func builtinAdd(args []parser.Expr) (parser.Expr, error) {
	if len(args) == 0 {
		return parser.Integer(0), nil
	}
	isFloat := false
	sumInt := int64(0)
	sumFloat := 0.0
	for _, arg := range args {
		switch v := arg.(type) {
		case parser.Integer:
			sumInt += int64(v)
			sumFloat += float64(v)
		case parser.Float:
			isFloat = true
			sumFloat += float64(v)
		default:
			return nil, fmt.Errorf("+: invalid argument type %T", arg)
		}
	}
	if isFloat {
		return parser.Float(sumFloat), nil
	}
	return parser.Integer(sumInt), nil
}

// NewGlobalEnv は組み込み手続きなどを登録したグローバル環境を返します。
func NewGlobalEnv() *Env {
	env := NewEnv(nil)
	env.Set("+", &Builtin{Fn: builtinAdd})
	env.Set("*", &Builtin{Fn: builtinMul})
	// 必要に応じて他の組み込み手続き（例: "-", "/" など）を追加してください
	return env
}
