package controllers

import (
	"godoc/i18n"
	"godoc/service"
	"runtime"
	"strings"

	"godoc/session"

	"godoc/rbac"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/iiinsomnia/yiigo"
)

type HomeController struct {
	*controller
}

type LoginForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
	Captcha  string `form:"captcha" binding:"required"`
}

func NewHomeController(r *gin.Engine) *HomeController {
	return &HomeController{
		construct(r, "main", "home"),
	}
}

func (h *HomeController) Index(c *gin.Context) {
	categoryService := service.NewCategoryService(c)
	categories, _ := categoryService.GetAll()

	h.render(c, "index", gin.H{
		"os":         runtime.GOOS,
		"cpu":        runtime.NumCPU(),
		"arch":       runtime.GOARCH,
		"go":         runtime.Version(),
		"copyright":  "IIInsomnia 2017",
		"categories": categories,
	})
}

func (h *HomeController) Login(c *gin.Context) {
	if c.Request.Method == "GET" {
		if !rbac.IsGuest(c) {
			h.redirect(c, "/")
			return
		}

		h.layout("normal").render(c, "login", gin.H{
			"captchaID": captcha.NewLen(5),
		})

		return
	}

	form := &LoginForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		h.json(c, false, i18n.I18NSlice(errors))

		return
	}

	captchaID := c.Param("captchaID")

	if !captcha.VerifyString(captchaID, form.Captcha) {
		h.json(c, false, "验证码错误", nil, "/login")
		return
	}

	authService := service.NewAuthService(c)
	err := authService.Login(c, form.Account, form.Password)

	if err != nil {
		h.json(c, false, err.Error(), nil, "/login")
		return
	}

	h.json(c, true, "登录成功", nil, "/")
}

func (h *HomeController) Logout(c *gin.Context) {
	session.Destroy(c)
	h.redirect(c, "/login")
}

func (h *HomeController) Captcha(c *gin.Context) {
	id := c.Param("id")

	captcha.Reload(id)
	err := captcha.WriteImage(c.Writer, id, 120, 34)

	if err != nil {
		yiigo.LogError(err.Error)
	}
}

func (h *HomeController) NotFound(c *gin.Context) {
	h.renderError(c, 404, "Page Not Found.")
}

func (h *HomeController) InternalServerError(c *gin.Context) {
	h.renderError(c, 500, "Internal Server Error.")
}
