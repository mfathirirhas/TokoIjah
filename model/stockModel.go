package model

import (
	"github.com/mfathirirhas/TokoIjah/domain"
)


func (db *DB) CreateStock(s *domain.Stock) {
	db.Create(s)
}

func (db *DB) GetAllStock() []domain.Stock {
	var allStock []domain.Stock
	db.Find(&allStock)
	return allStock
}

func (db *DB) GetStockByID(ID int) domain.Stock {
	var stock domain.Stock
	db.First(&stock, ID)
	return stock
}

func (db *DB) GetStockBySku(sku string) domain.Stock {
	var stock domain.Stock
	db.Where("sku = ?", sku).First(&stock)
	return stock
}

func (db *DB) UpdateStock(s domain.Stock) domain.Stock {
	db.Save(s)
	return s
}
