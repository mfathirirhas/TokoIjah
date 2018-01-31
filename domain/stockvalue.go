package domain

import (
	// "github.com/jinzhu/gorm"
)

type Stockvalue struct {
	ID				int		`gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Sku				string	`gorm:"not null" json:"sku"`
	Name			string	`json:"name"`
	Amount			int		`json:"amount"`
	BuyingPrice		int		`json:"buyingprice"`
	Total			int		`json:"total"`
}

type IStockvalue interface {
	CreateStockValue(*Stockvalue)
	GetAllStockValues() []Stockvalue
	GetStockValueByID(int) Stockvalue
	GetStockValuesBySku(string) Stockvalue
	UpdateStockValue(Stockvalue) Stockvalue
}