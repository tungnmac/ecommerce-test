package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// Load and execute init.sql
func InitDatabase() {
	sqlFile := "scripts/init.sql"

	// Read SQL file
	data, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	queries := strings.Split(string(data), ";")

	log.Println(queries)

	// Execute each SQL command
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := DB.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Error executing query: %v", err)
		}
	}

	fmt.Println("✅ Database initialized successfully!")
}
