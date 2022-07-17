package main

import (
	routes "CareerGuidance/routes"
	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())

	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.AdminRoutes(router)
	routes.MentorRoutes(router)
	//	integration with select_file.html to test upload cv func
	router.LoadHTMLGlob("template/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	})
	router.Run(":" + port)
}
