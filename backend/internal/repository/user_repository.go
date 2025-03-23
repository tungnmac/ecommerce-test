package repository

import (
	"context"
	"ecommerce-test/config"
	"ecommerce-test/internal/models"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) UserRepository {
	return UserRepository{
		DB: db,
	}
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(user models.User) error {
	now := time.Now()
	user.CreatedAt = models.Date(now)
	user.UpdatedAt = models.Date(now)
	_, err := r.DB.Exec(context.Background(),
		"INSERT INTO test_users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		&user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	var createdAt, updatedAt time.Time

	err := r.DB.QueryRow(context.Background(),
		"SELECT id, name, email, password, created_at, updated_at FROM test_users WHERE email=$1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &createdAt, &updatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.CreatedAt = models.Date(createdAt)
	user.UpdatedAt = models.Date(updatedAt)
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, name, email, created_at, updated_at FROM test_users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		user.CreatedAt = models.Date(createdAt)
		user.UpdatedAt = models.Date(updatedAt)
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := config.DB.QueryRow(context.Background(), "SELECT id, name, email FROM test_users WHERE id=$1", id).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
