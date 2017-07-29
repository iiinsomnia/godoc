package views

import rice "github.com/GeertJohan/go.rice"

var View *rice.Box

func LoadViews() {
	View = rice.MustFindBox("../views")
}
