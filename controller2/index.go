package controller2

import (
	"net/http"

	"github.com/goodplayer/Princess/utils/logging"

	"github.com/gin-gonic/gin"
)

var (
	LOGGER = logging.NewLogger("IndexAction")
)

func indexAction(c *gin.Context) {
	//TODO
	LOGGER.Info(c.Request.Header.Get("Accept-Language"))
	defaultLang(c)

	c.HTML(http.StatusOK, "index.html", struct{}{})
}
