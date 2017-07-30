package controllers

import (
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

type CategoryController struct {
	*controller
}

type CategoryForm struct {
	Name string `form:"name" binding:"required"`
}

func NewCategoryController(r *gin.Engine) *CategoryController {
	return &CategoryController{
		construct(r, "main", "category"),
	}
}

func (c *CategoryController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	categoryService := service.NewCategoryService(ctx)
	category, err := categoryService.GetDetail(_id)

	if err != nil {
		if err.Error() == "not found" {
			c.renderError(ctx, 404, "类别不存在")
			return
		}

		c.renderError(ctx, 500, "数据获取失败")

		return
	}

	categories, _ := categoryService.GetAll()

	projectService := service.NewProjectService(ctx)
	projects, _ := projectService.GetProjects(_id)

	c.render(ctx, "view", gin.H{
		"category":   category,
		"categories": categories,
		"projects":   projects,
	})
}

func (c *CategoryController) Add(ctx *gin.Context) {
	identity := rbac.GetIdentity(ctx)

	if identity.Role == 1 {
		c.renderError(ctx, 403, "无操作权限")
		return
	}

	form := &CategoryForm{}

	if validate := ctx.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		c.json(ctx, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name": form.Name,
	}

	categoryService := service.NewCategoryService(ctx)
	id, err := categoryService.Add(data)

	if err != nil {
		c.json(ctx, false, "添加失败")
		return
	}

	c.json(ctx, true, "添加成功", nil, fmt.Sprintf("/categories/view/%d", id))
}

func (c *CategoryController) Edit(ctx *gin.Context) {
	identity := rbac.GetIdentity(ctx)

	if identity.Role == 1 {
		c.renderError(ctx, 403, "无操作权限")
		return
	}

	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	form := &CategoryForm{}

	if validate := ctx.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		c.json(ctx, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name": form.Name,
	}

	categoryService := service.NewCategoryService(ctx)
	err := categoryService.Edit(_id, data)

	if err != nil {
		c.json(ctx, false, "编辑失败")
		return
	}

	c.json(ctx, true, "编辑成功", nil, fmt.Sprintf("/categories/view/%s", id))
}

func (c *CategoryController) Delete(ctx *gin.Context) {
	identity := rbac.GetIdentity(ctx)

	if identity.Role != 3 {
		c.renderError(ctx, 403, "无操作权限")
		return
	}

	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	categoryService := service.NewCategoryService(ctx)

	err := categoryService.Delete(_id)

	if err != nil {
		c.json(ctx, false, "删除失败")
		return
	}

	c.json(ctx, true, "删除成功", nil, "/")
}
