package sessionutil

import (
	"github.com/astaxie/beego/session"
	"github.com/gin-gonic/gin"
)

const contextSessionKey string = "___context_session___"

func GetContextSessionKey() string {
	return contextSessionKey
}

func InitContextSession(c *gin.Context, s session.Store) {
	c.Set(contextSessionKey, s)
}
