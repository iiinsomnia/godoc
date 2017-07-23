package controllers

import (
	"fmt"
	"godoc/params"
	"godoc/rbac"
	"godoc/views"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iiinsomnia/yiigo"
)

type controller struct {
	router  *gin.Engine
	view    *view
	version string
}

type view struct {
	deft *html
	curr *html
	err  map[string]string
	dir  string
	ext  string
}

type html struct {
	funcs  template.FuncMap
	layout string
}

// 构造函数
func construct(r *gin.Engine, layout string, tplDir string) *controller {
	version := yiigo.GetEnvString("app", "version", "1.0.0")

	return &controller{
		router: r,
		view: &view{
			deft: &html{
				funcs: template.FuncMap{
					"getRoleName": params.GetRoleName,
					"date":        views.Date,
				},
				layout: layout,
			},
			curr: &html{
				funcs: template.FuncMap{
					"getRoleName": params.GetRoleName,
					"date":        views.Date,
				},
				layout: layout,
			},
			err: map[string]string{
				"layout": "normal",
				"dir":    "error",
				"tpl":    "error",
			},
			dir: tplDir,
			ext: "html",
		},
		version: version,
	}
}

func (c *controller) addFuncs(funcs template.FuncMap) *controller {
	for k, v := range funcs {
		c.view.curr.funcs[k] = v
	}

	return c
}

func (c *controller) layout(layout string) *controller {
	c.view.curr.layout = layout

	return c
}

func (c *controller) render(ctx *gin.Context, tpl string, args ...gin.H) {
	defer c.recover(ctx)

	viewFiles := []string{}

	layoutFile := fmt.Sprintf("layouts/%s.%s", c.view.curr.layout, c.view.ext)
	tplFile := fmt.Sprintf("%s/%s.%s", c.view.dir, tpl, c.view.ext)

	viewFiles = append(viewFiles, layoutFile, tplFile)

	viewStrings := []string{}

	for _, v := range viewFiles {
		viewString := views.View.MustString(v)
		viewStrings = append(viewStrings, viewString)
	}

	html := strings.Join(viewStrings, "")

	templ := template.Must(template.New("").Funcs(c.view.curr.funcs).Parse(html))

	c.router.SetHTMLTemplate(templ)

	data := gin.H{}

	if len(args) > 0 {
		data = args[0]
	}

	data["version"] = c.version
	data["identity"] = rbac.GetIdentity(ctx)

	ctx.HTML(http.StatusOK, tpl, data)
}

func (c *controller) renderError(ctx *gin.Context, code int, msg string) {
	viewFiles := []string{
		fmt.Sprintf("layouts/%s.%s", c.view.err["layout"], c.view.ext),
		fmt.Sprintf("%s/%s.%s", c.view.err["dir"], c.view.err["tpl"], c.view.ext),
	}

	viewStrings := []string{}

	for _, v := range viewFiles {
		viewString := views.View.MustString(v)
		viewStrings = append(viewStrings, viewString)
	}

	html := strings.Join(viewStrings, "")

	templ := template.Must(template.New("").Parse(html))

	c.router.SetHTMLTemplate(templ)

	ctx.HTML(http.StatusOK, c.view.err["tpl"], gin.H{
		"code": code,
		"msg":  msg,
	})
}

func (c *controller) redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusFound, location)
}

func (c *controller) json(ctx *gin.Context, success bool, msg interface{}, resp ...interface{}) {
	obj := gin.H{
		"success": success,
		"msg":     msg,
	}

	if num := len(resp); num > 0 {
		switch num {
		case 1:
			obj["data"] = resp[0]
		case 2:
			obj["data"] = resp[0]
			obj["redirect"] = resp[1]
		}
	}

	ctx.JSON(http.StatusOK, obj)
}

func (c *controller) recover(ctx *gin.Context) {
	c.view.curr.funcs = c.view.deft.funcs
	c.view.curr.layout = c.view.deft.layout

	if err := recover(); err != nil {
		yiigo.LogErrorf("%s", err)
		ctx.String(http.StatusInternalServerError, "%s", err)
	}
}
