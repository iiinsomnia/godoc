package assets

import rice "github.com/GeertJohan/go.rice"

var Asset *rice.Box

func LoadAssets() {
	Asset = rice.MustFindBox("../assets")
}
