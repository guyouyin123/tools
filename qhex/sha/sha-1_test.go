package sha

import (
	"fmt"
	"testing"
)

func TestSHA1(t *testing.T) {
	a := SHA1("123456", "yan")
	fmt.Println(a)
}
func TestSHA256(t *testing.T) {
	a := SHA256("123456", "yan")
	fmt.Println(a)
}
