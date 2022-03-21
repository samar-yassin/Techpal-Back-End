package main

import (
	routes "CareerGuidance/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("Port")
	if port == "" {
		log.Fatal("No port Provided")
	}
	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Run(":" + port)
}
