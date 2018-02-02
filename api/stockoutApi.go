package api

import (
	"os"
	"net/http"
	"time"
	"io"
	"bufio"
	"strconv"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

func RemoveProduct(db domain.IStockout, dbStock domain.IStock) gin.HandlerFunc {
	return func(gc *gin.Context){

		var stock domain.Stock
		var stockout domain.Stockout

		if gc.BindJSON(&stockout) == nil {
			stockout.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			stockout.Total = stockout.OutAmount * stockout.SalePrice
			
			stock = dbStock.GetStockBySku(stockout.Sku)
			if stock.Sku != "" { // if already existed before, update the stock
				db.RemoveProduct(&stockout)
				stock.Amount -= stockout.OutAmount
				updatedStock := dbStock.UpdateStock(stock)
				gc.JSON(http.StatusOK, gin.H{
					"status": "true",
					"message": "Products removed successfully",
					"id": stockout.ID,
					"stock": updatedStock.Amount,
				})
				return
			} else { // if never exist before, then this product was never in stock before
				gc.JSON(http.StatusBadRequest, gin.H{
					"status": "false",
					"message": "Products were never in stock before!",
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

func StockoutExportToCSV(db domain.IStockout) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allstockout []domain.Stockout
		allstockout = db.GetAllOutProducts()

		csvdata := init2dArray(len(allstockout), 10)

		for i:=0; i<len(allstockout); i++ {
			csvdata[i][0] = strconv.Itoa(i+1)
			csvdata[i][1] = allstockout[i].Timestamp
			csvdata[i][2] = allstockout[i].Sku
			csvdata[i][3] = allstockout[i].Name
			csvdata[i][4] = strconv.Itoa(allstockout[i].OutAmount)
			csvdata[i][5] = strconv.Itoa(allstockout[i].SalePrice)
			csvdata[i][6] = strconv.Itoa(allstockout[i].Total)
			csvdata[i][7] = allstockout[i].Note
		}

		fileName := time.Now().Format("2006-01-02") + "-Stockout.csv"
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
			"message": "Stockout data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

func StockoutImportCSV(db domain.IStockout) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stockout []domain.Stockout

		file, _ := gc.FormFile("stockoutimport")
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

			stockoutoutamount, _ := strconv.Atoi(line[4])
			stockoutsaleprice, _ := strconv.Atoi(line[5])    
    		stockouttotal, _ := strconv.Atoi(line[6])
			stockout = append(stockout, domain.Stockout{
				Timestamp: line[1],  // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku: line[2],
				Name: line[3],
				OutAmount: stockoutoutamount,
				SalePrice: stockoutsaleprice,
				Total: stockouttotal,
				Note: line[7],
			})
		}

		if len(stockout) > 0 {
			for i:=0; i<len(stockout); i++ {
				db.RemoveProduct(&stockout[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "data csv migrated successfully to stockout table",
				"data": stockout,
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