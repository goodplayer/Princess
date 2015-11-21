package model

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"
)

import (
	"moetang.info/prod/Princess/repo"
)

var (
	_userUtil userUtil
)

func init() {
	_userUtil = userUtil{}
}

type userUtil struct {
}

func UserUtil() userUtil {
	return _userUtil
}

func (userUtil) CheckPasswordMatch(user *User, password string) bool {
	salt := user.Salt
	rawPassword := []byte(password)
	raw := append(salt, rawPassword...)
	finalArray := sha256.Sum256(raw)
	final := finalArray[:]
	finalPassword := base64.StdEncoding.EncodeToString(final)
	return finalPassword == user.Password
}

func (userUtil) PreparePassword(user *User) error {
	// for register
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	rawPassword := []byte(user.Password)
	raw := append(salt, rawPassword...)
	finalArray := sha256.Sum256(raw)
	final := finalArray[:]
	finalPassword := base64.StdEncoding.EncodeToString(final)
	user.Password = finalPassword
	user.Salt = salt
	return nil
}

func (userUtil) UserExists(username string) (bool, error) {
	id := int64(-1)
	r := repo.Run().QueryRow(`SELECT id FROM "user" where username = $1;`, username)
	e := r.Scan(&id)
	if e != nil {
		if e == sql.ErrNoRows {
			return false, nil
		} else {
			return false, e
		}
	} else {
		return true, nil
	}
}

type User struct {
	Id             int64
	Username       string
	Password       string
	Nickname       string
	Status         int64
	Email          string
	CreateTime     int64
	LastUpdateTime int64
	Salt           []byte
}

func (this *User) String() string {
	return fmt.Sprintf("id=[%d], username=[%s], password=[%s], nickname=[%s], status=[%d]", this.Id, this.Username, this.Password, this.Nickname, this.Status)
}

func NewUser() *User {
	return new(User)
}

func (this *User) Save() error {
	this.CreateTime = time.Now().UnixNano()
	this.LastUpdateTime = this.CreateTime
	_, err := repo.Run().Exec(`INSERT INTO "user"(username, password, nickname, status, email, createtime, lastupdatetime, salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		this.Username, this.Password, this.Nickname, this.Status, this.Email, this.CreateTime, this.LastUpdateTime, this.Salt)
	return err
}

func (this *User) FillUserById() error {
	r := repo.Run().QueryRow(`SELECT username, password, nickname, status, email, createtime, lastupdatetime, salt FROM "user" where id = $1;`, this.Id)
	e := r.Scan(&this.Username, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime, &this.Salt)
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
	r := repo.Run().QueryRow(`SELECT id, password, nickname, status, email, createtime, lastupdatetime, salt FROM "user" where username = $1;`, this.Username)
	e := r.Scan(&this.Id, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime, &this.Salt)
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
