package controllers

import (
	"log"
	"net/http"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/goodplayer/Princess/controller/utils"
	"github.com/goodplayer/Princess/model"
)

func ShowUsersAction(c *gin.Context) {
	result := utils.NewTemplateModel(c)
	if !utils.IsAdmin(c) {
		result["ErrorCode"] = "没有权限"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}

	postList, err := model.UserUtil().GetAllUsers()
	if err != nil {
		log.Println("[ERROR] admin / get all users error.", err)
		result["ErrorCode"] = "系统内部错误"
		c.HTML(http.StatusInternalServerError, "error.html", result)
		return
	}
	result["post_list"] = postList
	c.HTML(http.StatusOK, "admin_users.html", result)
}
