package evaluator

import (
	"fmt"

	"github.com/Warashi/lispish/parser"
)

// Env は変数とその値の束縛を保持する環境です。
// outer があれば、ネストした環境（静的スコープ）を実現します。
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

// Callable インターフェースは、関数オブジェクトとして呼び出し可能なものが実装すべきメソッドを定義します。
type Callable interface {
	// Call は引数を受け取り、その評価結果を返します。
	Call(args []parser.Expr) (parser.Expr, error)
}

// Builtin は組み込み関数を表す型です。
// 新たな組み込み関数を追加する場合、Name と実際の関数処理（Fn）を設定してインスタンス化してください。
type Builtin struct {
	Name string
	Fn   func(args []parser.Expr) (parser.Expr, error)
}

// Call により、組み込み関数を呼び出します。
func (b *Builtin) Call(args []parser.Expr) (parser.Expr, error) {
	return b.Fn(args)
}

// Closure はユーザ定義の関数（lambda式）のクロージャを表します。
type Closure struct {
	params []parser.Symbol
	body   parser.Expr
	env    *Env
}

// Call により、クロージャ内の式を引数付きで評価します。
func (c *Closure) Call(args []parser.Expr) (parser.Expr, error) {
	if len(args) != len(c.params) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(c.params), len(args))
	}
	newEnv := NewEnv(c.env)
	for i, param := range c.params {
		newEnv.Set(param, args[i])
	}
	return Eval(c.body, newEnv)
}

// Eval は AST（parser.Expr）を評価し、その結果を返します。
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
				// (define var expr) または (define (fun arg...) body...)
				if len(exp) < 3 {
					return nil, fmt.Errorf("define: too few arguments")
				}
				// 関数定義の短縮形の場合
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
					// 変数定義の場合: (define var expr)
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
				// (lambda (params...) body...) → クロージャを生成して返す
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

		// 引数は評価する
		var args []parser.Expr
		for _, arg := range exp[1:] {
			evaluatedArg, err := Eval(arg, env)
			if err != nil {
				return nil, err
			}
			args = append(args, evaluatedArg)
		}

		// op が Callable インターフェースを実装しているかチェック
		callable, ok := op.(Callable)
		if !ok {
			return nil, fmt.Errorf("not a function: %v", op)
		}
		return callable.Call(args)

	// コメントはそのまま返す（実行時には無視してもよい）
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

// --- 組み込み関数の実装例 ---

// builtinAdd は "+" を実装します。
// 整数・浮動小数点数に対して加算を行います。
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

// builtinMul は "*" を実装します。
// 引数が整数または浮動小数点数の場合に乗算を行います。
func builtinMul(args []parser.Expr) (parser.Expr, error) {
	if len(args) == 0 {
		return parser.Integer(1), nil
	}
	isFloat := false
	prodInt := int64(1)
	prodFloat := 1.0
	for _, arg := range args {
		switch v := arg.(type) {
		case parser.Integer:
			prodInt *= int64(v)
			prodFloat *= float64(v)
		case parser.Float:
			isFloat = true
			prodFloat *= float64(v)
		default:
			return nil, fmt.Errorf("*: invalid argument type %T", arg)
		}
	}
	if isFloat {
		return parser.Float(prodFloat), nil
	}
	return parser.Integer(prodInt), nil
}

// NewGlobalEnv は、組み込み関数などが登録されたグローバル環境を生成して返します。
// 新たな組み込み関数を追加する場合は、ここに env.Set() を追加してください。
func NewGlobalEnv() *Env {
	env := NewEnv(nil)
	env.Set("+", &Builtin{
		Name: "+",
		Fn:   builtinAdd,
	})
	env.Set("*", &Builtin{
		Name: "*",
		Fn:   builtinMul,
	})
	// 必要に応じて他の組み込み関数（例: "-", "/" など）を追加可能です。
	return env
}
