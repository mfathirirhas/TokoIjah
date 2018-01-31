package domain

import (
	// "github.com/jinzhu/gorm"
)

type Stockout struct {
	ID				int		`gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Timestamp		string	`gorm:"not null" json:"timestamp"`
	Sku				string	`gorm:"not null" json:"sku"`
	Name			string	`json:"name"`
	OutAmount		int		`json:"outamount"`
	SalePrice		int		`json:"saleprice"`
	Total			int		`json:"total"`
	Note			string	`json:"note"`
}

type IStockout interface {
	RemoveProduct(*Stockout)
	GetAllOutProducts() []Stockout
	GetOutProductsBySku(string) []Stockout
}