package controller

import (
	"net/http"
	"strconv"
)

import (
	"github.com/gin-gonic/gin"

	"moetang.info/prod/Princess/model"
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
	//TODO
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
