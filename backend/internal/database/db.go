package database

// import (
// 	"database/sql"
// 	"ecommerce-test/backend/config"
// 	"fmt"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func ConnectDB() {
// 	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		config.GetEnv("POSTGRES_HOST", "localhost"), config.GetEnv("POSTGRES_PORT", "5432"), config.GetEnv("POSTGRES_USER", "root"), config.GetEnv("POSTGRES_PASS", "postgres"), config.GetEnv("POSTGRES_DB", "default"))

// 	var err error
// 	DB, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %v", err)
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		log.Fatalf("Error connecting to the database: %v", err)
// 	}

// 	fmt.Println("Connected to PostgreSQL successfully!")
// }
