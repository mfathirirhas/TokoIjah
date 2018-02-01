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

func CreateSaleReport(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var salereport domain.Salereport
		if gc.BindJSON(&salereport) == nil {
			salereport.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			salereport.Total = (salereport.Amount * salereport.Saleprice)
			salereport.Profit = (salereport.Amount * salereport.Saleprice) - (salereport.Amount * salereport.Buyingprice)
			db.CreateSaleReport(&salereport)
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"message": "Sale report created successfully!",
				"id": salereport.ID,
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