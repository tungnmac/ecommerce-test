package main

import (
	"ecommerce-test/config"
	"ecommerce-test/internal/routes"
	"fmt"
)

// @title Golang Ecommerce Test API
// @version 1.0
// @description API documentation for the Golang Ecommerce Test
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db := config.ConnectDB()
	config.InitJWT()

	// Initialize database schema and seed data
	// config.InitDatabase()

	router := routes.SetupRouter(db)

	fmt.Println("Server running on port 8080")
	router.Run(":8080") // Gin handles logging and errors internally
}
