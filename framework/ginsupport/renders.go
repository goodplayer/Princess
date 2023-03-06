package ginsupport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRenderTemplateSuccess(template string, resCtx gin.H) Render {
	return &templateRender{
		status:   http.StatusOK,
		template: template,
		data:     resCtx,
	}
}

type templateRender struct {
	status   int
	template string
	data     gin.H
}

func (t *templateRender) renderContext(ctx *gin.Context) {
	ctx.HTML(t.status, t.template, t.data)
}
