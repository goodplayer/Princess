package controller2

import (
	"net/http"
	"time"

	"github.com/goodplayer/Princess/config"
	"github.com/goodplayer/Princess/model"
	"github.com/goodplayer/Princess/utils/password"

	"github.com/gin-gonic/gin"
)

func InitReg(r *gin.Engine) {
	r.GET("/reg", showRegAction)
	r.POST("/reg", regAction)
	r.GET("/login", showLoginAction)
	r.POST("/login", loginAction)
}

func showRegAction(c *gin.Context) {
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
	}

	c.HTML(http.StatusOK, "reg.html", ctx)
}

func showLoginAction(c *gin.Context) {
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
	}

	c.HTML(http.StatusOK, "login.html", ctx)
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
	user.DisplayName = &displayName
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

	// clear session
	err = clearLoginUserFromSession(c)
	if err != nil {
		LOGGER.Error("clear login user error.", err)
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

func loginAction(c *gin.Context) {
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	email := c.PostForm("email")
	pw := c.PostForm("password")

	LOGGER.Info("login info: ", email, " ", pw)

	user := new(model.User2)
	user.Email = email

	err := user.LoadUser2()
	if err != nil {
		LOGGER.Error("load user error.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "错误：" + err.Error(),
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	LOGGER.Error(user)

	match, err := password.VerifyScrypt(user.Password, pw)
	if err != nil {
		LOGGER.Error("verify user password error.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "错误：" + err.Error(),
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	if !match {
		LOGGER.Error("user login failed, password error.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "password error",
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	err = SaveLoginUserFromSession(c, user)
	if err != nil {
		LOGGER.Error("save login user to session.", err)
		ctx := map[string]interface{}{
			"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
			"Message":   "错误：" + err.Error(),
		}
		c.HTML(http.StatusInternalServerError, "error.html", ctx)
		return
	}

	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
		"Message":   "登录成功",
	}
	c.HTML(http.StatusOK, "info.html", ctx)
}
