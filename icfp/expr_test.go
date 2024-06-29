package icfp

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEfficiency1(t *testing.T) {
	s := `B$ L! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! B$ v! I" L! B+ B+ v! v! B+ v! v!`
	exprs := Parse(s)
	expr, rest := CombineToExpr(exprs)
	assert.Empty(t, rest)
	res := Eval(expr)
	assert.Equal(t, Integer{Int: big.NewInt(17592186044416)}, res)
}
