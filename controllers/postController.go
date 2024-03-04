package controllers

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type PostData struct {
	PostId uuid.UUID     `json:"id"`
	Text   string        `json:"text"`
	Files  []models.File `json:"files"`
	User   UserData      `json:"user"`
}

func GetPosts(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	var posts []models.Post

	offset := (page - 1) * 20

	result := initializers.DB.Offset(offset).Limit(20).Order("created_at desc").Find(&posts)
	if result.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var response []PostData
	for _, post := range posts {
		var user models.User
		result = initializers.DB.Preload("Avatar").First(&user, "id = ?", post.Author)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		respUser := PostData{
			PostId: post.ID,
			Text:   post.Text,
			Files:  post.Files,
			User: UserData{
				ID:          user.ID,
				Username:    user.Username,
				DisplayName: user.DisplayName,
				Avatar:      user.Avatar.FileName,
			},
		}
		response = append(response, respUser)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"posts": response,
	})
}

func GetPostById(c *gin.Context) {
	postId := c.Param("id")

	uniqueId, err := uuid.Parse(postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
		return
	}

	var postData PostData

	var post models.Post
	result := initializers.DB.Preload("Files").First(&post, "id = ?", uniqueId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	var user models.User
	result = initializers.DB.Preload("Avatar").First(&user, "id = ?", post.Author)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	postData = PostData{
		PostId: uniqueId,
		Text:   post.Text,
		Files:  post.Files,
		User: UserData{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar.FileName,
		},
	}

	c.JSON(http.StatusOK, postData)
}

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
