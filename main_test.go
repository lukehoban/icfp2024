package main

import (
	"fmt"
	"testing"

	"github.com/lukehoban/icfp2024/icfp"
	"github.com/stretchr/testify/assert"
)

func TestFirstMessage(t *testing.T) {
	s := "solve language_test 4w3s0m3"
	tok := icfp.StringToToken(s)
	fmt.Printf("Sending: %s\n", s)
	ret, err := icfp.CommunicateToken(string(tok))
	assert.NoError(t, err)

	expr, rest := icfp.CombineToExpr(ret)
	fmt.Printf("Expr: %v\n", expr)
	fmt.Printf("Rest: %v\n", rest)

	assert.Empty(t, rest)

	res := icfp.Eval(expr, map[int]icfp.Expr{})
	fmt.Printf("Res: %v\n", res)
}
