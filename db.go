package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	return db, nil
}
