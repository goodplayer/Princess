package controller2

import (
	"encoding/json"

	"github.com/goodplayer/Princess/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IsLoginFromSession(c *gin.Context) (bool, error) {
	u, err := GetLoginUserFromSession(c)
	if err != nil {
		return false, err
	} else {
		return u != nil, nil
	}
}

func GetLoginUserFromSession(c *gin.Context) (*model.User2, error) {
	s := sessions.Default(c)
	obj := s.Get("user")
	if obj != nil {
		u := new(model.User2)
		err := json.Unmarshal([]byte(obj.(string)), u)
		if err != nil {
			return nil, err
		}
		return u, nil
	} else {
		return nil, nil
	}
}

func SaveLoginUserFromSession(c *gin.Context, u *model.User2) error {
	s := sessions.Default(c)

	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	s.Set("user", string(data))
	return s.Save()
}

func clearLoginUserFromSession(c *gin.Context) error {
	s := sessions.Default(c)
	s.Delete("user")
	return s.Save()
}
