package icfp

import (
	"fmt"
	"math/big"
	"strings"
)

type Expr interface {
	IsExpr()
}

type Boolean bool
type Integer struct {
	*big.Int
}
type String string
type Unop struct {
	Op  string
	Arg Expr
}
type Binop struct {
	Op    string
	Left  Expr
	Right Expr
}
type If struct {
	Test Expr
	Then Expr
	Else Expr
}
type Lambda struct {
	Param int64
	Body  Expr
}
type Var struct {
	v int64
}

func (b Boolean) IsExpr() {}
func (i Integer) IsExpr() {}
func (s String) IsExpr()  {}
func (i If) IsExpr()      {}
func (u Unop) IsExpr()    {}
func (b Binop) IsExpr()   {}
func (l Lambda) IsExpr()  {}
func (v Var) IsExpr()     {}

const lookup = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\"#$%&'()*+,-./:;<=>?@[\\]^_`|~ \n"

func ParseInteger(s string) *big.Int {
	ret := big.NewInt(0)
	for i := 0; i < len(s); i++ {
		ret = ret.Mul(ret, big.NewInt(94))
		ret = ret.Add(ret, big.NewInt(int64(s[i])-33))
	}
	return ret
}

func ParseToken(token string) Expr {
	indicator := token[0]
	switch indicator {
	case 'T':
		return Boolean(true)
	case 'F':
		return Boolean(false)
	case 'I':
		return Integer{ParseInteger(token[1:])}
	case 'S':
		s := ""
		for i := 1; i < len(token); i++ {
			s += string(lookup[int(token[i])-33])
		}
		return String(s)
	case '?':
		return If{}
	case 'B':
		return Binop{Op: token[1:]}
	case 'U':
		return Unop{Op: token[1:]}
	case 'L':
		param := ParseInteger(token[1:])
		return Lambda{Param: param.Int64()}
	case 'v':
		v := ParseInteger(token[1:])
		return Var{v: v.Int64()}
	default:
		panic(fmt.Sprintf("Unknown token: %s", token))
	}

}

func Parse(s string) []Expr {
	tokens := strings.Split(s, " ")
	var ret []Expr
	for _, token := range tokens {
		ret = append(ret, ParseToken(token))
	}
	return ret
}

func CombineToExpr(exprs []Expr) (Expr, []Expr) {
	if len(exprs) == 0 {
		panic("Empty exprs")
	}
	expr := exprs[0]
	exprs = exprs[1:]
	switch v := expr.(type) {
	case Integer, Boolean, String:
		return v, exprs
	case If:
		test, exprs := CombineToExpr(exprs)
		then, exprs := CombineToExpr(exprs)
		els, exprs := CombineToExpr(exprs)
		return If{test, then, els}, exprs
	case Binop:
		left, exprs := CombineToExpr(exprs)
		right, exprs := CombineToExpr(exprs)
		return Binop{v.Op, left, right}, exprs
	case Lambda:
		body, exprs := CombineToExpr(exprs)
		return Lambda{Param: v.Param, Body: body}, exprs
	case Var:
		return v, exprs
	case Unop:
		arg, exprs := CombineToExpr(exprs)
		return Unop{v.Op, arg}, exprs
	default:
		panic(fmt.Sprintf("Unknown type: %T", expr))
	}
}

var cache map[Expr]map[Expr]Expr

func readCache(key1 Expr, key2 Expr) (Expr, bool) {
	if cache == nil {
		return nil, false
	}
	if _, ok := cache[key1]; !ok {
		return nil, false
	}
	if val, ok := cache[key1][key2]; ok {
		if val == nil {
			return nil, false
		}
		return val, true
	}
	return nil, false
}

func writeCache(key1 Expr, key2 Expr, val Expr) {
	if cache == nil {
		cache = make(map[Expr]map[Expr]Expr)
	}
	if _, ok := cache[key1]; !ok {
		cache[key1] = make(map[Expr]Expr)
	}
	cache[key1][key2] = val
}

