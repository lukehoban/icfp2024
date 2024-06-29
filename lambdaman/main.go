package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lukehoban/icfp2024/icfp"
)

const example = `###.#...
...L..##
.#######`

type Grid [][]byte

func parseGrid(s string) Grid {
	lines := strings.Split(s, "\n")
	g := make(Grid, len(lines))
	for i, line := range lines {
		g[i] = []byte(line)
	}
	return g
}

func FindPath(x0, y0, x1, y1 int, grid Grid) []byte {
	moves := []byte{}
	for {
		options := []byte{}
		if x0 == x1 && y0 == y1 {
			return moves
		}
		if x1 > x0 {
			options = append(options, 'R')
		}
		if x1 < x0 {
			options = append(options, 'L')
		}
		if y1 > y0 {
			options = append(options, 'D')
		}
		if y1 < y0 {
			options = append(options, 'U')
		}
		for _, option := range options {
			switch option {
			case 'L':
				if grid[y0][x0-1] != '#' {
					moves = append(moves, option)
					x0--
				}
			case 'R':
				if grid[y0][x0+1] != '#' {
					moves = append(moves, option)
					x0++
				}
			case 'U':
				if grid[y0-1][x0] != '#' {
					moves = append(moves, option)
					y0--
				}
			case 'D':
				if grid[y0+1][x0] != '#' {
					moves = append(moves, option)
					y0++
				}
			}
		}
	}
}

func run(example string) (string, error) {
	grid := parseGrid(example)
	fmt.Printf("%v\n", grid)
	width := len(grid[0])
	dots := map[int]struct{}{}
	x := 0
	y := 0
	for i, row := range grid {
		for j, cell := range row {
			if cell == '.' {
				dots[i*width+j] = struct{}{}
			}
			if cell == 'L' {
				y = i
				x = j
			}
		}
	}
	fmt.Printf("%v\n", dots)
	fmt.Printf("lambdaman is at (%d, %d)\n", x, y)
	var ret []byte
	for {
		closest := 100000000000
		closestX := 0
		closestY := 0
		closestIndex := 0
		if len(dots) == 0 {
			break
		}
		for dot := range dots {
			dotX := dot % width
			dotY := dot / width
			dist := (dotX-x)*(dotX-x) + (dotY-y)*(dotY-y)
			if dist < closest {
				closest = dist
				closestX = dotX
				closestY = dotY
				closestIndex = dot
			}
		}
		fmt.Printf("closest dot is at (%d, %d)\n", closestX, closestY)
		delete(dots, closestIndex)

		moves := FindPath(x, y, closestX, closestY, grid)
		fmt.Printf("moves: %s\n", string(moves))
		x = closestX
		y = closestY
		ret = append(ret, moves...)
	}

	fmt.Printf("moves: %s\n", string(ret))
	return string(ret), nil
}

func do() error {
	var wg sync.WaitGroup

	for i := 1; i <= 25; i++ {
		time.Sleep(10 * time.Second)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// if i < 16 {
			// 	return
			// }
			fmt.Printf("Getting lambdaman%d\n", i)
			in, err := icfp.CommunicateString(fmt.Sprintf("get lambdaman%d", i))
			if err != nil {
				fmt.Printf("Failed to get spaceship%d: %v\n", i, err)
				return
			}
			fmt.Printf("Got back %d\n", len(string(in)))
			s, err := run(string(in))
			if err != nil {
				fmt.Printf("Invalid response for lambdaman%d: %v\n", i, err)
				return
			}

			fmt.Printf("lambdaman%d: %s\n", i, s)

			answer, err := icfp.CommunicateString(fmt.Sprintf("solve lambdaman%d %s", i, s))
			if err != nil {
				fmt.Printf("Failed to get spaceship%d: %v\n", i, err)
				return
			}
			fmt.Printf("Response: %s\n", answer)
		}(i)
	}

	wg.Wait()

	return nil
}

func main() {
	err := do()
	if err != nil {
		panic(err)
	}
}
