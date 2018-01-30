package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status"	: "OK",
		"message" 	: "Homepage",	
	})
}