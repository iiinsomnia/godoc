package controllers

import (
	"godoc/helpers"
	"godoc/i18n"
	"godoc/rbac"
	"godoc/service"
	"html/template"
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
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
	Role  int    `form:"role" binding:"required"`
}

type PasswordForm struct {
	Password       string `form:"password" binding:"required"`
	PasswordRepeat string `form:"password_repeat" binding:"required"`
}

func NewUserController(r *gin.Engine) *UserController {
	return &UserController{
		construct(r, "main", "user"),
	}
}

func (u *UserController) Index(c *gin.Context) {
	query := c.Request.URL.Query()

	userService := service.NewUserService(c)
	count, curPage, totalPage, data, err := userService.GetUserList(query, 1)

	if err != nil {
		u.renderError(c, 500, "数据获取失败")
		return
	}

	u.addFuncs(template.FuncMap{
		"httpBuildQueryUrl": helpers.HttpBuildQueryUrl,
	}).addTpls("layouts/pagination").render(c, "index", gin.H{
		"count":     count,
		"totalPage": totalPage,
		"query":     query,
		"users":     data,
		"pagination": gin.H{
			"uri":      "/users",
			"curPage":  curPage,
			"prevPage": curPage - 1,
			"nextPage": curPage + 1,
			"lastPage": totalPage,
			"pages":    helpers.Pagination(curPage, totalPage),
		},
	})
}

func (u *UserController) Add(c *gin.Context) {
	if c.Request.Method == "GET" {
		u.render(c, "add", gin.H{
			"defaultPass": yiigo.GetEnvString("app", "defaultPass", "123"),
			"roles":       rbac.Roles,
		})

		return
	}

	form := &UserForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.json(c, false, i18n.I18NSlice(errors))

		return
	}

	userService := service.NewUserService(c)

	unique, err := userService.CheckUnique(form.Name, form.Email)

	if err != nil {
		u.json(c, false, "添加失败")
		return
	}

	if !unique {
		u.json(c, false, "用户名或邮箱已被注册")
		return
	}

	defaultPass := yiigo.GetEnvString("app", "defaultPass", "123")
	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"name":     form.Name,
		"email":    form.Email,
		"role":     form.Role,
		"password": helpers.MD5(defaultPass + salt),
		"salt":     salt,
	}

	_, err = userService.Add(data)

	if err != nil {
		u.json(c, false, "添加失败")
		return
	}

	u.json(c, true, "添加成功", nil, "/users")
}

func (u *UserController) Edit(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUserService(c)

	if c.Request.Method == "GET" {
		user, err := userService.GetDetail(_id)

		if err != nil {
			if err.Error() == "not found" {
				u.renderError(c, 404, "用户不存在")
				return
			}

			u.renderError(c, 500, "数据获取失败")

			return
		}

		u.render(c, "edit", gin.H{
			"user":  user,
			"roles": rbac.Roles,
		})

		return
	}

	form := &UserForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.json(c, false, i18n.I18NSlice(errors))

		return
	}

	unique, err := userService.CheckUnique(form.Name, form.Email, _id)

	if err != nil {
		u.json(c, false, "添加失败")
		return
	}

	if !unique {
		u.json(c, false, "用户名或邮箱已被使用")
		return
	}

	data := yiigo.X{
		"name":  form.Name,
		"email": form.Email,
		"role":  form.Role,
	}

	err = userService.Edit(_id, data)

	if err != nil {
		u.json(c, false, "编辑失败")
		return
	}

	u.json(c, true, "编辑成功", nil, "/users")
}

func (u *UserController) Password(c *gin.Context) {
	form := &PasswordForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.json(c, false, i18n.I18NSlice(errors))

		return
	}

	if form.Password != form.PasswordRepeat {
		u.json(c, false, "密码确认错误")
		return
	}

	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"password": helpers.MD5(form.Password + salt),
		"salt":     salt,
	}

	identity := rbac.GetIdentity(c)

	userService := service.NewUserService(c)
	err := userService.Edit(identity.ID, data)

	if err != nil {
		u.json(c, false, "修改失败")
		return
	}

	u.json(c, true, "修改成功", nil, "/users")
}

func (u *UserController) Reset(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	defaultPass := yiigo.GetEnvString("app", "defaultPass", "123")
	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"password": helpers.MD5(defaultPass + salt),
		"salt":     salt,
	}

	userService := service.NewUserService(c)
	err := userService.Edit(_id, data)

	if err != nil {
		u.json(c, false, "重置失败")
		return
	}

	u.json(c, true, "重置成功", nil, "/users")
}

func (u *UserController) Delete(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUserService(c)

	err := userService.Delete(_id)

	if err != nil {
		u.json(c, false, "删除失败")
		return
	}

	u.json(c, true, "删除成功", nil, "/users")
}
