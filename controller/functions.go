package controller

import (
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

func ShowReqUriFunction(c *gin.Context) string {
	return c.Request.Method + " " + c.Request.URL.RequestURI()
}

func StdDate(t interface{}) interface{} {
	timeInMillis, ok := t.(int64)
	if !ok {
		return ""
	}
	tttt := time.Unix(timeInMillis / 1000 / 1000 / 1000, timeInMillis % (1000 * 1000 * 1000))
	return tttt.Format("2006-01-02 15:04:05")
}
