package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveRel(t *testing.T) {

	a, newv, dist := moveRel(Point{1, -1}, Point{0, 0})
	assert.Equal(t, 3, a)
	assert.Equal(t, 0, dist)
	assert.Equal(t, Point{1, -1}, newv)

	a, newv, dist = moveRel(Point{0, -2}, Point{1, -1})
	assert.Equal(t, 1, a)
	assert.Equal(t, 0, dist)
	assert.Equal(t, Point{0, -2}, newv)

	a, newv, dist = moveRel(Point{1, -2}, Point{0, -2})
	assert.Equal(t, 6, a)
	assert.Equal(t, 0, dist)
	assert.Equal(t, Point{1, -2}, newv)

	a, newv, dist = moveRel(Point{0, -3}, Point{1, -2})
	assert.Equal(t, 1, a)
	assert.Equal(t, 0, dist)
	assert.Equal(t, Point{0, -3}, newv)

	a, newv, dist = moveRel(Point{1, -2}, Point{0, -3})
	assert.Equal(t, 9, a)
	assert.Equal(t, 0, dist)
	assert.Equal(t, Point{1, -2}, newv)
}

func Test1(t *testing.T) {
	points, err := readFile("1.txt")
	assert.NoError(t, err)

	actions := Walk(points)
	assert.Equal(t, []int{3, 1, 6, 1, 9}, actions)

}
