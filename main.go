package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nehjoshi/fetch-challenge/handlers"
)

func main() {
	router := gin.Default()
	router.GET("/", handlers.HomeRoute)
	router.GET("/receipts/:id/points", handlers.GetPoints)
	router.POST("/receipts/process", handlers.ProcessReceipt)

	router.Run(":5000")
}
