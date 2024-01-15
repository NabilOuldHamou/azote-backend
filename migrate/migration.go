package main

import (
	"azote-backend/initializers"
	"azote-backend/models"
	"log"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	log.Println("Migrating models to db...")

	err := initializers.DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Post{})
	if err != nil {
		log.Fatalf("Automatic migration has failed : %v", err)
	}

	log.Println("Migration successful.")
}
