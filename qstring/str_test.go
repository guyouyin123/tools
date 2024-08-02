package qstring

import (
	"fmt"
	"testing"
)

func TestContainsChinese(t *testing.T) {
	fmt.Println(ContainsChinese("aaa"))
	fmt.Println(ContainsChinese("aaa你好"))
}
