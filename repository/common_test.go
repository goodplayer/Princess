package repository

import (
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GlobalDb *gorm.DB

var initDbOnce *sync.Once

func init() {
	initDbOnce = new(sync.Once)
}

func Init() {
	initDbOnce.Do(func() {
		dsn := "sqlserver://sa:P@ssw0rdP@ssw0rd@192.168.31.207:1433?database=princess"
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(err)
		}
		GlobalDb = db
	})
}
