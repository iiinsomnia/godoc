package helpers

import "strconv"

// Int 字符串转int
func Int(s string) int {
	n, _ := strconv.Atoi(s)

	return n
}
