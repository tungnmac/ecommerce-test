package services

import (
	"ecommerce-test/internal/models"
	"ecommerce-test/internal/repository"
	"errors"
)

type ProductService struct {
	Repo repository.ProductRepository
}

func (s *ProductService) CreateProduct(product models.Product) error {
	// Check if product already exists
	existingProduct, _ := s.Repo.GetProductByReference(product.ProductReference)
	if existingProduct != nil {
		return errors.New("product with reference already exists")
	}

	// Save product to database
	return s.Repo.InsertProduct(product)
}

// GetPaginatedProducts retrieves paginated products
func (s *ProductService) GetPaginatedProducts(limit, offset int) ([]models.Product, error) {
	return s.Repo.GetProducts(limit, offset)
}

func (s *ProductService) GetProductByReference(reference string) (*models.Product, error) {
	return s.Repo.GetProductByReference(reference)
}

// Update an existing product
func (s *ProductService) UpdateProduct(product models.Product) error {
	return s.Repo.UpdateProduct(product)
}

// GetFilteredProducts returns products based on filters
func (s *ProductService) GetFilteredProducts(filters map[string][]string, limit, offset int) ([]models.Product, error) {
	return s.Repo.GetFilteredProducts(filters, limit, offset)
}

// GetProductCategoryStats returns product counts by category and supplier
func (s *ProductService) GetProductCategoryStatistics() (map[string]interface{}, error) {
	categoryStats, err := s.Repo.GetProductCategoryStats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"category_distribution": categoryStats,
	}, nil
}

// GetProductSupplierStatistics returns product counts by category and supplier
func (s *ProductService) GetProductSupplierStatistics() (map[string]interface{}, error) {

	supplierStats, err := s.Repo.GetProductSupplierStats()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"supplier_distribution": supplierStats,
	}, nil
}
