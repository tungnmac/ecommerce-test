package main

import (
	"go-mock/db"
	"go-mock/queue"

	"github.com/gin-gonic/gin"
)

const (
	apiURL    = "http://localhost:8088/data"
	batchSize = 10000
	totalRows = 1000000
)

func main() {
	db.InitDB()
	queue.InitQueue()
	go queue.StartWorker()
	go queue.StartScheduler(apiURL, batchSize, totalRows)

	// Create Gin router
	r := gin.Default()

	// Start the server
	r.Run(":8888") // Runs on http://localhost:8080
}
