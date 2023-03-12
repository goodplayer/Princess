package ginplugin

import (
	"github.com/gin-gonic/gin"
)

const defaultSessionManagerKey = "github.com/goodplayer/Princess"

func SessionHandler(sessionDomainService *SessionDomainService) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		sessionDomainService.initSessionManager(ctx)
		ctx.Next()
	}
}

func Session(ctx *gin.Context) SessionManager {
	return session(ctx)
}

type SessionManager interface {
	Set(key string, value []byte)
	Get(key string) []byte
	Delete(key string)
	SaveAndFreeze() error
}
