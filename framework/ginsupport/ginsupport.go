package ginsupport

import (
	"github.com/gin-gonic/gin"
)

type Render interface {
	renderContext(ctx *gin.Context)
}

type Action func(ctx *gin.Context) Render

func GroupGET(group *gin.RouterGroup, uri string, action Action) {
	group.GET(uri, func(context *gin.Context) {
		render := action(context)
		render.renderContext(context)
	})
}

func GroupPOST(group *gin.RouterGroup, uri string, action Action) {
	group.POST(uri, func(context *gin.Context) {
		render := action(context)
		render.renderContext(context)
	})
}
