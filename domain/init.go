package domain

import (
	"github.com/goodplayer/Princess/domain/user"
	"github.com/goodplayer/Princess/framework/app"
)

func Init(ac *app.ApplicationContainer) error {
	ac.Set(&user.UserDomainService{})

	return nil
}
