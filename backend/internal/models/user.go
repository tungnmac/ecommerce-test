package models

// User model
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	CreatedAt Date   `json:"created_at"`
	UpdatedAt Date   `json:"updated_at"`
}
