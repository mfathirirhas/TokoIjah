package domain

import (
	// "github.com/jinzhu/gorm"
)

type Stockin struct {
	ID				int		`gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Timestamp		string	`gorm:"not null" json:"timestamp"`
	Sku				string	`gorm:"not null" json:"sku"`
	Name			string	`json:"name"`
	OrderAmount		int		`json:"orderamount"`
	ReceivedAmount	int		`json:"receivedamount"`
	BuyingPrice		int		`json:"buyingprice"`
	Total			int		`json:"total"`
	Receipt			string	`json:"receipt"`
	Note			string	`json:"note"`
}

type IStockin interface {
	StoreProduct(*Stockin)
	GetAllStoredProducts() []Stockin
	GetStoredProductsBySku(string) []Stockin
}