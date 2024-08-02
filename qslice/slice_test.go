package qslice

import (
	"fmt"
	"testing"
)

func TestSliceRemoveDuplicates(t *testing.T) {
	//a := []int{1, 1, 1, 2, 2, 2}
	//b := []int{10}

	a := []string{"a", "b"}
	b := []string{"c", "d"}

	fmt.Println(DiffStrSlice(a, b))
}
