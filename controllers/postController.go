package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreatePost(c *gin.Context) {
	form, _ := c.MultipartForm()
	uploadedFiles := form.File["files"]
	var body struct {
		Text string
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
	post := models.Post{
		Text:   body.Text,
		Files:  files,
		Author: id,
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {

		for _, f := range files {
			var filePath string
			if os.Getenv("GIN_MODE") == "release" {
				filePath = filepath.Join(initializers.ReleaseBasePath, f.FileName)
			} else {
				filePath = filepath.Join(initializers.DebugBasePath, f.FileName)
			}
			if os.Remove(filePath) != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}

			if initializers.DB.Delete(&f).Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unknown DB error",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"post": post,
	})

}

func uploadFiles(c *gin.Context, uploadedFiles []*multipart.FileHeader) ([]models.File, error) {
	var files []models.File

	for _, file := range uploadedFiles {

		osFile, _ := file.Open()
		buffer := make([]byte, 261)
		_, err := osFile.Read(buffer)
		if err != nil {
			return nil, err
		}

		if !filetype.IsImage(buffer) && !filetype.IsAudio(buffer) && !filetype.IsVideo(buffer) {
			return nil, errors.New("file is not an image, audio or video")
		}

		splitName := strings.Split(file.Filename, ".")
		newName := uuid.New().String() + "." + splitName[len(splitName)-1]
		var path string
		if os.Getenv("GIN_MODE") == "release" {
			path = initializers.ReleaseBasePath + newName
		} else {
			path = initializers.DebugBasePath + newName
		}

		if err := c.SaveUploadedFile(file, path); err != nil {
			return nil, err
		}

		mFile := models.File{
			FileName: newName,
		}
		result := initializers.DB.Create(&mFile)
		if result.Error != nil {
			return nil, errors.New("database error")
		}

		files = append(files, mFile)
	}

	return files, nil
}
