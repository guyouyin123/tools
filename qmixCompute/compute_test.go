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
	a := fmt.Sprintf("%.2f", float64(10))
	b := fmt.Sprintf("%.2f", 10.000001)
	fmt.Println(a)
	fmt.Println(b)
}
