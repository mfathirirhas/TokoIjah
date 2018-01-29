package config

import (
	"database/sql"
	"log"
	
	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbname = "tokoijah.db"
	mainDB *sql.DB
)

func Db() *sql.DB {

	db, err := sql.Open("sqlite3", dbname)
	dbErr(err)
	mainDB = db

	return mainDB
}

func dbErr(err error) {
	if err != nil {
		log.Fatal("database error: %s",err)
	}	
}