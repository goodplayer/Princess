package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

import ()

func NoRouteHandler(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", NewTemplateModel(c))
}
