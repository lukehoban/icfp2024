package main

import (
	"fmt"
	"strconv"
)

/*
. . . . 0 . . . .
. B > . = . . . .
. v 1 . . > . . .
. . - . . . + S .
. . . . . ^ . . .
. . v . . 0 > . .
. . . . . . A + .
. 1 @ 6 . . < . .
. . 3 . 0 @ 3 . .
. . . . . 3 . . .
*/

func Parse() [][]string {
	return [][]string{
		{".", ".", ".", ".", "0", ".", ".", ".", "."},
		{".", "B", ">", ".", "=", ".", ".", ".", "."},
		{".", "v", "1", ".", ".", ">", ".", ".", "."},
		{".", ".", "-", ".", ".", ".", "+", "S", "."},
		{".", ".", ".", ".", ".", "^", ".", ".", "."},
		{".", ".", "v", ".", ".", "0", ">", ".", "."},
		{".", ".", ".", ".", ".", ".", "A", "+", "."},
		{".", "1", "@", "6", ".", ".", "<", ".", "."},
		{".", ".", "3", ".", "0", "@", "3", ".", "."},
		{".", ".", ".", ".", ".", "3", ".", ".", "."},
	}
}

type Val interface{ isVal() }
type Int int
type Op string

func (Int) isVal() {}
func (Op) isVal()  {}

func PrintBoard(board [][]Val) string {
	s := ""
	for _, row := range board {
		for _, cell := range row {
			if cell == nil {
				s += "."
			} else if i, ok := cell.(Int); ok {
				s += fmt.Sprintf("%d", i)
			} else {
				s += string(cell.(Op))
			}
			s += " "
		}
		s += "\n"
	}
	return s
}

func Step(in [][]Val) ([][]Val, Val) {
	out := make([][]Val, len(in))
	for y := range in {
		out[y] = make([]Val, len(in[y]))
		copy(out[y], in[y])
	}
	submitLocations := []struct{ x, y int }{}
	for y, row := range in {
		for x, cell := range row {
			switch c := cell.(type) {
			case Int:
				continue
			case Op:
				switch c {
				case "S":
					submitLocations = append(submitLocations, struct{ x, y int }{x, y})
				case ">":
					if x == 0 || in[y][x-1] == nil {
						continue
					}
					out[y][x+1] = in[y][x-1]
					out[y][x-1] = nil
				case "<":
					if x == len(in[y])-1 || in[y][x+1] == nil {
						continue
					}
					out[y][x-1] = in[y][x+1]
					out[y][x+1] = nil
				case "^":
					if y == len(in)-1 || in[y+1][x] == nil {
						continue
					}
					out[y-1][x] = in[y+1][x]
					out[y+1][x] = nil
				case "v":
					if y == 0 || in[y-1][x] == nil {
						continue
					}
					out[y+1][x] = in[y-1][x]
					out[y-1][x] = nil
				case "+", "-", "*", "/", "%":
					if x == 0 || in[y][x-1] == nil || y == 0 || in[y-1][x] == nil {
						continue
					}
					a, ok := in[y][x-1].(Int)
					if !ok {
						continue
					}
					b, ok := in[y-1][x].(Int)
					if !ok {
						continue
					}
					res := Int(0)
					switch c {
					case "+":
						res = a + b
					case "-":
						res = a - b
					case "*":
						res = a * b
					case "/":
						res = a / b
					case "%":
						res = a % b
					}
					out[y][x+1] = res
					out[y+1][x] = res
				case "=", "#":
					if x == 0 || in[y][x-1] == nil || y == 0 || in[y-1][x] == nil {
						continue
					}
					a := in[y][x-1]
					b := in[y-1][x]
					if a == nil || b == nil {
						continue
					}
					test := false
					switch c {
					case "=":
						test = a == b
					case "#":
						test = a != b
					default:
						panic("unreachable")
					}
					if test {
						out[y][x+1] = b
						out[y+1][x] = a
					}
				}
			default:
				if c == nil {
					continue
				}
				panic(fmt.Sprintf("NYI %#v", cell))
			}
		}
	}
	// If "S" was overwritten, submit it
	for _, loc := range submitLocations {
		if out[loc.y][loc.x] != Op("S") {
			return nil, out[loc.y][loc.x]
		}
	}
	return out, nil
}

func Run(a, b int) Val {
	program := Parse()
	m := [][]Val{}
	t := 1
	for _, row := range program {
		r := []Val{}
		for _, cell := range row {
			if cell == "." {
				r = append(r, nil)
			} else if i, err := strconv.Atoi(cell); err == nil {
				r = append(r, Int(i))
			} else if cell == "A" {
				r = append(r, Int(a))
			} else if cell == "B" {
				r = append(r, Int(b))
			} else {
				r = append(r, Op(cell))
			}
		}
		m = append(m, r)
	}
	fmt.Printf("[t=%d]\n%s\n", t, PrintBoard(m))
	for {
		newm, ret := Step(m)
		t++
		if ret != nil {
			return ret
		}
		m = newm
		fmt.Printf("[t=%d]\n%s\n", t, PrintBoard(m))
	}
}

func do() error {
	Run(3, 4)
	return nil
}

func main() {
	err := do()
	if err != nil {
		panic(err)
	}
}
