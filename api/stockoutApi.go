package api

import (
	// "log"
	"net/http"
	"time"
	// "strconv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func RemoveProduct(db domain.IStockout) gin.HandlerFunc {
	return func(gc *gin.Context){

		var stockout domain.Stockout
		if gc.BindJSON(&stockout) == nil {
			stockout.Timestamp = time.Now().String()
			stockout.Total = stockout.OutAmount * stockout.SalePrice
			db.RemoveProduct(&stockout)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Products stored successfully",
				"id": stockout.ID,
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

func GetAllOutProducts(db domain.IStockout) gin.HandlerFunc {
	return func(gc *gin.Context){

		var stockout []domain.Stockout
		stockout = db.GetAllOutProducts()
		if len(stockout) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockout,
			}) 
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "No records of products out from store yet!",
			}) 
			return
		}
		
		return
	}
}

func GetOutProductsBySku(db domain.IStockout) gin.HandlerFunc {
	return func(gc *gin.Context){

		var stockout []domain.Stockout
		stockout = db.GetOutProductsBySku(gc.Param("sku"))
		if stockout[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockout,
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