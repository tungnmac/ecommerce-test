package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var DB *pgx.Conn

func ConnectDB() *pgx.Conn {
	// Load environment variables from .env file
	err := godotenv.Load()
	log.Println(err)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get DB connection URL from environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	// Establish database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	DB = conn

	return DB
}

func CloseDB() {
	if DB != nil {
		DB.Close(context.Background())
		fmt.Println("🔌 Database connection closed")
	}
}
