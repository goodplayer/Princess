package utils

import (
	"github.com/gin-gonic/gin"

	"github.com/goodplayer/Princess/model"
	"github.com/goodplayer/Princess/session"
)

func IsAdmin(c *gin.Context) bool {
	sess := session.GetSession(c)
	if sess == nil {
		return false
	}
	userT := sess.Get("user")
	if userT == nil {
		return false
	}
	user := userT.(*model.User)
	if user.Authority == model.USER_AUTHORITY_ADMIN {
		return true
	}
	return false
}

func IsNormalUser(c *gin.Context) bool {
	sess := session.GetSession(c)
	if sess == nil {
		return false
	}
	userT := sess.Get("user")
	if userT == nil {
		return false
	}
	user := userT.(*model.User)
	if user.Authority == model.USER_AUTHORITY_NORMAL {
		return true
	}
	return false
}
