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

func StockvalueExportToCSV(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allstockvalue []domain.Stockvalue
		allstockvalue = db.GetAllStockValues()

		csvdata := init2dArray(len(allstockvalue), 10)

		for i:=0; i<len(allstockvalue); i++ {
			csvdata[i][0] = strconv.Itoa(i+1)
			csvdata[i][1] = allstockvalue[i].Sku
			csvdata[i][2] = allstockvalue[i].Name
			csvdata[i][3] = strconv.Itoa(allstockvalue[i].Amount)
			csvdata[i][4] = strconv.Itoa(allstockvalue[i].BuyingPrice)
			csvdata[i][5] = strconv.Itoa(allstockvalue[i].Total)
		}

		fileName := time.Now().Format("2006-02-01") + "-Stockvalue.csv"
		file, err := os.Create("./"+fileName)
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
			"message": "Stockvalue data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}