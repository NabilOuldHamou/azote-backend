package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetUsers Returns all users
func GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	type userResponse struct {
		ID          uuid.UUID `json:"id"`
		Username    string    `json:"username"`
		DisplayName string    `json:"display_name"`
		Avatar      string    `json:"profile_picture"`
	}

	var users []models.User

	offset := (page - 1) * 50

	result := initializers.DB.Offset(offset).Limit(50).Find(&users)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var response []userResponse
	for _, user := range users {
		respUser := userResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar.FileName,
		}
		response = append(response, respUser)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"users": response,
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
