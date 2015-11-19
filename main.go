package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
)

import (
	"github.com/gin-gonic/gin"

	"moetang.info/prod/Princess/config"
	"moetang.info/prod/Princess/controller"
	"moetang.info/prod/Princess/repo"
	"moetang.info/prod/Princess/session"
	"moetang.info/prod/Princess/session/sessionutil"
)

func init() {
	ROOT_PATH, _ := os.Getwd()
	config.GLOBAL_CONFIG = config.Load(path.Join(ROOT_PATH, "princess.conf"))

	fmt.Println(config.GLOBAL_CONFIG)
}

func main() {
	r := gin.New()
	initProcess(r)

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
	controller.RegisterRoute(r)
	repo.InitRepo(config.GLOBAL_CONFIG)
}

func initMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery(), gin.Logger())
	//init session
	session.Init(config.GLOBAL_CONFIG.Sessionkey)
	r.Use(func(c *gin.Context) {
		sess, err := session.GetSessionManager().SessionStart(c.Writer, c.Request)
		if err != nil {
			log.Println(err)
		}
		defer sess.SessionRelease(c.Writer)
		sessionutil.InitContextSession(c, sess)
		c.Next()
	})
}
