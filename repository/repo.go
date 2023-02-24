package repository

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// initialize uuid generator which will be used as primary key generator
var _ = uuid.Must(uuid.NewV7())

type Db struct {
	db *gorm.DB
}

func NewDb(db *gorm.DB) *Db {
	return &Db{db: db}
}
