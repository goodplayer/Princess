package model

import (
	"database/sql"
	"time"
)

import (
	"moetang.info/prod/Princess/repo"
)

type Post struct {
	Id             int64
	Title          string
	Abstract       string
	Content        string
	Status         int64
	PostTime       int64
	LastUpdateTime int64
	PostUser       *User
}

func NewPost() *Post {
	return new(Post)
}

func (this *Post) Save() error {
	this.PostTime = time.Now().UnixNano()
	this.LastUpdateTime = this.PostTime
	_, err := repo.Run().Exec(`
INSERT INTO post(title, abstract, content, status, posttime, lastupdatetime, userid)
VALUES ($1, $2, $3, $4, $5, $6, $7);
	`, this.Title, this.Abstract, this.Content, this.Status, this.PostTime, this.LastUpdateTime, this.PostUser.Id)
	return err
}

func (this *Post) GetPostById() error {
	var userId int64 = -1
	r := repo.Run().QueryRow(`
SELECT title, abstract, content, status, posttime, lastupdatetime, userid FROM post where id = $1;
	`, this.Id)
	e := r.Scan(&this.Title, &this.Abstract, &this.Content, &this.Status, &this.PostTime, &this.LastUpdateTime, &userId)
	if e != nil {
		if e == sql.ErrNoRows {
			return NO_SUCH_RECORD
		} else {
			return e
		}
	}
	if userId <= 0 {
		return NO_USER_RELATED
	}
	user := NewUser()
	user.Id = userId
	err := user.GetUserById()
	if err != nil {
		return err
	}
	this.PostUser = user
	return nil
}
