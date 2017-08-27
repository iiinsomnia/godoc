package controllers

import (
	"database/sql"
	"fmt"
	"godoc/i18n"
	"godoc/params"
	"godoc/rbac"
	"godoc/service"
	"html/template"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iiinsomnia/yiigo"
)

type DocController struct {
	*controller
}

type DocForm struct {
	Title    string `form:"title" binding:"required"`
	Label    string `form:"label"`
	Markdown string `form:"markdown" binding:"required"`
}

func NewDocController(r *gin.Engine) *DocController {
	return &DocController{
		construct(r, "main", "doc"),
	}
}

func (d *DocController) View(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	docService := service.NewDocService(c)
	doc, err := docService.GetDetail(_id)

	if err != nil {
		if err == sql.ErrNoRows {
			d.V(c).RenderErr(404, "项目不存在")
			return
		}

		d.V(c).RenderErr(500, "数据获取失败")

		return
	}

	docs, _ := docService.GetDocs(doc.ProjectID)

	historyService := service.NewHistoryService(c)
	history, _ := historyService.GetHistory(_id)

	d.V(c).F(template.FuncMap{
		"explode":        strings.Split,
		"getHistoryFlag": params.GetHistoryFlag,
	}).Render("view", gin.H{
		"doc":     doc,
		"docs":    docs,
		"history": history,
	})
}

func (d *DocController) Add(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role == 1 {
		d.V(c).RenderErr(403, "无操作权限")
		return
	}

	projectID := c.Param("project")
	_projectID, _ := strconv.Atoi(projectID)

	projectService := service.NewProjectService(c)
	project, err := projectService.GetDetail(_projectID)

	if err != nil {
		if err == sql.ErrNoRows {
			d.V(c).RenderErr(404, "项目不存在")
			return
		}

		d.V(c).RenderErr(500, "数据获取失败")

		return
	}

	if c.Request.Method == "GET" {
		docService := service.NewDocService(c)
		docs, _ := docService.GetDocs(_projectID)

		d.V(c).Render("add", gin.H{
			"project": project,
			"docs":    docs,
		})

		return
	}

	form := &DocForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		d.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"title":       form.Title,
		"category_id": project.CategoryID,
		"project_id":  project.ID,
		"label":       form.Label,
		"markdown":    form.Markdown,
	}

	history := yiigo.X{
		"category_id": project.CategoryID,
		"project_id":  project.ID,
	}

	docService := service.NewDocService(c)
	id, err := docService.Add(data, history)

	if err != nil {
		d.JSON(c, false, "添加失败")
		return
	}

	d.JSON(c, true, "添加成功", nil, fmt.Sprintf("/docs/view/%d", id))
}

func (d *DocController) Edit(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role == 1 {
		d.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	docService := service.NewDocService(c)

	if c.Request.Method == "GET" {
		doc, err := docService.GetDetail(_id)

		if err != nil {
			if err == sql.ErrNoRows {
				d.V(c).RenderErr(404, "文档不存在")
				return
			}

			d.V(c).RenderErr(500, "数据获取失败")

			return
		}

		docs, _ := docService.GetDocs(doc.ProjectID)

		d.V(c).Render("edit", gin.H{
			"doc":  doc,
			"docs": docs,
		})

		return
	}

	form := &DocForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		d.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"title":    form.Title,
		"label":    form.Label,
		"markdown": form.Markdown,
	}

	err := docService.Edit(_id, data)

	if err != nil {
		d.JSON(c, false, "编辑失败")
		return
	}

	d.JSON(c, true, "编辑成功", nil, fmt.Sprintf("/docs/view/%s", id))
}

func (d *DocController) Delete(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		d.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	docService := service.NewDocService(c)
	doc, err := docService.GetDetail(_id)

	if err != nil {
		if err == sql.ErrNoRows {
			d.JSON(c, false, "文档不存在")
			return
		}

		d.JSON(c, false, "删除失败")

		return
	}

	err = docService.Delete(_id)

	if err != nil {
		d.JSON(c, false, "删除失败")
		return
	}

	d.JSON(c, true, "删除成功", nil, fmt.Sprintf("/projects/view/%d", doc.ProjectID))
}
