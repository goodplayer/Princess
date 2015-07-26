package controller

import (
	"github.com/gin-gonic/gin"
)

func ShowReqUriFunction(c *gin.Context) string {
	return c.Request.Method + " " + c.Request.URL.RequestURI()
}
