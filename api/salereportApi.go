package api

import (
	"os"
	"net/http"
	"time"
	"io"
	"bufio"
	"strconv"
	"encoding/csv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mfathirirhas/TokoIjah/domain"
)

// SaleReportReqBody data structure for generate salereport json req body
type SaleReportReqBody struct {
	From		string	`json:"datefrom"`
	To			string	`json:"dateto"`
	Csvexport	string	`json:"exportcsv"`
}

// CreateSaleReport API to create one sale report instance and save it to salereports table
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

// GetAllSaleReports get all rows of salereports from salereports table
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

// GetSaleReportsBySKU get salereports by sku
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

// SalereportExportToCSV export salereports data into csv file
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

// GetSaleReportsByDate get salereports by range of date(timestamp)
func GetSaleReportsByDate(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var reqBody SaleReportReqBody
		var from string
		var to string
		var csvexport string
		decoder := json.NewDecoder(gc.Request.Body)
		err := decoder.Decode(&reqBody)
		from = reqBody.From
		to = reqBody.To
		csvexport = reqBody.Csvexport
		if err != nil {	
			from = "undefined"
			to = "undefined"
		}

		var salereports []domain.Salereport
		salereports = db.GetSaleReportsByDate(from, to)
		if len(salereports) > 0 {
			if csvexport == "1" {
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

				fileName := time.Now().Format("2006-01-02") + "-Salereport_"+ dateOnlyFormat(from) +"_"+ dateOnlyFormat(to) +".csv"
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
					"daterange": from + " - " + to,
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

// SalereportImportCSV import csv data into salereports table
func SalereportImportCSV(db domain.ISalereport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var salereport []domain.Salereport

		file, _ := gc.FormFile("salereportimport")
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

			salereportamount, _ := strconv.Atoi(line[5])
			salereportsaleprice, _ := strconv.Atoi(line[6])
			salereporttotal, _ := strconv.Atoi(line[7])
			salereportbuyingprice, _ := strconv.Atoi(line[8])
			salereportprofit, _ := strconv.Atoi(line[9])
			salereport = append(salereport, domain.Salereport{
				OrderID: line[1],
				Timestamp: line[2],  // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku: line[3],
				Name: line[4],
				Amount: salereportamount,
				Saleprice: salereportsaleprice,
				Total: salereporttotal,
				Buyingprice: salereportbuyingprice,
				Profit: salereportprofit,
			})
		}

		if len(salereport) > 0 {
			for i:=0; i<len(salereport); i++ {
				db.CreateSaleReport(&salereport[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status": true,
				"message": "data csv migrated successfully to salereport table",
				"data": salereport,
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

func dateOnlyFormat(t string) string {
    runes := []rune(t)
    return string(runes[0:10])
}