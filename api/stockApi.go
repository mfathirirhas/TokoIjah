package api

import (
	"net/http"
	"strconv"
	"time"
	"os"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func CreateStock(db domain.IStock) gin.HandlerFunc{
	return func(gc *gin.Context) {

		var stock domain.Stock
		if gc.BindJSON(&stock) == nil {
			db.CreateStock(&stock)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stock data created successfully!",
				"id": stock.ID,
			})
			return
		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"message": "Bad Request",
			})
			return
		}
		return
	}
}

func GetAllStock(db domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock []domain.Stock
		stock = db.GetAllStock()
		if len(stock) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stock,
			}) 
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stock is empty!",
			}) 
			return
		}
		
		return
	}
}

func GetStockByID(db domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock domain.Stock
		id,err := strconv.Atoi(gc.Param("id"))
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": "false",
				"message": "error!",
			})
			return
		}

		stock = db.GetStockByID(id)
		if stock.ID == id {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stock,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "Stock not found!",
			})
			return
		}
	}
}

func GetStockBySku(db domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock domain.Stock
		stock = db.GetStockBySku(gc.Param("sku"))
		if stock.Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data": stock,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data": "Stock not found!",
			})
			return
		}
	}
}

func UpdateStock(db domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock domain.Stock
		if gc.BindJSON(&stock) == nil {
			updatedata := db.GetStockBySku(stock.Sku)
			updatedata.Name = stock.Name
			updatedata.Amount = stock.Amount
			updated := db.UpdateStock(updatedata)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Stock updated successfully",
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

func StockExportToCSV(db domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allstock []domain.Stock
		allstock = db.GetAllStock()

		csvdata := init2dArray(len(allstock), 4)

		for i:=0; i<len(allstock); i++ {
			csvdata[i][0] = strconv.Itoa(i+1)
			csvdata[i][1] = allstock[i].Sku
			csvdata[i][2] = allstock[i].Name
			csvdata[i][3] = strconv.Itoa(allstock[i].Amount)
		}

		fileName := time.Now().Format("2006-01-02") + "-Stock.csv"
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
			"message": "Stock data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

func init2dArray(y int, x int) [][]string {
	a := make([][]string, y)
	for i := range a {
		a[i] = make([]string, x)
	}
	return a
}