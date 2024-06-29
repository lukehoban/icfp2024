package main

import (
	"fmt"
	"os"

	"github.com/lukehoban/icfp2024/icfp"
)

func main() {
	s := os.Args[1]
	tok := icfp.StringToToken(s)
	ret, err := icfp.CommunicateToken(string(tok))
	if err != nil {
		panic(err)
	}

	expr, rest := icfp.CombineToExpr(ret)
	if len(rest) > 0 {
		fmt.Printf("WARNING - didn't use all input! %v\n", rest)
	}

	res := icfp.Eval(expr)
	fmt.Printf("%v", res)
}
