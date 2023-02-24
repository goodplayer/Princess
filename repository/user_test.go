package repository_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/goodplayer/Princess/repository"

	"github.com/gofrs/uuid"
)

func TestDb_SaveUser(t *testing.T) {
	repository.Init()
	db := repository.NewDb(repository.GlobalDb)

	uid := uuid.Must(uuid.NewV7())

	user := &repository.User{
		UserId:     uid.String(),
		UserName:   strings.ReplaceAll(uid.String(), "-", ""),
		AliasName:  "demoalias01",
		Password:   "password01",
		UserStatus: 1,
	}
	if _, err := db.SaveUser(user); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SaveUser(user); err == nil {
		t.Fatal(errors.New("should report error due to primary key conflict"))
	}
}
