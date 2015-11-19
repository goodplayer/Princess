package model

import (
	"database/sql"
	"time"
	"fmt"
)

import (
	"moetang.info/prod/Princess/repo"
)

type User struct {
	Id             int64
	Username       string
	Password       string
	Nickname       string
	Status         int64
	Email          string
	CreateTime     int64
	LastUpdateTime int64
}

func (this *User) String() string{
	return fmt.Sprintf("id=[%d], username=[%s], password=[%s], nickname=[%s], status=[%d]", this.Id, this.Username, this.Password, this.Nickname, this.Status)
}

func NewUser() *User {
	return new(User)
}

func (this *User) Save() error {
	this.CreateTime = time.Now().UnixNano()
	this.LastUpdateTime = this.CreateTime
	_, err := repo.Run().Exec(`INSERT INTO "user"(username, password, nickname, status, email, createtime, lastupdatetime) VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		this.Username, this.Password, this.Nickname, this.Status, this.Email, this.CreateTime, this.LastUpdateTime)
	return err
}

func (this *User) FillUserById() error {
	r := repo.Run().QueryRow(`SELECT username, password, nickname, status, email, createtime, lastupdatetime FROM "user" where id = $1;`, this.Id)
	e := r.Scan(&this.Username, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime)
	if e != nil {
		if e == sql.ErrNoRows {
			return NO_SUCH_RECORD
		} else {
			return e
		}
	} else {
		return nil
	}
}

func (this *User) FillUserByUsername() error {
	r := repo.Run().QueryRow(`SELECT id, password, nickname, status, email, createtime, lastupdatetime FROM "user" where username = $1;`, this.Username)
	e := r.Scan(&this.Id, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime)
	if e != nil {
		if e == sql.ErrNoRows {
			return NO_SUCH_RECORD
		} else {
			return e
		}
	} else {
		return nil
	}
}
