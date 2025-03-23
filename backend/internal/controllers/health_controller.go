package controllers

import (
	"context"
	"ecommerce-test/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary API Health Check
// @Description Checks if the API and database are running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/health [get]
func HealthCheck(ctx *gin.Context) {
	// Check database connection
	err := config.DB.Ping(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Database connection failed",
		})
		return
	}

	// If everything is okay
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is running",
	})
}
