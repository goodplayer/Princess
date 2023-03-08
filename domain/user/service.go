package user

import "github.com/goodplayer/Princess/repository"

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

func (u *UserDomainService) UserAuth() {

}
