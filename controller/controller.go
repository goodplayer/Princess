package controller

import (
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
)

func IndexAction(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", NewTemplateModel(c))
}
