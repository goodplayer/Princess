package controller

import (
	"html/template"
	"time"
)

import (
	"github.com/astaxie/beego/session"
	"github.com/gin-gonic/gin"

	"github.com/goodplayer/Princess/model"
	princess_session "github.com/goodplayer/Princess/session"
)

const (
	_session_login_flag string = "___login_flag___"
)

func ShowReqUriFunction(c *gin.Context) string {
	return c.Request.Method + " " + c.Request.URL.RequestURI()
}

func StdDate(t interface{}) interface{} {
	timeInMillis, ok := t.(int64)
	if !ok {
		return ""
	}
	tttt := time.Unix(timeInMillis/1000/1000/1000, timeInMillis%(1000*1000*1000))
	return tttt.Format("2006-01-02 15:04:05")
}

func IsLogin(c *gin.Context) string {
	sess := princess_session.GetSession(c)
	i := sess.Get(_session_login_flag)
	if i != nil {
		return "true"
	} else {
		return "false"
	}
}

func MarkLogin(sess session.Store) {
	sess.Set(_session_login_flag, _session_login_flag)
}

func MarkNotLogin(sess session.Store) {
	sess.Delete(_session_login_flag)
}

func GetUserNickName(c *gin.Context) string {
	sess := princess_session.GetSession(c)
	if sess == nil {
		return "[未登录]"
	}
	userT := sess.Get("user")
	if userT == nil {
		return "[未登录]"
	}
	user := userT.(*model.User)
	return user.Nickname
}

func GetUser(c interface{}) interface{} {
	ctx, ok := c.(*gin.Context)
	if !ok {
		return nil
	}
	sess := princess_session.GetSession(ctx)
	if sess == nil {
		return nil
	}
	userT := sess.Get("user")
	if userT == nil {
		return nil
	}
	user := userT.(*model.User)
	return user
}

func TemplateRawOutput(str interface{}) interface{} {
	if nstr, ok := str.(string); ok {
		return template.HTML(nstr)
	} else {
		return str
	}
}

func IsAdmin(c *gin.Context) string {
	sess := princess_session.GetSession(c)
	if sess == nil {
		return "[未登录]"
	}
	userT := sess.Get("user")
	if userT == nil {
		return "[未登录]"
	}
	user := userT.(*model.User)
	if user.Authority == model.USER_AUTHORITY_ADMIN {
		return "true"
	}
	return "false"
}

func IsNormalUser(c *gin.Context) string {
	sess := princess_session.GetSession(c)
	if sess == nil {
		return "[未登录]"
	}
	userT := sess.Get("user")
	if userT == nil {
		return "[未登录]"
	}
	user := userT.(*model.User)
	if user.Authority == model.USER_AUTHORITY_NORMAL {
		return "true"
	}
	return "false"
}
