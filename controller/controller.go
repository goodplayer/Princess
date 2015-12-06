package controller

import (
	"log"
	"net/http"
	"strconv"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

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
	if len(username) == 0 || len(password) == 0 {
		result["ErrorCode"] = "请填写完整信息"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	user := model.NewUser()
	user.Username = username
	err := user.FillUserByUsername()
	log.Println(user)

	if err == model.NO_SUCH_RECORD {
		result["ErrorCode"] = "用户名或密码错误"
	} else if err != nil {
		log.Println(err)
		result["ErrorCode"] = "系统错误"
	} else if !model.UserUtil().CheckPasswordMatch(user, password) {
		// check password
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

func RegisterAction(c *gin.Context) {
	result := NewTemplateModel(c)

	username := c.PostForm("j_username")
	password := c.PostForm("j_password")
	nickname := c.PostForm("j_nickname")
	email := c.PostForm("j_email")

	if len(username) == 0 || len(password) < 8 || len(nickname) == 0 || len(email) == 0 {
		result["ErrorCode"] = "请填写完整信息"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}
	if len(username) > 32 || len(password) > 32 || len(nickname) > 32 || len(email) > 90 {
		result["ErrorCode"] = "信息填写过长"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}
	//TODO 检查参数格式

	exist, err := model.UserUtil().UserExists(username)
	if err != nil {
		result["ErrorCode"] = "系统错误，请重试"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		log.Println("[ERROR] check user exist:", err)
		return
	}
	if exist {
		result["ErrorCode"] = "用户名已经存在"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	user := model.NewUser()
	user.Username = username
	user.Password = password
	user.Nickname = nickname
	user.Status = 0
	user.Authority = 0
	user.Email = email
	err = model.UserUtil().PreparePassword(user)
	if err != nil {
		result["ErrorCode"] = "系统错误，请重试"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		log.Println("[ERROR] PreparePassword:", err)
		return
	}

	err = user.Save()
	if err != nil {
		result["ErrorCode"] = "系统错误，请重试"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		log.Println("[ERROR] user save:", err)
		return
	}
	result["Result"] = "用户注册成功，请登录"
	c.HTML(http.StatusOK, "info.html", result)
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

func ShowNewPostPageAction(c *gin.Context) {
	if IsLogin(c) == "false" {
		result := NewTemplateModel(c)
		result["ErrorCode"] = "请先登录"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	user := GetUser(c).(*model.User)
	if user.Authority <= 0 {
		result := NewTemplateModel(c)
		result["ErrorCode"] = "没有权限发表文章"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	c.HTML(http.StatusOK, "new_post.html", NewTemplateModel(c))
}

func NewPostAction(c *gin.Context) {
	if IsLogin(c) == "false" {
		result := NewTemplateModel(c)
		result["ErrorCode"] = "请先登录"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	if len(title) == 0 || len(content) == 0 {
		result := NewTemplateModel(c)
		result["ErrorCode"] = "请填写所有内容"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	user := GetUser(c).(*model.User)
	if user.Authority <= 0 {
		result := NewTemplateModel(c)
		result["ErrorCode"] = "没有权限发表文章"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	//	contentData := bluemonday.UGCPolicy().SanitizeBytes([]byte(content))
	content = string(blackfriday.MarkdownCommon([]byte(content)))
	title = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(title)))

	post := model.NewPost()
	post.Title = title
	abstractLen := len(content)
	if abstractLen > 6000 {
		abstractLen = 6000
	}
	post.Abstract = content[:abstractLen]
	post.Content = content
	post.PostUser = user
	post.Status = 0
	err := post.Save()
	if err != nil {
		log.Println("[ERROR] save post:", err)
		result := NewTemplateModel(c)
		result["ErrorCode"] = "系统错误"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	result := NewTemplateModel(c)
	result["Result"] = "发表成功"
	c.HTML(http.StatusOK, "info.html", result)
}
