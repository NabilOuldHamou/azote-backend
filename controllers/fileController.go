package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	splitName := strings.Split(file.Filename, ".")
	path := "assets/images/" + uuid.New().String() + "." + splitName[len(splitName)-1]
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown error"})
		return
	}

	f := models.File{
		Location: path,
	}

	result := initializers.DB.Create(&f)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An internal server error occurred",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"recipeId": "feur",
	})
}
