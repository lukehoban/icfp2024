package main

import (
	"fmt"
	"io"
	"os"

	"github.com/lukehoban/icfp2024/icfp"
)

func main() {
	s := ""
	if len(os.Args) >= 2 {
		s = os.Args[1]
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		s = string(b)
	}
	tok := icfp.StringToToken(s)
	ret, err := icfp.CommunicateToken(string(tok))
	if err != nil {
		panic(err)
	}

	expr, rest := icfp.CombineToExpr(ret)
	if len(rest) > 0 {
		fmt.Printf("WARNING - didn't use all input! %v\n", rest)
	}

	fmt.Printf("expr %v\n", expr)
	res := icfp.Eval(expr)
	fmt.Printf("%v", res)
}
