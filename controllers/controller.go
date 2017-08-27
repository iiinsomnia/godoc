package controllers

import (
	"godoc/helpers"
	"godoc/rbac"
	"godoc/views"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type controller struct {
	router *gin.Engine
	tpl    *tpl
}

type tpl struct {
	layout string
	dir    string
}

func construct(r *gin.Engine, layout string, dir string) *controller {
	return &controller{
		router: r,
		tpl: &tpl{
			layout: layout,
			dir:    dir,
		},
	}
}

func (c *controller) V(ctx *gin.Context) *views.View {
	v := &views.View{
		Router: c.router,
		Ctx:    ctx,
		HTML: &views.HTML{
			Layout: c.tpl.layout,
			Funcs: template.FuncMap{
				"getRoleName": rbac.GetRoleName,
				"date":        helpers.Date,
			},
			Err: map[string]string{
				"layout": "normal",
				"dir":    "error",
				"tpl":    "error",
			},
			Dir: c.tpl.dir,
			Ext: "html",
		},
	}

	return v
}

func (c *controller) Redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusFound, location)
}

func (c *controller) JSON(ctx *gin.Context, success bool, msg interface{}, resp ...interface{}) {
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
