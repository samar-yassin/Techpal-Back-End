package routes

import (
	controller "CareerGuidance/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//incomingRoutes.Use(middleware.Authenticate())
	//incomingRoutes.GET("/users", controller.GetUsers())
	//incomingRoutes.GET("/users/:userid", controller.GetUser())
	incomingRoutes.GET("api/user", controller.User())

}
