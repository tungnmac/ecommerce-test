package models

type Product struct {
	ID               int     `json:"id"`
	ProductReference string  `json:"product_reference"`
	ProductName      string  `json:"product_name"`
	Status           string  `json:"status"`
	ProductCategory  string  `json:"product_category"`
	Price            float64 `json:"price"`
	StockLocation    string  `json:"stock_location"`
	Supplier         string  `json:"supplier"`
	Quantity         int     `json:"quantity"`
	CreatedAt        Date    `json:"created_at"`
	UpdatedAt        Date    `json:"updated_at"`
	// DateAdded        Date    `json:"date_added"`
}
