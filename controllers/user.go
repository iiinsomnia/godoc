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

type UserController struct {
	*controller
}

type UserForm struct {
	Name string `form:"name" binding:"required"`
}

func NewUserController(r *gin.Engine) *UserController {
	return &UserController{
		construct(r, "main", "user"),
	}
}

func (c *UserController) View(ctx *gin.Context) {
	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUserService(ctx)
	user, err := userService.GetDetail(_id)

	if err != nil {
		if err.Error() == "not found" {
			c.renderError(ctx, 404, "类别不存在")
			return
		}

		c.renderError(ctx, 500, "数据获取失败")

		return
	}

	categories, _ := userService.GetAll()

	projectService := service.NewProjectService(ctx)
	projects, _ := projectService.GetProjects(_id)

	c.render(ctx, "view", gin.H{
		"user":       user,
		"categories": categories,
		"projects":   projects,
	})
}

func (c *UserController) Add(ctx *gin.Context) {
	form := &UserForm{}

	if validate := ctx.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		c.json(ctx, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name": form.Name,
	}

	userService := service.NewUserService(ctx)
	id, err := userService.Add(data)

	if err != nil {
		c.json(ctx, false, "添加失败")
		return
	}

	c.json(ctx, true, "添加成功", nil, fmt.Sprintf("/categories/view/%d", id))
}

func (c *UserController) Edit(ctx *gin.Context) {
	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	form := &UserForm{}

	if validate := ctx.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		c.json(ctx, false, i18n.I18NSlice(errors))

		return
	}

	data := yiigo.X{
		"name": form.Name,
	}

	userService := service.NewUserService(ctx)
	err := userService.Edit(_id, data)

	if err != nil {
		c.json(ctx, false, "编辑失败")
		return
	}

	c.json(ctx, true, "编辑成功", nil, fmt.Sprintf("/categories/view/%s", id))
}

func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUserService(ctx)

	err := userService.Delete(_id)

	if err != nil {
		c.json(ctx, false, "删除失败")
		return
	}

	c.json(ctx, true, "删除成功", nil, "/")
}
