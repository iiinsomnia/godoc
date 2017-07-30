package helpers

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// MD5 获取字符串md5值
func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}

// Int 字符串转int
func Int(s string) int {
	n, _ := strconv.Atoi(s)

	return n
}
