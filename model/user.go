package model

import (
	"database/sql"
)

import (
	"moetang.info/prod/Princess/repo"
)

type User struct {
	Id       int64
	Username string
	Password string
	Nickname string
	Status   int64
	Email    string
}

func NewUser() *User {
	return new(User)
}

func (this *User) Save() error {
	_, err := repo.DB.Exec(`INSERT INTO "user"(username, password, nickname, status, email) VALUES ($1, $2, $3, $4, $5);`,
		this.Username, this.Password, this.Nickname, this.Status, this.Email)
	return err
}

func (this *User) GetUserById() error {
	r := repo.DB.QueryRow(`SELECT username, password, nickname, status, email FROM "user" where id = $1;`, this.Id)
	e := r.Scan(&this.Username, &this.Password, &this.Nickname, &this.Status, &this.Email)
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

func (this *User) GetUserByUsername() error {
	r := repo.DB.QueryRow(`SELECT id, password, nickname, status, email FROM "user" where username = $1;`, this.Username)
	e := r.Scan(&this.Id, &this.Password, &this.Nickname, &this.Status, &this.Email)
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
