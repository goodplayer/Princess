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
	DB *sql.DB
)

func InitRepo(config *config.Config) {
	var err error
	DB, err = sql.Open("postgres", config.DbConfig.DbConnStr)
	if err != nil {
		log.Println("db error:", err)
		panic(err)
	}
	DB.SetMaxIdleConns(config.DbConfig.DbMinCount)
	DB.SetMaxOpenConns(config.DbConfig.DbMaxCount)
}
