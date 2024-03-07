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
	api.Api.GET("/validate", middleware.RequireAuth, controllers.ValidateToken)

	// Users
	api.Api.GET("/users", middleware.RequireAuth, controllers.GetUsers)
	api.Api.GET("/users/:id", middleware.RequireAuth, controllers.GetUserById)
	api.Api.PUT("/users", middleware.RequireAuth, controllers.UpdateUser)

	// Files
	if currentMode == "debug" {
		api.Router.Static("assets", initializers.DebugBasePath)
	} else {
		api.Router.Static("assets", initializers.ReleaseBasePath)
	}

	// Posts
	api.Api.POST("/posts", middleware.RequireAuth, controllers.CreatePost)
	api.Api.GET("/posts", middleware.RequireAuth, controllers.GetPosts)
	api.Api.GET("/posts/:id", middleware.RequireAuth, controllers.GetPostById)

	// Starting
	err := api.Router.Run()
	if err != nil {
		log.Fatal("Router could not be created!\n" + err.Error())
	}
}
