package controllers

import (
	"ecommerce-test/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.UserService
}

// GetUsers godoc
// @Summary Get all users
// @Description Returns a list of users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} map[string]string
// @Router /api/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.Service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Fetches a single user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.Service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
