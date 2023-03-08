package user

import (
	"errors"

	"github.com/goodplayer/Princess/repository"

	"github.com/sirupsen/logrus"
)

var userServiceLogger = logrus.New()

type UserDomainService struct {
	db *repository.Db
}

func (u *UserDomainService) DependOn(db *repository.Db) {
	u.db = db
}

func (u *UserDomainService) UserExistsByEmail(email string) (bool, error) {
	return u.db.CheckUserExistsByEmail(email)
}

func (u *UserDomainService) DoUserRegister(user *struct {
	Email    string
	Password string
}) (*User, error) {

	newUser := new(User).InitFill(user.Email)
	if err := newUser.UpdatePassword(user.Password); err != nil {
		return nil, err
	}
	if _, err := u.db.SaveUser(newUser.ToDbModel()); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *UserDomainService) UserAuthByPassword(username, password string) (*User, error) {
	if user, err := u.db.LoadUserByUsername(username); err != nil {
		userServiceLogger.Errorf("load username failed:%s", err)
		return nil, err
	} else if user == nil {
		userServiceLogger.Warnf("user not found:%s", username)
		return nil, errors.New("user not found or password not match")
	} else {
		u := new(User).FromDbModel(user)
		if u.VerifyPassword(password) {
			return u, nil
		} else {
			userServiceLogger.Warnf("user password not match:%s", username)
			return nil, errors.New("user not found or password not match")
		}
	}
}
