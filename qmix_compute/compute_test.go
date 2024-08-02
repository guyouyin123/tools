package qmix_compute

import (
	"fmt"
	"testing"
)

func TestMixCompute(t *testing.T) {
	result := MixCompute("(a-b)/b", map[rune]float64{
		'a': float64(100),
		'b': float64(200),
	})
	fmt.Println(result)
}
