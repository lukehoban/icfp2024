package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstMessage(t *testing.T) {
	s := "solve language_test 4w3s0m3"
	tok := StringToToken(s)
	fmt.Printf("Sending: %s\n", s)
	ret, err := Communicate(string(tok))
	assert.NoError(t, err)

	expr, rest := CombineToExpr(ret)
	fmt.Printf("Expr: %v\n", expr)
	fmt.Printf("Rest: %v\n", rest)

	assert.Empty(t, rest)

	res := Eval(expr, map[int]Expr{})
	fmt.Printf("Res: %v\n", res)
}
