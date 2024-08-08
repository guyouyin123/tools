package qmixCompute

import (
	"fmt"
	"testing"
)

func TestMixCompute(t *testing.T) {
	result := MixCompute("a*b+c", map[rune]float64{
		'a': float64(100),
		'b': 0.5,
		'c': float64(10),
	})
	fmt.Println(result)

	result2 := MixCompute("a/b", map[rune]float64{
		'a': float64(10),
		'b': 3,
	})
	fmt.Println(result2)
}
