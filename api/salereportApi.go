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

func CreateSaleReport(db domain.ISalereport, dbstockvalue domain.IStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var salereport domain.Salereport
		var stockvalue domain.Stockvalue

		if gc.BindJSON(&salereport) == nil {
			salereport.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			stockvalue = dbstockvalue.GetStockValuesBySku(salereport.Sku)
			if stockvalue.Sku != "" {
				salereport.Buyingprice = stockvalue.BuyingPrice	
			} else {
				gc.JSON(http.StatusBadRequest, gin.H{
					"status": false,
					"message": "There is no such product in stockvalue!",
				})
				return
			}
			salereport.Total = (salereport.Amount * salereport.Saleprice)
			salereport.Profit = (salereport.Amount * salereport.Saleprice) - (salereport.Amount * salereport.Buyingprice)
			db.CreateSaleReport(&salereport)
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "Sale report is created successfully!",
				"id": salereport.ID,
				"sale": salereport,
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

func SalereportExportToCSV(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allsalereport []domain.Salereport
		allsalereport = db.GetAllSaleReports()

		csvdata := init2dArray(len(allsalereport), 10)

		for i:=0; i<len(allsalereport); i++ {
			csvdata[i][0] = strconv.Itoa(i+1)
			csvdata[i][1] = allsalereport[i].OrderID
			csvdata[i][2] = allsalereport[i].Timestamp
			csvdata[i][3] = allsalereport[i].Sku
			csvdata[i][4] = allsalereport[i].Name
			csvdata[i][5] = strconv.Itoa(allsalereport[i].Amount)
			csvdata[i][6] = strconv.Itoa(allsalereport[i].Saleprice)
			csvdata[i][7] = strconv.Itoa(allsalereport[i].Total)
			csvdata[i][8] = strconv.Itoa(allsalereport[i].Buyingprice)
			csvdata[i][9] = strconv.Itoa(allsalereport[i].Profit)
		}

		fileName := time.Now().Format("2006-01-02") + "-Salereport.csv"
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
			"message": "Salereport data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

func GetSaleReportsByDate(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var salereports []domain.Salereport
		salereports = db.GetSaleReportsByDate(gc.PostForm("datefrom"), gc.PostForm("dateto"))
		if len(salereports) > 0 {
			if gc.PostForm("exportcsv") == "1" {
				// export to csv
				csvdata := init2dArray(len(salereports), 10)
				for i:=0; i<len(salereports); i++ {
					csvdata[i][0] = strconv.Itoa(i+1)
					csvdata[i][1] = salereports[i].OrderID
					csvdata[i][2] = salereports[i].Timestamp
					csvdata[i][3] = salereports[i].Sku
					csvdata[i][4] = salereports[i].Name
					csvdata[i][5] = strconv.Itoa(salereports[i].Amount)
					csvdata[i][6] = strconv.Itoa(salereports[i].Saleprice)
					csvdata[i][7] = strconv.Itoa(salereports[i].Total)
					csvdata[i][8] = strconv.Itoa(salereports[i].Buyingprice)
					csvdata[i][9] = strconv.Itoa(salereports[i].Profit)
				}

				fileName := time.Now().Format("2006-01-02") + "-Salereport_"+gc.PostForm("datefrom")+"---"+ gc.PostForm("dateto")+".csv"
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
					"message": "Salereport exported successfully!",
					"filename": fileName,
				})
				return

			} else {
				// just print the json response
				omzet := 0
				grossprofit := 0
				totalsale := 0
				totalproducts := 0
				for i:=0; i<len(salereports); i++ {
					omzet += salereports[i].Total
					grossprofit += salereports[i].Profit
					if salereports[i].OrderID != "" {
						totalsale++ 
					}
					totalproducts += salereports[i].Amount
				}
				gc.JSON(http.StatusOK, gin.H{
					"status": true,
					"datereport": time.Now().Format("2006-01-02"),
					"daterange": gc.PostForm("datefrom") + " - " + gc.PostForm("dateto"),
					"omzet": omzet,
					"grossprofit": grossprofit,
					"totalsale": totalsale,
					"totalsoldproducts": totalproducts,
					"data": salereports,
				})
				return
			}
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": false,
				"message": "No sale reports by those dates!",
			})
			return
		}

		// in case nothing to be returned
		gc.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"message": "something wrong!",
			"data": salereports,
		})
		return
	}
}