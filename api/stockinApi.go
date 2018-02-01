package api

import (
	"os"
	"net/http"
	"time"
	"strconv"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func StoreProduct(db domain.IStockin, dbStock domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock domain.Stock
		var stockin domain.Stockin
		
		if gc.BindJSON(&stockin) == nil {
			stockin.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			stockin.Total = stockin.OrderAmount * stockin.BuyingPrice
			db.StoreProduct(&stockin)

			stock = dbStock.GetStockBySku(stockin.Sku)
			if stock.Sku != "" { // if product already existed, update the amount
				stock.Amount += stockin.ReceivedAmount
				updatedStock := dbStock.UpdateStock(stock)
				gc.JSON(http.StatusOK, gin.H{
					"status": "true",
					"message": "Products stored successfully",
					"id": stockin.ID,
					"stock": updatedStock.Amount,
				})
				return
			} else { // if never exist before, create new stock
				stock.Sku = stockin.Sku
				stock.Name = stockin.Name
				stock.Amount = stockin.ReceivedAmount
				dbStock.CreateStock(&stock)
				gc.JSON(http.StatusOK, gin.H{
					"status": "true",
					"message": "New Products stored successfully",
					"id": stockin.ID,
					"stock": stock.Amount,
				})
				return
			}

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

func StockinExportToCSV(db domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allstockin []domain.Stockin
		allstockin = db.GetAllStoredProducts()

		csvdata := init2dArray(len(allstockin), 10)

		for i:=0; i<len(allstockin); i++ {
			csvdata[i][0] = strconv.Itoa(i+1)
			csvdata[i][1] = allstockin[i].Timestamp
			csvdata[i][2] = allstockin[i].Sku
			csvdata[i][3] = allstockin[i].Name
			csvdata[i][4] = strconv.Itoa(allstockin[i].OrderAmount)
			csvdata[i][5] = strconv.Itoa(allstockin[i].ReceivedAmount)
			csvdata[i][6] = strconv.Itoa(allstockin[i].BuyingPrice)
			csvdata[i][7] = strconv.Itoa(allstockin[i].Total)
			csvdata[i][8] = allstockin[i].Receipt
			csvdata[i][9] = allstockin[i].Note
		}

		fileName := time.Now().Format("2006-01-02") + "-Stockin.csv"
		file, err := os.Create("./csv/"+fileName)
		if err != nil {
			gc.JSON(http.StatusConflict, gin.H{
				"status": false,
				"message": "Failed to export file!",
			})
			return
		}
    	defer file.Close()

    	writer := csv.NewWriter(file)
    	defer writer.Flush()

    	for _, value := range csvdata {
        	err := writer.Write(value)
        	if err != nil {
				gc.JSON(http.StatusConflict, gin.H{
					"status": false,
					"message": "Failed to export file!",
				})
				return
			}
		}
		
		gc.JSON(http.StatusOK, gin.H{
			"status": true,
			"message": "Stockin data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}