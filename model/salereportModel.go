package model

import (
	"github.com/mfathirirhas/TokoIjah/domain"
)

func (db *DB) CreateSaleReport(s *domain.Salereport) {
	db.Create(s)
}

func (db *DB) GetAllSaleReports() []domain.Salereport {
	var allSaleReports []domain.Salereport
	db.Find(&allSaleReports)
	return allSaleReports
}

func (db *DB) GetSaleReportsBySKU(sku string) []domain.Salereport {
	var saleReports []domain.Salereport
	db.Where("sku = ?", sku).Find(&saleReports)
	return saleReports
}