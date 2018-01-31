package model

import (
	"github.com/mfathirirhas/TokoIjah/domain"
)


func (db *DB) StoreProduct(s *domain.Stockin) {
	db.Create(s)
}

func (db *DB) GetAllStoredProducts() []domain.Stockin {
	var allStoredProducts []domain.Stockin
	db.Find(&allStoredProducts)
	return allStoredProducts
}

func (db *DB) GetStoredProductsBySku(sku string) []domain.Stockin {
	var storedProduct []domain.Stockin
	db.Where("sku = ?", sku).Find(&storedProduct)
	return storedProduct
}