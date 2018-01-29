package config

import (
	"database/sql"
	"github.com/mfathirirhas/TokoIjah/libs"
	
	// import sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbname = "tokoijah.db"
	mainDB *sql.DB
)

func Db() *sql.DB {

	db, err := sql.Open("sqlite3", dbname)
	libs.DbErr(err)
	mainDB = db

	return mainDB
}