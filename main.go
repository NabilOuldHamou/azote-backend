package main

import (
	"azote-backend/api"
	"azote-backend/controllers"
	"azote-backend/initializers"
	"azote-backend/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var currentMode string

func init() {
	initializers.LoadEnv()
	currentMode = os.Getenv("GIN_MODE")
	if currentMode == "debug" {
		initializers.CreateAssetsFolder(initializers.DebugBasePath)
	} else {
		initializers.CreateAssetsFolder(initializers.ReleaseBasePath)
	}
	initializers.ConnectToDB()
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	api.CreateRouter()

	// Auth
	api.Api.POST("/signup", controllers.Signup)
	api.Api.POST("/login", controllers.Login)

	// Users
	api.Api.GET("/users", controllers.GetUsers)
	api.Api.GET("/users/:id", controllers.GetUserById)

	// Files
	if currentMode == "debug" {
		api.Router.Static("assets", initializers.DebugBasePath+"/images")
	} else {
		api.Router.Static("assets", initializers.ReleaseBasePath+"/images")
	}

	// Posts
	api.Api.POST("/posts", middleware.RequireAuth, controllers.CreatePost)

	// Starting
	err := api.Router.Run()
	if err != nil {
		log.Fatal("Router could not be created!\n" + err.Error())
	}
}
