package session

import (
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

func GetSessionManager() *session.Manager {
	return globalSessions
}

func Init(cookieName string) {
	var err error
	cf := new(session.ManagerConfig)
	cf.CookieName = cookieName
	cf.Gclifetime = 3600
	cf.Maxlifetime = 3600 * 24 * 30
	globalSessions, err = session.NewManager("memory", cf)
	if err != nil {
		panic(err)
	}
	go globalSessions.GC()
}
