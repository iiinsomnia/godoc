package params

// GetRoleName 获取角色名称
func GetRoleName(role int) string {
	roleName := ""

	switch role {
	case 1:
		roleName = "超级管理员"
	case 2:
		roleName = "高级管理员"
	case 3:
		roleName = "普通管理员"
	}

	return roleName
}
