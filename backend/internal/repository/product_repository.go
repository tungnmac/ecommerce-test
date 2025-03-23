package repository

import (
	"context"
	"ecommerce-test/config"
	"ecommerce-test/internal/models"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type ProductRepository struct{}

// InsertProduct inserts a new product with an auto-generated product_reference
func (r *ProductRepository) InsertProduct(product models.Product) (models.Product, error) {
	currentTime := time.Now()
	yearMonth := currentTime.Format("200601") // Format: YYYYMM

	// Get last product reference for the current month
	var lastProductRef string
	err := config.DB.QueryRow(context.Background(),
		"SELECT product_reference FROM test_products WHERE product_reference LIKE $1 ORDER BY product_reference DESC LIMIT 1",
		fmt.Sprintf("PROD-%s-%%", yearMonth)).Scan(&lastProductRef)

	log.Println(lastProductRef)
	// Extract sequence number
	var newSeq int
	if err != nil {
		// No previous product found for this month, start at 1
		newSeq = 1
	} else {
		// Extract last sequence number (PROD-YYYYMM-XXX)
		fmt.Sscanf(lastProductRef, "PROD-%s-%03d", &yearMonth, &newSeq)
		newSeq++
	}

	log.Println(fmt.Sprintf("PROD-%s-%03d", yearMonth, newSeq))

	// Generate new product reference
	product.ProductReference = fmt.Sprintf("PROD-%s-%03d", yearMonth, newSeq)
	product.CreatedAt = models.Date(currentTime)

	// Insert new product
	// _, err = config.DB.Exec(context.Background(),
	// 	"INSERT INTO test_products (product_reference, product_name, created_at, status, product_category, price, stock_location, supplier, quantity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
	// 	product.ProductReference, product.ProductName, product.CreatedAt, product.Status, product.ProductCategory, product.Price, product.StockLocation, product.Supplier, product.Quantity)

	return product, err
}

// GetProducts retrieves products with pagination
func (r *ProductRepository) GetProducts(limit, offset int) ([]models.Product, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT id, product_reference, product_name, created_at, status, product_category, price, stock_location FROM test_products ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		var createdAt time.Time
		err := rows.Scan(&product.ID, &product.ProductReference, &product.ProductName, &createdAt, &product.Status, &product.ProductCategory, &product.Price, &product.StockLocation)
		if err != nil {
			return nil, err
		}
		product.CreatedAt = models.Date(createdAt)
		products = append(products, product)
	}
	return products, nil
}

// GetProductByReference retrieves a single product by product_reference
func (r *ProductRepository) GetProductByReference(productReference string) (*models.Product, error) {
	var product models.Product

	query := `
		SELECT product_reference, product_name, created_at, status, product_category, price, 
		       stock_location, supplier, quantity 
		FROM test_products WHERE product_reference = $1
	`

	row := config.DB.QueryRow(context.Background(), query, productReference)
	err := row.Scan(
		&product.ProductReference, &product.ProductName, &product.CreatedAt, &product.Status,
		&product.ProductCategory, &product.Price, &product.StockLocation, &product.Supplier, &product.Quantity,
	)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return &product, nil
}

// UpdateProduct updates an existing product based on `product_reference`
func (r *ProductRepository) UpdateProduct(product models.Product) error {
	_, err := config.DB.Exec(context.Background(),
		`UPDATE test_products 
		SET product_name = $1, status = $3, product_category = $4, price = $5, 
		    stock_location = $6, supplier = $7, quantity = $8
		WHERE product_reference = $9`,
		product.ProductName, product.Status, product.ProductCategory, product.Price,
		product.StockLocation, product.Supplier, product.Quantity, product.ProductReference)

	return err
}

// GetFilteredProducts retrieves products based on dynamic filters
func (r *ProductRepository) GetFilteredProducts(filters map[string][]string, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply filters dynamically
	for key, values := range filters {
		if key == "min_price" {
			conditions = append(conditions, fmt.Sprintf("price >= $%d", argIndex))
			args = append(args, values[0])
			argIndex++
		} else if key == "max_price" {
			conditions = append(conditions, fmt.Sprintf("price <= $%d", argIndex))
			args = append(args, values[0])
			argIndex++
		} else if key == "product_name" {
			conditions = append(conditions, fmt.Sprintf("product_name ILIKE $%d", argIndex)) // Case-insensitive search
			args = append(args, "%"+values[0]+"%")                                           // Partial match
			argIndex++
		} else if key == "status" || key == "product_category" {
			// Handle multiple values for status and category
			placeholders := []string{}
			for _, value := range values {
				placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
				args = append(args, value)
				argIndex++
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = $%d", key, argIndex))
			args = append(args, values[0])
			argIndex++
		}
	}

	// Build the query dynamically
	query := "SELECT id, product_reference, product_name, created_at, updated_at, status, product_category, price, stock_location, supplier, quantity FROM test_products"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	// Execute the query
	rows, err := config.DB.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		var createdAt time.Time
		var updatedAt time.Time
		err := rows.Scan(
			&product.ID, &product.ProductReference, &product.ProductName, &createdAt, &updatedAt, &product.Status,
			&product.ProductCategory, &product.Price, &product.StockLocation, &product.Supplier, &product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		product.CreatedAt = models.Date(createdAt)
		product.UpdatedAt = models.Date(updatedAt)
		products = append(products, product)
	}

	return products, nil
}

// ProductCategoryStats holds category statistics
type ProductCategoryStats struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// ProductSupplierStats holds supplier statistics
type ProductSupplierStats struct {
	Supplier string `json:"supplier"`
	Count    int    `json:"count"`
}

// GetProductCategoryStats retrieves the count of products per category
func (r *ProductRepository) GetProductCategoryStats() ([]ProductCategoryStats, error) {
	query := "SELECT product_category, COUNT(*) FROM test_products GROUP BY product_category"
	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []ProductCategoryStats
	for rows.Next() {
		var stat ProductCategoryStats
		if err := rows.Scan(&stat.Category, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// GetProductSupplierStats retrieves the count of products per supplier
func (r *ProductRepository) GetProductSupplierStats() ([]ProductSupplierStats, error) {
	query := "SELECT supplier, COUNT(*) FROM test_products GROUP BY supplier"
	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []ProductSupplierStats
	for rows.Next() {
		var stat ProductSupplierStats
		if err := rows.Scan(&stat.Supplier, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// DeleteProductByID deletes a product by ID
func (r *ProductRepository) DeleteProductByID(id int) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM test_products
		WHERE id = $1`,
		id)

	return err
}

// DeleteProductByReference deletes a product by product_reference
func (r *ProductRepository) DeleteProductByReference(reference string) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM test_products
		WHERE product_reference = $1`,
		reference)

	return err
}

// DeleteMultipleProducts deletes multiple products by ID
func (r *ProductRepository) DeleteMultipleProducts(ids []int) error {
	_, err := config.DB.Exec(context.Background(),
		`DELETE FROM test_products
		WHERE id IN $1`,
		ids)

	return err
}
