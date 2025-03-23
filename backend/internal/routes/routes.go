package routes

import (
	"ecommerce-test/internal/controllers"
	"ecommerce-test/internal/middleware"
	"ecommerce-test/internal/repository"
	"ecommerce-test/internal/services"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"

	_ "ecommerce-test/docs" // Import generated docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *pgx.Conn) *gin.Engine {
	r := gin.Default()

	// Health Check API
	r.GET("/api/health", controllers.HealthCheck)

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userService := services.UserService{Repo: userRepo}
	userController := controllers.UserController{Service: userService}

	// Initialize dependencies
	authService := services.AuthService{UserRepo: userRepo}
	authController := controllers.AuthController{AuthService: authService}

	// Initialize dependencies
	productRepo := repository.ProductRepository{}
	productService := services.ProductService{Repo: productRepo}
	productController := controllers.ProductController{Service: productService}

	// Public routes
	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)

	// Protected routes (JWT required)
	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	{
		// User
		protected.GET("/users", userController.GetUsers)

		// Product
		protected.POST("/products", productController.CreateProduct)
		protected.GET("/products/:product_reference", productController.GetProductByReference)
		protected.PUT("/products/:product_reference", productController.UpdateProduct) // Update product
		protected.GET("/products", productController.GetFilteredProducts)              // Dynamic filtering
		protected.GET("/products/pdf", productController.GenerateProductPDF)           // Generate PDF
		// Distance API
		protected.GET("/distance", controllers.GetDistanceAPI)

		// Product category statistics route
		protected.GET("/statistics/products-per-category", productController.GetProductCategoryStatistics)
		// Product supplier statistics route
		protected.GET("/statistics/products-per-supplier", productController.GetProductSupplierStatistics)

	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Service is running"})
	})

	return r
}
