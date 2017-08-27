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
	categoryService := service.NewCategory(c)
	categories, _ := categoryService.GetAll()

	h.V(c).Render("index", gin.H{
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
			h.Redirect(c, "/")
			return
		}

		h.V(c).L("normal").Render("login", gin.H{
			"captchaID": captcha.NewLen(5),
		})

		return
	}

	form := &LoginForm{}

	if validate := c.ShouldBindWith(form, binding.Form); validate != nil {
		errors := strings.Split(validate.Error(), "\n")
		h.JSON(c, false, i18n.I18NSlice(errors))

		return
	}

	captchaID := c.Param("captchaID")

	if !captcha.VerifyString(captchaID, form.Captcha) {
		h.JSON(c, false, "验证码错误", nil, "/login")
		return
	}

	authService := service.NewAuth(c)
	err := authService.Login(c, form.Account, form.Password)

	if err != nil {
		h.JSON(c, false, err.Error(), nil, "/login")
		return
	}

	h.JSON(c, true, "登录成功", nil, "/")
}

func (h *HomeController) Logout(c *gin.Context) {
	session.Destroy(c)
	h.Redirect(c, "/login")
}

func (h *HomeController) Captcha(c *gin.Context) {
	id := c.Param("id")

	captcha.Reload(id)
	err := captcha.WriteImage(c.Writer, id, 120, 34)

	if err != nil {
		yiigo.Err(err.Error)
	}
}

func (h *HomeController) NotFound(c *gin.Context) {
	h.V(c).RenderErr(404, "Page Not Found.")
}

func (h *HomeController) InternalServerError(c *gin.Context) {
	h.V(c).RenderErr(500, "Internal Server Error.")
}
