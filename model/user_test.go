package model

import (
	"log"
	//	"testing"
)

import (
	"github.com/goodplayer/Princess/config"
	"github.com/goodplayer/Princess/repo"
)

func initDb() {
	conf := new(config.Config)
	conf.DbConfig.DbConnStr = "postgres://localhost/blog?sslmode=disable"
	conf.DbConfig.DbMinCount = 5
	conf.DbConfig.DbMaxCount = 10
	repo.InitRepo(conf)
}

func ExampleUserSave() {
	initDb()

	user := NewUser()
	user.Username = "testuser1"
	user.Password = "password1"
	user.Nickname = "nickname1"
	user.Email = "email1"
	user.Status = 0

	err := user.Save()

	if err != nil {
		log.Println(err)
	}
}

func ExampleUserGetById() {
	initDb()

	user := NewUser()
	user.Id = 8

	err := user.FillUserById()
	if err != nil {
		log.Println(err)
	}

	log.Println(user)
}

func ExampleUserGetByUsername() {
	initDb()

	user := NewUser()
	user.Username = "testuser1"

	err := user.FillUserByUsername()
	if err != nil {
		log.Println(err)
	}

	log.Println(user)
}
