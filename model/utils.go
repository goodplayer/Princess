package model

func UserStatusString(e interface{}) interface{} {
	i, ok := e.(int64)
	if !ok {
		return "获取错误!"
	}
	switch i {
	case USER_STATUS_NORMAL:
		return "正常"
	case USER_STATUS_DELETED:
		return "已删除"
	}
	return "未知"
}

func UserAuthorityString(e interface{}) interface{} {
	i, ok := e.(int64)
	if !ok {
		return "获取错误!"
	}
	switch i {
	case USER_AUTHORITY_NORMAL:
		return "普通用户"
	case USER_AUTHORITY_ADMIN:
		return "管理员"
	}
	return "未知"
}
