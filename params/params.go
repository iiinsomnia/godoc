package params

func GetHistoryFlag(flag int) string {
	flagName := ""

	switch flag {
	case 1:
		flagName = "创建"
	case 2:
		flagName = "修改"
	}

	return flagName
}
