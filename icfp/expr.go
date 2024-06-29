package icfp

import (
	"fmt"
	"strings"
)

type Expr interface {
	IsExpr()
}

type Boolean bool
type Integer int64
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

func ParseToken(token string) Expr {
	indicator := token[0]
	switch indicator {
	case 'T':
		return Boolean(true)
	case 'F':
		return Boolean(false)
	case 'I':
		ret := int64(0)
		for i := 1; i < len(token); i++ {
			ret = ret*94 + int64(token[i]) - 33
		}
		return Integer(ret)
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
		ret := int64(0)
		for i := 1; i < len(token); i++ {
			ret = ret*94 + int64(token[i]) - 33
		}
		return Lambda{Param: ret}
	case 'v':
		ret := int64(0)
		for i := 1; i < len(token); i++ {
			ret = ret*94 + int64(token[i]) - 33
		}
		return Var{v: ret}
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
			return Boolean(left == right)
		case "$":
			lambda := Eval(v.Left).(Lambda)
			reduced := Substitute(lambda.Body, lambda.Param, v.Right)
			return Eval(reduced)
		case "T":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(String)
			return right[0:left]
		case "D":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(String)
			return right[left:]
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
			return Boolean(left < right)
		case ">":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Boolean(left > right)
		case "%":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Integer(left % right)
		case "/":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Integer(left / right)
		case "*":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Integer(left * right)
		case "+":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Integer(left + right)
		case "-":
			left := Eval(v.Left).(Integer)
			right := Eval(v.Right).(Integer)
			return Integer(left - right)
		default:
			panic(fmt.Sprintf("Unknown binop: %s", v.Op))
		}
	case Unop:
		switch v.Op {
		case "-":
			arg := Eval(v.Arg).(Integer)
			return Integer(-arg)
		case "!":
			arg := Eval(v.Arg).(Boolean)
			return Boolean(!arg)
		case "$":
			i := Eval(v.Arg).(Integer)
			s := ""
			for i != 0 {
				d := i % 94
				i /= 94
				s = string(lookup[d]) + s
			}
			return String(s)
		case "#":
			s := Eval(v.Arg).(String)
			i := 0
			for _, c := range s {
				i = i*94 + int(strings.Index(lookup, string(c)))
			}
			return Integer(i)
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
