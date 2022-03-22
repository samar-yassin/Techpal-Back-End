package main

import (
	routes "CareerGuidance/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	//	integration with select_file.html to test upload cv func
	/*
		router.LoadHTMLGlob("template/*")
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "select_file.html", gin.H{})
		})
	*/
	router.Run(":" + port)
}
