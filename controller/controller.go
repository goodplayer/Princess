package controller

import (
	"net/http"
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

}
