package views

import (
	"time"

	rice "github.com/GeertJohan/go.rice"
)

var View *rice.Box

func LoadViews() {
	View = rice.MustFindBox("../views")
}

// Date 时间戳格式化日期
func Date(timestamp int64, format ...string) string {
	layout := "2006-01-02 15:04:05"

	if len(format) > 0 {
		layout = format[0]
	}

	date := time.Unix(timestamp, 0).Format(layout)

	return date
}
