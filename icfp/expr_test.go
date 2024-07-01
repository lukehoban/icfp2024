package icfp

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func evalString(t *testing.T, s string) Expr {
	exprs := Parse(s)
	expr, rest := CombineToExpr(exprs)
	assert.Empty(t, rest)
	return Eval(expr, nil)
}

func TestEfficiency1(t *testing.T) {
	s := `B$ L! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! I" L! B+ B+ v! v! B+ v! v!`
	v := evalString(t, s)
	assert.Equal(t, Integer{Int: big.NewInt(17592186044416)}, v)
}

func TestEfficiency2(t *testing.T) {
	s := `B+ I7c B* B$ B$ L" B$ L# B$ v" B$ v# v# L# B$ v" B$ v# v# L$ L% ? B= v% I! I" B+ I" B$ v$ B- v% I" I":c1+0 I!`
	v := evalString(t, s)
	assert.Equal(t, Integer{Int: big.NewInt(17592186044416)}, v)
}
