package controller2

import "github.com/gin-gonic/gin"

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
	//
	//r.Static("/statics", config.GLOBAL_CONFIG.StaticPath)

}
