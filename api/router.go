package api

import (
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/mfathirirhas/TokoIjah/model"
)


func InitRouter(db *model.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	setRoutes(db, router)

	return router
}

func setRoutes(db *model.DB, r *gin.Engine) {
	r.GET("/", home)

	// stock apis
	r.POST("/stock", CreateStock(db))
	r.GET("/stock", GetAllStock(db))
	r.GET("/stockbyid/:id", GetStockByID(db))
	r.GET("/stockbysku/:sku", GetStockBySku(db))
	r.POST("/stockupdate", UpdateStock(db))

	// stockin apis
	r.POST("/stockin", StoreProduct(db))
	r.GET("/stockin", GetAllStoredProducts(db))
	r.GET("/stockinbysku/:sku", GetStoredProductsBySku(db))

	// stockout apis
	r.POST("/stockout", RemoveProduct(db))
	r.GET("/stockout", GetAllOutProducts(db))
	r.GET("/stockoutbysku/:sku", GetOutProductsBySku(db))


	// salereport apis
	r.POST("/salereport", CreateSaleReport(db))
	r.GET("/salereport", GetAllSaleReports(db))
	r.GET("/salereportbysku/:sku", GetSaleReportsBySKU(db))

	// stockvalue apis
	r.POST("/stockvalue", CreateStockValue(db))
	r.GET("/stockvalue", GetAllStockValues(db))
	r.GET("/stockvaluebyid/:id", GetStockValueByID(db))
	r.GET("/stockvaluebysku/:sku", GetStockValuesBySku(db))
	r.POST("/stockvalueupdate", UpdateStockValue(db))

	// export to csv
	r.GET("/stockexport", StockExportToCSV(db))
}

