package api

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func CreateSaleReport(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var salereport domain.Salereport
		if gc.BindJSON(&salereport) == nil {
			salereport.Timestamp = time.Now().String()
			salereport.Total = (salereport.Amount * salereport.Saleprice)
			salereport.Profit = (salereport.Amount * salereport.Saleprice) - (salereport.Amount * salereport.Buyingprice)
			db.CreateSaleReport(&salereport)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Sale report created successfully!",
				"id": salereport.ID,
			})
			return
		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"message": "Bad Request",
			})
			return
		}
		return
	}
}

func GetAllSaleReports(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var salereport []domain.Salereport
		salereport = db.GetAllSaleReports()
		if len(salereport) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": salereport,
			}) 
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "No sale report yet!",
			}) 
			return
		}
		
		return
	}
}

func GetSaleReportsBySKU(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var salereport []domain.Salereport
		salereport = db.GetSaleReportsBySKU(gc.Param("sku"))
		if salereport[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": salereport,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "No salereports not found!",
			})
			return
		}
	}
}