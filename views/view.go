package views

import (
	"fmt"
	"godoc/rbac"
	"html/template"
	"net/http"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

var viewBox *rice.Box

type View struct {
	Router *gin.Engine
	Ctx    *gin.Context
	HTML   *HTML
}

type HTML struct {
	Layout string
	Funcs  template.FuncMap
	Err    map[string]string
	Extra  []string
	Dir    string
	Ext    string
}

// LoadViews init load view box
func LoadViews() {
	viewBox = rice.MustFindBox("../views")
}

// L set view layout
func (v *View) L(layout string) *View {
	v.HTML.Layout = layout

	return v
}

// F add view func
func (v *View) F(funcs template.FuncMap) *View {
	for k, f := range funcs {
		v.HTML.Funcs[k] = f
	}

	return v
}

// T add view tpl, eg: "layouts/pagination"
func (v *View) T(tpls ...string) *View {
	extra := []string{}

	for _, t := range tpls {
		tpl := fmt.Sprintf("%s.%s", t, v.HTML.Ext)
		extra = append(extra, tpl)
	}

	v.HTML.Extra = extra

	return v
}

// Render render a view
func (v *View) Render(tpl string, args ...gin.H) {
	defer v.recover()

	tpls := []string{}

	// 布局模板
	layout := fmt.Sprintf("layouts/%s.%s", v.HTML.Layout, v.HTML.Ext)
	tpls = append(tpls, layout)

	// 额外模板
	tpls = append(tpls, v.HTML.Extra...)

	// 当前模板
	curTpl := fmt.Sprintf("%s/%s.%s", v.HTML.Dir, tpl, v.HTML.Ext)
	tpls = append(tpls, curTpl)

	tplStrs := []string{}

	for _, t := range tpls {
		tplStrs = append(tplStrs, viewBox.MustString(t))
	}

	html := strings.Join(tplStrs, "")
	templ := template.Must(template.New("").Funcs(v.HTML.Funcs).Parse(html))

	v.Router.SetHTMLTemplate(templ)

	data := gin.H{}

	if len(args) > 0 {
		data = args[0]
	}

	data["version"] = yiigo.EnvString("app", "version", "1.0.0")
	data["identity"] = rbac.GetIdentity(v.Ctx)

	v.Ctx.HTML(http.StatusOK, tpl, data)
}

// RenderErr render a error view
func (v *View) RenderErr(code int, msg string) {
	tpls := []string{
		fmt.Sprintf("layouts/%s.%s", v.HTML.Err["layout"], v.HTML.Ext),
		fmt.Sprintf("%s/%s.%s", v.HTML.Err["dir"], v.HTML.Err["tpl"], v.HTML.Ext),
	}

	tplStrs := []string{}

	for _, t := range tpls {
		tplStrs = append(tplStrs, viewBox.MustString(t))
	}

	html := strings.Join(tplStrs, "")
	templ := template.Must(template.New("").Parse(html))

	v.Router.SetHTMLTemplate(templ)

	v.Ctx.HTML(http.StatusOK, v.HTML.Err["tpl"], gin.H{
		"code": code,
		"msg":  msg,
	})
}

func (v *View) recover() {
	if err := recover(); err != nil {
		yiigo.Error(err)
		v.RenderErr(500, "内部服务器错误")
	}
}
