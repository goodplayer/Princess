package model

import (
	"context"
	"github.com/goodplayer/Princess/repo"
)

type User2 struct {
	Id          int64
	Email       string
	Password    string
	DisplayName string

	Config string

	TimeCreated int64
	TimeUpdated int64
}

func (this *User2) SaveUser2() error {
	_, err := repo.Run().ExecContext(context.Background(), `
INSERT INTO princess_user
(user_name, user_password, display_name, status, config, user_type, time_created)
VALUES($1, $2, $3, 0, $4, 1, $5);
`, this.Email, this.Password, this.DisplayName, this.Config, this.TimeCreated)
	return err
}
