package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/goodplayer/Princess/api/controller"
	"github.com/goodplayer/Princess/config"
	"github.com/goodplayer/Princess/domain"
	"github.com/goodplayer/Princess/framework/app"
	"github.com/goodplayer/Princess/framework/ginsupport"
	"github.com/goodplayer/Princess/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func init() {
	ROOT_PATH, _ := os.Getwd()
	config.InitConfig(config.Load(path.Join(ROOT_PATH, "princess.conf")))

	fmt.Println(config.GlobalConfig())
}

func main() {
	ac := app.NewAppContainer()
	r := gin.New()

	dsn := "sqlserver://sa:P@ssw0rdP@ssw0rd@192.168.31.232:1433?database=princess"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	ac.Set(repository.NewDb(db))

	if err := initProcess(r, ac); err != nil {
		panic(err)
	}

	// start container
	ac.StartContainer()

	if gin.IsDebugging() {
		log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", config.GlobalConfig().HttpConfig.Bind)
	}
	defer func() {
		if err != nil && gin.IsDebugging() {
			log.Printf("[GIN-debug] [ERROR] %v\n", err)
		}
	}()
	server := &http.Server{
		Addr:    config.GlobalConfig().HttpConfig.Bind,
		Handler: r,
	}
	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initProcess(r *gin.Engine, ac *app.ApplicationContainer) error {
	if err := domain.Init(ac); err != nil {
		return err
	}

	// init middleware before registering controllers
	// including custom function in order for template to use
	initMiddleware(r)
	if gin.IsDebugging() {
		ginsupport.InitHtmlTemplateFromFolderWithDebugReloading(r, "templates", "html")
	} else {
		ginsupport.InitHtmlTemplateFromFolder(r, "templates", "html")
	}
	// load controller
	if err := controller.InitController(r, ac); err != nil {
		return err
	}

	return nil
}

func initMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery(), gin.Logger())
	//init session
	store := memstore.NewStore([]byte(config.GlobalConfig().Sessionkey))
	r.Use(sessions.Sessions(config.GlobalConfig().Sessionkey, store))
}
