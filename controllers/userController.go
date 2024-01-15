package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
)

// GetUsers Returns all users
func GetUsers(c *gin.Context) {
	var users []models.User

	result := initializers.DB.Find(&users)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"users": users,
	})
}

func GetUserById(c *gin.Context) {
	userId := c.Param("id")

	uniqueId, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "id = ?", uniqueId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