func Eval(expr Expr) Expr {
	switch v := expr.(type) {
	case Integer, Boolean, String:
		return v
	case If:
		test := Eval(v.Test).(Boolean)
		if test {
			return Eval(v.Then)
		} else {
			return Eval(v.Else)
		}
	case Binop:
		switch v.Op {
		case "=":
			left := Eval(v.Left)
			right := Eval(v.Right)
			i, oki := left.(Integer)
			j, okj := right.(Integer)
			if oki && okj {
				return Boolean(i.Cmp(j.Int) == 0)
			}
			return Boolean(left == right)
		case "$":
			// fmt.Printf("Beta-reduction: %v %v\n", v.Left, v.Right)
			// fmt.Printf(".")
			if val, ok := readCache(v.Left, v.Right); ok {
				return val
			}
			lambda := Eval(v.Left).(Lambda)
			reduced := Substitute(lambda.Body, lambda.Param, v.Right)
			res := Eval(reduced)
			writeCache(v.Left, v.Right, res)
			return res
		case "T":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(String)
			return right[0:left.Int64()]
		case "D":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(String)
			return right[left.Int64():]
		case ".":
			left := Eval(v.Left).(String)
			right := Eval(v.Right).(String)
			return left + right
		case "&":
			left := Eval(v.Left).(Boolean)
			right := Eval(v.Right).(Boolean)
			return Boolean(bool(left) && bool(right))
		case "|":
			left := Eval(v.Left).(Boolean)
			right := Eval(v.Right).(Boolean)
			return Boolean(bool(left) || bool(right))
		case "<":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			cmp := left.Cmp(right.Int)
			return Boolean(cmp == -1)
		case ">":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			cmp := left.Cmp(right.Int)
			return Boolean(cmp == 1)
		case "%":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)

			z := big.NewInt(0).Rem(left.Int, right.Int)
			return Integer{Int: z}
		case "/":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			z := big.NewInt(0).Quo(left.Int, right.Int)
			return Integer{Int: z}
		case "*":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			z := big.NewInt(0).Mul(left.Int, right.Int)
			return Integer{Int: z}
		case "+":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			z := big.NewInt(0).Add(left.Int, right.Int)
			return Integer{Int: z}
		case "-":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			z := big.NewInt(0).Sub(left.Int, right.Int)
			return Integer{Int: z}
		default:
			panic(fmt.Sprintf("Unknown binop: %s", v.Op))
		}
	case Unop:
		switch v.Op {
		case "-":
			arg := Eval(v.Arg).(Integer)
			z := big.NewInt(0).Neg(arg.Int)
			return Integer{Int: z}
		case "!":
			arg := Eval(v.Arg).(Boolean)
			return Boolean(!arg)
		case "$":
			i := Eval(v.Arg).(Integer)
			s := ""
			for i.Cmp(big.NewInt(0)) != 0 {
				d := big.NewInt(0).Mod(i.Int, big.NewInt(94))
				i = Integer{Int: big.NewInt(0).Div(i.Int, big.NewInt(94))}
				s = string(lookup[d.Int64()]) + s
			}
			return String(s)
		case "#":
			s := Eval(v.Arg).(String)
			i := big.NewInt(0)
			for _, c := range s {
				i.Mul(i, big.NewInt(94))
				i.Add(i, big.NewInt(int64(strings.Index(lookup, string(c)))))
			}
			return Integer{Int: i}
		default:
			panic(fmt.Sprintf("Unknown unop: %s", v.Op))
		}
	case Lambda:
		return Lambda{v.Param, v.Body}
	case Var:
		// Note: This should have been substituted in a beta reduction
		// before trying to evaluate if it was in scope.
		panic(fmt.Sprintf("Variable %d not found", v))
	default:
		panic(fmt.Sprintf("Unknown type: %T", expr))
	}
}

func Substitute(expr Expr, v int64, val Expr) Expr {
	switch e := expr.(type) {
	case Integer, Boolean, String:
		return e
	case If:
		return If{Substitute(e.Test, v, val), Substitute(e.Then, v, val), Substitute(e.Else, v, val)}
	case Binop:
		return Binop{e.Op, Substitute(e.Left, v, val), Substitute(e.Right, v, val)}
	case Unop:
		return Unop{e.Op, Substitute(e.Arg, v, val)}
	case Lambda:
		if e.Param == v {
			return e
		}
		// if e.Param.Cmp(v) == 0 {
		// 	return e
		// }
		return Lambda{e.Param, Substitute(e.Body, v, val)}
	case Var:
		if e.v == v {
			return val
		}
		return e
	default:
		panic(fmt.Sprintf("Unknown type: %T", expr))
	}
}

func StringToToken(s string) String {
	var ret String = "S"
	for _, c := range s {
		ret += String(byte(strings.Index(lookup, string(c)) + 33))
	}
	return ret
}

func RenderAsLambda(e Expr) string {
	varLookup := []string{"x", "y", "z", "w", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	switch v := e.(type) {
	case Integer:
		return fmt.Sprintf("%d", v)
	case Boolean:
		return fmt.Sprintf("%t", v)
	case String:
		return fmt.Sprintf("%q", string(v))
	case If:
		return fmt.Sprintf("(if %s %s %s)", RenderAsLambda(v.Test), RenderAsLambda(v.Then), RenderAsLambda(v.Else))
	case Binop:
		if v.Op == "$" {
			return fmt.Sprintf("(%s %s)", RenderAsLambda(v.Left), RenderAsLambda(v.Right))
		}
		return fmt.Sprintf("(%s %s %s)", v.Op, RenderAsLambda(v.Left), RenderAsLambda(v.Right))
	case Unop:
		return fmt.Sprintf("(%s %s)", v.Op, RenderAsLambda(v.Arg))
	case Lambda:
		return fmt.Sprintf("(λ%s.%s)", varLookup[v.Param], RenderAsLambda(v.Body))
	case Var:
		return varLookup[v.v]
	default:
		panic(fmt.Sprintf("Unknown type: %T", e))
	}
}

// (((λy.((λz.(y (z z))) (λz.(y (z z))))) (λy.(λz.(if (= z 0) "[" (. (y (/ z 39)) (T 1 (D (% z 39) "[t=1, x6y]\n.0S^2>v+/-345Crashed:TickLmE"))))))) 6501638242769916696)

// Y = (λy.((λz.(y (z z))) (λz.(y (z z)))))
// f = (λy.(λz.(if (= z 0) "[" (. (y (/ z 39)) (T 1 (D (% z 39) "[t=1, x6y]\n.0S^2>v+/-345Crashed:TickLmE"))))))
// ((Y f) 6501638242769916696)
