package model

import ()

import (
	"moetang.info/prod/Princess/repo"
)

type Post struct {
	Id       int64
	Title    string
	Content  string
	PostTime int64
	PostUser User
}

func (this *Post) Save() error {
	_, err := repo.DB.Exec("")
	return err
}
