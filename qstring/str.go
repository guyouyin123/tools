package qstring

import "unicode"

// 判断字符串是否包含中文字符
func ContainsChinese(str string) bool {
	for _, char := range str {
		if unicode.Is(unicode.Scripts["Han"], char) {
			return true
		}
	}
	return false
}
