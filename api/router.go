package api

import (
	// "net/http"
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
}

