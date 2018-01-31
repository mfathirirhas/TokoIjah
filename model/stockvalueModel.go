package model

import (
	"github.com/mfathirirhas/TokoIjah/domain"
)

func (db *DB) CreateStockValue(s *domain.Stockvalue) {
	db.Create(s)
}

func (db *DB) GetAllStockValues() []domain.Stockvalue {
	var allStockvalue []domain.Stockvalue
	db.Find(&allStockvalue)
	return allStockvalue
}

func (db *DB) GetStockValueByID(ID int) domain.Stockvalue {
	var stockvalue domain.Stockvalue
	db.First(&stockvalue, ID)
	return stockvalue
}

func (db *DB) GetStockValuesBySku(sku string) domain.Stockvalue {
	var stockvalue domain.Stockvalue
	db.Where("sku = ?", sku).First(&stockvalue)
	return stockvalue
}

func (db *DB) UpdateStockValue(s domain.Stockvalue) domain.Stockvalue {
	db.Save(s)
	return s
}