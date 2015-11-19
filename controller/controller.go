package controller

import (
	"log"
	"net/http"
	"strconv"
)

import (
	"github.com/gin-gonic/gin"

	"moetang.info/prod/Princess/model"
	"moetang.info/prod/Princess/session"
)

const (
	_session_user_id = "___princess_user_id___"

	session_key = "session"
)

func IndexAction(c *gin.Context) {
	posts, err := model.PostUtil().GetIndexPosts()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", NewTemplateModel(c))
	} else {
		result := NewTemplateModel(c)
		result["posts"] = posts
		c.HTML(http.StatusOK, "index.html", result)
	}
}

func LoginAction(c *gin.Context) {
	result := NewTemplateModel(c)
	success := false

	username := c.PostForm("j_username")
	password := c.PostForm("j_password")
	user := model.NewUser()
	user.Username = username
	err := user.FillUserByUsername()
	log.Println(user)

	if err == model.NO_SUCH_RECORD {
		result["ErrorCode"] = "用户名或密码错误"
	} else if err != nil {
		log.Println(err)
		result["ErrorCode"] = "系统错误"
	} else if user.Password != password {
		result["ErrorCode"] = "用户名或密码错误"
	} else {
		success = true
	}

	if !success {
		// 错误情况下
		c.HTML(http.StatusInternalServerError, "error.html", result)
	} else {
		// 正确情况下
		sess := session.GetSession(c)
		MarkLogin(sess)
		sess.Set("user", user) // fill user info
		sess.Set(_session_user_id, user.Id)
		log.Println("user_id:", sess.Get(_session_user_id), "login")
		c.Redirect(http.StatusFound, "/")
	}
}

func LogoutAction(c *gin.Context) {
	sess := session.GetSession(c)
	MarkNotLogin(sess)
	log.Println("user_id:", sess.Get(_session_user_id), "logout")
	sess.Flush()
	c.Redirect(http.StatusFound, "/")
}

func ShowPostAction(c *gin.Context) {
	postId := c.Param("postId")
	id, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", NewTemplateModel(c))
		return
	}

	post, ok, err := model.PostUtil().GetPostById(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", NewTemplateModel(c))
		return
	} else if !ok {
		NoRouteHandler(c)
		return
	}

	result := NewTemplateModel(c)
	result["post"] = post
	c.HTML(http.StatusOK, "post.html", result)
}
