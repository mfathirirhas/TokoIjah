package model

import (
	"github.com/mfathirirhas/TokoIjah/domain"
)

func (db *DB) RemoveProduct(s *domain.Stockout) {
	db.Create(s)
}

func (db *DB) GetAllOutProducts() []domain.Stockout {
	var allStockout []domain.Stockout
	db.Find(&allStockout)
	return allStockout
}

func (db *DB) GetOutProductsBySku(sku string) []domain.Stockout {
	var stockout []domain.Stockout
	db.Where("sku = ?", sku).Find(&stockout)
	return stockout
}