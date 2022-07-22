package routers

import (
	"strings"
)

// 索引加一函数
func Add(a, b int) int {
	return a + b
}

// 字符串截取函数
func StrCut(s string, length int, dot ...string) string {
	if len(s) < length {
		return s
	}
	if len(dot) <= 0 {
		dot = []string{"..."}
	}

	n := 0
	l := len(s)
	for i, b := range s {
		if n >= length*2 {
			if l > i {
				return s[:i] + strings.Join(dot, "")
			}
			return s[:i]
		}
		if b < 128 {
			n++
		} else {
			n += 2
		}
	}

	return s
}
