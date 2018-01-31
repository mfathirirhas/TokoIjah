package api

import (
	// "log"
	"net/http"
	// "time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func CreateStockValue(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stockvalue domain.Stockvalue
		if gc.BindJSON(&stockvalue) == nil {
			stockvalue.Total = stockvalue.BuyingPrice * stockvalue.Amount
			db.CreateStockValue(&stockvalue)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stock value created successfully",
				"id": stockvalue.ID,
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

func GetAllStockValues(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var stockvalue []domain.Stockvalue
		stockvalue = db.GetAllStockValues()
		if len(stockvalue) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockvalue,
			}) 
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stockvalues is empty!",
			}) 
			return
		}
		
		return
	}
}

func GetStockValueByID(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var stockvalue domain.Stockvalue
		id,err := strconv.Atoi(gc.Param("id"))
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"message": "error!",
			})
			return
		}
		stockvalue = db.GetStockValueByID(id)
		if stockvalue.ID == id {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockvalue,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "Stockvalue not found!",
			})
			return
		}
	}
}


func GetStockValuesBySku(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var stockvalue domain.Stockvalue
		stockvalue = db.GetStockValuesBySku(gc.Param("sku"))
		if stockvalue.Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stockvalue,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "Stockvalue not found!",
			})
			return
		}
	}
}

func UpdateStockValue(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {
		
		var stockvalue domain.Stockvalue
		if gc.BindJSON(&stockvalue) == nil {
			updatedata := db.GetStockValuesBySku(stockvalue.Sku)
			updatedata.Name = stockvalue.Name
			updatedata.Amount = stockvalue.Amount
			updatedata.BuyingPrice = stockvalue.BuyingPrice
			updatedata.Total = stockvalue.BuyingPrice * stockvalue.Amount
			updated := db.UpdateStockValue(updatedata)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stockvalue updated successfully",
				"Updated Data": updated,
			})
			return
		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"message": "Check data carefully",
			})
			return
		}
	}
}