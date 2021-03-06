package controller2

import (
	"html/template"
	"net/http"

	"github.com/goodplayer/Princess/config"

	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// home page for redirecting
	r.GET("/", indexAction)
	//
	//// login
	//r.POST("/login", LoginAction)
	//r.POST("/register", RegisterAction)
	//r.GET("/logout", LogoutAction)
	//
	//r.GET("/new_post", ShowNewPostPageAction)
	//r.POST("/posts", NewPostAction)
	//
	//// admin
	//r.GET("/admin/users", controllers.ShowUsersAction)

	//TODO remove
	InitRecording(r)

	InitReg(r)

	templ := template.New("")
	registerFunction(templ)
	templ = template.Must(templ.ParseGlob(config.GLOBAL_CONFIG.TemplatePath))
	r.SetHTMLTemplate(templ)

	r.Static("/statics", config.GLOBAL_CONFIG.StaticPath)

	r.NoRoute(NoRouteHandler)
}

func NoRouteHandler(c *gin.Context) {
	ctx := map[string]interface{}{
		"site_name": config.GLOBAL_CONFIG.SiteConfig.DefaultSiteName,
	}

	c.HTML(http.StatusNotFound, "404.html", ctx)
}

func registerFunction(t *template.Template) {

}
