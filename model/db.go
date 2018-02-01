package model

import (
	"log"
	"github.com/jinzhu/gorm"
	"github.com/mfathirirhas/TokoIjah/domain"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbEngine	= "sqlite3"
	dbName		= "./test.db"
)

type DB struct {
	*gorm.DB
}

//NewDB initializes the database
func InitDB() *DB {
	db, err := gorm.Open(dbEngine, dbName)
	if err != nil {
		log.Fatal("failed to initialize database: ",err)
	}

	db.AutoMigrate(&domain.Stock{}, &domain.Stockin{}, &domain.Stockout{}, &domain.Salereport{}, &domain.Stockvalue{})

	return &DB{db}
}
