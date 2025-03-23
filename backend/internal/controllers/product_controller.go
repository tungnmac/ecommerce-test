package controllers

import (
	"ecommerce-test/internal/models"
	"ecommerce-test/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jung-kurt/gofpdf"
	"github.com/oschwald/geoip2-golang"
	"github.com/umahmood/haversine"
)

type ProductController struct {
	Service services.ProductService
}

// CreateProduct godoc
// @Summary Add a new product
// @Description Creates a new product with auto-generated product_reference
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data (product_reference is auto-generated)"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paramProduct models.Product
	data, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	json.Unmarshal(data, &paramProduct)

	product, err := c.Service.CreateProduct(paramProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product added successfully", "product_reference": product.ProductReference})
}

// GetProductByReference godoc
// @Summary Get product by reference
// @Description Retrieves a product by its reference
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param reference path string true "Product Reference"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /api/products/{reference} [get]
func (c *ProductController) GetProductByReference(ctx *gin.Context) {
	reference := ctx.Param("product_reference")

	product, err := c.Service.GetProductByReference(reference)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Updates product details based on product_reference
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_reference path string true "Product Reference"
// @Param product body models.ParamProduct true "Updated Product Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/products/{product_reference} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	productRef := ctx.Param("product_reference")

	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	product.ProductReference = productRef // Ensure correct reference is used

	err := c.Service.UpdateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// GetFilteredProducts godoc
// @Summary Retrieve filtered products
// @Description Dynamic filtering and search for products, including multiple statuses and categories
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_reference query string false "Filter by product_reference"
// @Param product_name query string false "Search by product name"
// @Param status query string false "Filter by status (comma-separated values)"
// @Param product_category query string false "Filter by category (comma-separated values)"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param stock_location query string false "Filter by stock location"
// @Param supplier query string false "Filter by supplier"
// @Param limit query int false "Number of results per page" default(10)
// @Param offset query int false "Pagination offset" default(0)
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /api/products [get]
func (c *ProductController) GetFilteredProducts(ctx *gin.Context) {
	filters := map[string][]string{}

	// Extract query parameters
	if value := ctx.Query("product_reference"); value != "" {
		filters["product_reference"] = []string{value}
	}
	if value := ctx.Query("product_name"); value != "" {
		filters["product_name"] = []string{value}
	}
	if value := ctx.Query("status"); value != "" {
		filters["status"] = strings.Split(value, ",")
	}
	if value := ctx.Query("product_category"); value != "" {
		filters["product_category"] = strings.Split(value, ",")
	}
	if value := ctx.Query("stock_location"); value != "" {
		filters["stock_location"] = []string{value}
	}
	if value := ctx.Query("supplier"); value != "" {
		filters["supplier"] = []string{value}
	}
	if value := ctx.Query("min_price"); value != "" {
		filters["min_price"] = []string{value}
	}
	if value := ctx.Query("max_price"); value != "" {
		filters["max_price"] = []string{value}
	}

	// Pagination
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	// Fetch data
	products, err := c.Service.GetFilteredProducts(filters, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// GenerateProductPDF godoc
// @Summary Generate a PDF report of products
// @Description Generates a formatted PDF file containing product details
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce application/pdf
// @Success 200 {file} pdf
// @Failure 500 {object} map[string]string
// @Router /api/products/pdf [get]
func (c *ProductController) GenerateProductPDF(ctx *gin.Context) {
	products, err := c.Service.GetFilteredProducts(nil, 100, 0) // Retrieve up to 100 products
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Create a new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "B", 16)
	pdf.AddPage()

	// Title
	pdf.Cell(190, 10, "Product Report")
	pdf.Ln(12)

	// Table Headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Ref", "Name", "Category", "Price", "Status", "Stock"}
	colWidths := []float64{35, 50, 35, 20, 25, 25}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Table Data
	pdf.SetFont("Arial", "", 10)
	for _, product := range products {
		data := []string{
			product.ProductReference,
			product.ProductName,
			product.ProductCategory,
			fmt.Sprintf("$%.2f", product.Price),
			product.Status,
			fmt.Sprintf("%d", product.Quantity),
		}

		for i, text := range data {
			pdf.CellFormat(colWidths[i], 10, text, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Save PDF
	filePath := "reports/product_report.pdf"
	err = pdf.OutputFileAndClose(filePath)
	log.Println(err)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	// Return the generated PDF file
	ctx.File(filePath)
}

// Predefined product locations (For demonstration, replace with a database lookup)
var productLocations = map[string]haversine.Coord{
	"New York": {Lat: 40.7128, Lon: -74.0060},
	"London":   {Lat: 51.5074, Lon: -0.1278},
	"Tokyo":    {Lat: 35.6895, Lon: 139.6917},
	"Berlin":   {Lat: 52.5200, Lon: 13.4050},
	"Sydney":   {Lat: -33.8688, Lon: 151.2093},
}

// Get user location from IP
func getUserLocation(ip string) (haversine.Coord, error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb") // Ensure you have the GeoLite2 database
	if err != nil {
		return haversine.Coord{}, err
	}
	defer db.Close()

	parsedIP := net.ParseIP(ip)
	record, err := db.City(parsedIP)
	if err != nil {
		return haversine.Coord{}, err
	}

	return haversine.Coord{
		Lat: record.Location.Latitude,
		Lon: record.Location.Longitude,
	}, nil
}

// GetDistanceAPI godoc
// @Summary Calculate distance between user and product location
// @Description Returns the distance in km between user's IP location and product's city
// @Tags Distance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product_city query string true "City where product was produced"
// @Param ip query string false "User's IP address (optional, auto-detect)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/distance [get]
func GetDistanceAPI(ctx *gin.Context) {
	productCity := ctx.Query("product_city")
	ip := ctx.ClientIP() // Get user IP (auto-detect)

	if productCity == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Product city is required"})
		return
	}

	// Get product location
	productCoord, exists := productLocations[productCity]
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product city"})
		return
	}

	// Get user location from IP
	userCoord, err := getUserLocation(ip)
	if err != nil {
		log.Println("Failed to get user location:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not determine user location"})
		return
	}

	// Calculate distance using Haversine formula
	_, km := haversine.Distance(productCoord, userCoord)

	ctx.JSON(http.StatusOK, gin.H{
		"user_ip":        ip,
		"user_location":  fmt.Sprintf("%.4f, %.4f", userCoord.Lat, userCoord.Lon),
		"product_city":   productCity,
		"product_coords": fmt.Sprintf("%.4f, %.4f", productCoord.Lat, productCoord.Lon),
		"distance_km":    fmt.Sprintf("%.2f km", km),
	})
}

// GetProductCategoryStatistics godoc
// @Summary Get product category distribution statistics
// @Description Returns product counts per category
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/statistics/products-per-category [get]
func (c *ProductController) GetProductCategoryStatistics(ctx *gin.Context) {
	stats, err := c.Service.GetProductCategoryStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetProductSupplierStatistics godoc
// @Summary Get product supplier distribution statistics
// @Description Returns product counts per supplier
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/statistics/products-per-supplier [get]
func (c *ProductController) GetProductSupplierStatistics(ctx *gin.Context) {
	stats, err := c.Service.GetProductSupplierStatistics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// DeleteProduct godoc
// @Summary Delete a single product by ID
// @Description Remove a product from the database using its ID
// @Tags Products
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := c.Service.DeleteProductByID(int(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// DeleteProductByReference godoc
// @Summary Delete a product by reference
// @Description Remove a product from the database using its product reference
// @Tags Products
// @Security BearerAuth
// @Param reference path string true "Product Reference"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/reference/{reference} [delete]
func (c *ProductController) DeleteProductByReference(ctx *gin.Context) {
	reference := ctx.Param("reference")

	if err := c.Service.DeleteProductByReference(reference); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// DeleteMultipleProducts godoc
// @Summary Delete multiple products
// @Description Remove multiple products by providing a list of IDs
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.DeleteMultipleProductsRequest true "List of Product IDs"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/products/bulk [delete]
func (c *ProductController) DeleteMultipleProducts(ctx *gin.Context) {
	var request models.DeleteMultipleProductsRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(request.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No product IDs provided"})
		return
	}

	if err := c.Service.DeleteMultipleProducts(request.IDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Products deleted successfully"})
}
