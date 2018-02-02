package api

import (
	"os"
	"net/http"
	"time"
	"strconv"
	"math"
	"encoding/csv"
	"io"
	"bufio"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

// CreateStockValue create one instance of stock value into stockvalues table
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

// GetAllStockValues get all records from stockvalues table
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

// GetStockValueByID get a stockvalue data by id
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

// GetStockValuesBySku get a stock value data by sku
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

// UpdateStockValue update an already existing stock value data 
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

// StockvalueExportToCSV export all records from stockvalues table
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

		fileName := time.Now().Format("2006-01-02") + "-Stockvalue.csv"
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
			"message": "Stockvalue data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

// GenerateStockValue generate report from stock value in form of json data or csv
func GenerateStockValue(db domain.IStockvalue, dbstock domain.IStock, dbstockin domain.IStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock []domain.Stock
		var stockin []domain.Stockin
		var tempStockValue domain.Stockvalue

		stock = dbstock.GetAllStock() // get all products from actual stock
		if len(stock) > 0 {
			for i:=0; i<len(stock); i++ { // loop through products in stock
				var stockvalue domain.Stockvalue // if there're products in stock, initialize stockvalue
				stockvalue.Sku = stock[i].Sku
				stockvalue.Name = stock[i].Name
				stockvalue.Amount = stock[i].Amount
				stockin = dbstockin.GetStoredProductsBySku(stockvalue.Sku) // look stockin records for each products in stock by sku to calculate average price
				if len(stockin) > 0 { // if they exist in stockin records, then calculate average price
					sumTotal := 0 // sum of products total price in stockin records by sku
					sumReceivedAmount := 0 // sum of products receivedamount in stockin records by sku
					for j:=0; j<len(stockin); j++ {
						sumTotal += stockin[j].Total
						sumReceivedAmount += stockin[j].ReceivedAmount
					}
					averageValue := float64(sumTotal) / float64(sumReceivedAmount) // stockvalue average price = sum of stockin total / sum of stockin received amount
					stockvalue.BuyingPrice = int(Round(averageValue, .5, 0))
					stockvalue.Total = stockvalue.Amount * stockvalue.BuyingPrice

					// check for existing products inside stockvalues, if exist update the amount, buying price, and total. If not exist, create new one
					tempStockValue = db.GetStockValuesBySku(stockvalue.Sku)
					if tempStockValue.Sku != "" { // already exist in stockvalue table, update amount, price and total
						tempStockValue.Amount = stockvalue.Amount
						tempStockValue.BuyingPrice = stockvalue.BuyingPrice
						tempStockValue.Total = stockvalue.Total
						updatedStockValues := db.UpdateStockValue(tempStockValue)
						_ = updatedStockValues
					} else { // not exist and never recorded in stockvalue before, then create new one
						db.CreateStockValue(&stockvalue)
					}
					
				} else { // if the products by that sku dont exist in stockin, then the products never recorded into stockin
					gc.JSON(http.StatusBadRequest, gin.H{
						"status": false,
						"message": "Products in stock were never in Stockin records!",
					})
					return
				} //end if

			} // end loop stock
			// if for loop in stock for stock value finish, calculating value of stock is done, then calculate report
			var allStockValue []domain.Stockvalue
			allStockValue = db.GetAllStockValues()
			sumOfStockValueAmounts := 0
			sumOfStockValueTotals := 0
			for k:=0; k<len(allStockValue); k++ {
				sumOfStockValueAmounts += allStockValue[k].Amount
				sumOfStockValueTotals += allStockValue[k].Total
			}
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "Calculating stock values is done!",
				"Date": time.Now().Format("2006-01-02"),
				"Total SKU": len(allStockValue),
				"Total of Products": sumOfStockValueAmounts,
				"Total Value": sumOfStockValueTotals,
				"Stock Values": allStockValue,
			})
			return
		} else {
			gc.JSON(http.StatusNoContent, gin.H{
				"status": false,
				"message": "Stock is empty!",
			})
			return
		}

		return
	}
}

// StockvalueImportCSV import csv data into stockvalues table
func StockvalueImportCSV(db domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stockvalue []domain.Stockvalue

		file, _ := gc.FormFile("stockvalueimport")
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

			stockvalueamount, _ := strconv.Atoi(line[3])
			stockvaluebuyingprice, _ := strconv.Atoi(line[4])
			stockvaluetotal, _ := strconv.Atoi(line[5])
			stockvalue = append(stockvalue, domain.Stockvalue{
				Sku: line[1],  // start from sku column as we ignore id column(assume the csv include the IDs)
				Name: line[2],
				Amount: stockvalueamount,
				BuyingPrice: stockvaluebuyingprice,
				Total: stockvaluetotal,
			})
		}

		if len(stockvalue) > 0 {
			for i:=0; i<len(stockvalue); i++ {
				var tableStockvalue domain.Stockvalue
				tableStockvalue = db.GetStockValuesBySku(stockvalue[i].Sku)
				if tableStockvalue.Sku != "" { // data already exist in stock table, update the data then
					tableStockvalue.Name = stockvalue[i].Name
					tableStockvalue.Amount += stockvalue[i].Amount
					tableStockvalue.BuyingPrice = stockvalue[i].BuyingPrice
					tableStockvalue.Total = stockvalue[i].Total
					updatedStockvalue := db.UpdateStockValue(tableStockvalue)
					_ = updatedStockvalue // LOL
				} else { // data didn't exist in table, create new one then
					db.CreateStockValue(&stockvalue[i])
				}
			}
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "data csv migrated successfully to stockvalue table",
				"data": stockvalue,
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

// Round Go convert float to int without checking the round up/down value, it will always round down, hence we need this function
func Round(val float64, roundOn float64, places int ) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}