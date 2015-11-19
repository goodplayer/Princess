package session

import (
	"strconv"
)

import (
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

func GetSessionManager() *session.Manager {
	return globalSessions
}

func Init(cookieName string) {
	var err error
	globalSessions, err = session.NewManager("memory", `{"cookieName":"`+cookieName+`","gclifetime":3600,"gclifetime":`+strconv.Itoa(3600*24*30)+`}`)
	if err != nil {
		panic(err)
	}
	go globalSessions.GC()
}
