package session

import (
	"github.com/astaxie/beego/session"
	"github.com/gin-gonic/gin"

	"github.com/goodplayer/Princess/session/sessionutil"
)

func GetSession(c *gin.Context) session.Store {
	return c.MustGet(sessionutil.GetContextSessionKey()).(session.Store)
}
