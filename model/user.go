package model

import (
	"container/list"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

import (
	"github.com/goodplayer/Princess/repo"
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

func (userUtil) GetAllUsers() ([]*User, error) {
	rows, err := repo.Run().Query(`SELECT id, username, password, nickname, status, email, createtime, lastupdatetime, salt, authority FROM "user";`)
	if err != nil {
		return nil, err
	}
	defer func() {
		e := rows.Close()
		if e != nil {
			log.Println(e)
		}
	}()
	l := list.New()
	for rows.Next() {
		user, err := wrapUser(rows)
		if err != nil {
			return nil, err
		}
		l.PushBack(user)
	}
	if l.Len() > 0 {
		result := make([]*User, l.Len())
		idx := 0
		for e := l.Front(); e != nil; e = e.Next() {
			result[idx] = e.Value.(*User)
			idx++
		}
		return result, nil
	} else {
		return []*User{}, nil
	}
}

func wrapUser(r *sql.Rows) (*User, error) {
	user := new(User)
	e := r.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.Status, &user.Email, &user.CreateTime, &user.LastUpdateTime, &user.Salt, &user.Authority)
	if e != nil {
		return nil, e
	}
	return user, nil
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
	Authority      int64
}

func (this *User) String() string {
	return fmt.Sprintf("id=[%d], username=[%s], password=[%s], nickname=[%s], status=[%d]", this.Id, this.Username, this.Password, this.Nickname, this.Status)
}

func NewUser() *User {
	return new(User)
}

func (this *User) CanPost() bool {
	return this.Authority > 0
}

func (this *User) Save() error {
	this.CreateTime = time.Now().UnixNano()
	this.LastUpdateTime = this.CreateTime
	_, err := repo.Run().Exec(`INSERT INTO "user"(username, password, nickname, status, email, createtime, lastupdatetime, salt, authority) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		this.Username, this.Password, this.Nickname, this.Status, this.Email, this.CreateTime, this.LastUpdateTime, this.Salt, this.Authority)
	return err
}

func (this *User) FillUserById() error {
	r := repo.Run().QueryRow(`SELECT username, password, nickname, status, email, createtime, lastupdatetime, salt, authority FROM "user" where id = $1;`, this.Id)
	e := r.Scan(&this.Username, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime, &this.Salt, &this.Authority)
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
	r := repo.Run().QueryRow(`SELECT id, password, nickname, status, email, createtime, lastupdatetime, salt, authority FROM "user" where username = $1;`, this.Username)
	e := r.Scan(&this.Id, &this.Password, &this.Nickname, &this.Status, &this.Email, &this.CreateTime, &this.LastUpdateTime, &this.Salt, &this.Authority)
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
