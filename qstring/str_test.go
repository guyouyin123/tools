package qstring

import (
	"fmt"
	"testing"
)

func TestContainsChinese(t *testing.T) {
	fmt.Println(ContainsChinese("aaa"))
	fmt.Println(ContainsChinese("aaa你好"))

	a := RsplitN("aa/bb/cc.txt", "/", 2)
	fmt.Println(a[1])
}
