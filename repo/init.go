package repo

import (
	"database/sql"
	"log"
)

import (
	_ "github.com/lib/pq"

	"github.com/goodplayer/Princess/config"
)

var (
	db     *sql.DB
	isOpen bool

	db_obj      *sql.DB
	isDbObjOpen bool
)

func init() {
	isOpen = false
	isDbObjOpen = false
}

func InitRepo(config *config.Config) {
	{
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
	{
		var err error
		db_obj, err = sql.Open("postgres", config.DbObjConfig.DbConnStr)
		if err != nil {
			log.Println("db_obj error:", err)
			panic(err)
		}
		isDbObjOpen = true
		db_obj.SetMaxIdleConns(config.DbObjConfig.DbMinCount)
		db_obj.SetMaxOpenConns(config.DbObjConfig.DbMaxCount)
	}
}
