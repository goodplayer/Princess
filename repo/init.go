package repo

import (
	"database/sql"
	"log"
)

import (
	_ "github.com/lib/pq"

	"moetang.info/prod/Princess/config"
)

var (
	db     *sql.DB
	isOpen bool
)

func init() {
	isOpen = false
}

func InitRepo(config *config.Config) {
	var err error
	db, err = sql.Open("postgres", config.DbConfig.DbConnStr)
	if err != nil {
		log.Println("db error:", err)
		panic(err)
	}
	isOpen = true
	db.SetMaxIdleConns(config.DbConfig.DbMinCount)
	db.SetMaxOpenConns(config.DbConfig.DbMaxCount)
}
