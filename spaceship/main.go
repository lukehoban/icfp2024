package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*

1 -1
1 -3
2 -5
2 -8
3 -10

*/

func do() error {
	points, err := readFile(os.Args[1])
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", points)
	return nil
}

type Point struct {
	x, y int
}

func moveRel(delta Point, vel Point) (int, Point) {
	candidates := []Point{
		{-1, -1},
		{0, -1},
		{1, -1},
		{-1, 0},
		{0, 0},
		{+1, 0},
		{-1, 1},
		{0, 1},
		{1, 1},
	}
	best := 0
	bestrem := 1000000000
	for i, step := range candidates {
		xrem := vel.x + step.x - delta.x
		if xrem < 0 {
			xrem = -xrem
		}
		yrem := vel.y + step.y - delta.y
		if yrem < 0 {
			yrem = -yrem
		}
		maxrem := xrem
		if yrem < maxrem {
			maxrem = yrem
		}
		if maxrem < bestrem {
			bestrem = maxrem
			best = i
		}
	}
	// Start at 1 on keypad
	return best + 1, Point{x: vel.x + candidates[best].x, y: vel.y + candidates[best].y}
}

// func moveRelative(delta Point) []Point {
// 	// assume we need start at velocity 0,0 and need to end at velocity 0,0
// 	max := delta.x
// 	if delta.y > max {
// 		max = delta.y
// 	}
// 	keypad := [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
// 	for i := 0; i < max; i++ {
// 		var action int
// 		var keyPadRow [3]int
// 		if delta.y > 0 {
// 			keyPadRow = keypad[2]
// 		} else if delta.y < 0 {
// 			keyPadRow = keypad[0]
// 		} else {
// 			keyPadRow = keypad[1]
// 		}
// 		if delta.x > 0 {
// 			action = keyPadRow[2]
// 		} else if delta.x < 0 {
// 			action = keyPadRow[0]
// 		} else {
// 			action = keyPadRow[1]
// 		}
// 		fmt.Printf("%d\n", action)

// 	}

// }

func readFile(s string) ([][2]int, error) {
	byts, err := os.ReadFile(s)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(byts), "\n")
	var points [][2]int
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		b, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		points = append(points, [2]int{a, b})
	}
	return points, nil
}

func main() {
	err := do()
	if err != nil {
		panic(err)
	}
}
