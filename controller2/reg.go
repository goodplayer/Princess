package controller2

import (
	"github.com/goodplayer/Princess/model"
	"github.com/goodplayer/Princess/utils/password"
	"net/http"
	"time"

	"github.com/goodplayer/Princess/config"

	"github.com/gin-gonic/gin"
)

func InitReg(r *gin.Engine) {
	r.GET("/reg", showRegAction)
	r.POST("/reg", regAction)
}

func showRegAction(c *gin.Context) {
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
	}

	c.HTML(http.StatusOK, "reg.html", ctx)
}

func regAction(c *gin.Context) {
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	email := c.PostForm("email")
	pw := c.PostForm("password")
	displayName := c.PostForm("display_name")
	pw, err := password.EncryptScrypt(pw)
	if err != nil {
		LOGGER.Error("encrypt password error.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "错误：" + err.Error(),
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	LOGGER.Info("reg info: ", email, " ", pw, " ", displayName)

	user := new(model.User2)
	user.Email = email
	user.Password = pw
	user.DisplayName = displayName
	now := time.Now()
	user.TimeCreated = now.Unix()*1000 + now.UnixNano()/1000/1000%1000
	err = user.SaveUser2()
	if err != nil {
		LOGGER.Error("save user to db error.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "错误：" + err.Error(),
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
		"Message":   "注册成功",
	}
	c.HTML(http.StatusOK, "info.html", ctx)
}
