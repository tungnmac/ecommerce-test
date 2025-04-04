package main

import (
	"ecommerce-test/config"
	"ecommerce-test/internal/routes"
	"fmt"
	"log"
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
	// to change the flags on the default logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db := config.ConnectDB()
	config.InitJWT()

	// Initialize database schema and seed data
	// config.InitDatabase()

	router := routes.SetupRouter(db)

	fmt.Println("Server running on port 8080")
	router.Run(":8080") // Gin handles logging and errors internally
}
