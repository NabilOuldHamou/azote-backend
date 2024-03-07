package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserData struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Avatar      string    `json:"profile_picture"`
}

// GetUsers Returns all users
func GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	username := c.Query("username")

	var users []models.User

	offset := (page - 1) * 20

	var result *gorm.DB

	if username != "" {
		result = initializers.DB.Preload("Avatar").Where("username LIKE ?", username+"%").Offset(offset).Limit(20).Find(&users)
	} else {
		result = initializers.DB.Preload("Avatar").Offset(offset).Limit(20).Find(&users)
	}

	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var response []UserData
	for _, user := range users {
		respUser := UserData{
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
	result := initializers.DB.Preload("Avatar").Preload("Posts.Files").First(&user, "id = ?", uniqueId)

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

func UpdateUser(c *gin.Context) {
	form, _ := c.MultipartForm()
	uploadedFiles := form.File["Avatar"]
	var body struct {
		Email       string
		Password    string
		DisplayName string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(uploadedFiles) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No files were uploaded.",
		})
		return
	}

	files, err := uploadFiles(c, uploadedFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	currentUserId := c.GetString("userId")
	id, _ := uuid.Parse(currentUserId)

	initializers.DB.Model(&files[0]).Update("user_id", id.String())

	var user models.User
	result := initializers.DB.Preload("Avatar").First(&user, "id = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}

		initializers.DB.Model(&user).Updates(models.User{
			Email:       body.Email,
			Password:    string(hashedPassword),
			DisplayName: body.DisplayName,
		})
	} else {
		initializers.DB.Model(&user).Updates(models.User{
			Email:       body.Email,
			DisplayName: body.DisplayName,
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"user": user,
	})

}
