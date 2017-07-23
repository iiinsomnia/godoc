package controllers

import (
	"fmt"
	"godoc/i18n"
	"godoc/service"
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
		if err.Error() == "not found" {
			d.renderError(c, 404, "项目不存在")
			return
		}

		d.renderError(c, 500, "数据获取失败")

		return
	}

	docs, _ := docService.GetDocs(doc.ProjectID)

	d.render(c, "view", gin.H{
		"doc":  doc,
		"docs": docs,
	})
}

func (d *DocController) Add(c *gin.Context) {
	projectID := c.Param("project")
	_projectID, _ := strconv.Atoi(projectID)

	projectService := service.NewProjectService(c)
	project, err := projectService.GetDetail(_projectID)

	if err != nil {
		if err.Error() == "not found" {
			d.renderError(c, 404, "项目不存在")
			return
		}

		d.renderError(c, 500, "数据获取失败")

		return
	}

	if c.Request.Method == "GET" {
		docService := service.NewDocService(c)
		docs, _ := docService.GetDocs(_projectID)

		d.render(c, "add", gin.H{
			"project": project,
			"docs":    docs,
		})

		return
	}

	form := &DocForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		d.json(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"title":       form.Title,
		"category_id": project.CategoryID,
		"project_id":  project.ID,
		"label":       form.Label,
		"markdown":    form.Markdown,
	}

	docService := service.NewDocService(c)
	id, err := docService.Add(data)

	if err != nil {
		d.json(c, false, "添加失败")
		return
	}

	d.json(c, true, "添加成功", nil, fmt.Sprintf("/docs/view/%d", id))
}

func (d *DocController) Edit(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	docService := service.NewDocService(c)

	if c.Request.Method == "GET" {
		doc, err := docService.GetDetail(_id)

		if err != nil {
			if err.Error() == "not found" {
				d.renderError(c, 404, "文档不存在")
				return
			}

			d.renderError(c, 500, "数据获取失败")

			return
		}

		docs, _ := docService.GetDocs(doc.ProjectID)

		d.render(c, "edit", gin.H{
			"doc":  doc,
			"docs": docs,
		})

		return
	}

	form := &DocForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		d.json(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"title":    form.Title,
		"label":    form.Label,
		"markdown": form.Markdown,
	}

	err := docService.Edit(_id, data)

	if err != nil {
		d.json(c, false, "编辑失败")
		return
	}

	d.json(c, true, "编辑成功", nil, fmt.Sprintf("/docs/view/%s", id))
}

func (d *DocController) Delete(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	docService := service.NewDocService(c)
	doc, err := docService.GetDetail(_id)

	if err != nil {
		if err.Error() == "not found" {
			d.json(c, false, "文档不存在")
			return
		}

		d.json(c, false, "删除失败")

		return
	}

	err = docService.Delete(_id)

	if err != nil {
		d.json(c, false, "删除失败")
		return
	}

	d.json(c, true, "删除成功", nil, fmt.Sprintf("/projects/view/%d", doc.ProjectID))
}
