package models

type RegisterRequest struct {
	Usermame string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// DeleteMultipleProductsRequest represents the request body for deleting multiple products
type DeleteMultipleProductsRequest struct {
	IDs []int `json:"ids"`
}
