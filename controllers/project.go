package controllers

import (
	"database/sql"
	"fmt"
	"godoc/i18n"
	"godoc/rbac"
	"godoc/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iiinsomnia/yiigo"
)

type ProjectController struct {
	*controller
}

type ProjectForm struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
}

func NewProjectController(r *gin.Engine) *ProjectController {
	return &ProjectController{
		construct(r, "main", "project"),
	}
}

func (p *ProjectController) View(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	projectService := service.NewProjectService(c)
	project, err := projectService.GetDetail(_id)

	if err != nil {
		if err == sql.ErrNoRows {
			p.V(c).RenderErr(404, "项目不存在")
			return
		}

		p.V(c).RenderErr(500, "数据获取失败")

		return
	}

	docService := service.NewDocService(c)
	docs, _ := docService.GetDocs(project.ID)

	p.V(c).Render("view", gin.H{
		"project": project,
		"docs":    docs,
	})
}

func (p *ProjectController) Add(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role == 1 {
		p.V(c).RenderErr(403, "无操作权限")
		return
	}

	categoryID := c.Param("category")
	_categoryID, _ := strconv.Atoi(categoryID)

	categoryService := service.NewCategoryService(c)
	category, err := categoryService.GetDetail(_categoryID)

	if err != nil {
		if err == sql.ErrNoRows {
			p.V(c).RenderErr(404, "类别不存在")
			return
		}

		p.V(c).RenderErr(500, "数据获取失败")

		return
	}

	if c.Request.Method == "GET" {
		categories, _ := categoryService.GetAll()

		p.V(c).Render("add", gin.H{
			"category":   category,
			"categories": categories,
		})

		return
	}

	form := &ProjectForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		p.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name":        form.Name,
		"category_id": category.ID,
		"description": form.Description,
	}

	projectService := service.NewProjectService(c)
	_, err = projectService.Add(data)

	if err != nil {
		p.JSON(c, false, "添加失败")
		return
	}

	p.JSON(c, true, "添加成功", nil, fmt.Sprintf("/categories/view/%s", categoryID))
}

func (p *ProjectController) Edit(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role == 1 {
		p.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	projectService := service.NewProjectService(c)

	if c.Request.Method == "GET" {
		project, err := projectService.GetDetail(_id)

		if err != nil {
			if err == sql.ErrNoRows {
				p.V(c).RenderErr(404, "项目不存在")
				return
			}

			p.V(c).RenderErr(500, "数据获取失败")

			return
		}

		docService := service.NewDocService(c)
		docs, _ := docService.GetDocs(project.ID)

		p.V(c).Render("edit", gin.H{
			"project": project,
			"docs":    docs,
		})

		return
	}

	form := &ProjectForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		p.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name":        form.Name,
		"description": form.Description,
	}

	err := projectService.Edit(_id, data)

	if err != nil {
		p.JSON(c, false, "编辑失败")
		return
	}

	p.JSON(c, true, "编辑成功", nil, fmt.Sprintf("/projects/view/%s", id))
}

func (p *ProjectController) Delete(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		p.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	projectService := service.NewProjectService(c)
	project, err := projectService.GetDetail(_id)

	if err != nil {
		if err == sql.ErrNoRows {
			p.JSON(c, false, "项目不存在")
			return
		}

		p.JSON(c, false, "删除失败")

		return
	}

	err = projectService.Delete(_id)

	if err != nil {
		p.JSON(c, false, "删除失败")
		return
	}

	p.JSON(c, true, "删除成功", nil, fmt.Sprintf("/categories/view/%d", project.CategoryID))
}
