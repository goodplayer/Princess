package repo

import (
	"database/sql"
	"log"
)

func Run() *sql.DB {
	if !isOpen {
		//TODO need deal with closed situation
		return nil
	}
	return db
}

func RunObj() *sql.DB {
	if !isDbObjOpen {
		//TODO need deal with closed situation
		return nil
	}
	return db_obj
}

func Tx(f func(*sql.Tx) (interface{}, error)) (interface{}, error) {
	tx, err := Run().Begin()
	if err != nil {
		return nil, err
	}
	var rErr error = nil
	var result interface{} = nil
	defer func() {
		if rErr != nil {
			errFinal := tx.Rollback()
			if errFinal != nil {
				log.Println(errFinal)
			}
		} else {
			errFinal := tx.Commit()
			if errFinal != nil {
				log.Println(errFinal)
			}
		}
	}()
	result, rErr = f(tx)
	return result, rErr
}
