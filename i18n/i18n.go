package i18n

var I18N = map[string]string{
	"Key: 'LoginForm.Account' Error:Field validation for 'Account' failed on the 'required' tag":           "登录账号不可为空",
	"Key: 'LoginForm.Password' Error:Field validation for 'Password' failed on the 'required' tag":         "登录密码不可为空",
	"Key: 'LoginForm.Captcha' Error:Field validation for 'Captcha' failed on the 'required' tag":           "验证码不可为空",
	"Key: 'CategoryForm.Name' Error:Field validation for 'Name' failed on the 'required' tag":              "类别名称不可为空",
	"Key: 'ProjectForm.Name' Error:Field validation for 'Name' failed on the 'required' tag":               "项目名称不可为空",
	"Key: 'ProjectForm.Description' Error:Field validation for 'Description' failed on the 'required' tag": "项目描述不可为空",
	"Key: 'DocForm.Name' Error:Field validation for 'Name' failed on the 'required' tag":                   "文档名称不可为空",
	"Key: 'DocForm.Markdown' Error:Field validation for 'Markdown' failed on the 'required' tag":           "文档正文不可为空",
}

func I18NString(msg string) string {
	i18n := ""

	if v, ok := I18N[msg]; ok {
		i18n = v
	}

	return i18n
}

func I18NSlice(msgs []string) []string {
	i18n := []string{}

	for _, msg := range msgs {
		if v, ok := I18N[msg]; ok {
			i18n = append(i18n, v)
		} else {
			i18n = append(i18n, "")
		}
	}

	return i18n
}
