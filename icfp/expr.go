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
	Env   Env
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

type Thunk struct {
	Expr      Expr
	Env       Env
	Value     Expr
	Evaluated bool
}

type Env map[int64]*Thunk

func copyEnv(e Env) Env {
	newEnv := make(Env)
	for k, v := range e {
		newEnv[k] = v
	}
	return newEnv
}

func Eval(expr Expr, env Env) Expr {
	fmt.Printf(".")
	switch v := expr.(type) {
	case Integer, Boolean, String:
		return v
	case Lambda:
		return Lambda{Param: v.Param, Body: v.Body, Env: copyEnv(env)}
	case Var:
		thunk, ok := env[v.v]
		if !ok {
			panic(fmt.Sprintf("Unknown variable: %d", v.v))
		}
		if !thunk.Evaluated {
			// fmt.Printf("Evaluating thunk: %s with env %v\n", RenderAsLambda(thunk.Expr), thunk.Env)
			thunk.Value = Eval(thunk.Expr, thunk.Env)
			thunk.Evaluated = true
		}
		return thunk.Value
	case Binop:
		var left, right Expr
		if v.Op != "$" {
			left = Eval(v.Left, env)
			right = Eval(v.Right, env)
		}
		switch v.Op {
		case "$":
			// fmt.Printf("Beta-reducing: %s with env %v\n", RenderAsLambda(v), env)
			lambda := Eval(v.Left, env).(Lambda)
			// fmt.Printf("Creating thunk for arg: %s with env %v\n", RenderAsLambda(v.Right), env)
			argThunk := &Thunk{
				Expr:      v.Right,
				Env:       env,
				Value:     nil,
				Evaluated: false,
			}
			newEnv := copyEnv(lambda.Env)
			newEnv[lambda.Param] = argThunk
			// fmt.Printf("Calling lambda: %s with env %v\n", RenderAsLambda(lambda.Body), newEnv)
			return Eval(lambda.Body, newEnv)
		case "=":
			i, oki := left.(Integer)
			j, okj := right.(Integer)
			if oki && okj {
				return Boolean(i.Cmp(j.Int) == 0)
			}
			return Boolean(left == right)
		case "T":
			return right.(String)[0:left.(Integer).Int64()]
		case "D":
			return right.(String)[left.(Integer).Int64():]
		case ".":
			return left.(String) + right.(String)
		case "&":
			return Boolean(bool(left.(Boolean)) && bool(right.(Boolean)))
		case "|":
			return Boolean(bool(left.(Boolean)) || bool(right.(Boolean)))
		case "<":
			cmp := left.(Integer).Cmp(right.(Integer).Int)
			return Boolean(cmp == -1)
		case ">":
			cmp := left.(Integer).Cmp(right.(Integer).Int)
			return Boolean(cmp == 1)
		case "%":
			z := big.NewInt(0).Rem(left.(Integer).Int, right.(Integer).Int)
			return Integer{Int: z}
		case "/":
			z := big.NewInt(0).Quo(left.(Integer).Int, right.(Integer).Int)
			return Integer{Int: z}
		case "*":
			if right.(Integer).Cmp(big.NewInt(0)) == 0 {
				return Integer{Int: big.NewInt(0)}
			}
			z := big.NewInt(0).Mul(left.(Integer).Int, right.(Integer).Int)
			return Integer{Int: z}
		case "+":
			z := big.NewInt(0).Add(left.(Integer).Int, right.(Integer).Int)
			return Integer{Int: z}
		case "-":
			z := big.NewInt(0).Sub(left.(Integer).Int, right.(Integer).Int)
			return Integer{Int: z}
		default:
			panic(fmt.Sprintf("Unknown binop: %s", v.Op))
		}
	case Unop:
		arg := Eval(v.Arg, env)
		switch v.Op {
		case "-":
			z := big.NewInt(0).Neg(arg.(Integer).Int)
			return Integer{Int: z}
		case "!":
			return Boolean(!arg.(Boolean))
		case "$":
			i := arg.(Integer)
			s := ""
			for i.Cmp(big.NewInt(0)) != 0 {
				d := big.NewInt(0).Mod(i.Int, big.NewInt(94))
				i = Integer{Int: big.NewInt(0).Div(i.Int, big.NewInt(94))}
				s = string(lookup[d.Int64()]) + s
			}
			return String(s)
		case "#":
			s := arg.(String)
			i := big.NewInt(0)
			for _, c := range s {
				i.Mul(i, big.NewInt(94))
				i.Add(i, big.NewInt(int64(strings.Index(lookup, string(c)))))
			}
			return Integer{Int: i}
		default:
			panic(fmt.Sprintf("Unknown unop: %s", v.Op))
		}
	case If:
		test := Eval(v.Test, env).(Boolean)
		if test {
			return Eval(v.Then, env)
		} else {
			return Eval(v.Else, env)
		}
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
		var s string
		if v.Param >= int64(len(varLookup)) {
			s = fmt.Sprintf("v%d", v.Param)
		} else {
			s = varLookup[v.Param]
		}
		return fmt.Sprintf("(λ%s.%s)", s, RenderAsLambda(v.Body))
	case Var:
		if v.v >= int64(len(varLookup)) {
			return fmt.Sprintf("v%d", v.v)
		}
		return varLookup[v.v]
	default:
		panic(fmt.Sprintf("Unknown type: %T", e))
	}
}
