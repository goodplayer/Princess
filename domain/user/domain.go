package user

import (
	"strings"

	"github.com/goodplayer/Princess/repository"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	user *repository.User

	UserId    string
	Username  string
	AliasName string

	password string

	Status int
}

func (u *User) InitFill(username string) *User {
	u.Username = username
	u.AliasName = username
	u.Status = 0
	u.UserId = strings.ReplaceAll(uuid.Must(uuid.NewV4()).String(), "-", "")
	return u
}

func (u *User) ToDbModel() *repository.User {
	dbUser := &repository.User{
		UserId:     u.UserId,
		UserName:   u.Username,
		AliasName:  u.AliasName,
		Password:   u.password,
		UserStatus: u.Status,
	}
	u.user = dbUser

	return dbUser
}

func (u *User) FromDbModel(dbUser *repository.User) *User {
	u.Username = dbUser.UserName
	u.password = dbUser.Password
	u.user = dbUser
	u.AliasName = dbUser.AliasName
	u.UserId = dbUser.UserId
	u.Status = dbUser.UserStatus
	return u
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}

func (u *User) UpdatePassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.password = string(hash)
	return nil
}
