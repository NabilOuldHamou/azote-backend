package main

import (
	"azote-backend/api"
	"azote-backend/controllers"
	"azote-backend/initializers"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func init() {
	initializers.LoadEnv()
	initializers.CreateAssetsFolder()
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
	api.Router.Static("assets", "./assets/images")

	// Posts

	// Starting
	err := api.Router.Run()
	if err != nil {
		log.Fatal("Router could not be created!\n" + err.Error())
	}
}
