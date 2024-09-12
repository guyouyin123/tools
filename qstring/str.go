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

// RsplitN 切割字符串，从右边开始切割指定次数，并保留切割字符
func RsplitN(s string, delimiter string, cuts int) []string {
	if cuts == 0 {
		return []string{s}
	}
	// 从右向左切割字符串
	parts := strings.Split(s, delimiter)
	lenth := len(parts)
	// 如果切割次数超过了部分数量，限制为部分数量
	if cuts > lenth {
		cuts = len(parts)
	}
	tmp := lenth - cuts
	// 反转切割后的部分以从右边开始取
	result := make([]string, 0)
	st := ""
	for k, v := range parts {
		if k < tmp {
			st += v + delimiter
		} else {
			if st != "" {
				result = append(result, st)
				st = ""
			}
			result = append(result, v)
		}
	}
	return result
}
