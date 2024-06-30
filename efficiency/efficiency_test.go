package efficiency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEfficiency3(t *testing.T) {
	// 9345873499
	// (+ 2134 (* (((λy.((λz.(y (z z))) (λz.(y (z z))))) (λw.(λa.(if (= a 0) 1 (+ 1 (w (- a 1))))))) 9345873499) 1))
	var f func(int) int
	f = func(a int) int {
		if a == 0 {
			return 1
		}
		return 1 + f(a-1)
	}
	res := f(9345873499) + 1
	assert.Equal(t, res, 1)

}
