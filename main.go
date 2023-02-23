package main

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"

	"github.com/goodplayer/Princess/config"
	"github.com/goodplayer/Princess/controller2"
	"github.com/goodplayer/Princess/repo"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
	ROOT_PATH, _ := os.Getwd()
	config.GLOBAL_CONFIG = config.Load(path.Join(ROOT_PATH, "princess.conf"))

	fmt.Println(config.GLOBAL_CONFIG)
}

func main() {
	r := gin.New()
	initProcess(r)

	dsn := "sqlserver://sa:P@ssw0rdP@ssw0rd@192.168.31.207:1433?database=princess"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// start fastcgi
	go func() {
		if config.GLOBAL_CONFIG.HttpConfig.Enable {
			var err error = nil
			if gin.IsDebugging() {
				log.Printf("[GIN-debug] Listening and serving HTTP on %s\n", config.GLOBAL_CONFIG.HttpConfig.Bind)
			}
			defer func() {
				if err != nil && gin.IsDebugging() {
					log.Printf("[GIN-debug] [ERROR] %v\n", err)
				}
			}()

			server := &http.Server{
				Addr:    config.GLOBAL_CONFIG.HttpConfig.Bind,
				Handler: r,
			}

			err = server.ListenAndServe()
		}
	}()

	startFastCgi(r)
}

func startFastCgi(r *gin.Engine) {
	addr, err := net.ResolveTCPAddr("tcp", config.GLOBAL_CONFIG.Bind)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("[Fastcgi] starting fastcgi on", config.GLOBAL_CONFIG.Bind)
	fcgi.Serve(listener, r)
}

func initProcess(r *gin.Engine) {
	initMiddleware(r)
	controller2.Init(r)
	repo.InitRepo(config.GLOBAL_CONFIG)
}

func initMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery(), gin.Logger())
	//init session
	store := memstore.NewStore([]byte(config.GLOBAL_CONFIG.Sessionkey))
	r.Use(sessions.Sessions(config.GLOBAL_CONFIG.Sessionkey, store))
}
