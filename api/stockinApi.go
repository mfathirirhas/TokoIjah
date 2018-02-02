package api

import (
	"os"
	"net/http"
	"time"
	"strconv"
	"io"
	"bufio"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

// StoreProduct make note for new products imported into stock
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

// GetAllStoredProducts get all data records from stockin table 
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

// GetStoredProductsBySku get records data by sku
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

// StockinExportToCSV export all records from stockins table into csv file
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

// StockinImportCSV import csv data into stockins table
func StockinImportCSV(db domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stockin []domain.Stockin

		file, _ := gc.FormFile("stockinimport")
		dst := "./csv/"+ file.Filename
		gc.SaveUploadedFile(file, dst)
		// csvfile, err := os.Open("./csv/import_stock.csv")
		csvfile, err := os.Open("./csv/"+file.Filename)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"message": "error opening file, check file again",
			})
		}


		reader := csv.NewReader(bufio.NewReader(csvfile))
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				gc.JSON(http.StatusBadRequest, gin.H{
					"status": false,
					"message": "something's wrong!",
				})
			}

			stockinorderamount, _ := strconv.Atoi(line[4])
			stockinreceivedamount, _ := strconv.Atoi(line[5])    
    		stockinbuyingprice, _ := strconv.Atoi(line[6])
    		stockintotal, _ := strconv.Atoi(line[7])
			stockin = append(stockin, domain.Stockin{
				Timestamp: line[1],  // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku: line[2],
				Name: line[3],
				OrderAmount: stockinorderamount,
				ReceivedAmount: stockinreceivedamount,
				BuyingPrice: stockinbuyingprice,
				Total: stockintotal,
				Receipt: line[8],
				Note: line[9],
			})
		}

		if len(stockin) > 0 {
			for i:=0; i<len(stockin); i++ {
				db.StoreProduct(&stockin[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "data csv migrated successfully to stock table",
				"data": stockin,
			})
			return

		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status": false,
				"message": "error reading csv file, check file again for correct format!",
			})
			return

		}

		gc.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"message": "something's wrong!",
		})
		return

	}
}