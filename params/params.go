package params

func GetHistoryFlag(flag int) string {
	name := ""

	switch flag {
	case 1:
		name = "创建"
	case 2:
		name = "修改"
	}

	return name
}

func RoleDesc(role int) string {
	desc := ""

	switch role {
	case 1:
		desc = "只能查看文档"
	case 2:
		desc = "拥有添加和修改文档的权限"
	case 3:
		desc = "拥有所有权限"
	}

	return desc
}
