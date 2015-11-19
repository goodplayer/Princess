package session

import (
	"github.com/astaxie/beego/session"
	"github.com/gin-gonic/gin"

	"moetang.info/prod/Princess/session/sessionutil"
)

func GetSession(c *gin.Context) session.SessionStore {
	return c.MustGet(sessionutil.GetContextSessionKey()).(session.SessionStore)
}
