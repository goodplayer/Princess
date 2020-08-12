package model

import (
	"context"
	"errors"

	"github.com/goodplayer/Princess/repo"
)

type User2 struct {
	Id          int64   `json:"id"`
	Email       string  `json:"email"`
	Password    string  `json:"-"`
	DisplayName *string `json:"display_name"`

	Status   *int `json:"status"`
	UserType *int `json:"user_type"`

	Config *string `json:"config"`

	TimeCreated int64  `json:"time_created"`
	TimeUpdated *int64 `json:"time_updated"`
}

func (this *User2) SaveUser2() error {
	_, err := repo.Run().ExecContext(context.Background(), `
INSERT INTO princess_user
(user_name, user_password, display_name, status, config, user_type, time_created)
VALUES($1, $2, $3, 0, $4, 1, $5);
`, this.Email, this.Password, this.DisplayName, this.Config, this.TimeCreated)
	return err
}

func (this *User2) LoadUser2() error {
	rows, err := repo.Run().Query("select user_id, user_name, user_password, display_name, status, config, user_type, time_created, time_updated from princess_user where user_name = $1", this.Email)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&this.Id, &this.Email, &this.Password, &this.DisplayName, &this.Status, &this.Config, &this.UserType, &this.TimeCreated, &this.TimeUpdated)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return errors.New("no user found")
	}
}
