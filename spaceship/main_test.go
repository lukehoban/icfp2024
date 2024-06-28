package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveRel(t *testing.T) {
	// a := moveRel(Point{1, 0}, Point{0, 0})
	// assert.Equal(t, 7, a)

	a := moveRel(Point{0, 0}, Point{1, -1})
	assert.Equal(t, 1, a)
}
