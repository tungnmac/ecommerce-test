package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
)

// Define the structure of fake data
type Data struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	City  string `json:"city"`
}

// Generate fake data
func generateFakeData(batchSize int) []Data {
	var data []Data
	for i := 0; i < batchSize; i++ {
		item := Data{
			ID:    rand.Intn(1000000),
			Name:  faker.Name(),
			Email: faker.Email(),
			Phone: faker.Phonenumber(),
			// City:  faker.City(),
		}
		data = append(data, item)
	}
	return data
}

// API Endpoint
func getFakeData(c *gin.Context) {
	// Get `size` parameter from request (default to 1000)
	size, err := strconv.Atoi(c.DefaultQuery("size", "1000"))
	if err != nil || size < 1 {
		size = 1000
	}

	// Generate the data
	data := generateFakeData(size)

	// Return JSON response
	c.JSON(http.StatusOK, data)
}

func main() {
	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Create Gin router
	r := gin.Default()

	// Define the endpoint
	r.GET("/data", getFakeData)

	// Start the server
	r.Run(":8088") // Runs on http://localhost:8080
}
