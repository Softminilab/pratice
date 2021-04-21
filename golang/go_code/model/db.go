package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	db *sql.DB
)

func InitDB(dsn string) error {
	var (
		err error
	)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return errors.Wrap(err, "mysql db open failed")
	}

	err = db.Ping()
	if err != nil {
		return errors.Wrap(err, "mysql PING is failed")
	}
	return err
}

func Close() error {
	err := db.Close()
	if err != nil {
		return errors.Wrap(err, "mysql DB Close is failed")
	}
	return err
}