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

// RegisterUserRequest represents the user registration request
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32,password"`
}

// CreateProductRequest represents the request for creating a product
type CreateProductRequest struct {
	ProductName     string  `json:"product_name" validate:"required,min=3,max=100"`
	Status          string  `json:"status" validate:"required,oneof=active inactive"`
	ProductCategory string  `json:"product_category" validate:"required"`
	Price           float64 `json:"price" validate:"required,gt=0"`
	StockLocation   string  `json:"stock_location" validate:"required"`
	Supplier        string  `json:"supplier" validate:"required"`
	Quantity        int     `json:"quantity" validate:"required,min=1"`
}
