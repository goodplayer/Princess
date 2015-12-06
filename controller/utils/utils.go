package utils

import ()

import (
	"github.com/gin-gonic/gin"
)

var (
	g_TEMPLATE_COMMON_MODEL = make(map[string]interface{})
)

func Init(f func(map[string]interface{})) {
	f(g_TEMPLATE_COMMON_MODEL)
}

func NewTemplateModel(c *gin.Context) gin.H {
	// build in template obj
	r := gin.H{
		"CONTEXT": c,
		"SITE":    g_TEMPLATE_COMMON_MODEL,
	}
	return r
}
