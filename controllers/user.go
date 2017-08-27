package controllers

import (
	"database/sql"
	"godoc/helpers"
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
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		u.V(c).RenderErr(403, "无操作权限")
		return
	}

	query := c.Request.URL.Query()

	userService := service.NewUser(c)
	count, curPage, totalPage, data, err := userService.GetUserList(query)

	if err != nil && err != sql.ErrNoRows {
		u.V(c).RenderErr(500, "数据获取失败")
		return
	}

	u.V(c).F(template.FuncMap{
		"int":               helpers.Int,
		"httpBuildQueryUrl": helpers.HttpBuildQueryUrl,
	}).T("layouts/pagination").Render("index", gin.H{
		"count":     count,
		"totalPage": totalPage,
		"query":     query,
		"roles":     rbac.Roles,
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
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		u.V(c).RenderErr(403, "无操作权限")
		return
	}

	if c.Request.Method == "GET" {
		u.V(c).F(template.FuncMap{
			"roleDesc": params.RoleDesc,
		}).Render("add", gin.H{
			"defaultPass": yiigo.EnvString("app", "defaultPass", "123"),
			"roles":       rbac.Roles,
		})

		return
	}

	form := &UserForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	userService := service.NewUser(c)

	unique, err := userService.CheckUnique(form.Name, form.Email)

	if err != nil {
		u.JSON(c, false, "添加失败")
		return
	}

	if !unique {
		u.JSON(c, false, "用户名或邮箱已被注册")
		return
	}

	defaultPass := yiigo.EnvString("app", "defaultPass", "123")
	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"name":     form.Name,
		"email":    form.Email,
		"role":     form.Role,
		"password": yiigo.MD5(defaultPass + salt),
		"salt":     salt,
	}

	_, err = userService.Add(data)

	if err != nil {
		u.JSON(c, false, "添加失败")
		return
	}

	u.JSON(c, true, "添加成功", nil, "/users")
}

func (u *UserController) Edit(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		u.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUser(c)

	if c.Request.Method == "GET" {
		user, err := userService.GetDetail(_id)

		if err != nil {
			if err == sql.ErrNoRows {
				u.V(c).RenderErr(404, "用户不存在")
				return
			}

			u.V(c).RenderErr(500, "数据获取失败")

			return
		}

		u.V(c).F(template.FuncMap{
			"roleDesc": params.RoleDesc,
		}).Render("edit", gin.H{
			"user":  user,
			"roles": rbac.Roles,
		})

		return
	}

	form := &UserForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	unique, err := userService.CheckUnique(form.Name, form.Email, _id)

	if err != nil {
		u.JSON(c, false, "添加失败")
		return
	}

	if !unique {
		u.JSON(c, false, "用户名或邮箱已被使用")
		return
	}

	data := yiigo.X{
		"name":  form.Name,
		"email": form.Email,
		"role":  form.Role,
	}

	err = userService.Edit(_id, data)

	if err != nil {
		u.JSON(c, false, "编辑失败")
		return
	}

	u.JSON(c, true, "编辑成功", nil, "/users")
}

func (u *UserController) Password(c *gin.Context) {
	form := &PasswordForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		u.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	if form.Password != form.PasswordRepeat {
		u.JSON(c, false, "密码确认错误")
		return
	}

	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"password": yiigo.MD5(form.Password + salt),
		"salt":     salt,
	}

	identity := rbac.GetIdentity(c)

	userService := service.NewUser(c)
	err := userService.Edit(identity.ID, data)

	if err != nil {
		u.JSON(c, false, "修改失败")
		return
	}

	u.JSON(c, true, "修改成功", nil, "/users")
}

func (u *UserController) Reset(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		u.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	defaultPass := yiigo.EnvString("app", "defaultPass", "123")
	salt := rbac.GenerateSalt()

	data := yiigo.X{
		"password": yiigo.MD5(defaultPass + salt),
		"salt":     salt,
	}

	userService := service.NewUser(c)
	err := userService.Edit(_id, data)

	if err != nil {
		u.JSON(c, false, "重置失败")
		return
	}

	u.JSON(c, true, "重置成功", nil, "/users")
}

func (u *UserController) Delete(c *gin.Context) {
	identity := rbac.GetIdentity(c)

	if identity.Role != 3 {
		u.V(c).RenderErr(403, "无操作权限")
		return
	}

	id := c.Param("id")
	_id, _ := strconv.Atoi(id)

	userService := service.NewUser(c)

	err := userService.Delete(_id)

	if err != nil {
		u.JSON(c, false, "删除失败")
		return
	}

	u.JSON(c, true, "删除成功", nil, "/users")
}
