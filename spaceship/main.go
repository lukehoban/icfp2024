package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lukehoban/icfp2024/icfp"
)

func do() error {
	var wg sync.WaitGroup

	for i := 1; i <= 25; i++ {
		time.Sleep(100 * time.Second)
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// if i < 16 {
			// 	return
			// }
			in, err := icfp.CommunicateString(fmt.Sprintf("get spaceship%d", i))
			if err != nil {
				fmt.Printf("Failed to get spaceship%d: %v\n", i, err)
				return
			}
			// fmt.Printf("Got back %s\n", string(in))
			points, err := parse(string(in))
			if err != nil {
				fmt.Printf("Invalid response for spaceship%d: %v\n", i, err)
				return
			}
			actions := Walk(points)
			s := ""
			for _, a := range actions {
				s += strconv.Itoa(a)
			}
			fmt.Printf("Spaceship%d: %s\n", i, s)

			answer, err := icfp.CommunicateString(fmt.Sprintf("solve spaceship%d %s", i, s))
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

type Point struct {
	x, y int
}

func moveRel(delta Point, vel Point) (int, Point, int) {
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
		if yrem > maxrem {
			maxrem = yrem
		}
		if maxrem < bestrem {
			bestrem = maxrem
			best = i
		}
	}
	// Start at 1 on keypad
	return best + 1, Point{x: vel.x + candidates[best].x, y: vel.y + candidates[best].y}, bestrem
}

func Walk(points []Point) []int {
	pos := Point{0, 0}
	vel := Point{0, 0}
	var ret []int
	for _, point := range points {
		for {
			a, newv, _ := moveRel(Point{x: point.x - pos.x, y: point.y - pos.y}, vel)
			vel = newv
			ret = append(ret, a)
			pos.x += vel.x
			pos.y += vel.y
			if pos.x == point.x && pos.y == point.y {
				break
			}
		}
	}
	return ret
}

func parse(s string) ([]Point, error) {
	lines := strings.Split(s, "\n")
	var points []Point
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if line == "" {
			continue
		}
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
		points = append(points, Point{a, b})
	}
	return points, nil
}

func readFile(s string) ([]Point, error) {
	byts, err := os.ReadFile(s)
	if err != nil {
		return nil, err
	}
	return parse(string(byts))

}

func main() {
	err := do()
	if err != nil {
		panic(err)
	}
}
