package helpers

import (
	"crypto/md5"
	"fmt"
)

// MD5 获取字符串md5值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))

	return fmt.Sprintf("%x", h.Sum(nil))
}
