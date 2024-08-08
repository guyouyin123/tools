package qstring

import (
	"strings"
	"unicode"
)

// 判断字符串是否包含中文字符
func ContainsChinese(str string) bool {
	for _, char := range str {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}

// 右切
func RsplitN(s, sep string, n int) []string {
	st := reverseString(s)
	st2 := strings.SplitN(st, sep, n)
	for i := 0; i < len(st2)/2; i++ {
		st2[i], st2[len(st2)-i-1] = st2[len(st2)-i-1], st2[i]
	}
	for k, v := range st2 {
		st2[k] = reverseString(v)
	}
	return st2
}

func reverseString(str string) string {
	st := ""
	for i := len(str) - 1; i >= 0; i-- {
		st += string(str[i])
	}
	return st
}
