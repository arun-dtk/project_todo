package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.GET("/todos", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	server.Run(":8080") // localhost:8080
}
