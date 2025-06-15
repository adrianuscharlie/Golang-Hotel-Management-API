package main

import (
	"hms-backend/config"
	"hms-backend/model"
	"hms-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&model.Room{},
		&model.RoomType{},
		&model.Booking{},
		&model.Guest{},
		&model.Transaction{})
	r := gin.Default()
	routes.RegisterRoutes(r, config.DB)
	r.Run(":4000")
}
