package api

import (
	// "log"
	"net/http"
	"time"
	// "strconv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func StoreProduct(db domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stockin domain.Stockin
		if gc.BindJSON(&stockin) == nil {
			stockin.Timestamp = time.Now().String()
			stockin.Total = stockin.OrderAmount * stockin.BuyingPrice
			db.StoreProduct(&stockin)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Products stored successfully",
				"id": stockin.ID,
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

func GetAllStoredProducts(db domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context){

		var stockin []domain.Stockin
		stockin = db.GetAllStoredProducts()
		if len(stockin) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockin,
			}) 
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Store records is empty!",
			}) 
			return
		}
		
		return
	}
}

func GetStoredProductsBySku(db domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var stockin []domain.Stockin
		stockin = db.GetStoredProductsBySku(gc.Param("sku"))
		if stockin[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockin,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "No records not found!",
			})
			return
		}
	}
}