package controller2

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func defaultLang(c *gin.Context) string {
	host := c.Request.Host

	cookie, err := c.Cookie("lang")
	if err != http.ErrNoCookie || len(cookie) == 0 {
		c.SetCookie("lang", LangEn, 3600*24*30, "/", host, false, true)
		return LangEn
	} else {
		return cookie
	}
}
