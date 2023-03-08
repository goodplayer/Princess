package controller

import (
	"github.com/goodplayer/Princess/framework/app"

	"github.com/gin-gonic/gin"
)

func InitController(r *gin.Engine, ac *app.ApplicationContainer) error {
	if err := initHomepage(r, ac); err != nil {
		return err
	}
	if err := initUser(r, ac); err != nil {
		return err
	}

	return nil
}
