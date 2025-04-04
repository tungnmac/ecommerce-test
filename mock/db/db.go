package db

import (
	"context"
	"log"
	"time"

	"go-mock/request"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func InitDB() {
	var err error
	dbPool, err = pgxpool.New(context.Background(), "postgres://root:postgres@localhost:5432/pod_io_dev")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
}

// Batch insert data into PostgreSQL
func InsertBatch(data []request.Data) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Println("Transaction begin error:", err)
		return
	}

	stmt := `INSERT INTO records (id, name, email, phone) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`

	for _, d := range data {
		_, err := tx.Exec(ctx, stmt, d.ID, d.Name, d.Email, d.Phone)
		if err != nil {
			log.Println("Insert error:", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit error:", err)
	}
	log.Println("Transaction commit success!")
}
