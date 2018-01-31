package domain

import (
	// "github.com/jinzhu/gorm"
)

type Salereport struct {
	ID				int		`gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	OrderID			string	`gorm:"not null" json:"orderid"`
	Timestamp		string	`gorm:"not null" json:"timestamp"`
	Sku				string	`json:"sku"`
	Name			string	`json:"name"`
	Amount			int		`json:"amount"`
	Saleprice		int		`json:"saleprice"`
	Total			int		`json:"total"`
	Buyingprice		int		`json:"buyingprice"`
	Profit			int		`json:"profit"`
}

type ISalereport interface {
	CreateSaleReport(*Salereport)
	GetAllSaleReports() []Salereport
	GetSaleReportsBySKU(string) []Salereport
}